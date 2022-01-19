package network

import (
	"encoding/hex"
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/qbft"
)

type MessageID []byte

func MessageIDForValidatorPKAndRole(pk []byte, role beacon.RoleType) MessageID {
	panic("implement")
}

func (msgID MessageID) String() string {
	return hex.EncodeToString(msgID)
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
	// GetID returns a unique msg ID that is used to identify to which validator should the message be sent for processing
	GetID() MessageID
	// GetData returns msg data
	GetData() []byte

	// GetSigningData returns signing data for the message to be used for signing the message
	GetSigningData() []byte
}

type SignedMessage interface {
	GetSigner() qbft.NodeID
	GetSig() []byte
	GetMessage() Message
}
