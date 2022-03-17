package ssv

import (
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/utils/threshold"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

//func TestValidator_processPostConsensusSig(t *testing.T) {
//	threshold.Init()
//	sk1 := &bls.SecretKey{}
//	sk1.SetByCSPRNG()
//
//	t.Run("valid no quorum", func(t *testing.T) {
//		v := newTestingValidator()
//		dr := newTestingDutyRunner()
//	})
//}

func TestValidator_signPostConsensusMsg(t *testing.T) {
	km := NewTestingKeyManager()

	t.Run("valid sig", func(t *testing.T) {
		v := newTestingValidator()
		v.signer = km
		msg := &PostConsensusMessage{
			Height:          1,
			DutySignature:   []byte{1, 2, 3, 4},
			DutySigningRoot: []byte{1, 1, 1, 1},
			Signers:         []types.OperatorID{1},
		}
		sig, err := v.signPostConsensusMsg(msg)
		require.NoError(t, err)
		require.NotNil(t, sig.GetSignature())
		require.NotNil(t, sig.GetSigners())

		require.NoError(t, sig.GetSignature().Verify(msg, types.PrimusTestnet, types.PostConsensusSigType, testingSK1.GetPublicKey().Serialize()))
	})
}

func TestValidator_verifyBeaconPartialSignature(t *testing.T) {
	threshold.Init()
	t.Run("valid", func(t *testing.T) {
		sk := &bls.SecretKey{}
		sk.SetByCSPRNG()

		v := newTestingValidator()
		v.share.Committee[0] = &types.Operator{
			OperatorID: 1,
			PubKey:     sk.GetPublicKey().Serialize(),
		}

		root := []byte{1, 2, 3, 4, 5, 6}
		sig := sk.SignByte(root).Serialize()

		err := v.verifyBeaconPartialSignature(&PostConsensusMessage{
			Height:          1,
			DutySignature:   sig,
			DutySigningRoot: root,
			Signers:         []types.OperatorID{1},
		})
		require.NoError(t, err)
	})

	t.Run("invalid", func(t *testing.T) {
		sk := &bls.SecretKey{}
		sk.SetByCSPRNG()

		v := newTestingValidator()
		v.share.Committee[0] = &types.Operator{
			OperatorID: 1,
			PubKey:     sk.GetPublicKey().Serialize(),
		}

		root := []byte{1, 2, 3, 4, 5, 6}
		sig := sk.SignByte(root).Serialize()

		err := v.verifyBeaconPartialSignature(&PostConsensusMessage{
			Height:          1,
			DutySignature:   sig,
			DutySigningRoot: []byte{1, 2, 3, 4, 5, 5},
			Signers:         []types.OperatorID{1},
		})
		require.EqualError(t, err, "could not verify signature from iBFT member 1")
	})
}

func TestValidator_validatePostConsensusMsg(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		sk := &bls.SecretKey{}
		sk.SetByCSPRNG()

		v := newTestingValidator()
		v.share.Committee[0] = &types.Operator{
			OperatorID: 1,
			PubKey:     sk.GetPublicKey().Serialize(),
		}

		root := []byte{1, 2, 3, 4, 5, 6}
		sig := sk.SignByte(root).Serialize()
		msg := &PostConsensusMessage{
			Height:          1,
			DutySignature:   sig,
			DutySigningRoot: root,
			Signers:         []types.OperatorID{1},
		}
		r, err := types.ComputeSigningRoot(msg, types.ComputeSignatureDomain(types.PrimusTestnet, types.PostConsensusSigType))
		require.NoError(t, err)
		signedMsg := &SignedPostConsensusMessage{
			message:   msg,
			signature: sk.SignByte(r).Serialize(),
			signers:   []types.OperatorID{1},
		}

		executionState := &DutyExecutionState{
			PostConsensusSigRoot: root,
		}
		require.NoError(t, v.validatePostConsensusMsg(executionState, signedMsg))
	})

	t.Run("SignedPostConsensusMessage wrong signer count", func(t *testing.T) {
		sk := &bls.SecretKey{}
		sk.SetByCSPRNG()

		v := newTestingValidator()
		v.share.Committee[0] = &types.Operator{
			OperatorID: 1,
			PubKey:     sk.GetPublicKey().Serialize(),
		}

		root := []byte{1, 2, 3, 4, 5, 6}
		sig := sk.SignByte(root).Serialize()
		msg := &PostConsensusMessage{
			Height:          1,
			DutySignature:   sig,
			DutySigningRoot: root,
			Signers:         []types.OperatorID{1},
		}
		r, err := types.ComputeSigningRoot(msg, types.ComputeSignatureDomain(types.PrimusTestnet, types.PostConsensusSigType))
		require.NoError(t, err)
		signedMsg := &SignedPostConsensusMessage{
			message:   msg,
			signature: sk.SignByte(r).Serialize(),
			signers:   []types.OperatorID{1},
		}

		executionState := &DutyExecutionState{
			PostConsensusSigRoot: root,
		}

		signedMsg.signers = []types.OperatorID{1, 2}
		require.EqualError(t, v.validatePostConsensusMsg(executionState, signedMsg), "SignedPostConsensusMessage allows 1 signer")
		signedMsg.signers = []types.OperatorID{}
		require.EqualError(t, v.validatePostConsensusMsg(executionState, signedMsg), "SignedPostConsensusMessage allows 1 signer")
	})

	t.Run("invalid SignedPostConsensusMessage sig", func(t *testing.T) {
		sk := &bls.SecretKey{}
		sk.SetByCSPRNG()

		v := newTestingValidator()
		v.share.Committee[0] = &types.Operator{
			OperatorID: 1,
			PubKey:     sk.GetPublicKey().Serialize(),
		}

		root := []byte{1, 2, 3, 4, 5, 6}
		sig := sk.SignByte(root).Serialize()
		msg := &PostConsensusMessage{
			Height:          1,
			DutySignature:   sig,
			DutySigningRoot: root,
			Signers:         []types.OperatorID{1},
		}

		signedMsg := &SignedPostConsensusMessage{
			message:   msg,
			signature: sk.SignByte([]byte{1, 2, 3, 4}).Serialize(),
			signers:   []types.OperatorID{1},
		}

		executionState := &DutyExecutionState{
			PostConsensusSigRoot: root,
		}
		require.EqualError(t, v.validatePostConsensusMsg(executionState, signedMsg), "failed to verify DutySignature: failed to verify signature")
	})

	t.Run("wrong post consensus root", func(t *testing.T) {
		sk := &bls.SecretKey{}
		sk.SetByCSPRNG()

		v := newTestingValidator()
		v.share.Committee[0] = &types.Operator{
			OperatorID: 1,
			PubKey:     sk.GetPublicKey().Serialize(),
		}

		root := []byte{1, 2, 3, 4, 5, 6}
		sig := sk.SignByte(root).Serialize()
		msg := &PostConsensusMessage{
			Height:          1,
			DutySignature:   sig,
			DutySigningRoot: root,
			Signers:         []types.OperatorID{1},
		}
		r, err := types.ComputeSigningRoot(msg, types.ComputeSignatureDomain(types.PrimusTestnet, types.PostConsensusSigType))
		require.NoError(t, err)
		signedMsg := &SignedPostConsensusMessage{
			message:   msg,
			signature: sk.SignByte(r).Serialize(),
			signers:   []types.OperatorID{1},
		}

		executionState := &DutyExecutionState{
			PostConsensusSigRoot: []byte("wrong root"),
		}
		require.EqualError(t, v.validatePostConsensusMsg(executionState, signedMsg), "post consensus message signing root is wrong")
	})

	t.Run("too many signers PostConsensusMessage", func(t *testing.T) {
		sk := &bls.SecretKey{}
		sk.SetByCSPRNG()

		v := newTestingValidator()
		v.share.Committee[0] = &types.Operator{
			OperatorID: 1,
			PubKey:     sk.GetPublicKey().Serialize(),
		}

		root := []byte{1, 2, 3, 4, 5, 6}
		sig := sk.SignByte(root).Serialize()
		msg := &PostConsensusMessage{
			Height:          1,
			DutySignature:   sig,
			DutySigningRoot: root,
			Signers:         []types.OperatorID{1, 2},
		}
		r, err := types.ComputeSigningRoot(msg, types.ComputeSignatureDomain(types.PrimusTestnet, types.PostConsensusSigType))
		require.NoError(t, err)
		signedMsg := &SignedPostConsensusMessage{
			message:   msg,
			signature: sk.SignByte(r).Serialize(),
			signers:   []types.OperatorID{1},
		}

		executionState := &DutyExecutionState{
			PostConsensusSigRoot: root,
		}
		require.EqualError(t, v.validatePostConsensusMsg(executionState, signedMsg), "PostConsensusMessage allows 1 signer")
	})

	t.Run("no signers PostConsensusMessage", func(t *testing.T) {
		sk := &bls.SecretKey{}
		sk.SetByCSPRNG()

		v := newTestingValidator()
		v.share.Committee[0] = &types.Operator{
			OperatorID: 1,
			PubKey:     sk.GetPublicKey().Serialize(),
		}

		root := []byte{1, 2, 3, 4, 5, 6}
		sig := sk.SignByte(root).Serialize()
		msg := &PostConsensusMessage{
			Height:          1,
			DutySignature:   sig,
			DutySigningRoot: root,
			Signers:         []types.OperatorID{},
		}
		r, err := types.ComputeSigningRoot(msg, types.ComputeSignatureDomain(types.PrimusTestnet, types.PostConsensusSigType))
		require.NoError(t, err)
		signedMsg := &SignedPostConsensusMessage{
			message:   msg,
			signature: sk.SignByte(r).Serialize(),
			signers:   []types.OperatorID{1},
		}

		executionState := &DutyExecutionState{
			PostConsensusSigRoot: root,
		}
		require.EqualError(t, v.validatePostConsensusMsg(executionState, signedMsg), "PostConsensusMessage allows 1 signer")
	})

	t.Run("beacon partial sig invalid", func(t *testing.T) {
		sk := &bls.SecretKey{}
		sk.SetByCSPRNG()

		wrongSK := &bls.SecretKey{}
		wrongSK.SetByCSPRNG()

		v := newTestingValidator()
		v.share.Committee[0] = &types.Operator{
			OperatorID: 1,
			PubKey:     sk.GetPublicKey().Serialize(),
		}

		root := []byte{1, 2, 3, 4, 5, 6}
		sig := wrongSK.SignByte(root).Serialize()
		msg := &PostConsensusMessage{
			Height:          1,
			DutySignature:   sig,
			DutySigningRoot: []byte{1, 2, 3, 4, 5, 6},
			Signers:         []types.OperatorID{1},
		}
		r, err := types.ComputeSigningRoot(msg, types.ComputeSignatureDomain(types.PrimusTestnet, types.PostConsensusSigType))
		require.NoError(t, err)
		signedMsg := &SignedPostConsensusMessage{
			message:   msg,
			signature: sk.SignByte(r).Serialize(),
			signers:   []types.OperatorID{1},
		}

		executionState := &DutyExecutionState{
			PostConsensusSigRoot: root,
		}
		require.EqualError(t, v.validatePostConsensusMsg(executionState, signedMsg), "could not verify beacon partial signature: could not verify signature from iBFT member 1")
	})
}
