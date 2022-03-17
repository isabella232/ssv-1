package spectest

import (
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/utils/threshold"
	"github.com/herumi/bls-eth-go-binary/bls"
)

var AllTests = []*SpecTest{
	happyFullFlow(),
}

type SpecTest struct {
	Name                    string
	DutyRunner              *ssv.DutyRunner
	Messages                []*types.SSVMessage
	PostDutyRunnerStateRoot string
	ExpectedError           error
}

var baseRunner = func() *ssv.DutyRunner {
	return ssv.NewDutyRunner(
		beacon.RoleTypeAttester,
		testingValidatorPK[:],
		testingShare,
		ssv.NewTestingQBFTController([]byte{1, 2, 3, 4}),
		ssv.NewTestingStorage(),
	)
}
var BaseValidator = func() *ssv.Validator {
	ret := ssv.NewValidator(
		ssv.NewTestingNetwork(),
		ssv.NewTestingBeaconNode(),
		ssv.NewTestingStorage(),
		testingShare,
		ssv.NewTestingKeyManager(),
		func(data []byte) error {
			return nil
		},
	)
	return ret
}

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
	OperatorID:      1,
	ValidatorPubKey: testingValidatorPK[:],
	SharePubKey:     testingSK1.GetPublicKey().Serialize(),
	DomainType:      types.PrimusTestnet,
	Quorum:          3,
	PartialQuorum:   2,
	Committee:       testingCommittee,
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

var ssvMsg = func(qbftMsg *qbft.SignedMessage, postMsg *ssv.PostConsensusMessage) *types.SSVMessage {
	var msgType types.MsgType
	var data []byte
	if qbftMsg != nil {
		msgType = types.SSVConsensusMsgType
		data, _ = qbftMsg.Encode()
	} else if postMsg != nil {
		msgType = types.SSVPostConsensusMsgType
		data, _ = postMsg.Encode()
	} else {
		panic("msg type undefined")
	}

	msgID := types.MessageIDForValidatorPKAndRole(testingValidatorPK[:], beacon.RoleTypeAttester)

	return &types.SSVMessage{
		MsgType: msgType,
		MsgID:   msgID,
		Data:    data,
	}
}
