package ssv

import (
	"encoding/json"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/go-ssz"
)

// DutyExecutionState holds all the relevant progress the duty execution progress
type DutyExecutionState struct {
	RunningInstance qbft.IInstance

	DecidedValue *consensusData

	SignedAttestation *spec.Attestation
	SignedProposal    *spec.SignedBeaconBlock

	CollectedPartialSigs map[types.OperatorID][]byte
	PostConsensusSigRoot []byte
	// Quorum is the number of min signatures needed for quorum
	Quorum uint64

	Finished bool
}

func (pcs *DutyExecutionState) AddPartialSig(sigMsg *PostConsensusMessage) error {
	if len(sigMsg.Signers) != 1 {
		return errors.New("PostConsensusMessage has != 1 Signers")
	}

	if pcs.CollectedPartialSigs[sigMsg.Signers[0]] == nil {
		pcs.CollectedPartialSigs[sigMsg.Signers[0]] = sigMsg.DutySignature
	}
	return nil
}

// ReconstructAttestationSig aggregates collected partial sigs, reconstructs a valid sig and returns an attestation obj with the reconstructed sig
func (pcs *DutyExecutionState) ReconstructAttestationSig() (*spec.Attestation, error) {
	panic("implement")
}

func (pcs *DutyExecutionState) HasPostConsensusSigQuorum() bool {
	return uint64(len(pcs.CollectedPartialSigs)) >= pcs.Quorum
}

// SetFinished will mark this execution state as finished
func (pcs *DutyExecutionState) SetFinished() {
	pcs.Finished = true
}

// IsFinished returns true if this execution state is finished
func (pcs *DutyExecutionState) IsFinished() bool {
	return pcs.Finished
}

// Encode returns the encoded struct in bytes or error
func (pcs *DutyExecutionState) Encode() ([]byte, error) {
	m := make(map[string]interface{})

	if pcs.RunningInstance != nil {
		byts, err := pcs.RunningInstance.Encode()
		if err != nil {
			return nil, errors.Wrap(err, "could not encode RunningInstance")
		}
		m["running_instance"] = byts
	}

	if pcs.DecidedValue != nil {
		byts, err := pcs.DecidedValue.Encode()
		if err != nil {
			return nil, errors.Wrap(err, "could not encode DecidedValue")
		}
		m["decided_value"] = byts
	}

	if pcs.SignedAttestation != nil {
		byts, err := ssz.Marshal(pcs.SignedAttestation)
		if err != nil {
			return nil, errors.Wrap(err, "could not encode SignedAttestation")
		}
		m["signed_att"] = byts
	}

	if pcs.SignedProposal != nil {
		byts, err := ssz.Marshal(pcs.SignedProposal)
		if err != nil {
			return nil, errors.Wrap(err, "could not encode SignedProposal")
		}
		m["signed_proposal"] = byts
	}

	m["collected_partial_sigs"] = pcs.CollectedPartialSigs
	m["post_consensus_root"] = pcs.PostConsensusSigRoot
	m["quorum"] = pcs.Quorum
	m["finished"] = pcs.Finished
	return json.Marshal(m)
}

// Decode returns error if decoding failed
func (pcs *DutyExecutionState) Decode(data []byte) error {
	panic("implement")
}
