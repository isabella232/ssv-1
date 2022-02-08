package qbft

import (
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/utils"
	"github.com/bloxapp/ssv/utils/threadsafe"
	"github.com/pkg/errors"
	"sync"
)

type IInstance interface {
	// Start will start the new instance with a specific value and height
	Start(value []byte, height uint64)
	// ProcessMsg will process a signed msg
	ProcessMsg(msg SignedMessage) (decided bool, decidedValue []byte, err error)
	// IsDecided will return true and a non-empty byte slice if instance decided.
	IsDecided() (bool, []byte)
	// GetHeight returns the instance's height
	GetHeight() uint64
}

// Instance is a single QBFT instance that starts with a Start call (including a value).
// Every new msg the ProcessMsg function needs to be called
type Instance struct {
	state      State
	p2pNetwork types.QBFTNetwork
	valueCheck types.ValueCheck

	proposeContainer     MsgContainer
	prepareContainer     MsgContainer
	commitContainer      MsgContainer
	roundChangeContainer MsgContainer

	decided      *threadsafe.SafeBool
	decidedValue *threadsafe.SafeBytes
	processMsgF  *utils.ThreadSafeF
	startOnce    sync.Once
	startValue   []byte
}

// Start is an interface implementation
func (i *Instance) Start(value []byte, height uint64) {
	i.startOnce.Do(func() {
		i.startValue = value
		i.state.SetRound(FirstRound)
		i.state.SetHeight(height)

		// propose if this node is the proposer
		if proposer(i.state, FirstRound) == i.state.GetConfig().GetID() {
			proposal := createProposal(i.state, i.startValue)
			if err := i.state.GetConfig().GetP2PNetwork().BroadcastSignedMessage(proposal); err != nil {
				// TODO - log
			}
		}
	})
}

// ProcessMsg processes a new QBFT msg, returns non nil error on msg processing error
func (i *Instance) ProcessMsg(msg SignedMessage) (decided bool, decidedValue []byte, err error) {
	res := i.processMsgF.Run(func() interface{} {
		switch msg.GetMessage().MsgType {
		case ProposalType:
			return uponProposal(i.state, msg, i.proposeContainer)
		case PrepareType:
			return uponPrepare(i.state, msg, i.prepareContainer, i.commitContainer)
		case CommitType:
			decided, decidedValue, err = uponCommit(i.state, msg, i.commitContainer)
			i.decided.Set(decided)
			if decided {
				i.decidedValue.Set(decidedValue)
			}
			return err
		case RoundChangeType:
			return uponRoundChange(i.state, msg, i.roundChangeContainer, i.valueCheck)
		default:
			return errors.New("signed message type not supported")
		}
	})
	if res != nil {
		return false, nil, res.(error)
	}
	return i.decided.Get(), i.decidedValue.Get(), nil
}

// IsDecided interface implementation
func (i *Instance) IsDecided() (bool, []byte) {
	return i.decided.Get(), i.decidedValue.Get()
}

// GetHeight interface implementation
func (i *Instance) GetHeight() uint64 {
	return i.state.GetHeight()
}
