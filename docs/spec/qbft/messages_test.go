package qbft

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSignedMessage_DeepCopy(t *testing.T) {
	expected, err := testingSignedMsg.GetRoot()
	require.NoError(t, err)

	c := testingSignedMsg.DeepCopy()
	r, err := c.GetRoot()
	require.NoError(t, err)
	require.EqualValues(t, expected, r)
}
