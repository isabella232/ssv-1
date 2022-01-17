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

	switch msg.GetType() {
	case network.Consensus:
		panic("")
	case network.PostConsensusSignature:
		panic("")
	case network.Sync:
		panic("")
	default:
		return errors.New("unknown msg")
	}
}

func (v *Validator) validateMessage(msg network.Message) error {
	if !v.id.MessageIDBelongs(msg.GetID()) {
		return errors.New("msg ID doesn't match validator ID")
	}

	panic("implement")
}
