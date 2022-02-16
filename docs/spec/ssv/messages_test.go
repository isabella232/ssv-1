package ssv

import (
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSignedPostConsensusMessage_MatchedSigners(t *testing.T) {
	t.Run("matched same order", func(t *testing.T) {
		msg := &SignedPostConsensusMessage{}
		msg.signers = []types.NodeID{1, 2, 3, 4}
		require.True(t, msg.MatchedSigners([]types.NodeID{1, 2, 3, 4}))
	})

	t.Run("matched different order", func(t *testing.T) {
		msg := &SignedPostConsensusMessage{}
		msg.signers = []types.NodeID{1, 2, 3, 4}
		require.True(t, msg.MatchedSigners([]types.NodeID{2, 1, 4, 3}))
	})

	t.Run("matched same order with duplicate", func(t *testing.T) {
		msg := &SignedPostConsensusMessage{}
		msg.signers = []types.NodeID{3, 1, 2, 3}
		require.True(t, msg.MatchedSigners([]types.NodeID{3, 1, 2, 3}))
	})

	t.Run("matched different duplicate", func(t *testing.T) {
		msg := &SignedPostConsensusMessage{}
		msg.signers = []types.NodeID{1, 2, 3, 3}
		require.True(t, msg.MatchedSigners([]types.NodeID{3, 1, 2, 3}))
	})

	t.Run("not matched same order", func(t *testing.T) {
		msg := &SignedPostConsensusMessage{}
		msg.signers = []types.NodeID{1, 2, 3, 4, 4}
		require.False(t, msg.MatchedSigners([]types.NodeID{1, 2, 3, 4}))
	})

	t.Run("not matched", func(t *testing.T) {
		msg := &SignedPostConsensusMessage{}
		msg.signers = []types.NodeID{1, 2, 3, 3}
		require.False(t, msg.MatchedSigners([]types.NodeID{1, 2, 3, 4}))
	})

	t.Run("not matched", func(t *testing.T) {
		msg := &SignedPostConsensusMessage{}
		msg.signers = []types.NodeID{1, 2, 3}
		require.False(t, msg.MatchedSigners([]types.NodeID{1, 2, 3, 4}))
	})
}

//func TestSignedPostConsensusMessage_Aggregate(t *testing.T) {
//	threshold.Init()
//	sk1 := &bls.SecretKey{}
//	sk1.SetByCSPRNG()
//	sk2 := &bls.SecretKey{}
//	sk2.SetByCSPRNG()
//	sk3 := &bls.SecretKey{}
//	sk3.SetByCSPRNG()
//	sk4 := &bls.SecretKey{}
//	sk4.SetByCSPRNG()
//
//	t.Run("valid aggregate", func(t *testing.T) {
//		sig := sk1.SignByte([]byte{1, 2, 3, 4})
//		sig.Add(sk2.SignByte([]byte{1, 2, 3, 4}))
//		msg1 := &PostConsensusMessage{
//			DutySignature:   sk1.SignByte([]byte{1, 2, 3, 4}).Serialize(),
//			Signers:         []types.NodeID{1},
//			DutySigningRoot: []byte{1, 2, 3, 4},
//		}
//		msg2 := &PostConsensusMessage{
//			DutySignature:   sk2.SignByte([]byte{1, 2, 3, 4}).Serialize(),
//			Signers:         []types.NodeID{2},
//			DutySigningRoot: []byte{1, 2, 3, 4},
//		}
//
//		require.NoError(t, msg1.Aggregate(msg2))
//		msg1.MatchedSigners([]types.NodeID{1, 2})
//	})
//
//	t.Run("partially matching Signers", func(t *testing.T) {
//		sig := sk1.SignByte([]byte{1, 2, 3, 4})
//		sig.Add(sk2.SignByte([]byte{1, 2, 3, 4}))
//		msg1 := &PostConsensusMessage{
//			Signers:         []types.NodeID{1, 2},
//			DutySigningRoot: []byte{1, 2, 3, 4},
//		}
//		msg2 := &PostConsensusMessage{
//			Signers:         []types.NodeID{2},
//			DutySigningRoot: []byte{1, 2, 3, 4},
//		}
//
//		require.EqualError(t, msg1.Aggregate(msg2), "signer IDs partially/ fully match")
//	})
//
//	t.Run("different roots", func(t *testing.T) {
//		msg1 := &PostConsensusMessage{
//			DutySigningRoot: []byte{1, 2, 3, 4},
//		}
//		msg2 := &PostConsensusMessage{
//			DutySigningRoot: []byte{1, 2, 3, 3},
//		}
//
//		require.EqualError(t, msg1.Aggregate(msg2), "can't aggregate msgs with different roots")
//	})
//}

func TestSignedPostConsensusMessage_Marshaling(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		signed := &SignedPostConsensusMessage{
			message: &PostConsensusMessage{
				Height:          1,
				DutySignature:   []byte{1, 2, 3, 4},
				DutySigningRoot: []byte{1, 1, 1, 1},
				Signers:         []types.NodeID{1},
			},
			signers:   []types.NodeID{1},
			signature: []byte{1, 2, 3, 4},
		}

		byts, err := signed.Encode()
		require.NoError(t, err)

		decoded := &SignedPostConsensusMessage{}
		require.NoError(t, decoded.Decode(byts))
	})
}
