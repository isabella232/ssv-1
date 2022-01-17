package network

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/qbft"
)

type MessageID []byte

func MessageIDForValidatorPKAndRole(pk []byte, role beacon.RoleType) MessageID {
	panic("implement")
}

type Type uint64

const (
	// Consensus are all QBFT consensus related messages
	Consensus Type = iota
	// Sync are all QBFT sync messages
	Sync
	// PostConsensusSignature are all partial signatures sent after consensus
	PostConsensusSignature
)

type Message interface {
	GetType() Type
	// GetID returns the a unique msg ID
	GetID() MessageID
	GetData() []byte

	GetSignedMessage() (qbft.SignedMessage, error)
}
