package ssv

import "github.com/bloxapp/ssv/docs/spec/qbft"

type PostConsensusSigMessage interface {
	GetSignerID() qbft.NodeID
	GetSig() []byte
}
