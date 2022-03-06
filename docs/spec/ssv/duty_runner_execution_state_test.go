package ssv

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDutyExecutionState_Marshaling(t *testing.T) {
	es := &DutyExecutionState{
		RunningInstance: newTestingQBFTInstance(),
	}

	byts, err := es.Encode()
	require.NoError(t, err)

	decoded := &DutyExecutionState{}
	require.NoError(t, decoded.Decode(byts))
}
