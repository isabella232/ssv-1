package ssv

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

// ProcessMessage processes network Message of all types
func (v *Validator) ProcessMessage(msg types.SSVMessage) error {
	if err := v.validateMessage(msg); err != nil {
		return errors.Wrap(err, "Message invalid")
	}

	dutyRunner := v.DutyRunners.DutyRunnerForMsgID(msg.GetID())
	if dutyRunner == nil {
		return errors.Errorf("could not get duty runner for msg ID %s", msg.GetID().String())
	}

	switch msg.GetType() {
	case types.SSVConsensusMsgType:
		signedMsg := &qbft.SignedMessage{}
		if err := signedMsg.Decode(msg.GetData()); err != nil {
			return errors.Wrap(err, "could not get post consensus Message from network Message")
		}
		return v.processConsensusMsg(dutyRunner, signedMsg)
	case types.SSVPostConsensusMsgType:
		signedMsg := &SignedPostConsensusMessage{}
		if err := signedMsg.Decode(msg.GetData()); err != nil {
			return errors.Wrap(err, "could not get post consensus Message from network Message")
		}
		return v.processPostConsensusSig(dutyRunner, signedMsg)
	case types.SSVSyncMsgType:
		panic("implement")
	default:
		return errors.New("unknown msg")
	}
}

func (v *Validator) validateMessage(msg types.SSVMessage) error {
	if !v.share.ValidatorPubKey.MessageIDBelongs(msg.GetID()) {
		return errors.New("msg ID doesn't match validator ID")
	}

	dutyRunner := v.DutyRunners.DutyRunnerForMsgID(msg.GetID())
	if dutyRunner == nil {
		return errors.Errorf("could not get duty runner for msg ID %s", msg.GetID().String())
	}

	if msg.GetType() > 2 {
		return errors.New("msg type not supported")
	}

	if len(msg.GetData()) == 0 {
		return errors.New("msg data is invalid")
	}

	return nil
}
