package ssv

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/pkg/errors"
)

// ExecuteDuty executes a duty for the validator
func (v *Validator) ExecuteDuty(duty *beacon.Duty) error {
	dutyRunner := v.dutyRunners[duty.Type]
	if dutyRunner == nil {
		return errors.Errorf("duty type %s not supported", duty.Type.String())
	}

}
