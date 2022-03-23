package testingutils

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv"
	"github.com/bloxapp/ssv/docs/spec/types"
)

var BaseRunner = func() *ssv.DutyRunner {
	return ssv.NewDutyRunner(
		beacon.RoleTypeAttester,
		TestingShare,
		NewTestingQBFTController([]byte{1, 2, 3, 4}),
		NewTestingStorage(),
	)
}

var DecidedRunner = func() *ssv.DutyRunner {
	msgs := []*types.SSVMessage{
		SSVMsg(SignQBFTMsg(TestingSK1, 1, &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       ProposalDataBytes(TestConsensusDataByts, nil, nil),
		}), nil),
		SSVMsg(SignQBFTMsg(TestingSK1, 1, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       PrepareDataBytes(TestConsensusDataByts),
		}), nil),
		SSVMsg(SignQBFTMsg(TestingSK2, 2, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       PrepareDataBytes(TestConsensusDataByts),
		}), nil),
		SSVMsg(SignQBFTMsg(TestingSK3, 3, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       PrepareDataBytes(TestConsensusDataByts),
		}), nil),
		SSVMsg(SignQBFTMsg(TestingSK1, 1, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       CommitDataBytes(TestConsensusDataByts),
		}), nil),
		SSVMsg(SignQBFTMsg(TestingSK2, 2, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       CommitDataBytes(TestConsensusDataByts),
		}), nil),
		SSVMsg(SignQBFTMsg(TestingSK3, 3, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       CommitDataBytes(TestConsensusDataByts),
		}), nil),
	}

	v := BaseValidator()
	if err := v.DutyRunners[beacon.RoleTypeAttester].StartNewInstance([]byte{1, 2, 3, 4}); err != nil {
		panic(err.Error())
	}
	for _, msg := range msgs {
		if err := v.ProcessMessage(msg); err != nil {
			panic(err.Error())
		}
	}
	return v.DutyRunners[beacon.RoleTypeAttester]
}
