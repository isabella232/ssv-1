package ssv

import "github.com/bloxapp/ssv/beacon"

// Storage is a persistent storage for the SSV
type Storage interface {
	// SaveHighestDecided saves the decided value as highest for a validator PK and role
	SaveHighestDecided(validatorPK []byte, role beacon.RoleType, decidedValue consensusInputData) error
	// GetHighestDecided returns the saved decided value (highest) for a validator PK and role
	GetHighestDecided(validatorPK []byte, role beacon.RoleType) (consensusInputData, error)
}
