package qbft

import "github.com/bloxapp/ssv/docs/spec/types"

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
	Share                           types.Share
	ID                              []byte // instance identifier
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
