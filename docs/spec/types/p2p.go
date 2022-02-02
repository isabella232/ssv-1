package types

// QBFTNetwork is a collection of funcs for the QBFT network
type QBFTNetwork interface {
	BroadcastSignedMessage(msg MessageEncoder) error
}

// SSVNetwork holds all SSV relevant networking functions
type SSVNetwork interface {
}
