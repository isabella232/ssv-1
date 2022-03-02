package qbft

import (
	"encoding/json"
	"github.com/bloxapp/ssv/docs/spec/types"
)

type MessageType int

const (
	ProposalMsgType MessageType = iota
	PrepareMsgType
	CommitMsgType
	RoundChangeMsgType
	DecidedMsgType
)

type ProposalData interface {
	// GetData returns the data for which this QBFT instance tries to decide, can be any arbitrary data
	GetData() []byte
	// GetRoundChangeJustification returns a signed message with quorum as justification for the round change
	GetRoundChangeJustification() []*SignedMessage
	// GetPrepareJustification returns a signed message with quorum as justification for a prepared round change
	GetPrepareJustification() []*SignedMessage
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
	GetRoundChangeJustification() []*SignedMessage
}

type Message struct {
	MsgType    MessageType
	Height     uint64 // QBFT instance height
	Round      Round  // QBFT round for which the msg is for
	Identifier []byte // instance identifier this msg belongs to
	Data       []byte
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
	return json.Marshal(msg)
}

// Decode returns error if decoding failed
func (msg *Message) Decode(data []byte) error {
	return json.Unmarshal(data, &msg)
}

// GetRoot returns the root used for signing and verification
func (msg *Message) GetRoot() ([]byte, error) {
	panic("implement")
}

// DeepCopy returns a new instance of Message, deep copied
func (msg *Message) DeepCopy() *Message {
	panic("implement")
}

type SignedMessage struct {
	Signature types.Signature
	Signers   []types.OperatorID
	Message   *Message // message for which this signature is for
}

func (signedMsg *SignedMessage) GetSignature() types.Signature {
	return signedMsg.Signature
}
func (signedMsg *SignedMessage) GetSigners() []types.OperatorID {
	return signedMsg.Signers
}

// MatchedSigners returns true if the provided signer ids are equal to GetSignerIds() without order significance
func (signedMsg *SignedMessage) MatchedSigners(ids []types.OperatorID) bool {
	panic("implement")
}

// Aggregate will aggregate the signed message if possible (unique signers, same digest, valid)
func (signedMsg *SignedMessage) Aggregate(sig types.MessageSignature) error {
	panic("implement")
}

// Encode returns a msg encoded bytes or error
func (signedMsg *SignedMessage) Encode() ([]byte, error) {
	return json.Marshal(signedMsg)
}

// Decode returns error if decoding failed
func (signedMsg *SignedMessage) Decode(data []byte) error {
	return json.Unmarshal(data, &signedMsg)
}

// GetRoot returns the root used for signing and verification
func (signedMsg *SignedMessage) GetRoot() ([]byte, error) {
	return signedMsg.Message.GetRoot()
}

// DeepCopy returns a new instance of SignedMessage, deep copied
func (signedMsg *SignedMessage) DeepCopy() *SignedMessage {
	panic("implement")
}
