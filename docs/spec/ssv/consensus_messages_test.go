package ssv

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidator_ProcessConsensusMsg(t *testing.T) {
	t.Run("non decided qbft msg", func(t *testing.T) {
		v := newTestingValidator()
		dr := newTestingDutyRunner()

		require.NoError(t, v.processConsensusMsg(dr, qbft.TestingSignedMsg))
	})

	t.Run("failed msg processing", func(t *testing.T) {
		t.Fail()
		//v := newTestingValidator()
		//dr := newTestingDutyRunner()
		//dr.State.QBFTController.(*testingQBFTController).failProcessMsg = true
		//
		//require.EqualError(t, v.processConsensusMsg(dr, &qbft.SignedMessage{}), "failed to process consensus msg: failed process msg")
	})

	t.Run("decided", func(t *testing.T) {
		t.Fail()
		//v := newTestingValidator()
		//dr := newTestingDutyRunner()
		//require.NoError(t, dr.StartNewInstance(testConsensusDataByts))
		//dr.State.QBFTController.(*testingQBFTController).returnDecided = true
		//dr.State.QBFTController.(*testingQBFTController).returnDecidedValue = testConsensusDataByts
		//
		//require.NoError(t, v.processConsensusMsg(dr, &qbft.SignedMessage{}))
		//require.NotNil(t, dr.State.DutyExecutionState.DecidedValue)
	})
}
