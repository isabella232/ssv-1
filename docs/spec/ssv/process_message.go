package ssv

import (
	"github.com/bloxapp/ssv/docs/spec/network"
	"github.com/pkg/errors"
)

// ProcessMessage processes network message of all types
func (v *Validator) ProcessMessage(msg network.Message) error {
	if err := v.validateMessage(msg); err != nil {
		return errors.Wrap(err, "message invalid")
	}

}

func (v *Validator) validateMessage(msg network.Message) error {
	panic("implement")
}
