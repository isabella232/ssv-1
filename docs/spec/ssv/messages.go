package ssv

import (
	"github.com/bloxapp/ssv/docs/spec/types"
)

type PostConsensusSigMessage interface {
	types.MessageEncoder
	types.MessageSignature
	// GetHeight returns the QBFT for which this post consensus sig message belongs
	GetHeight() uint64
}
