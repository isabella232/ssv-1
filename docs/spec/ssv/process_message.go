package ssv

import (
	"github.com/bloxapp/ssv/docs/spec/network"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/pkg/errors"
)

func consensusMsgFromNetworkMsg(msg network.Message) (qbft.SignedMessage, error) {
	panic("implement")
}

func postConsensusSigMsgFromNetworkMsg(msg network.Message) (PostConsensusSigMessage, error) {
	panic("implement")
}

// ProcessMessage processes network message of all types
func (v *Validator) ProcessMessage(signedMsg network.SignedMessage) error {
	if err := v.validateMessage(signedMsg); err != nil {
		return errors.Wrap(err, "message invalid")
	}

	msg := signedMsg.GetMessage()

	dutyRunner := v.dutyRunners.DutyRunnerForMsgID(msg.GetID())
	if dutyRunner == nil {
		return errors.Errorf("could not get duty runner for msg ID %s", msg.GetID().String())
	}

	switch msg.GetType() {
	case network.Consensus:
		consensusMsg, err := consensusMsgFromNetworkMsg(msg)
		if err != nil {
			return errors.Wrap(err, "could not get consensus message from network message")
		}
		return v.processConsensusMsg(dutyRunner, consensusMsg)
	case network.PostConsensusSignature:
		sigMsg, err := postConsensusSigMsgFromNetworkMsg(msg)
		if err != nil {
			return errors.Wrap(err, "could not get post consensus message from network message")
		}
		return v.processPostConsensusSig(dutyRunner, sigMsg)
	case network.Sync:
		panic("")
	default:
		return errors.New("unknown msg")
	}
}

func (v *Validator) validateMessage(signedMsg network.SignedMessage) error {
	if !v.id.MessageIDBelongs(signedMsg.GetMessage().GetID()) {
		return errors.New("msg ID doesn't match validator ID")
	}

	dutyRunner := v.dutyRunners.DutyRunnerForMsgID(signedMsg.GetMessage().GetID())
	if dutyRunner == nil {
		return errors.Errorf("could not get duty runner for msg ID %s", signedMsg.GetMessage().GetID().String())
	}

	if !v.verifySignedMsgSignature(signedMsg) {
		return errors.New("failed to verify msg signature")
	}

	panic("implement")
}

func (v *Validator) verifySignedMsgSignature(signedMsg network.SignedMessage) bool {
	panic("implement")
}
