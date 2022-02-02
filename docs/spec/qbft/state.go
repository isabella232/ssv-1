package qbft

import "github.com/bloxapp/ssv/docs/spec/types"

type Config interface {
	// GetNodes returns instance nodes
	GetNodes() []types.Node
	// GetID returns the specific node's ID
	GetID() types.NodeID
	// HasQuorum returns true if at least 2f+1 signatures present
	HasQuorum([]SignedMessage) bool
	// HasPartialQuorum returns true if at least f+1 signatures present
	HasPartialQuorum([]SignedMessage) bool
	// GetValueCheck returns value check instance
	GetValueCheck() types.ValueCheck
	// GetNetwork returns a p2p Network instance
	GetP2PNetwork() types.QBFTNetwork
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
	GetProposalAcceptedForCurrentRound() SignedMessage
	// SetProposalAcceptedForCurrentRound sets the variable
	SetProposalAcceptedForCurrentRound(msg SignedMessage)
}
