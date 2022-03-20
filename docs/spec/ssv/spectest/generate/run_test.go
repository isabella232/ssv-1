package main

import (
	"encoding/hex"
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAll(t *testing.T) {
	for _, test := range spectest.AllTests {
		t.Run(test.Name, func(t *testing.T) {
			runTest(t, test)
		})
	}
}

func runTest(t *testing.T, test *spectest.SpecTest) {
	v := testingutils.BaseValidator()
	v.DutyRunners[beacon.RoleTypeAttester] = test.DutyRunner

	for _, msg := range test.Messages {
		require.NoError(t, v.ProcessMessage(*msg))
	}

	postRoot, err := test.DutyRunner.DutyExecutionState.GetRoot()
	require.NoError(t, err)

	require.EqualValues(t, test.PostDutyRunnerStateRoot, hex.EncodeToString(postRoot))
}
