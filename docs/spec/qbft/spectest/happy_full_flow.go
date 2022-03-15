package spectest

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/types"
)

func happyFullFlow() *SpecTest {
	pre := baseInstance()
	msgs := []*qbft.SignedMessage{
		signMsg(testingSK1, types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       proposalDataByts([]byte{1, 2, 3, 4}, nil, nil),
		}),
		signMsg(testingSK1, types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       prepareDataByts([]byte{1, 2, 3, 4}, nil, nil),
		}),
		signMsg(testingSK2, types.OperatorID(2), &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       prepareDataByts([]byte{1, 2, 3, 4}, nil, nil),
		}),
		signMsg(testingSK3, types.OperatorID(3), &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       prepareDataByts([]byte{1, 2, 3, 4}, nil, nil),
		}),
		signMsg(testingSK1, types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       commitDataByts([]byte{1, 2, 3, 4}, nil, nil),
		}),
		signMsg(testingSK2, types.OperatorID(2), &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       commitDataByts([]byte{1, 2, 3, 4}, nil, nil),
		}),
		signMsg(testingSK3, types.OperatorID(3), &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       commitDataByts([]byte{1, 2, 3, 4}, nil, nil),
		}),
	}
	return &SpecTest{
		Name:     "happy full flow",
		Pre:      pre,
		PostRoot: "edd470ed8a6b435581662a55cead5d963eac8646b4ee8650cc1c30416e413e39",
		Messages: msgs,
	}
}
