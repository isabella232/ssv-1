package ssv

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDutyExecutionState_AddPartialSig(t *testing.T) {
	t.Run("add to empty", func(t *testing.T) {
		s := NewTestingDutyExecutionState()
		s.AddPartialSig(&testingPostConsensusSigMessage{
			signerID: 1,
		})

		require.Len(t, s.collectedPartialSigs, 1)
	})

	t.Run("add multiple", func(t *testing.T) {
		s := NewTestingDutyExecutionState()
		s.AddPartialSig(&testingPostConsensusSigMessage{
			signerID: 1,
		})
		s.AddPartialSig(&testingPostConsensusSigMessage{
			signerID: 2,
		})
		s.AddPartialSig(&testingPostConsensusSigMessage{
			signerID: 3,
		})

		require.Len(t, s.collectedPartialSigs, 3)
	})

	t.Run("add duplicate", func(t *testing.T) {
		s := NewTestingDutyExecutionState()
		s.AddPartialSig(&testingPostConsensusSigMessage{
			signerID: 1,
		})
		s.AddPartialSig(&testingPostConsensusSigMessage{
			signerID: 1,
		})
		s.AddPartialSig(&testingPostConsensusSigMessage{
			signerID: 3,
		})

		require.Len(t, s.collectedPartialSigs, 2)
	})
}
