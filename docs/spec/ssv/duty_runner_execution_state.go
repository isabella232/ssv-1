package ssv

import (
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

// DutyExecutionState holds all the relevant progress the duty execution progress
type DutyExecutionState struct {
	Height          uint64
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
