package qbft

import "github.com/bloxapp/ssv/docs/spec/types"

type MessageType int

const (
	ProposalType MessageType = iota
	PrepareType
	CommitType
	RoundChangeType
)

type Message interface {
	types.MessageEncoder
	types.MessageRoot

	// GetType returns the msg type
	GetType() MessageType
	// GetHeight returns the QBFT instance height (similar to block number)
	GetHeight() uint64
	// GetRound returns the round for which the msg was sent
	GetRound() Round
	// GetInstanceIdentifier returns the instance identifier this msg belongs to
	GetInstanceIdentifier() []byte

	// GetProposalData returns proposal specific data
	GetProposalData() ProposalData
	// GetPrepareData returns prepare specific data
	GetPrepareData() PrepareData
	// GetCommitData returns commit specific data
	GetCommitData() PrepareData
	// GetRoundChangeData returns round change specific data
	GetRoundChangeData() RoundChangeData
}

type SignedMessage interface {
	types.MessageEncoder
	types.MessageSignature
	// GetMessage returns the message for which this signature is for
	GetMessage() Message
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
