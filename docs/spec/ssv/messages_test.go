package ssv

import (
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/utils/threshold"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPostConsensusSigMessage_MatchedSigners(t *testing.T) {
	t.Run("matched same order", func(t *testing.T) {
		msg := &PostConsensusSigMessage{}
		msg.signers = []types.NodeID{1, 2, 3, 4}
		require.True(t, msg.MatchedSigners([]types.NodeID{1, 2, 3, 4}))
	})

	t.Run("matched different order", func(t *testing.T) {
		msg := &PostConsensusSigMessage{}
		msg.signers = []types.NodeID{1, 2, 3, 4}
		require.True(t, msg.MatchedSigners([]types.NodeID{2, 1, 4, 3}))
	})

	t.Run("matched same order with duplicate", func(t *testing.T) {
		msg := &PostConsensusSigMessage{}
		msg.signers = []types.NodeID{3, 1, 2, 3}
		require.True(t, msg.MatchedSigners([]types.NodeID{3, 1, 2, 3}))
	})

	t.Run("matched different duplicate", func(t *testing.T) {
		msg := &PostConsensusSigMessage{}
		msg.signers = []types.NodeID{1, 2, 3, 3}
		require.True(t, msg.MatchedSigners([]types.NodeID{3, 1, 2, 3}))
	})

	t.Run("not matched same order", func(t *testing.T) {
		msg := &PostConsensusSigMessage{}
		msg.signers = []types.NodeID{1, 2, 3, 4, 4}
		require.False(t, msg.MatchedSigners([]types.NodeID{1, 2, 3, 4}))
	})

	t.Run("not matched", func(t *testing.T) {
		msg := &PostConsensusSigMessage{}
		msg.signers = []types.NodeID{1, 2, 3, 3}
		require.False(t, msg.MatchedSigners([]types.NodeID{1, 2, 3, 4}))
	})

	t.Run("not matched", func(t *testing.T) {
		msg := &PostConsensusSigMessage{}
		msg.signers = []types.NodeID{1, 2, 3}
		require.False(t, msg.MatchedSigners([]types.NodeID{1, 2, 3, 4}))
	})
}

func TestPostConsensusSigMessage_IsValidSignature(t *testing.T) {
	threshold.Init()
	sk1 := &bls.SecretKey{}
	sk1.SetByCSPRNG()
	sk2 := &bls.SecretKey{}
	sk2.SetByCSPRNG()
	sk3 := &bls.SecretKey{}
	sk3.SetByCSPRNG()
	sk4 := &bls.SecretKey{}
	sk4.SetByCSPRNG()

	nodes := []*types.Node{
		{NodeID: 1, PubKey: sk1.GetPublicKey().Serialize()},
		{NodeID: 2, PubKey: sk2.GetPublicKey().Serialize()},
		{NodeID: 3, PubKey: sk3.GetPublicKey().Serialize()},
		{NodeID: 4, PubKey: sk4.GetPublicKey().Serialize()},
	}

	t.Run("valid sig single signer", func(t *testing.T) {
		msg := &PostConsensusSigMessage{
			signature: sk1.SignByte([]byte{1, 2, 3, 4}).Serialize(),
			signers:   []types.NodeID{1},
			root:      []byte{1, 2, 3, 4},
		}

		require.True(t, msg.IsValidSignature(nodes))
	})

	t.Run("valid sig multi signer", func(t *testing.T) {
		sig := sk1.SignByte([]byte{1, 2, 3, 4})
		sig.Add(sk2.SignByte([]byte{1, 2, 3, 4}))
		msg := &PostConsensusSigMessage{
			signature: sig.Serialize(),
			signers:   []types.NodeID{1, 2},
			root:      []byte{1, 2, 3, 4},
		}

		require.True(t, msg.IsValidSignature(nodes))
	})

	t.Run("invalid sig", func(t *testing.T) {
		sig := sk1.SignByte([]byte{1, 2, 3, 4})
		sig.Add(sk2.SignByte([]byte{1, 2, 3, 4}))
		msg := &PostConsensusSigMessage{
			signature: sig.Serialize(),
			signers:   []types.NodeID{1, 2, 3},
			root:      []byte{1, 2, 3, 4},
		}

		require.False(t, msg.IsValidSignature(nodes))
	})
}

func TestPostConsensusSigMessage_Aggregate(t *testing.T) {
	threshold.Init()
	sk1 := &bls.SecretKey{}
	sk1.SetByCSPRNG()
	sk2 := &bls.SecretKey{}
	sk2.SetByCSPRNG()
	sk3 := &bls.SecretKey{}
	sk3.SetByCSPRNG()
	sk4 := &bls.SecretKey{}
	sk4.SetByCSPRNG()

	t.Run("valid aggregate", func(t *testing.T) {
		sig := sk1.SignByte([]byte{1, 2, 3, 4})
		sig.Add(sk2.SignByte([]byte{1, 2, 3, 4}))
		msg1 := &PostConsensusSigMessage{
			signature: sk1.SignByte([]byte{1, 2, 3, 4}).Serialize(),
			signers:   []types.NodeID{1},
			root:      []byte{1, 2, 3, 4},
		}
		msg2 := &PostConsensusSigMessage{
			signature: sk2.SignByte([]byte{1, 2, 3, 4}).Serialize(),
			signers:   []types.NodeID{2},
			root:      []byte{1, 2, 3, 4},
		}

		require.NoError(t, msg1.Aggregate(msg2))
		msg1.MatchedSigners([]types.NodeID{1, 2})
	})

	t.Run("partially matching signers", func(t *testing.T) {
		sig := sk1.SignByte([]byte{1, 2, 3, 4})
		sig.Add(sk2.SignByte([]byte{1, 2, 3, 4}))
		msg1 := &PostConsensusSigMessage{
			signers: []types.NodeID{1, 2},
			root:    []byte{1, 2, 3, 4},
		}
		msg2 := &PostConsensusSigMessage{
			signers: []types.NodeID{2},
			root:    []byte{1, 2, 3, 4},
		}

		require.EqualError(t, msg1.Aggregate(msg2), "signer IDs partially/ fully match")
	})

	t.Run("different roots", func(t *testing.T) {
		msg1 := &PostConsensusSigMessage{
			root: []byte{1, 2, 3, 4},
		}
		msg2 := &PostConsensusSigMessage{
			root: []byte{1, 2, 3, 3},
		}

		require.EqualError(t, msg1.Aggregate(msg2), "can't aggregate msgs with different roots")
	})
}
