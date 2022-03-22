package ssv_test

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
	"github.com/stretchr/testify/require"
	"testing"
)

var testingSignedQBFTMsg = func() *qbft.SignedMessage {
	msg := &qbft.Message{
		MsgType:    qbft.ProposalMsgType,
		Height:     qbft.FirstHeight,
		Round:      qbft.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ProposalDataBytes([]byte{1, 2, 3, 4}, nil, nil),
	}
	return testingutils.SignQBFTMsg(testingutils.TestingSK1, 1, msg)
}()

func TestValidator_ProcessConsensusMsg(t *testing.T) {
	t.Run("non decided qbft msg", func(t *testing.T) {
		v := testingutils.BaseValidator()
		dr := v.DutyRunners[beacon.RoleTypeAttester]
		require.NoError(t, dr.StartNewInstance([]byte{1, 2, 3, 4}))
		require.NoError(t, v.ProcessMessage(testingutils.SSVMsg(testingSignedQBFTMsg, nil)))
	})

	t.Run("failed msg processing (wrong identifier)", func(t *testing.T) {
		v := testingutils.BaseValidator()
		dr := v.DutyRunners[beacon.RoleTypeAttester]
		dr.QBFTController.Identifier = []byte{1, 2, 3, 3}
		require.NoError(t, dr.StartNewInstance([]byte{1, 2, 3, 4}))
		require.EqualError(t, v.ProcessMessage(testingutils.SSVMsg(testingSignedQBFTMsg, nil)), "failed to process consensus msg: message doesn't belong to Identifier 01020303")
	})

	t.Run("decided", func(t *testing.T) {
		v := testingutils.BaseValidator()
		dr := v.DutyRunners[beacon.RoleTypeAttester]
		require.NoError(t, dr.StartNewInstance([]byte{1, 2, 3, 4}))

		msgs := []*qbft.SignedMessage{
			testingutils.SignQBFTMsg(testingutils.TestingSK1, 1, &qbft.Message{
				MsgType:    qbft.ProposalMsgType,
				Height:     qbft.FirstHeight,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ProposalDataBytes(testingutils.TestConsensusDataByts, nil, nil),
			}),
			testingutils.SignQBFTMsg(testingutils.TestingSK1, 1, &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     qbft.FirstHeight,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.PrepareDataBytes(testingutils.TestConsensusDataByts),
			}),
			testingutils.SignQBFTMsg(testingutils.TestingSK2, 2, &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     qbft.FirstHeight,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.PrepareDataBytes(testingutils.TestConsensusDataByts),
			}),
			testingutils.SignQBFTMsg(testingutils.TestingSK3, 3, &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     qbft.FirstHeight,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.PrepareDataBytes(testingutils.TestConsensusDataByts),
			}),
			testingutils.SignQBFTMsg(testingutils.TestingSK1, 1, &qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     qbft.FirstHeight,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.CommitDataBytes(testingutils.TestConsensusDataByts),
			}),
			testingutils.SignQBFTMsg(testingutils.TestingSK2, 2, &qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     qbft.FirstHeight,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.CommitDataBytes(testingutils.TestConsensusDataByts),
			}),
			testingutils.SignQBFTMsg(testingutils.TestingSK3, 3, &qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     qbft.FirstHeight,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.CommitDataBytes(testingutils.TestConsensusDataByts),
			}),
		}

		for _, msg := range msgs {
			require.NoError(t, v.ProcessMessage(testingutils.SSVMsg(msg, nil)))
		}

		decidedVal, err := dr.DutyExecutionState.DecidedValue.Encode()
		require.NoError(t, err)
		require.EqualValues(t, decidedVal, testingutils.TestConsensusDataByts)
	})
}
