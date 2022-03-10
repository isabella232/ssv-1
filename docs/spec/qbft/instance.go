package qbft

import (
	"encoding/json"
	"github.com/bloxapp/ssv/docs/spec/utils"
	"github.com/pkg/errors"
	"sync"
)

type proposedValueCheck func(data []byte) error

// Instance is a single QBFT instance that starts with a Start call (including a value).
// Every new msg the ProcessMsg function needs to be called
type Instance struct {
	State      State
	config     IConfig
	valueCheck proposedValueCheck

	ProposeContainer     *MsgContainer
	PrepareContainer     *MsgContainer
	CommitContainer      *MsgContainer
	RoundChangeContainer *MsgContainer

	Decided      bool
	DecidedValue []byte
	processMsgF  *utils.ThreadSafeF
	startOnce    sync.Once
	StartValue   []byte
}

// Start is an interface implementation
func (i *Instance) Start(value []byte, height uint64) {
	i.startOnce.Do(func() {
		i.StartValue = value
		i.State.Round = FirstRound
		i.State.Height = height

		// propose if this node is the proposer
		if proposer(i.State, FirstRound) == i.State.Share.OperatorID {
			proposal, err := createProposal(i.State, i.config, i.StartValue, nil, nil)
			if err != nil {
				// TODO log
			}
			if err := i.config.GetNetwork().Broadcast(proposal); err != nil {
				// TODO - log
			}
		}
	})
}

// ProcessMsg processes a new QBFT msg, returns non nil error on msg processing error
func (i *Instance) ProcessMsg(msg *SignedMessage) (decided bool, decidedValue []byte, aggregatedCommit *SignedMessage, err error) {
	res := i.processMsgF.Run(func() interface{} {
		switch msg.Message.MsgType {
		case ProposalMsgType:
			return uponProposal(i.State, i.config, msg, i.ProposeContainer)
		case PrepareMsgType:
			return uponPrepare(i.State, i.config, msg, i.PrepareContainer, i.CommitContainer)
		case CommitMsgType:
			decided, decidedValue, aggregatedCommit, err = uponCommit(i.State, i.config, msg, i.CommitContainer)
			i.Decided = decided
			if decided {
				i.DecidedValue = decidedValue
			}

			// TODO - Roberto comment: we should send a Decided msg here
			return err
		case RoundChangeMsgType:
			return uponRoundChange(i.State, i.config, msg, i.RoundChangeContainer, i.valueCheck)
		default:
			return errors.New("signed message type not supported")
		}
	})
	if res != nil {
		return false, nil, nil, res.(error)
	}
	return i.Decided, i.DecidedValue, aggregatedCommit, nil
}

// IsDecided interface implementation
func (i *Instance) IsDecided() (bool, []byte) {
	return i.Decided, i.DecidedValue
}

// GetHeight interface implementation
func (i *Instance) GetHeight() uint64 {
	return i.State.Height
}

// Encode implementation
func (i *Instance) Encode() ([]byte, error) {
	return json.Marshal(i)
}

// Decode implementation
func (i *Instance) Decode(data []byte) error {
	return json.Unmarshal(data, &i)
}
