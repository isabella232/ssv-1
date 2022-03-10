package qbft

import (
	"encoding/hex"
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
	config     Config
	valueCheck proposedValueCheck

	proposeContainer     *MsgContainer
	prepareContainer     *MsgContainer
	commitContainer      *MsgContainer
	roundChangeContainer *MsgContainer

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
		i.State.Round = FirstRound
		i.State.Height = height

		// propose if this node is the proposer
		if proposer(i.State, FirstRound) == i.State.Share.OperatorID {
			proposal, err := createProposal(i.State, i.config, i.startValue, nil, nil)
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
			return uponProposal(i.State, i.config, msg, i.proposeContainer)
		case PrepareMsgType:
			return uponPrepare(i.State, i.config, msg, i.prepareContainer, i.commitContainer)
		case CommitMsgType:
			decided, decidedValue, aggregatedCommit, err = uponCommit(i.State, i.config, msg, i.commitContainer)
			i.Decided = decided
			if decided {
				i.decidedValue = decidedValue
			}

			// TODO - Roberto comment: we should send a Decided msg here
			return err
		case RoundChangeMsgType:
			return uponRoundChange(i.State, i.config, msg, i.roundChangeContainer, i.valueCheck)
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
	return i.State.Height
}

// Encode implementation
func (i *Instance) Encode() ([]byte, error) {
	m := make(map[string]interface{})

	byts, err := i.State.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode State")
	}
	m["State"] = hex.EncodeToString(byts)

	if i.proposeContainer != nil {
		byts, err = i.proposeContainer.Encode()
		if err != nil {
			return nil, errors.Wrap(err, "could not encode proposeContainer")
		}
		m["propose_container"] = hex.EncodeToString(byts)
	}

	if i.prepareContainer != nil {
		byts, err = i.prepareContainer.Encode()
		if err != nil {
			return nil, errors.Wrap(err, "could not encode prepareContainer")
		}
		m["prepare_container"] = hex.EncodeToString(byts)
	}

	if i.commitContainer != nil {
		byts, err = i.commitContainer.Encode()
		if err != nil {
			return nil, errors.Wrap(err, "could not encode commitContainer")
		}
		m["commit_container"] = hex.EncodeToString(byts)
	}

	if i.roundChangeContainer != nil {
		byts, err = i.roundChangeContainer.Encode()
		if err != nil {
			return nil, errors.Wrap(err, "could not encode roundChangeContainer")
		}
		m["round_change_container"] = hex.EncodeToString(byts)
	}

	m["Decided"] = i.Decided
	m["decided_value"] = hex.EncodeToString(i.decidedValue)
	m["start_value"] = hex.EncodeToString(i.startValue)
	return json.Marshal(m)
}

// Decode implementation
func (i *Instance) Decode(data []byte) error {
	m := make(map[string]interface{})
	if err := json.Unmarshal(data, &m); err != nil {
		return errors.Wrap(err, "could not unmarshal instance data")
	}

	if m["State"] != nil {
		i.State = State{}
		data, err := hex.DecodeString(m["State"].(string))
		if err != nil {
			return errors.New("could not decode instance state hex")
		}
		if err := i.State.Decode(data); err != nil {
			return errors.Wrap(err, "could not decode instance state")
		}
	}

	if m["propose_container"] != nil {
		i.proposeContainer = &MsgContainer{}
		data, err := hex.DecodeString(m["propose_container"].(string))
		if err != nil {
			return errors.New("could not decode propose_container hex")
		}
		if err := i.proposeContainer.Decode(data); err != nil {
			return errors.Wrap(err, "could not decode instance proposeContainer")
		}
	}

	if m["prepare_container"] != nil {
		i.prepareContainer = &MsgContainer{}
		data, err := hex.DecodeString(m["prepare_container"].(string))
		if err != nil {
			return errors.New("could not decode prepare_container hex")
		}
		if err := i.prepareContainer.Decode(data); err != nil {
			return errors.Wrap(err, "could not decode instance prepareContainer")
		}
	}

	if m["commit_container"] != nil {
		i.commitContainer = &MsgContainer{}
		data, err := hex.DecodeString(m["commit_container"].(string))
		if err != nil {
			return errors.New("could not decode commit_container hex")
		}
		if err := i.commitContainer.Decode(data); err != nil {
			return errors.Wrap(err, "could not decode instance commitContainer")
		}
	}

	if m["round_change_container"] != nil {
		i.roundChangeContainer = &MsgContainer{}
		data, err := hex.DecodeString(m["round_change_container"].(string))
		if err != nil {
			return errors.New("could not decode round_change_container hex")
		}
		if err := i.roundChangeContainer.Decode(data); err != nil {
			return errors.Wrap(err, "could not decode instance roundChangeContainer")
		}
	}

	if m["Decided"] != nil {
		i.Decided = m["Decided"].(bool)
	}

	if m["decided_value"] != nil {
		data, err := hex.DecodeString(m["decided_value"].(string))
		if err != nil {
			return errors.New("could not decode decided_value hex")
		}
		i.decidedValue = data
	}

	if m["start_value"] != nil {
		data, err := hex.DecodeString(m["start_value"].(string))
		if err != nil {
			return errors.New("could not decode start_value hex")
		}
		i.startValue = data
	}

	return nil
}
