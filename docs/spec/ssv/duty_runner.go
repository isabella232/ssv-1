package ssv

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

// PostConsensusSigCollectionSlotTimeout represents for how many slots the post consensus sig collection has timeout for
const PostConsensusSigCollectionSlotTimeout spec.Slot = 32

// DutyRunner is manages the execution of a duty from start to finish, it can only execute 1 duty at a time.
// Prev duty must finish before the next one can start.
type DutyRunner struct {
	BeaconRoleType beacon.RoleType
	ValidatorPK    []byte
	Share          *types.Share
	// DutyExecutionState holds all relevant params for a full duty execution (consensus & post consensus)
	DutyExecutionState *DutyExecutionState
	QBFTController     *qbft.Controller
	storage            Storage
}

// CanStartNewDuty returns nil if:
// - no running instance exists or
// - a QBFT instance Decided and all post consensus sigs collectd or
// - a QBFT instance Decided and 32 slots passed from Decided duty
// else returns an error
// Will return error if not same role type
func (dr *DutyRunner) CanStartNewDuty(duty *beacon.Duty) error {
	if dr.DutyExecutionState == nil {
		return nil
	}

	if dr.BeaconRoleType != duty.Type {
		return errors.New("duty runner role != duty.MsgType")
	}
	if !bytes.Equal(dr.ValidatorPK, duty.PubKey[:]) {
		return errors.New("duty runner validator pk != duty.PubKey")
	}

	if decided, _ := dr.DutyExecutionState.RunningInstance.IsDecided(); !decided {
		return errors.New("consensus on duty is running")
	}

	if !dr.DutyExecutionState.HasPostConsensusSigQuorum() &&
		dr.DutyExecutionState.DecidedValue.Duty.Slot+PostConsensusSigCollectionSlotTimeout >= duty.Slot { // if 32 slots (1 epoch) passed from running duty, start a new duty
		return errors.New("post consensus sig collection is running")
	}
	return nil
}

// StartNewInstance starts a new QBFT instance for value
func (dr *DutyRunner) StartNewInstance(value []byte) error {
	if value == nil {
		return errors.New("new instance value nil")
	}
	if err := dr.QBFTController.StartNewInstance(value); err != nil {
		return errors.Wrap(err, "could not start new QBFT instance")
	}
	newInstance := dr.QBFTController.InstanceForHeight(dr.QBFTController.Height)

	dr.DutyExecutionState = &DutyExecutionState{
		RunningInstance: newInstance,
		Quorum:          dr.Share.Quorum,
	}
	return dr.QBFTController.StartNewInstance(value)
}

// PostConsensusStateForHeight returns a DutyExecutionState instance for a specific Height
func (dr *DutyRunner) PostConsensusStateForHeight(height uint64) *DutyExecutionState {
	if dr.DutyExecutionState != nil && dr.DutyExecutionState.RunningInstance.GetHeight() == height {
		return dr.DutyExecutionState
	}
	return nil
}

// DecideRunningInstance sets the Decided duty and partially signs the Decided data, returns a PostConsensusMessage to be broadcasted or error
func (dr *DutyRunner) DecideRunningInstance(decidedValue *consensusData, signer types.KeyManager) (*PostConsensusMessage, error) {
	ret := &PostConsensusMessage{
		Height:  dr.DutyExecutionState.RunningInstance.GetHeight(),
		Signers: []types.OperatorID{dr.Share.OperatorID},
	}
	switch dr.BeaconRoleType {
	case beacon.RoleTypeAttester:
		signedAttestation, r, err := signer.SignAttestation(decidedValue.AttestationData, decidedValue.Duty, decidedValue.Duty.PubKey[:])
		if err != nil {
			return nil, errors.Wrap(err, "failed to sign attestation")
		}

		dr.DutyExecutionState.DecidedValue = decidedValue
		dr.DutyExecutionState.SignedAttestation = signedAttestation
		dr.DutyExecutionState.PostConsensusSigRoot = ensureRoot(r)
		dr.DutyExecutionState.CollectedPartialSigs = map[types.OperatorID][]byte{}

		ret.DutySigningRoot = dr.DutyExecutionState.PostConsensusSigRoot
		ret.DutySignature = dr.DutyExecutionState.SignedAttestation.Signature[:]

		return ret, nil
	default:
		return nil, errors.Errorf("unknown duty %s", decidedValue.Duty.Type.String())
	}
}

// GetRoot returns the root used for signing and verification
func (dr *DutyRunner) GetRoot() ([]byte, error) {
	marshaledRoot, err := dr.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode DutyRunnerState")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret[:], nil
}

// Encode returns the encoded struct in bytes or error
func (dr *DutyRunner) Encode() ([]byte, error) {
	return json.Marshal(dr)
}

// Decode returns error if decoding failed
func (dr *DutyRunner) Decode(data []byte) error {
	return json.Unmarshal(data, &dr)
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
