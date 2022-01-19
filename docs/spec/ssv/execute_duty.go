package ssv

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/beacon"
	"github.com/pkg/errors"
)

type consensusInputData interface {
	// Marshal parses the obj to bytes
	Marshal() []byte
	// UnMarshal takes bytes and populates the obj
	UnMarshal(data []byte) error

	GetDuty() (*beacon.Duty, error)
	GetAttestationData() (*phase0.AttestationData, error)
	GetBlockData() (*phase0.BeaconBlock, error)
}

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

	var input consensusInputData
	switch duty.Type {
	case beacon.RoleTypeAttester:
		_, err := v.beacon.GetAttestationData(duty.Slot, duty.CommitteeIndex)
		if err != nil {
			return errors.Wrap(err, "failed to get attestation data")
		}

		// TODO set input with duty and attestation data
	default:
		return errors.Errorf("duty type %s unkwon", duty.Type.String())
	}

	if err := dutyRunner.qbftController.StartNewInstance(input.Marshal()); err != nil {
		return errors.Wrap(err, "can't start new QBFT instance for duty")
	}

	return nil
}
