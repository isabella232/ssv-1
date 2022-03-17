package ssv

import (
	"encoding/hex"
	"encoding/json"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/go-ssz"
)

// DutyRunners is a map of duty runners mapped by msg id hex.
type DutyRunners map[beacon.RoleType]*DutyRunner

// DutyRunnerForMsgID returns a DutyRunner from the provided msg ID, or nil if not found
func (ci DutyRunners) DutyRunnerForMsgID(msgID types.MessageID) *DutyRunner {
	role := msgID.GetRoleType()
	return ci[role]
}

type Network interface {
	Broadcast(message types.Encoder) error
}

// Storage is a persistent storage for the SSV
type Storage interface {
}

type BeaconNode interface {
	// GetAttestationData returns attestation data by the given slot and committee index
	GetAttestationData(slot phase0.Slot, committeeIndex phase0.CommitteeIndex) (*phase0.AttestationData, error)
	// SubmitAttestation submit the attestation to the node
	SubmitAttestation(attestation *phase0.Attestation) error
}

// consensusData holds all relevant duty and data Decided on by consensus
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
