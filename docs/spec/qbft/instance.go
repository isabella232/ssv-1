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
	state      State
	config     Config
	valueCheck proposedValueCheck

	proposeContainer     MsgContainer
	prepareContainer     MsgContainer
	commitContainer      MsgContainer
	roundChangeContainer MsgContainer

	Decided      bool
	decidedValue []byte
	processMsgF  *utils.ThreadSafeF
	startOnce    sync.Once
	startValue   []byte
}

// Start is an interface implementation
func (i *Instance) Start(value []byte, height uint64) {
	i.startOnce.Do(func() {
		i.startValue = value
		i.state.Round = FirstRound
		i.state.Height = height

		// propose if this node is the proposer
		if proposer(i.state, FirstRound) == i.state.Share.OperatorID {
			proposal, err := createProposal(i.state, i.config, i.startValue, nil, nil)
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
			return uponProposal(i.state, i.config, msg, i.proposeContainer)
		case PrepareMsgType:
			return uponPrepare(i.state, i.config, msg, i.prepareContainer, i.commitContainer)
		case CommitMsgType:
			decided, decidedValue, aggregatedCommit, err = uponCommit(i.state, i.config, msg, i.commitContainer)
			i.Decided = decided
			if decided {
				i.decidedValue = decidedValue
			}

			// TODO - Roberto comment: we should send a Decided msg here
			return err
		case RoundChangeMsgType:
			return uponRoundChange(i.state, i.config, msg, i.roundChangeContainer, i.valueCheck)
		default:
			return errors.New("signed message type not supported")
		}
	})
	if res != nil {
		return false, nil, nil, res.(error)
	}
	return i.Decided, i.decidedValue, aggregatedCommit, nil
}

// IsDecided interface implementation
func (i *Instance) IsDecided() (bool, []byte) {
	return i.Decided, i.decidedValue
}

// GetHeight interface implementation
func (i *Instance) GetHeight() uint64 {
	return i.state.Height
}

// Encode implementation
func (i *Instance) Encode() ([]byte, error) {
	m := make(map[string]interface{})

	byts, err := i.state.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode state")
	}
	m["state"] = byts

	byts, err = i.proposeContainer.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode proposeContainer")
	}
	m["propose_container"] = byts

	byts, err = i.prepareContainer.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode prepareContainer")
	}
	m["prepare_container"] = byts

	byts, err = i.commitContainer.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode commitContainer")
	}
	m["commit_container"] = byts

	byts, err = i.roundChangeContainer.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode roundChangeContainer")
	}
	m["round_change_container"] = byts

	m["Decided"] = i.Decided
	m["decided_value"] = i.decidedValue
	m["start_value"] = i.startValue
	return json.Marshal(m)
}

// Decode implementation
func (i *Instance) Decode(data []byte) error {
	panic("implement")
}
