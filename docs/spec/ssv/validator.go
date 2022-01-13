package ssv

import (
	"bytes"
	"github.com/bloxapp/ssv/docs/spec/network"
)

type ValidatorID []byte

func (vid ValidatorID) MatchMessageID(msgID network.MessageID) bool {
	return bytes.Equal(vid, msgID)
}

type Share interface {
	// Sign will return a signature using the share secret key
	Sign(data []byte) ([]byte, error)
	// GetValidatorPubKey returns the validator public key to which the share is associated with
	GetValidatorPubKey() []byte
}

// Validator represents an SSV ETH consensus validator share assigned, coordinates duty execution and more.
type Validator struct {
	share Share
	id    ValidatorID
	consensusInstances
}
