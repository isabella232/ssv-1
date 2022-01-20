package eth1

import (
	"crypto/rsa"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"time"
)

// Options configurations related to eth1
type Options struct {
	ETH1Addr              string        `yaml:"ETH1Addr" env:"ETH_1_ADDR" env-required:"true" env-description:"ETH1 node WebSocket address"`
	ETH1SyncOffset        string        `yaml:"ETH1SyncOffset" env:"ETH_1_SYNC_OFFSET" env-description:"block number to start the sync from"`
	ETH1ConnectionTimeout time.Duration `yaml:"ETH1ConnectionTimeout" env:"ETH_1_CONNECTION_TIMEOUT" env-default:"10s" env-description:"eth1 node connection timeout"`
	RegistryContractAddr  string        `yaml:"RegistryContractAddr" env:"REGISTRY_CONTRACT_ADDR_KEY" env-default:"0x9573C41F0Ed8B72f3bD6A9bA6E3e15426A0aa65B" env-description:"registry contract address"`
	RegistryContractABI   string        `yaml:"RegistryContractABI" env:"REGISTRY_CONTRACT_ABI" env-description:"registry contract abi json file"`
	CleanRegistryData     bool          `yaml:"CleanRegistryData" env:"CLEAN_REGISTRY_DATA" env-default:"false" env-description:"cleans registry contract data (validator shares) and forces re-sync"`
	AbiVersion            Version       `yaml:"AbiVersion" env:"ABI_VERSION" env-default:"0" env-description:"smart contract abi version (format)"`
}

// Event represents an eth1 event log in the system
type Event struct {
	// Log is the raw event log
	Log types.Log
	// Data is the parsed event
	Data interface{}
	// Decrypted indicates whether the share keys where decrypted during parsing
	Decrypted bool
}

// SyncResult used to pass stats and results of an events sync
type SyncResult struct {
	// Total is the amount of total events
	Total int
	// Successful is the amount of successfully parsed events
	Successful int
	// LastBlock is the last block where the sync applies to
	LastBlock uint64
}

// NewSyncResult creates a new SyncResult object
func NewSyncResult(n int, nsuccess int, lastBlock uint64) SyncResult {
	return SyncResult{
		Total:      n,
		Successful: nsuccess,
		LastBlock:  lastBlock,
	}
}

// ShareEncryptionKeyProvider is a function that returns the operator private key
type ShareEncryptionKeyProvider = func() (*rsa.PrivateKey, bool, error)

// Client represents the required interface for eth1 client
type Client interface {
	Subscribe(handler EventHandler) error
	Sync(fromBlock *big.Int, handler EventHandler) (SyncResult, error)
}

// EventHandler handles incoming event from registry contract
type EventHandler func(*Event) error
