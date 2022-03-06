package qbft

import (
	"encoding/json"
	"github.com/bloxapp/ssv/docs/spec/types"
)

type signing interface {
	// GetSigner returns a signer instance
	GetSigner() types.SSVSigner
	// GetSigningPubKey returns the public key used to sign all QBFT messages
	GetSigningPubKey() []byte
	// GetSignatureDomainType returns the domain type used for signatures
	GetSignatureDomainType() types.DomainType
}

type Config interface {
	signing
	// GetValueCheck returns value check instance
	GetValueCheck() proposedValueCheck
	// GetNetwork returns a p2p Network instance
	GetNetwork() Network
	// GetTimer returns round timer
	GetTimer() Timer
}

type State struct {
	Share                           *types.Share
	ID                              []byte // instance Identifier
	Round                           Round
	Height                          uint64
	LastPreparedRound               Round
	LastPreparedValue               []byte
	ProposalAcceptedForCurrentRound *SignedMessage
}

// GetRoot returns the state's deterministic root
func (s *State) GetRoot() []byte {
	panic("implement")
}

// Encode returns a msg encoded bytes or error
func (s *State) Encode() ([]byte, error) {
	return json.Marshal(s)
}

// Decode returns error if decoding failed
func (s *State) Decode(data []byte) error {
	return json.Unmarshal(data, &s)
}
