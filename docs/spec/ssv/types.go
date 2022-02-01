package ssv

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/network"
	"github.com/bloxapp/ssv/docs/spec/qbft"
)

type ValidatorID []byte

// ValidatorIDFromPubKey returns ValidatorID from a validator public key
func ValidatorIDFromPubKey(pk []byte) ValidatorID {
	return pk
}

// MessageIDBelongs returns true if message ID belongs to validator
func (vid ValidatorID) MessageIDBelongs(msgID network.MessageID) bool {
	panic("implement")
}

// DutyRunners is a map of duty runners mapped by msg id hex.
type DutyRunners map[beacon.RoleType]*DutyRunner

// DutyRunnerForMsgID returns a DutyRunner from the provided msg ID, or nil if not found
func (ci DutyRunners) DutyRunnerForMsgID(msgID network.MessageID) *DutyRunner {
	panic("implement")
}

//// ControllerFromDuty returns a QBFT controller from the provided duty, nil if not found
//func (ci DutyRunners) ControllerFromDuty(duty *beacon.Duty) *qbft.Controller {
//	panic("implement")
//}

type Share interface {
	// Sign will return a signature using the share secret key
	Sign(data []byte) ([]byte, error)
	// GetValidatorPubKey returns the validator public key to which the share is associated with
	GetValidatorPubKey() []byte
	GetQBFTCommittee() []qbft.NodeID
	GetQuorumCount() uint64
}

type Network interface {
	BroadcastMessage(message network.Message) error
}

// consensusData holds all relevant duty and data decided on by consensus
type consensusData struct {
	Duty            *beacon.Duty
	AttestationData *phase0.AttestationData
	BlockData       *phase0.BeaconBlock
}

func (cid *consensusData) Marshal() []byte {
	panic("implement")
}

func (cid *consensusData) UnMarshal(data []byte) error {
	panic("implement")
}
