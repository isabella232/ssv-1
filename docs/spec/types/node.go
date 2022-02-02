package types

type NodeID uint64

// Node represents an SSV operator node
type Node interface {
	// GetPublicKey returns the public key with which the node is identified with
	GetPublicKey() []byte
	// GetID returns the node's ID
	GetID() NodeID
}
