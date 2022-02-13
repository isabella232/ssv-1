package types

import (
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
)

// BeaconDataCheck validates beacon duty data (AttestationData, Block, etc.), including slashing protection
type BeaconDataCheck struct {
	KeyManager KeyManager
}

func (vc *BeaconDataCheck) CheckAttestationData(data *spec.AttestationData) error {
	// validate not slashable
	return vc.KeyManager.IsAttestationSlashable(data)
}
