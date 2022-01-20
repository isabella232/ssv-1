package goeth

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/bloxapp/ssv/eth1"
	"github.com/bloxapp/ssv/monitoring/metrics"
	"github.com/bloxapp/ssv/utils/tasks"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
	"log"
	"math/big"
	"strings"
	"time"
)

const (
	healthCheckTimeout        = 10 * time.Second
	blocksInBatch      uint64 = 100000
)

type eth1NodeStatus int32

var (
	metricsEth1NodeStatus = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ssv:eth1:node_status",
		Help: "Status of the connected eth1 node",
	})
	statusUnknown eth1NodeStatus = 0
	statusSyncing eth1NodeStatus = 1
	statusOK      eth1NodeStatus = 2
)

func init() {
	if err := prometheus.Register(metricsEth1NodeStatus); err != nil {
		log.Println("could not register prometheus collector")
	}
}

// ClientOptions are the options for the client
type ClientOptions struct {
	Ctx                        context.Context
	Logger                     *zap.Logger
	NodeAddr                   string
	RegistryContractAddr       string
	ContractABI                string
	ConnectionTimeout          time.Duration
	ShareEncryptionKeyProvider eth1.ShareEncryptionKeyProvider

	AbiVersion eth1.Version
}

// eth1Client is the internal implementation of Client
type eth1Client struct {
	ctx    context.Context
	conn   *ethclient.Client
	logger *zap.Logger

	shareEncryptionKeyProvider eth1.ShareEncryptionKeyProvider

	nodeAddr             string
	registryContractAddr common.Address
	contractABI          string
	connectionTimeout    time.Duration

	abiVersion eth1.Version
}

// verifies that the client implements HealthCheckAgent
var _ metrics.HealthCheckAgent = &eth1Client{}

// NewEth1Client creates a new instance
func NewEth1Client(opts ClientOptions) (eth1.Client, error) {
	logger := opts.Logger.With(zap.String("component", "eth1GoETH"),
		zap.String("contract address", opts.RegistryContractAddr))
	logger.Info("creating a new eth1 client", zap.String("address", opts.NodeAddr))

	ec := eth1Client{
		ctx:                        opts.Ctx,
		logger:                     logger,
		shareEncryptionKeyProvider: opts.ShareEncryptionKeyProvider,
		nodeAddr:                   opts.NodeAddr,
		registryContractAddr:       common.HexToAddress(opts.RegistryContractAddr),
		connectionTimeout:          opts.ConnectionTimeout,
		abiVersion:                 opts.AbiVersion,
	}

	if err := ec.connect(); err != nil {
		logger.Error("could not connect to eth1 client", zap.Error(err))
		return nil, err
	}

	return &ec, nil
}

// Subscribe subscribes to events from the contract
func (ec *eth1Client) Subscribe(handler eth1.EventHandler) error {
	contractAbi, err := abi.JSON(strings.NewReader(ec.contractABI))
	if err != nil {
		return errors.Wrap(err, "failed to parse ABI interface")
	}
	abiParser := eth1.NewParser(ec.logger, ec.abiVersion)

	in, err := ec.subscribeAndListen()
	if err != nil {
		return errors.Wrap(err, "could not subscribe to events")
	}
	// read events from the results channel, invoke the given handler with successfully parsed
	for vLog := range in {
		if err = ec.processEventLog(vLog, contractAbi, abiParser, handler); err != nil {
			ec.logger.Warn("failed to process event log", zap.Error(err))
		}
	}
	return nil
}

// Sync reads events history
func (ec *eth1Client) Sync(fromBlock *big.Int, handler eth1.EventHandler) (eth1.SyncResult, error) {
	from := fromBlock.Uint64()
	ec.logger.Debug("syncing smart contract events", zap.Uint64("fromBlock", from))

	res, err := ec.syncEvents(fromBlock, handler)
	if err != nil {
		return res, err
	}
	ec.logger.Debug("finished syncing registry contract",
		zap.Int("total events", res.Total), zap.Int("successfully parsed", res.Successful),
		zap.Uint64("last block", res.LastBlock), zap.Uint64("from block", from))
	return res, err
}

// HealthCheck provides health status of eth1 node
func (ec *eth1Client) HealthCheck() []string {
	if ec.conn == nil {
		return []string{"not connected to eth1 node"}
	}
	ctx, cancel := context.WithTimeout(ec.ctx, healthCheckTimeout)
	defer cancel()
	sp, err := ec.conn.SyncProgress(ctx)
	if err != nil {
		metricsEth1NodeStatus.Set(float64(statusUnknown))
		return []string{"could not get eth1 node sync progress"}
	}
	if sp != nil {
		metricsEth1NodeStatus.Set(float64(statusSyncing))
		return []string{fmt.Sprintf("eth1 node is currently syncing: starting=%d, current=%d, highest=%d",
			sp.StartingBlock, sp.CurrentBlock, sp.HighestBlock)}
	}
	// eth1 node is connected and synced
	metricsEth1NodeStatus.Set(float64(statusOK))
	return []string{}
}

