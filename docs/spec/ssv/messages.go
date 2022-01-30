package ssv

import "github.com/bloxapp/ssv/docs/spec/qbft"

type PostConsensusSigMessage interface {
	// GetHeight returns the QBFT for which this post consensus sig message belongs
	GetHeight() uint64
	GetSignerID() qbft.NodeID
	GetSig() []byte
}
