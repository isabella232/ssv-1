package types

import "github.com/bloxapp/ssv/docs/spec/qbft"

// QBFTNetwork holds all QBFT relevant networking functions
type QBFTNetwork interface {
	BroadcastSignedMessage(msg qbft.SignedMessage) error
}

// SSVNetwork holds all SSV relevant networking functions
type SSVNetwork interface {
}
