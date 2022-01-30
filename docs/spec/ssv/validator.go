package ssv

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/types"
)

// Validator represents an SSV ETH consensus validator share assigned, coordinates duty execution and more.
// Every validator has a validatorID which is validator's public key.
// Each validator has multiple DutyRunners, for each duty type.
type Validator struct {
	network     Network
	beacon      beacon.Beacon
	storage     Storage
	share       Share
	id          ValidatorID
	dutyRunners DutyRunners
	signer      beacon.Signer
	valCheck    types.ValueCheck
}
