package network

type MessageID []byte
type Type uint64

const (
	// Consensus are all QBFT consensus related messages
	Consensus Type = iota
	// Sync are all QBFT sync messages
	Sync
	// Signature are all partial signatures sent after consensus
	Signature
)

type Message interface {
	GetType() Type
	GetID() MessageID
	GetData() []byte
}
