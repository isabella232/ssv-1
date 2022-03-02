package ssv

import (
	"encoding/hex"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/utils/threshold"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

var testDuty = &beacon.Duty{
	Type:                    beacon.RoleTypeAttester,
	PubKey:                  spec.BLSPubKey{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8},
	Slot:                    12,
	ValidatorIndex:          1,
	CommitteeIndex:          22,
	CommitteesAtSlot:        36,
	ValidatorCommitteeIndex: 11,
}
var testAttData = &spec.AttestationData{
	Slot:            1,
	Index:           3,
	BeaconBlockRoot: spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
	Source: &spec.Checkpoint{
		Epoch: 1,
		Root:  spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
	},
	Target: &spec.Checkpoint{
		Epoch: 2,
		Root:  spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
	},
}
var testConsensusData = &consensusData{
	Duty:            testDuty,
	AttestationData: testAttData,
}
var testConsensusDataByts, _ = testConsensusData.Encode()
var testingValidatorPK = spec.BLSPubKey{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4}
var committee = []*types.Operator{
	{
		OperatorID: 1,
	},
	{
		OperatorID: 2,
	},
	{
		OperatorID: 3,
	},
	{
		OperatorID: 4,
	},
}

func newTestingValidator() *Validator {
	signer := newTestingKeyManager()
	return &Validator{
		valCheck: &types.BeaconDataCheck{KeyManager: signer},
		signer:   signer,
		share: &types.Share{
			PubKey:     testingValidatorPK[:],
			Committee:  committee,
			Quorum:     3,
			DomainType: types.PrimusTestnet,
		},
		network: &testingNetwork{},
	}
}

func newTestingDutyExecutionState() *DutyExecutionState {
	return &DutyExecutionState{
		CollectedPartialSigs: make(map[types.OperatorID][]byte),
	}
}

type testingQBFTController struct {
	instances  map[uint64]*testingQBFTInstance
	height     uint64
	identifier []byte

	failProcessMsg     bool
	returnDecided      bool
	returnDecidedValue []byte
}

func newTestingQBFTController(identifier []byte) *testingQBFTController {
	return &testingQBFTController{
		identifier: identifier,
		height:     0,
		instances:  make(map[uint64]*testingQBFTInstance),
	}
}

// StartNewInstance will start a new QBFT instance, if can't will return error
func (tContr *testingQBFTController) StartNewInstance(value []byte) error {
	inst := newTestingQBFTInstance()
	tContr.height++
	inst.height = tContr.height
	tContr.instances[inst.height] = inst
	return nil
}

// ProcessMsg processes a new msg, returns true if decided, non nil byte slice if decided (decided value) and error
// decided returns just once per instance as true, following messages (for example additional commit msgs) will not return decided true
func (tContr *testingQBFTController) ProcessMsg(msg *qbft.SignedMessage) (bool, []byte, error) {
	if tContr.failProcessMsg {
		return false, nil, errors.New("failed process msg")
	}
	if tContr.returnDecided {
		return true, tContr.returnDecidedValue, nil
	}
	return false, nil, nil
}

// InstanceForHeight returns an instance for a specific Height, nil if not found
func (tContr *testingQBFTController) InstanceForHeight(height uint64) qbft.IInstance {
	if inst, found := tContr.instances[height]; found {
		return inst
	}
	return nil
}

// GetHeight returns the current running instance Height or, if not started, the last decided Height
func (tContr *testingQBFTController) GetHeight() uint64 {
	return tContr.height
}

// GetIdentifier returns QBFT identifier, used to identify messages
func (tContr *testingQBFTController) GetIdentifier() []byte {
	return tContr.identifier
}

type testingQBFTInstance struct {
	height  uint64
	decided bool
}

func newTestingQBFTInstance() *testingQBFTInstance {
	return &testingQBFTInstance{
		height:  1,
		decided: false,
	}
}

// Start implementation
func (tInstance *testingQBFTInstance) Start(value []byte, height uint64) {

}

// ProcessMsg implementation
func (tInstance *testingQBFTInstance) ProcessMsg(msg *qbft.SignedMessage) (decided bool, decidedValue []byte, aggCommit *qbft.SignedMessage, err error) {
	return false, nil, nil, nil

}

// IsDecided implementation
func (tInstance *testingQBFTInstance) IsDecided() (bool, []byte) {
	return tInstance.decided, nil
}

// GetHeight implementation
func (tInstance *testingQBFTInstance) GetHeight() uint64 {
	return tInstance.height
}

type testingStorage struct {
	storage map[string]map[beacon.RoleType]*consensusData
}

func newTestingStorage() *testingStorage {
	return &testingStorage{
		storage: make(map[string]map[beacon.RoleType]*consensusData),
	}
}

// SaveHighestDecided saves the decided value as highest for a validator PK and role
func (s *testingStorage) SaveHighestDecided(validatorPK []byte, role beacon.RoleType, decidedValue *consensusData) error {
	if s.storage[hex.EncodeToString(validatorPK)] == nil {
		s.storage[hex.EncodeToString(validatorPK)] = make(map[beacon.RoleType]*consensusData)
	}
	s.storage[hex.EncodeToString(validatorPK)][role] = decidedValue
	return nil
}

// GetHighestDecided returns the saved decided value (highest) for a validator PK and role
func (s *testingStorage) GetHighestDecided(validatorPK []byte, role beacon.RoleType) (*consensusData, error) {
	if s.storage[hex.EncodeToString(validatorPK)] == nil {
		return nil, errors.New("can't find validator PK")
	}
	if value, found := s.storage[hex.EncodeToString(validatorPK)][role]; found {
		return value, nil
	}
	return nil, errors.New("can't find role")
}

func newTestingDutyRunner() *DutyRunner {
	return &DutyRunner{
		State: &DutyRunnerState{
			BeaconRoleType: beacon.RoleTypeAttester,
			Share: &types.Share{
				PubKey:     testingValidatorPK[:],
				Committee:  committee,
				Quorum:     3,
				OperatorID: 1,
			},
			QBFTController: newTestingQBFTController([]byte{1, 2, 3, 4}),
		},
		storage: newTestingStorage(),
	}
}

type testingNetwork struct {
}

func (net *testingNetwork) Broadcast(message types.MessageEncoder) error {
	return nil
}

type testingKeyManager struct {
	sk     *bls.SecretKey
	domain types.DomainType
}

func newTestingKeyManager() types.KeyManager {
	threshold.Init()
	sk := &bls.SecretKey{}
	sk.SetByCSPRNG()
	return &testingKeyManager{
		sk:     sk,
		domain: types.PrimusTestnet,
	}
}

// IsAttestationSlashable returns error if attestation data is slashable
func (km *testingKeyManager) IsAttestationSlashable(data *spec.AttestationData) error {
	return nil
}

func (km *testingKeyManager) SignRoot(data types.Root, sigType types.SignatureType, pk []byte) (types.Signature, error) {
	root, _ := types.ComputeSigningRoot(data, types.ComputeSignatureDomain(km.domain, sigType))
	return km.sk.SignByte(root).Serialize(), nil
}

func (km *testingKeyManager) SignAttestation(data *spec.AttestationData, duty *beacon.Duty, pk []byte) (*spec.Attestation, []byte, error) {
	sig := spec.BLSSignature{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6}
	att := &spec.Attestation{
		Data:      data,
		Signature: sig,
	}
	return att, sig[:], nil
}

func (km *testingKeyManager) AddShare(shareKey *bls.SecretKey) error {
	return nil
}
