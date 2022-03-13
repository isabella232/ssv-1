package qbft

import (
	"github.com/stretchr/testify/require"
	"testing"
)

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
