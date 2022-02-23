package types

// NodeID is a unique ID for the node, used to create shares and verify msgs
type NodeID uint64

// Node represents an SSV operator node
type Node struct {
	NodeID NodeID
	PubKey []byte
}

// GetPublicKey returns the public key with which the node is identified with
func (n *Node) GetPublicKey() []byte {
	return n.PubKey
}

// GetID returns the node's ID
func (n *Node) GetID() NodeID {
	return n.NodeID
}
