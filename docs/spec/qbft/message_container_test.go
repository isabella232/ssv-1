package qbft

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMsgContainer_AddIfDoesntExist(t *testing.T) {
	t.Run("same msg and signers", func(t *testing.T) {
		c := &MsgContainer{
			Msgs: map[Round][]*SignedMessage{},
		}

		added, err := c.AddIfDoesntExist(TestingSignedMsg)
		require.NoError(t, err)
		require.True(t, added)

		added, err = c.AddIfDoesntExist(TestingSignedMsg)
		require.NoError(t, err)
		require.False(t, added)
	})

	t.Run("same msg different signers", func(t *testing.T) {
		c := &MsgContainer{
			Msgs: map[Round][]*SignedMessage{},
		}

		added, err := c.AddIfDoesntExist(TestingSignedMsg)
		require.NoError(t, err)
		require.True(t, added)

		added, err = c.AddIfDoesntExist(signMsg(TestingSK, 2, testingMessage))
		require.NoError(t, err)
		require.True(t, added)
	})
}

func TestMsgContainer_Marshaling(t *testing.T) {
	c := &MsgContainer{
		Msgs: map[Round][]*SignedMessage{},
	}
	c.Msgs[1] = []*SignedMessage{TestingSignedMsg}

	byts, err := c.Encode()
	require.NoError(t, err)

	decoded := &MsgContainer{}
	require.NoError(t, decoded.Decode(byts))

	decodedByts, err := decoded.Encode()
	require.NoError(t, err)
	require.EqualValues(t, byts, decodedByts)
}
