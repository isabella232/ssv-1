package spectest

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/utils/threshold"
	"github.com/herumi/bls-eth-go-binary/bls"
)

var AllTests = []*SpecTest{
	happyFullFlow(),
}

type SpecTest struct {
	Name     string
	Pre      *qbft.Instance
	PostRoot string
	Messages []*qbft.SignedMessage
	err      error
}

var signMsg = func(sk *bls.SecretKey, id types.OperatorID, msg *qbft.Message) *qbft.SignedMessage {
	domain := types.PrimusTestnet
	sigType := types.QBFTSigType

	r, _ := types.ComputeSigningRoot(msg, types.ComputeSignatureDomain(domain, sigType))
	sig := sk.SignByte(r)

	return &qbft.SignedMessage{
		Message:   msg,
		Signers:   []types.OperatorID{id},
		Signature: sig.Serialize(),
	}
}
var proposalDataByts = func(data []byte, rcj, pj []*qbft.SignedMessage) []byte {
	d := &qbft.ProposalData{
		Data:                     data,
		RoundChangeJustification: rcj,
		PrepareJustification:     pj,
	}
	ret, _ := d.Encode()
	return ret
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
var baseInstance = func() *qbft.Instance {
	ret := qbft.NewInstance(testingConfig)
	ret.State = &qbft.State{
		Share: &types.Share{
			OperatorID:    1,
			PubKey:        testingSK1.GetPublicKey().Serialize(),
			DomainType:    types.PrimusTestnet,
			Quorum:        3,
			PartialQuorum: 2,
			Committee: []*types.Operator{
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
			},
		},
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

var testingConfig = &qbft.Config{
	Signer:    testingKeyManager(),
	SigningPK: testingSK1.GetPublicKey().Serialize(),
	Domain:    types.PrimusTestnet,
	ValueCheck: func(data []byte) error {
		return nil
	},
	Storage: &testingStorage{},
	Network: &testingNetwork{},
}

func testingKeyManager() types.SSVSigner {
	ret := types.NewSSVKeyManager(types.PrimusTestnet)
	ret.AddShare(testingSK1)
	ret.AddShare(testingSK2)
	ret.AddShare(testingSK3)
	ret.AddShare(testingSK4)
	return ret
}

type testingStorage struct {
}

// SaveHighestDecided saves the Decided value as highest for a validator PK and role
func (s *testingStorage) SaveHighestDecided(signedMsg *qbft.SignedMessage) error {
	return nil
}

type testingNetwork struct {
}

func (net *testingNetwork) Broadcast(msg types.Encoder) error {
	return nil
}

func (net *testingNetwork) BroadcastDecided(msg types.Encoder) error {
	return nil
}
