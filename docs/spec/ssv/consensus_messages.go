package ssv

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

func (v *Validator) processConsensusMsg(dutyRunner *DutyRunner, msg qbft.SignedMessage) error {
	decided, decidedValueByts, err := dutyRunner.qbftController.ProcessMsg(msg)
	if err != nil {
		return errors.Wrap(err, "failed to process consensus msg")
	}

	if err := v.valCheck.Check(decidedValueByts); err != nil {
		return errors.Wrap(err, "post consensus value check fail")
	}

	if !decided {
		return nil
	}

	/**
	Decided returns true only once so if it is true it must be for the current running instance
	*/
	decidedValue := &consensusData{}
	if err := decidedValue.Decode(decidedValueByts); err != nil {
		return errors.Wrap(err, "failed to parse decided value to consensusData")
	}

	postConsensusMsg, err := dutyRunner.DecideRunningInstance(decidedValue, v.signer)
	if err != nil {
		return errors.Wrap(err, "failed to decide duty at runner")
	}

	data, err := postConsensusMsg.Encode()
	if err != nil {
		return errors.Wrap(err, "failed to encode PostConsensusSigMessage")
	}

	msgToBroadcast := &types.SSVMessage{
		MsgType: types.PostConsensusSignature,
		MsgID:   types.MessageIDForValidatorPKAndRole(v.share.GetValidatorPubKey(), dutyRunner.beaconRoleType),
		Data:    data,
	}

	if err := v.network.BroadcastMessage(msgToBroadcast); err != nil {
		return errors.Wrap(err, "can't broadcast partial sig")
	}
	return nil
}
