package qbft

import "github.com/bloxapp/ssv/docs/spec/types"

type Round uint64

const (
	NoRound    = 0 // NoRound represents a nil/ zero round
	FirstRound = 1 // FirstRound value is the first round in any QBFT instance start
)

// Network is a collection of funcs for the QBFT network
type Network interface {
	Broadcast(msg types.MessageEncoder) error
}

type Storage interface {
}
