package ssv

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

func (v *Validator) processConsensusMsg(dutyRunner *DutyRunner, msg *qbft.SignedMessage) error {
	decided, decidedValueByts, err := dutyRunner.qbftController.ProcessMsg(msg)
	if err != nil {
		return errors.Wrap(err, "failed to process consensus msg")
	}

	/**
	Decided returns true only once so if it is true it must be for the current running instance
	*/
	if !decided {
		return nil
	}

	decidedValue := &consensusData{}
	if err := decidedValue.Decode(decidedValueByts); err != nil {
		return errors.Wrap(err, "failed to parse decided value to consensusData")
	}

	if err := v.checkDecidedValue(decidedValue); err != nil {
		return errors.Wrap(err, "decided value is invalid")
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

	if err := v.network.Broadcast(msgToBroadcast); err != nil {
		return errors.Wrap(err, "can't broadcast partial sig")
	}
	return nil
}

func (v *Validator) checkDecidedValue(decidedValue *consensusData) error {
	switch decidedValue.Duty.Type {
	case beacon.RoleTypeAttester:
		return v.valCheck.CheckAttestationData(decidedValue.AttestationData)
	default:
		return errors.New("can't validate unknown decided value")
	}
}
