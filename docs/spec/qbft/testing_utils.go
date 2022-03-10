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
	id := types.OperatorID(1)
	domain := types.PrimusTestnet
	sigType := types.QBFTSigType

	r, _ := types.ComputeSigningRoot(testingMessage, types.ComputeSignatureDomain(domain, sigType))
	sig := TestingSK.SignByte(r)

	return &SignedMessage{
		Message:   testingMessage,
		Signers:   []types.OperatorID{id},
		Signature: sig.Serialize(),
	}
}()
var TestingSK = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.SetByCSPRNG()
	return ret
}()

var testingInstanceStruct = &Instance{
	State: State{
		Share: &types.Share{
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
		},
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
	Height:     1,
	StoredInstances: [HistoricalInstanceCapacity]*Instance{
		testingInstanceStruct,
	},
}
