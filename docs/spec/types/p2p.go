package types

import "github.com/bloxapp/ssv/docs/spec/qbft"

type P2P interface {
	BroadcastSignedMessage(msg qbft.SignedMessage) error
}
