package qbft

import (
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/utils/threshold"
	"github.com/herumi/bls-eth-go-binary/bls"
)

var testingMessage = &Message{
	MsgType:    ProposalMsgType,
	Height:     1,
	Round:      1,
	Identifier: []byte{1, 2, 3, 4},
	Data:       []byte{1, 2, 3, 4},
}
var TestingSignedMsg = func() *SignedMessage {
	return signMsg(TestingSK, 1, testingMessage)
}()
var signMsg = func(sk *bls.SecretKey, id types.OperatorID, msg *Message) *SignedMessage {
	domain := types.PrimusTestnet
	sigType := types.QBFTSigType

	r, _ := types.ComputeSigningRoot(testingMessage, types.ComputeSignatureDomain(domain, sigType))
	sig := sk.SignByte(r)

	return &SignedMessage{
		Message:   testingMessage,
		Signers:   []types.OperatorID{id},
		Signature: sig.Serialize(),
	}
}
var TestingSK = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.SetByCSPRNG()
	return ret
}()
var testingShare = &types.Share{
	OperatorID:    1,
	PubKey:        TestingSK.GetPublicKey().Serialize(),
	DomainType:    types.PrimusTestnet,
	Quorum:        3,
	PartialQuorum: 2,
	Committee: []*types.Operator{
		{
			OperatorID: 1,
			PubKey:     TestingSK.GetPublicKey().Serialize(),
		},
	},
}
var testingInstanceStruct = &Instance{
	State: &State{
		Share:                           testingShare,
		ID:                              []byte{1, 2, 3, 4},
		Round:                           1,
		Height:                          1,
		LastPreparedRound:               1,
		LastPreparedValue:               []byte{1, 2, 3, 4},
		ProposalAcceptedForCurrentRound: TestingSignedMsg,
	},
	ProposeContainer: &MsgContainer{
		Msgs: map[Round][]*SignedMessage{
			1: {
				TestingSignedMsg,
			},
		},
	},
	PrepareContainer: &MsgContainer{
		Msgs: map[Round][]*SignedMessage{
			1: {
				TestingSignedMsg,
			},
		},
	},
	CommitContainer: &MsgContainer{
		Msgs: map[Round][]*SignedMessage{
			1: {
				TestingSignedMsg,
			},
		},
	},
	RoundChangeContainer: &MsgContainer{
		Msgs: map[Round][]*SignedMessage{
			1: {
				TestingSignedMsg,
			},
		},
	},
	Decided:      false,
	DecidedValue: []byte{1, 2, 3, 4},
}
var testingControllerStruct = &Controller{
	Identifier: []byte{1, 2, 3, 4},
	Height:     Height(1),
	Share:      testingShare,
	StoredInstances: [HistoricalInstanceCapacity]*Instance{
		testingInstanceStruct,
	},
}
