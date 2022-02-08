package qbft

import "github.com/bloxapp/ssv/docs/spec/types"

type MessageType int

const (
	ProposalType MessageType = iota
	PrepareType
	CommitType
	RoundChangeType
)

type SignedMessage interface {
	types.MessageEncoder
	types.MessageSignature
	// GetMessage returns the message for which this signature is for
	GetMessage() *Message
}

type ProposalData interface {
	// GetData returns the data for which this QBFT instance tries to decide, can be any arbitrary data
	GetData() []byte
	// GetRoundChangeJustification returns a signed message with quorum as justification for the round change
	GetRoundChangeJustification() []SignedMessage
	// GetPrepareJustification returns a signed message with quorum as justification for a prepared round change
	GetPrepareJustification() []SignedMessage
}

type PrepareData interface {
	// GetData returns the data for which this QBFT instance tries to decide, can be any arbitrary data
	GetData() []byte
}

type CommitData interface {
	// GetData returns the data for which this QBFT instance tries to decide, can be any arbitrary data
	GetData() []byte
}

type RoundChangeData interface {
	GetPreparedValue() []byte
	GetPreparedRound() Round
	// GetNextProposalData returns NOT nil byte array if the signer is the next round's proposal.
	GetNextProposalData() []byte
	// GetRoundChangeJustification returns signed prepare messages for the last prepared state
	GetRoundChangeJustification() []SignedMessage
}

type Message struct {
	MsgType    MessageType
	Height     uint64 // QBFT instance height
	Round      Round  // QBFT round for which the msg is for
	Identifier []byte // instance identifier this msg belongs to
}

// GetProposalData returns proposal specific data
func (msg *Message) GetProposalData() ProposalData {
	panic("implement")
}

// GetPrepareData returns prepare specific data
func (msg *Message) GetPrepareData() PrepareData {
	panic("implement")
}

// GetCommitData returns commit specific data
func (msg *Message) GetCommitData() PrepareData {
	panic("implement")
}

// GetRoundChangeData returns round change specific data
func (msg *Message) GetRoundChangeData() RoundChangeData {
	panic("implement")
}

// Encode returns a msg encoded bytes or error
func (msg *Message) Encode() ([]byte, error) {
	panic("implement")
}

// Decode returns error if decoding failed
func (msg *Message) Decode(data []byte) error {
	panic("implement")
}

// GetRoot returns the root used for signing and verification
func (msg *Message) GetRoot() []byte {
	panic("implement")
}
