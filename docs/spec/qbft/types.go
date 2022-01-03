package qbft

type Round uint64

const (
	NoRound    = 0 // NoRound represents a nil/ zero round
	FirstRound = 1 // FirstRound value is the first round in any QBFT instance start
)
