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

type quorum interface {
	// HasQuorum returns true if at least 2f+1 signatures present
	// https://github.com/ConsenSys/qbft-formal-spec-and-verification/blob/main/dafny/spec/L1/node_auxiliary_functions.dfy#L259
	HasQuorum([]*SignedMessage) bool
	// HasPartialQuorum returns true if at least f+1 signatures present
	// https://github.com/ConsenSys/qbft-formal-spec-and-verification/blob/main/dafny/spec/L1/node_auxiliary_functions.dfy#L244
	HasPartialQuorum([]*SignedMessage) bool
}

type Config interface {
	signing
	quorum
	// GetOperators returns the operator committee for this instance
	GetOperators() []*types.Operator
	// GetID returns the specific node's ID
	GetID() types.OperatorID
	// GetValueCheck returns value check instance
	GetValueCheck() proposedValueCheck
	// GetNetwork returns a p2p Network instance
	GetNetwork() Network
	// GetTimer returns round timer
	GetTimer() Timer
}

type State struct {
	ID                              []byte // instance identifier
	Round                           Round
	Height                          uint64
	LastPreparedRound               Round
	LastPreparedValue               []byte
	ProposalAcceptedForCurrentRound *SignedMessage
}
