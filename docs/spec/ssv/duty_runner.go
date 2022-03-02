package ssv

import (
	"bytes"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

// PostConsensusSigCollectionSlotTimeout represents for how many slots the post consensus sig collection has timeout for
const PostConsensusSigCollectionSlotTimeout spec.Slot = 32

// DutyRunner is manages the execution of a duty from start to finish, it can only execute 1 duty at a time.
// Prev duty must finish before the next one can start.
type DutyRunner struct {
	State   *DutyRunnerState
	storage Storage
}

// CanStartNewDuty returns nil if:
// - no running instance exists or
// - a QBFT instance decided and all post consensus sigs collectd or
// - a QBFT instance decided and 32 slots passed from decided duty
// else returns an error
// Will return error if not same role type
func (dr *DutyRunner) CanStartNewDuty(duty *beacon.Duty) error {
	if dr.State.DutyExecutionState == nil {
		return nil
	}

	if dr.State.BeaconRoleType != duty.Type {
		return errors.New("duty runner role != duty.MsgType")
	}
	if !bytes.Equal(dr.State.Share.PubKey, duty.PubKey[:]) {
		return errors.New("duty runner validator pk != duty.PubKey")
	}

	if decided, _ := dr.State.DutyExecutionState.runningInstance.IsDecided(); !decided {
		return errors.New("consensus on duty is running")
	}

	if !dr.State.DutyExecutionState.HasPostConsensusSigQuorum() &&
		dr.State.DutyExecutionState.decidedValue.Duty.Slot+PostConsensusSigCollectionSlotTimeout >= duty.Slot { // if 32 slots (1 epoch) passed from running duty, start a new duty
		return errors.New("post consensus sig collection is running")
	}
	return nil
}

// StartNewInstance starts a new QBFT instance for value
func (dr *DutyRunner) StartNewInstance(value []byte) error {
	if value == nil {
		return errors.New("new instance value nil")
	}
	if err := dr.State.QBFTController.StartNewInstance(value); err != nil {
		return errors.Wrap(err, "could not start new QBFT instance")
	}
	newInstance := dr.State.QBFTController.InstanceForHeight(dr.State.QBFTController.GetHeight())

	dr.State.DutyExecutionState = &dutyExecutionState{
		runningInstance: newInstance,
		height:          dr.State.QBFTController.GetHeight(),
		quorumCount:     dr.State.Share.Quorum,
	}
	return dr.State.QBFTController.StartNewInstance(value)
}

// PostConsensusStateForHeight returns a dutyExecutionState instance for a specific Height
func (dr *DutyRunner) PostConsensusStateForHeight(height uint64) *dutyExecutionState {
	if dr.State.DutyExecutionState != nil && dr.State.DutyExecutionState.runningInstance.GetHeight() == height {
		return dr.State.DutyExecutionState
	}
	return nil
}

// DecideRunningInstance sets the decided duty and partially signs the decided data, returns a PostConsensusMessage to be broadcasted or error
func (dr *DutyRunner) DecideRunningInstance(decidedValue *consensusData, signer types.KeyManager) (*PostConsensusMessage, error) {
	ret := &PostConsensusMessage{
		Height:  dr.State.DutyExecutionState.height,
		Signers: []types.OperatorID{dr.State.Share.OperatorID},
	}
	switch dr.State.BeaconRoleType {
	case beacon.RoleTypeAttester:
		signedAttestation, r, err := signer.SignAttestation(decidedValue.AttestationData, decidedValue.Duty, decidedValue.Duty.PubKey[:])
		if err != nil {
			return nil, errors.Wrap(err, "failed to sign attestation")
		}

		dr.State.DutyExecutionState.decidedValue = decidedValue
		dr.State.DutyExecutionState.signedAttestation = signedAttestation
		dr.State.DutyExecutionState.postConsensusSigRoot = ensureRoot(r)
		dr.State.DutyExecutionState.collectedPartialSigs = map[types.OperatorID][]byte{}

		ret.DutySigningRoot = dr.State.DutyExecutionState.postConsensusSigRoot
		ret.DutySignature = dr.State.DutyExecutionState.signedAttestation.Signature[:]

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
