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

type IConfig interface {
	signing
	// GetValueCheck returns value check instance
	GetValueCheck() proposedValueCheck
	// GetNetwork returns a p2p Network instance
	GetNetwork() Network
	// GetTimer returns round timer
	GetTimer() Timer
}

type config struct {
	signer     types.SSVSigner
	signingPK  []byte
	domain     types.DomainType
	valueCheck proposedValueCheck
	storage    Storage
	network    Network
}

// GetSigner returns a signer instance
func (c *config) GetSigner() types.SSVSigner {
	return c.signer
}

// GetSigningPubKey returns the public key used to sign all QBFT messages
func (c *config) GetSigningPubKey() []byte {
	return c.signingPK
}

// GetSignatureDomainType returns the domain type used for signatures
func (c *config) GetSignatureDomainType() types.DomainType {
	return c.domain
}

// GetValueCheck returns value check instance
func (c *config) GetValueCheck() proposedValueCheck {
	return c.valueCheck
}

// GetNetwork returns a p2p Network instance
func (c *config) GetNetwork() Network {
	return c.network
}

// GetTimer returns round timer
func (c *config) GetTimer() Timer {
	return nil
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
