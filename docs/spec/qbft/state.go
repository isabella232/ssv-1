package qbft

type NodeID uint64

type Node interface {
	// GetPublicKey returns the public key with which the node is identified with
	GetPublicKey() []byte
	// GetID returns the node's ID
	GetID() NodeID
}

type Config interface {
	// GetNodes returns instance nodes
	GetNodes() []Node
	// GetID returns the specific node's ID
	GetID() NodeID
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
	// GetConfig returns instance config
	GetConfig() Config
	// GetLastPreparedRound returns the highest prepared round, NoRound if not prepared
	GetLastPreparedRound() Round
	// GetLastPreparedValue returns the highest prepared value, nil if not prepared
	GetLastPreparedValue() []byte

	// CreateRoundChangeMessage returns a signed change round msg for the current state
	CreateRoundChangeMessage() SignedMessage
}
