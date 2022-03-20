package ssv

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

// StartDuty starts a duty for the validator
func (v *Validator) StartDuty(duty *beacon.Duty) error {
	dutyRunner := v.DutyRunners[duty.Type]
	if dutyRunner == nil {
		return errors.Errorf("duty type %s not supported", duty.Type.String())
	}

	if err := dutyRunner.CanStartNewDuty(duty); err != nil {
		return errors.Wrap(err, "can't start new duty")
	}

	input := &types.ConsensusData{}
	switch dutyRunner.BeaconRoleType {
	case beacon.RoleTypeAttester:
		attData, err := v.beacon.GetAttestationData(duty.Slot, duty.CommitteeIndex)
		if err != nil {
			return errors.Wrap(err, "failed to get attestation data")
		}
		input.Duty = duty
		input.AttestationData = attData
	default:
		return errors.Errorf("duty type %s unkwon", duty.Type.String())
	}

	byts, err := input.Encode()
	if err != nil {
		return errors.Wrap(err, "could not encode input")
	}

	// validate input
	if err := v.valCheck(byts); err != nil {
		return errors.Wrap(err, "StartDuty input data invalid")
	}

	if err := dutyRunner.StartNewInstance(byts); err != nil {
		return errors.Wrap(err, "can't start new duty runner instance for duty")
	}

	return nil
}
