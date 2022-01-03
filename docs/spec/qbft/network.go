package qbft

type Network interface {
	BroadcastSignedMessage(msg SignedMessage) error
}
