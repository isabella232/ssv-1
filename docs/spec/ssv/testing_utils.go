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
var testingCommittee = []*types.Operator{
	{
		OperatorID: 1,
		PubKey:     testingSK1.GetPublicKey().Serialize(),
	},
	{
		OperatorID: 2,
		PubKey:     testingSK2.GetPublicKey().Serialize(),
	},
	{
		OperatorID: 3,
		PubKey:     testingSK3.GetPublicKey().Serialize(),
	},
	{
		OperatorID: 4,
		PubKey:     testingSK4.GetPublicKey().Serialize(),
	},
}
var testingShare = &types.Share{
	OperatorID:    1,
	PubKey:        testingSK1.GetPublicKey().Serialize(),
	DomainType:    types.PrimusTestnet,
	Quorum:        3,
	PartialQuorum: 2,
	Committee:     testingCommittee,
}
var testingSK1 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("453d8351a0880da481659260660be34fa8ac7f89fdd75af278c027474fcee285")
	return ret
}()
var testingSK2 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("6bfbd05d17d29c28d62e5c27f9646bc87b012ffd9d60575264def4f0291140dd")
	return ret
}()
var testingSK3 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("600a1293a76ab9dd141cd59c3705b268dbc2590c6bd5b66039ff27bac87e759f")
	return ret
}()
var testingSK4 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("538f613374f85d8d880b30031cab9c4ef82d2c6100850c7d4b097abc5a898790")
	return ret
}()
var testingQBFTConfig = &qbft.Config{
	Signer:    newTestingKeyManager(),
	SigningPK: testingSK1.GetPublicKey().Serialize(),
	Domain:    types.PrimusTestnet,
	ValueCheck: func(data []byte) error {
		return nil
	},
	Storage: &testingStorage{},
	Network: &testingNetwork{},
}

func newTestingValidator() *Validator {
	signer := newTestingKeyManager()
	dutyRunner := newTestingDutyRunner()
	return &Validator{
		valCheck: func(data []byte) error {
			return nil
		},
		signer:  signer,
		share:   testingShare,
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

func newTestingQBFTController(identifier []byte) *qbft.Controller {
	ret := qbft.NewController(
		[]byte{1, 2, 3, 4},
		testingShare,
		types.PrimusTestnet,
		newTestingKeyManager(),
		func(data []byte) error {
			return nil
		},
		newTestingStorage(),
		&testingNetwork{},
	)
	ret.Height = 0
	ret.Identifier = identifier
	ret.Domain = types.PrimusTestnet
	return ret
}

func newTestingQBFTInstance() *qbft.Instance {
	ret := qbft.NewInstance(testingQBFTConfig, nil, nil)
	ret.State = &qbft.State{
		Share:                           testingShare,
		ID:                              []byte{1, 2, 3, 4},
		Round:                           qbft.FirstRound,
		Height:                          qbft.FirstHeight,
		LastPreparedRound:               qbft.NoRound,
		LastPreparedValue:               nil,
		ProposalAcceptedForCurrentRound: nil,
	}
	ret.ProposeContainer = &qbft.MsgContainer{
		Msgs: map[qbft.Round][]*qbft.SignedMessage{},
	}
	ret.PrepareContainer = &qbft.MsgContainer{
		Msgs: map[qbft.Round][]*qbft.SignedMessage{},
	}
	ret.CommitContainer = &qbft.MsgContainer{
		Msgs: map[qbft.Round][]*qbft.SignedMessage{},
	}
	ret.RoundChangeContainer = &qbft.MsgContainer{
		Msgs: map[qbft.Round][]*qbft.SignedMessage{},
	}
	return ret
}

func newTestingDutyRunner() *DutyRunner {
	return &DutyRunner{
		BeaconRoleType: beacon.RoleTypeAttester,
		ValidatorPK:    testingValidatorPK[:],
		Share:          testingShare,
		QBFTController: newTestingQBFTController([]byte{1, 2, 3, 4}),
		storage:        newTestingStorage(),
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

type testingNetwork struct {
}

func (net *testingNetwork) Broadcast(message types.Encoder) error {
	return nil
}

func (net *testingNetwork) BroadcastDecided(msg types.Encoder) error {
	return nil
}

func newTestingKeyManager() types.KeyManager {
	ret := types.NewSSVKeyManager(types.PrimusTestnet)
	ret.AddShare(testingSK1)
	ret.AddShare(testingSK2)
	ret.AddShare(testingSK3)
	ret.AddShare(testingSK4)
	return ret
}
