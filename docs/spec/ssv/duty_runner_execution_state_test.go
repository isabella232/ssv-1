package ssv_test

import (
	"github.com/bloxapp/ssv/docs/spec/ssv"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDutyExecutionState_Marshaling(t *testing.T) {
	es := &ssv.DutyExecutionState{
		RunningInstance: testingutils.BaseInstance(),
	}

	byts, err := es.Encode()
	require.NoError(t, err)

	decoded := &ssv.DutyExecutionState{}
	require.NoError(t, decoded.Decode(byts))
}
