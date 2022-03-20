package main

import (
	"encoding/hex"
	"github.com/bloxapp/ssv/docs/spec/qbft/spectest"
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
	for _, msg := range test.Messages {
		_, _, _, err := test.Pre.ProcessMsg(msg)
		require.NoError(t, err)
	}

	postRoot, err := test.Pre.State.GetRoot()
	require.NoError(t, err)

	require.EqualValues(t, test.PostRoot, hex.EncodeToString(postRoot), "post root not valid")
}
