package ssv

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/network"
	"github.com/bloxapp/ssv/docs/spec/qbft"
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

	var decidedValue consensusInputData
	// TODO - parse decidedValueByts to decidedValue

	if err := dutyRunner.setAndSignDuty(decidedValue, v.signer); err != nil {
		return errors.Wrap(err, "failed to sign duty")
	}

	var msgToBroadcast network.Message
	switch dutyRunner.runningDuty.Type {
	case beacon.RoleTypeAttester:
		msgToBroadcast = v.createPartialSigMsg(dutyRunner.postConsensusState.attestation.Signature[:])
	default:
		return errors.Errorf("duty type %s unkwon", dutyRunner.runningDuty.Type.String())
	}

	if err := v.network.BroadcastMessage(msgToBroadcast); err != nil {
		return errors.Wrap(err, "can't broadcast partial sig")
	}
	return nil
}