// connect connects to eth1 client
func (ec *eth1Client) connect() error {
	if ec.conn != nil {
		ec.conn.Close()
	}
	// create an IPC based RPC connection to a remote node
	ec.logger.Info("dialing eth1 node...")
	ctx, cancel := context.WithTimeout(context.Background(), ec.connectionTimeout)
	defer cancel()
	conn, err := ethclient.DialContext(ctx, ec.nodeAddr)
	if err != nil {
		ec.logger.Error("could not connect to the eth1 client", zap.Error(err))
		return err
	}
	ec.logger.Info("successfully connected to eth1 goETH")
	ec.conn = conn
	return nil
}

// reconnect tries to reconnect multiple times with an exponent interval
func (ec *eth1Client) reconnect() {
	limit := 64 * time.Second
	tasks.ExecWithInterval(func(lastTick time.Duration) (stop bool, cont bool) {
		ec.logger.Info("reconnecting to eth1 node")
		if err := ec.connect(); err != nil {
			// continue until reaching to limit, and then panic as eth1 connection is required
			if lastTick >= limit {
				ec.logger.Panic("failed to reconnect to eth1 node", zap.Error(err))
			} else {
				ec.logger.Warn("could not reconnect to eth1 node, still trying", zap.Error(err))
			}
			return false, false
		}
		return true, false
	}, 1*time.Second, limit+(1*time.Second))

	ec.logger.Debug("managed to reconnect to eth1 node")
}

// subscribeAndListen subscribes to get event logs, notifies new events on the returned channel
func (ec *eth1Client) subscribeAndListen() (chan types.Log, error) {
	sub, logs, err := ec.subscribeEventLogs()
	if err != nil {
		return nil, errors.Wrap(err, "could not subscribe to events")
	}
	// notifies to a buffered channel, so we won't bottleneck the incoming events channel
	in := make(chan types.Log, 256)
	go func() {
		defer close(in)
		var err error
		//var lastLog types.Log
		for {
		listenToSubscription:
			for {
				select {
				case err = <-sub.Err():
					break listenToSubscription
				case vLog := <-logs:
					//lastLog = vLog
					in <- vLog
				case <-ec.ctx.Done():
					return
				}
			}
			if err != nil {
				// if failed to read from subscription -> re-subscribe
				ec.logger.Warn("failed to read logs from subscription", zap.Error(err))
				if sub, logs, err = ec.reSubscribeEventLogs(); err != nil {
					ec.logger.Panic("could not re-subscribe to contract events", zap.Error(err))
				}
				// TODO trigger a sync for the blocks that were created during reconnection
				//currentBlock, err := ec.conn.BlockNumber(ec.ctx)
				//if err != nil {
				//	ec.logger.Warn("could not get current block", zap.Error(err))
				//}
				//if currentBlock > lastLog.BlockNumber {
				//	_, err := ec.syncEvents(big.NewInt(int64(lastLog.BlockNumber)), func(event *eth1.Event) error {
				//		in <- event.Log
				//		return err
				//	})
				//	if err != nil {
				//		ec.logger.Panic("could not sync missing blocks", zap.Error(err))
				//	}
				//}
			}
		}
	}()

	return in, nil
}

// subscribeEventLogs subscribes to contract's event logs
func (ec *eth1Client) subscribeEventLogs() (ethereum.Subscription, chan types.Log, error) {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{ec.registryContractAddr},
	}
	logs := make(chan types.Log)
	sub, err := ec.conn.SubscribeFilterLogs(ec.ctx, query, logs)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to subscribe to logs")
	}
	ec.logger.Debug("subscribed to contract events")

	return sub, logs, nil
}

// subscribeEventLogs subscribes to contract's event logs
func (ec *eth1Client) reSubscribeEventLogs() (sub ethereum.Subscription, logs chan types.Log, err error) {
	err = tasks.Retry(func() error {
		ec.reconnect()
		sub, logs, err = ec.subscribeEventLogs()
		if err != nil {
			ec.logger.Warn("could not subscribe to contract events after reconnection", zap.Error(err))
		}
		return err
	}, 3)
	return sub, logs, err
}

