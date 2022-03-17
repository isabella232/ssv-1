package ssv

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidator_ProcessConsensusMsg(t *testing.T) {
	t.Run("non decided qbft msg", func(t *testing.T) {
		v := newTestingValidator()
		dr := v.dutyRunners[beacon.RoleTypeAttester]
		require.NoError(t, dr.StartNewInstance([]byte{1, 2, 3, 4}))
		require.NoError(t, v.processConsensusMsg(dr, testingSignedQBFTMsg))
	})

	t.Run("failed msg processing (wrong identifier)", func(t *testing.T) {
		v := newTestingValidator()
		dr := v.dutyRunners[beacon.RoleTypeAttester]
		dr.QBFTController.Identifier = []byte{1, 2, 3, 3}
		require.NoError(t, dr.StartNewInstance([]byte{1, 2, 3, 4}))
		require.EqualError(t, v.processConsensusMsg(dr, testingSignedQBFTMsg), "failed to process consensus msg: message doesn't belong to Identifier 01020303")
	})

	t.Run("decided", func(t *testing.T) {
		v := newTestingValidator()
		dr := newTestingDutyRunner()
		require.NoError(t, dr.StartNewInstance([]byte{1, 2, 3, 4}))

		msgs := []*qbft.SignedMessage{
			qbft.SignMsg(testingSK1, 1, &qbft.Message{
				MsgType:    qbft.ProposalMsgType,
				Height:     qbft.FirstHeight,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       proposalDataByts(testConsensusDataByts, nil, nil),
			}),
			qbft.SignMsg(testingSK1, 1, &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     qbft.FirstHeight,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       prepareDataByts(testConsensusDataByts),
			}),
			qbft.SignMsg(testingSK2, 2, &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     qbft.FirstHeight,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       prepareDataByts(testConsensusDataByts),
			}),
			qbft.SignMsg(testingSK3, 3, &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     qbft.FirstHeight,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       prepareDataByts(testConsensusDataByts),
			}),
			qbft.SignMsg(testingSK1, 1, &qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     qbft.FirstHeight,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       commitDataByts(testConsensusDataByts),
			}),
			qbft.SignMsg(testingSK2, 2, &qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     qbft.FirstHeight,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       commitDataByts(testConsensusDataByts),
			}),
			qbft.SignMsg(testingSK3, 3, &qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     qbft.FirstHeight,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       commitDataByts(testConsensusDataByts),
			}),
		}

		for _, msg := range msgs {
			require.NoError(t, v.processConsensusMsg(dr, msg))
		}

		decidedVal, err := dr.DutyExecutionState.DecidedValue.Encode()
		require.NoError(t, err)
		require.EqualValues(t, decidedVal, testConsensusDataByts)
	})
}
