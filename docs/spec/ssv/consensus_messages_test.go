package ssv

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"testing"
)

func TestValidator_ProcessConsensusMsg(t *testing.T) {
	t.Run("non decided qbft msg", func(t *testing.T) {
		v := newTestingValidator()
		dr := newTestingDutyRunner()

		v.processConsensusMsg(dr, &qbft.SignedMessage())
	})
}
