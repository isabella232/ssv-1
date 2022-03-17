package spectest

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv"
	"github.com/bloxapp/ssv/docs/spec/types"
)

func happyFullFlow() *SpecTest {
	dr := baseRunner()
	dr.StartNewInstance([]byte{1, 2, 3, 4})

	msgs := []*types.SSVMessage{
		ssvMsg(qbft.SignMsg(testingSK1, 1, &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       ssv.ProposalDataByts(ssv.TestConsensusDataByts, nil, nil),
		}), nil),
	}

	return &SpecTest{
		Name:                    "happy full flow",
		DutyRunner:              dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "",
	}
}
