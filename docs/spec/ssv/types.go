package ssv

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/go-ssz"
)

// ValidatorID is an eth2 validator public key
type ValidatorID []byte

// MessageIDBelongs returns true if message ID belongs to validator
func (vid ValidatorID) MessageIDBelongs(msgID types.MessageID) bool {
	toMatch := msgID[:len(vid)]
	return bytes.Equal(vid, toMatch)
}

// DutyRunners is a map of duty runners mapped by msg id hex.
type DutyRunners map[beacon.RoleType]*DutyRunner

// DutyRunnerForMsgID returns a DutyRunner from the provided msg ID, or nil if not found
func (ci DutyRunners) DutyRunnerForMsgID(msgID types.MessageID) *DutyRunner {
	role := msgID.GetRoleType()
	return ci[role]
}

type Network interface {
	Broadcast(message *types.SSVMessage) error
}

// Storage is a persistent storage for the SSV
type Storage interface {
	// SaveHighestDecided saves the decided value as highest for a validator PK and role
	SaveHighestDecided(validatorPK []byte, role beacon.RoleType, decidedValue *consensusData) error
	// GetHighestDecided returns the saved decided value (highest) for a validator PK and role
	GetHighestDecided(validatorPK []byte, role beacon.RoleType) (*consensusData, error)
}

type Share struct {
	pubKey     []byte
	committee  []*types.Node
	quorum     uint64
	domainType types.DomainType
}

// GetValidatorPubKey returns the validator public key to which the share is associated with
func (share *Share) GetValidatorPubKey() []byte {
	return share.pubKey
}

func (share *Share) GetQBFTCommittee() []*types.Node {
	return share.committee
}

func (share *Share) GetQuorumCount() uint64 {
	return share.quorum
}

func (share *Share) GetDomainType() types.DomainType {
	return share.domainType
}

func (share *Share) Encode() ([]byte, error) {
	return json.Marshal(share)
}

func (share *Share) Decode(data []byte) error {
	return json.Unmarshal(data, share)
}

// consensusData holds all relevant duty and data decided on by consensus
type consensusData struct {
	Duty            *beacon.Duty
	AttestationData *phase0.AttestationData
	BlockData       *phase0.BeaconBlock
}

func (cid *consensusData) Encode() ([]byte, error) {
	m := make(map[string]string)
	if cid.Duty != nil {
		d, err := json.Marshal(cid.Duty)
		if err != nil {
			return nil, errors.Wrap(err, "duty marshaling failed")
		}
		m["duty"] = hex.EncodeToString(d)
	}

	if cid.AttestationData != nil {
		d, err := ssz.Marshal(cid.AttestationData)
		if err != nil {
			return nil, errors.Wrap(err, "attestation data marshaling failed")
		}
		m["att_data"] = hex.EncodeToString(d)
	}

	if cid.BlockData != nil {
		d, err := ssz.Marshal(cid.BlockData)
		if err != nil {
			return nil, errors.Wrap(err, "block data marshaling failed")
		}
		m["block_data"] = hex.EncodeToString(d)
	}
	return json.Marshal(m)
}

func (cid *consensusData) Decode(data []byte) error {
	m := make(map[string]string)
	if err := json.Unmarshal(data, &m); err != nil {
		return errors.Wrap(err, "could not unmarshal consensusData")
	}

	if val, ok := m["duty"]; ok {
		cid.Duty = &beacon.Duty{}
		d, err := hex.DecodeString(val)
		if err != nil {
			return errors.Wrap(err, "Duty decode string failed")
		}
		if err := json.Unmarshal(d, cid.Duty); err != nil {
			cid.Duty = nil
			return errors.Wrap(err, "could not unmarshal duty")
		}
	}

	if val, ok := m["att_data"]; ok {
		cid.AttestationData = &phase0.AttestationData{}
		d, err := hex.DecodeString(val)
		if err != nil {
			return errors.Wrap(err, "AttestationData decode string failed")
		}
		if err := ssz.Unmarshal(d, cid.AttestationData); err != nil {
			cid.AttestationData = nil
			return errors.Wrap(err, "could not unmarshal AttestationData")
		}
	}

	if val, ok := m["block_data"]; ok {
		cid.BlockData = &phase0.BeaconBlock{}
		d, err := hex.DecodeString(val)
		if err != nil {
			return errors.Wrap(err, "BlockData decode string failed")
		}
		if err := ssz.Unmarshal(d, cid.BlockData); err != nil {
			cid.BlockData = nil
			return errors.Wrap(err, "could not unmarshal BeaconBlock")
		}
	}
	return nil
}

// Signer interface holds all relevant signature signing funcs for SSV
type Signer interface {
	// SignPostConsensusMessage returns a PostConsensusSigMessage instance or error
	SignPostConsensusMessage(height uint64, root, pk []byte) (*PostConsensusSigMessage, error)
}
