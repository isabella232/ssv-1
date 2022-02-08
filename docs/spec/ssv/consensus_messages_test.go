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

		require.NoError(t, v.processConsensusMsg(dr, &qbft.SignedMessage{}))
	})

	t.Run("failed msg processing", func(t *testing.T) {
		t.Fail()
	})
}
