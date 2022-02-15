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
	// GetNodes returns instance nodes
	GetNodes() []*types.Node
	// GetID returns the specific node's ID
	GetID() types.NodeID
	// GetValueCheck returns value check instance
	GetValueCheck() proposedValueCheck
	// GetNetwork returns a p2p Network instance
	GetNetwork() Network
	// GetStorage returns a p2p Network instance
	GetStorage() Storage
	// GetTimer returns round timer
	GetTimer() Timer
}

type State interface {
	// GetInstanceIdentifier returns the instance identifier this msg belongs to
	GetInstanceIdentifier() []byte
	// GetRound returns the round for which the msg was sent
	GetRound() Round
	// SetRound will set round
	SetRound(newRound Round)
	// GetHeight returns the current height
	GetHeight() uint64
	// SetHeight sets the current height
	SetHeight(height uint64) uint64
	// GetConfig returns instance config
	GetConfig() Config
	// GetLastPreparedRound returns the latest prepared round, NoRound if not prepared
	GetLastPreparedRound() Round
	// SetLastPreparedRound sets the latest prepared round
	SetLastPreparedRound(round Round)
	// GetLastPreparedValue returns the latest prepared value, nil if not prepared
	GetLastPreparedValue() []byte
	// SetLastPreparedValue sets the latest prepared value
	SetLastPreparedValue(value []byte)

	// GetProposalAcceptedForCurrentRound returns the set proposal for current round, nil if not set
	GetProposalAcceptedForCurrentRound() *SignedMessage
	// SetProposalAcceptedForCurrentRound sets the variable
	SetProposalAcceptedForCurrentRound(msg *SignedMessage)
}