// syncEvents is syncing  events history of the given contract
func (ec *eth1Client) syncEvents(fromBlock *big.Int, handler eth1.EventHandler) (eth1.SyncResult, error) {
	var nSuccess, nTotal int
	var toBlock *big.Int
	currentBatchSize := int64(blocksInBatch)
	currentBlock, err := ec.conn.BlockNumber(ec.ctx)
	if err != nil {
		return eth1.SyncResult{}, errors.Wrap(err, "failed to get current block")
	}
	for {
		if currentBlock-fromBlock.Uint64() > uint64(currentBatchSize) {
			toBlock = big.NewInt(int64(fromBlock.Uint64()) + currentBatchSize)
		} else { // no more batches are required -> setting toBlock to nil
			toBlock = big.NewInt(int64(currentBlock))
		}
		logger := ec.logger.With(zap.Uint64("fromBlock", fromBlock.Uint64()),
			zap.Uint64("toBlock", toBlock.Uint64()))
		logs, err := ec.fetchEventLogs(fromBlock, toBlock)
		if err != nil {
			// in case request exceeded limit, try again with lower batch size
			// will stop after log(blocksInBatch) tries
			if isLimitExceededErr(err) {
				currentBatchSize /= 2
				ec.logger.Debug("using a lower batch size", zap.Int64("currentBatchSize", currentBatchSize))
				if currentBatchSize > 1 {
					continue
				}
			}
		}
		logger.Debug("got event logs", zap.Int("results", len(logs)))
		processed, err := ec.processEventLogs(logs, handler)
		if err != nil {
			return eth1.NewSyncResult(nTotal, nSuccess, toBlock.Uint64()), errors.Wrap(err, "failed to get events")
		}
		nSuccess += processed
		nTotal += len(logs)
		if toBlock.Uint64() == currentBlock { // finished
			break
		}
		// resetting batch size and from/current block
		currentBatchSize = int64(blocksInBatch)
		if currentBlock, err = ec.conn.BlockNumber(ec.ctx); err != nil {
			return eth1.NewSyncResult(nTotal, nSuccess, toBlock.Uint64()), errors.Wrap(err, "failed to get current block")
		}
		fromBlock = toBlock
	}

	return eth1.NewSyncResult(nTotal, nSuccess, toBlock.Uint64()), nil
}

func (ec *eth1Client) fetchEventLogs(fromBlock, toBlock *big.Int) ([]types.Log, error) {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{ec.registryContractAddr},
		FromBlock: fromBlock,
	}
	if toBlock != nil {
		query.ToBlock = toBlock
	}
	ec.logger.Debug("fetching event logs", zap.Any("query", query))
	return ec.conn.FilterLogs(ec.ctx, query)
}

func (ec *eth1Client) processEventLogs(logs []types.Log, handler eth1.EventHandler) (int, error) {
	contractAbi, err := abi.JSON(strings.NewReader(ec.contractABI))
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse ABI interface")
	}
	abiParser := eth1.NewParser(ec.logger, ec.abiVersion)

	nSuccess := 0

	for _, vLog := range logs {
		if err = ec.processEventLog(vLog, contractAbi, abiParser, handler); err != nil {
			ec.logger.Warn("failed to process event log during sync", zap.Error(err))
			continue
		}
		nSuccess++
	}

	return nSuccess, nil
}

//var valAddedCount int64
func (ec *eth1Client) processEventLog(vLog types.Log, contractAbi abi.ABI, abiParser eth1.AbiParser, handler eth1.EventHandler) error {
	if vLog.Removed {
		return nil
	}
	e, err := ec.parseEventLog(vLog, contractAbi, abiParser)
	if err != nil {
		return errors.Wrap(err, "failed to parse event during sync")
	}
	if e == nil {
		return nil
	}
	if err = handler(e); err != nil {
		return errors.Wrap(err, "failed to handle event during sync")
	}
	return nil
}

// parseEventLog parses the given event
func (ec *eth1Client) parseEventLog(vLog types.Log, contractAbi abi.ABI, abiParser eth1.AbiParser) (*eth1.Event, error) {
	eventType, err := contractAbi.EventByID(vLog.Topics[0])
	if err != nil { // unknown event -> ignored
		ec.logger.Warn("failed to find event type", zap.Error(err), zap.String("txHash", vLog.TxHash.Hex()))
		return nil, nil
	}

	shareEncryptionKey, found, err := ec.shareEncryptionKeyProvider()
	if !found {
		return nil, errors.New("failed to find operator private key")
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get operator private key")
	}

	switch eventName := eventType.Name; eventName {
	case "OperatorAdded":
		parsed, isDecrypted, err := abiParser.ParseOperatorAddedEvent(shareEncryptionKey, vLog.Data, vLog.Topics, contractAbi)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse OperatorAdded event")
		}
		return &eth1.Event{Log: vLog, Data: parsed, Decrypted: isDecrypted}, nil
	case "ValidatorAdded":
		parsed, isDecrypted, err := abiParser.ParseValidatorAddedEvent(shareEncryptionKey, vLog.Data, contractAbi)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse ValidatorAdded event")
		}
		//ec.logger.Debug("validator was parsed successfully",
		//	zap.String("pubKey", hex.EncodeToString(parsed.PublicKey)),
		//	zap.Uint64("blockNumber", vLog.BlockNumber),
		//	zap.Int64("valAddedCount", atomic.AddInt64(&valAddedCount, 1)))
		if isDecrypted {
			ec.logger.Debug("validator is assigned to this operator",
				zap.String("pubKey", hex.EncodeToString(parsed.PublicKey)))
		}
		return &eth1.Event{Log: vLog, Data: parsed, Decrypted: isDecrypted}, nil
	default:
		ec.logger.Debug("unknown contract event was received", zap.String("hash", vLog.TxHash.Hex()), zap.String("eventName", eventName))
	}
	return nil, nil
}

// isLimitExceededErr returns true if the given error is of type limit exceeded
func isLimitExceededErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "websocket: read limit exceeded")
}
