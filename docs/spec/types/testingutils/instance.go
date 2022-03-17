package testingutils

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/types"
)

var TestingConfig = &qbft.Config{
	Signer:    NewTestingKeyManager(),
	SigningPK: TestingSK1.GetPublicKey().Serialize(),
	Domain:    types.PrimusTestnet,
	ValueCheck: func(data []byte) error {
		return nil
	},
	Storage: NewTestingStorage(),
	Network: NewTestingNetwork(),
}
var TestingValidatorPubKey = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8}
var BaseInstance = func() *qbft.Instance {
	ret := qbft.NewInstance(TestingConfig, nil, nil)
	ret.State = &qbft.State{
		Share: &types.Share{
			OperatorID:      1,
			ValidatorPubKey: TestingValidatorPubKey,
			SharePubKey:     TestingSK1.GetPublicKey().Serialize(),
			DomainType:      types.PrimusTestnet,
			Quorum:          3,
			PartialQuorum:   2,
			Committee: []*types.Operator{
				{
					OperatorID: 1,
					PubKey:     TestingSK1.GetPublicKey().Serialize(),
				},
				{
					OperatorID: 2,
					PubKey:     TestingSK2.GetPublicKey().Serialize(),
				},
				{
					OperatorID: 3,
					PubKey:     TestingSK3.GetPublicKey().Serialize(),
				},
				{
					OperatorID: 4,
					PubKey:     TestingSK4.GetPublicKey().Serialize(),
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
