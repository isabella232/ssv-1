package ssv

import (
	"bytes"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

// PostConsensusSigCollectionSlotTimeout represents for how many slots the post consensus sig collection has timeout for
const PostConsensusSigCollectionSlotTimeout spec.Slot = 32

// dutyExecutionState holds all the relevant progress the duty execution progress
type dutyExecutionState struct {
	height          uint64
	runningInstance qbft.IInstance

	decidedValue *consensusData

	signedAttestation *spec.Attestation
	signedProposal    *spec.SignedBeaconBlock

	collectedPartialSigs map[types.NodeID][]byte
	postConsensusSigRoot []byte
	// quorumCount is the number of min signatures needed for quorum
	quorumCount uint64

	finised bool
}

func (pcs *dutyExecutionState) AddPartialSig(sigMsg types.MessageSignature) error {
	if len(sigMsg.GetSigners()) != 1 {
		return errors.New("PostConsensusMessage has != 1 Signers")
	}

	if pcs.collectedPartialSigs[sigMsg.GetSigners()[0]] == nil {
		pcs.collectedPartialSigs[sigMsg.GetSigners()[0]] = sigMsg.GetSignature()
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

// DutyRunner is manages the execution of a duty from start to finish, it can only execute 1 duty at a time.
// Prev duty must finish before the next one can start.
type DutyRunner struct {
	beaconRoleType beacon.RoleType
	validatorPK    []byte
	storage        Storage
	// dutyExecutionState holds all relevant params for a full duty execution (consensus & post consensus)
	dutyExecutionState *dutyExecutionState
	qbftController     qbft.IController
	nodeID             types.NodeID
	share              *Share
}

// CanStartNewDuty returns nil if:
// - no running instance exists or
// - a QBFT instance decided and all post consensus sigs collectd or
// - a QBFT instance decided and 32 slots passed from decided duty
// else returns an error
// Will return error if not same role type
func (dr *DutyRunner) CanStartNewDuty(duty *beacon.Duty) error {
	if dr.dutyExecutionState == nil {
		return nil
	}

	if dr.beaconRoleType != duty.Type {
		return errors.New("duty runner role != duty.Type")
	}
	if !bytes.Equal(dr.validatorPK, duty.PubKey[:]) {
		return errors.New("duty runner validator pk != duty.PubKey")
	}

	if decided, _ := dr.dutyExecutionState.runningInstance.IsDecided(); !decided {
		return errors.New("consensus on duty is running")
	}

	if !dr.dutyExecutionState.HasPostConsensusSigQuorum() &&
		dr.dutyExecutionState.decidedValue.Duty.Slot+PostConsensusSigCollectionSlotTimeout >= duty.Slot { // if 32 slots (1 epoch) passed from running duty, start a new duty
		return errors.New("post consensus sig collection is running")
	}
	return nil
}

// StartNewInstance starts a new QBFT instance for value
func (dr *DutyRunner) StartNewInstance(value []byte) error {
	if value == nil {
		return errors.New("new instance value nil")
	}
	if err := dr.qbftController.StartNewInstance(value); err != nil {
		return errors.Wrap(err, "could not start new QBFT instance")
	}
	newInstance := dr.qbftController.InstanceForHeight(dr.qbftController.GetHeight())

	dr.dutyExecutionState = &dutyExecutionState{
		runningInstance: newInstance,
		height:          dr.qbftController.GetHeight(),
		quorumCount:     dr.share.GetQuorumCount(),
	}
	return dr.qbftController.StartNewInstance(value)
}

// PostConsensusStateForHeight returns a dutyExecutionState instance for a specific Height
func (dr *DutyRunner) PostConsensusStateForHeight(height uint64) *dutyExecutionState {
	if dr.dutyExecutionState != nil && dr.dutyExecutionState.runningInstance.GetHeight() == height {
		return dr.dutyExecutionState
	}
	return nil
}

// DecideRunningInstance sets the decided duty and partially signs the decided data, returns a PostConsensusMessage to be broadcasted or error
func (dr *DutyRunner) DecideRunningInstance(decidedValue *consensusData, signer types.KeyManager) (*PostConsensusMessage, error) {
	ret := &PostConsensusMessage{
		Height:  dr.dutyExecutionState.height,
		Signers: []types.NodeID{dr.nodeID},
	}
	switch dr.beaconRoleType {
	case beacon.RoleTypeAttester:
		signedAttestation, r, err := signer.SignAttestation(decidedValue.AttestationData, decidedValue.Duty, decidedValue.Duty.PubKey[:])
		if err != nil {
			return nil, errors.Wrap(err, "failed to sign attestation")
		}

		dr.dutyExecutionState.decidedValue = decidedValue
		dr.dutyExecutionState.signedAttestation = signedAttestation
		dr.dutyExecutionState.postConsensusSigRoot = ensureRoot(r)
		dr.dutyExecutionState.collectedPartialSigs = map[types.NodeID][]byte{}

		ret.DutySigningRoot = dr.dutyExecutionState.postConsensusSigRoot
		ret.DutySignature = dr.dutyExecutionState.signedAttestation.Signature[:]

		return ret, nil
	default:
		return nil, errors.Errorf("unknown duty %s", decidedValue.Duty.Type.String())
	}
}

// ensureRoot ensures that DutySigningRoot will have sufficient allocated memory
// otherwise we get panic from bls:
// github.com/herumi/bls-eth-go-binary/bls.(*Sign).VerifyByte:738
func ensureRoot(root []byte) []byte {
	n := len(root)
	if n == 0 {
		n = 1
	}
	tmp := make([]byte, n)
	copy(tmp[:], root[:])
	return tmp[:]
}
