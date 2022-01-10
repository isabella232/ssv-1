package qbft

import (
	"github.com/bloxapp/ssv/docs/spec/utils"
	"github.com/bloxapp/ssv/utils/threadsafe"
	"github.com/pkg/errors"
	"sync"
)

// Instance is a single QBFT instance that starts with a Start call (including a value).
// Every new msg the ProcessMsg function needs to be called
type Instance struct {
	state      State
	network    Network
	valueCheck ValueCheck

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

func (i *Instance) Start(value []byte) {
	i.startOnce.Do(func() {
		i.startValue = value

		// propose if this node is the proposer
		if proposer(i.state) == i.state.GetConfig().GetID() {
			proposal := createProposal(i.state, i.startValue)
			if err := i.state.GetConfig().GetNetwork().BroadcastSignedMessage(proposal); err != nil {
				// TODO log
			}
		}
	})
}

// ProcessMsg processes a new QBFT msg, returns non nil error on msg processing error
func (i *Instance) ProcessMsg(msg SignedMessage) error {
	res := i.processMsgF.Run(func() interface{} {
		switch msg.GetMessage().GetType() {
		case ProposalType:
			return uponProposal(i.state, msg, i.proposeContainer)
		case PrepareType:
			return uponPrepare(i.state, msg, i.prepareContainer, i.commitContainer)
		case CommitType:
			decided, decidedValue, err := uponCommit(i.state, msg, i.commitContainer)
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
		return res.(error)
	}
	return nil
}

// IsDecided returns true and a non nil byte slice of the decided value if decided.
func (i *Instance) IsDecided() (bool, []byte) {
	return i.decided.Get(), i.decidedValue.Get()
}
