package ssv

import (
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

// dutyExecutionState holds all the relevant progress the duty execution progress
type dutyExecutionState struct {
	height          uint64
	runningInstance qbft.IInstance

	decidedValue *consensusData

	signedAttestation *spec.Attestation
	signedProposal    *spec.SignedBeaconBlock

	collectedPartialSigs map[types.OperatorID][]byte
	postConsensusSigRoot []byte
	// quorumCount is the number of min signatures needed for quorum
	quorumCount uint64

	finised bool
}

func (pcs *dutyExecutionState) AddPartialSig(sigMsg *PostConsensusMessage) error {
	if len(sigMsg.Signers) != 1 {
		return errors.New("PostConsensusMessage has != 1 Signers")
	}

	if pcs.collectedPartialSigs[sigMsg.Signers[0]] == nil {
		pcs.collectedPartialSigs[sigMsg.Signers[0]] = sigMsg.DutySignature
	}
	return nil
}

// ReconstructAttestationSig aggregates collected partial sigs, reconstructs a valid sig and returns an attestation obj with the reconstructed sig
func (pcs *dutyExecutionState) ReconstructAttestationSig() (*spec.Attestation, error) {
	panic("implement")
}

func (pcs *dutyExecutionState) HasPostConsensusSigQuorum() bool {
	return uint64(len(pcs.collectedPartialSigs)) >= pcs.quorumCount
}

// SetFinished will mark this execution state as finished
func (pcs *dutyExecutionState) SetFinished() {
	pcs.finised = true
}

// IsFinished returns true if this execution state is finished
func (pcs *dutyExecutionState) IsFinished() bool {
	return pcs.finised
}
