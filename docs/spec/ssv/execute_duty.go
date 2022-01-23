package ssv

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/pkg/errors"
)

// StartDuty starts a duty for the validator
func (v *Validator) StartDuty(duty *beacon.Duty) error {
	dutyRunner := v.dutyRunners[duty.Type]
	if dutyRunner == nil {
		return errors.Errorf("duty type %s not supported", duty.Type.String())
	}

	if err := dutyRunner.CanStartNewInstance(); err != nil {
		return errors.Wrap(err, "can't start new duty")
	}

	dutyRunner.resetForNewDuty()

	input := consensusInputData{}
	switch duty.Type {
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

	if err := dutyRunner.qbftController.StartNewInstance(input.Marshal()); err != nil {
		return errors.Wrap(err, "can't start new QBFT instance for duty")
	}

	return nil
}
