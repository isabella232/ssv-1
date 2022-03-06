package ssv

import (
	"encoding/hex"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/utils/threshold"
	"github.com/herumi/bls-eth-go-binary/bls"
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
		PubKey:     qbft.TestingSK.GetPublicKey().Serialize(),
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
	dutyRunner := newTestingDutyRunner()
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
		dutyRunners: DutyRunners{
			beacon.RoleTypeAttester: dutyRunner,
		},
	}
}

func newTestingDutyExecutionState() *DutyExecutionState {
	return &DutyExecutionState{
		CollectedPartialSigs: make(map[types.OperatorID][]byte),
	}
}

//type testingQBFTController struct {
//	Instances  map[uint64]*testingQBFTInstance
//	Height     uint64
//	Identifier []byte
//
//	failProcessMsg     bool
//	returnDecided      bool
//	returnDecidedValue []byte
//}

func newTestingQBFTController(identifier []byte) *qbft.Controller {
	ret := qbft.NewController(
		newTestingKeyManager(),
		nil,
		newTestingStorage(),
		nil,
	)
	ret.Height = 0
	ret.Identifier = identifier
	ret.StartNewInstance([]byte{1, 2, 3, 4})
	return ret
}

func newTestingQBFTInstance() *qbft.Instance {
	return &qbft.Instance{
		State: qbft.State{
			Height: 1,
		},
		Decided: false,
	}
}

type testingStorage struct {
	storage map[string]*qbft.SignedMessage
}

func newTestingStorage() *testingStorage {
	return &testingStorage{
		storage: make(map[string]*qbft.SignedMessage),
	}
}

// SaveHighestDecided saves the Decided value as highest for a validator PK and role
func (s *testingStorage) SaveHighestDecided(signedMsg *qbft.SignedMessage) error {
	s.storage[hex.EncodeToString(signedMsg.Message.Identifier)] = signedMsg
	return nil
}

//// GetHighestDecided returns the saved Decided value (highest) for a validator PK and role
//func (s *testingStorage) GetHighestDecided(validatorPK []byte, role beacon.RoleType) (*consensusData, error) {
//	if s.storage[hex.EncodeToString(validatorPK)] == nil {
//		return nil, errors.New("can't find validator PK")
//	}
//	if value, found := s.storage[hex.EncodeToString(validatorPK)][role]; found {
//		return value, nil
//	}
//	return s.storage[hex.EncodeToString(signedMsg.Message.Identifier)], errors.New("can't find role")
//}

func newTestingDutyRunner() *DutyRunner {
	return &DutyRunner{
		State: &DutyRunnerState{
			BeaconRoleType: beacon.RoleTypeAttester,
			Share: types.Share{
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
