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
	threshold.Init()
	sk1 := &bls.SecretKey{}
	sk1.SetByCSPRNG()

	km := &testingKeyManager{
		sk:     sk1,
		domain: types.PrimusTestnet,
	}

	t.Run("valid sig", func(t *testing.T) {
		v := newTestingValidator()
		v.signer = km
		msg := &PostConsensusMessage{
			Height:          1,
			DutySignature:   []byte{1, 2, 3, 4},
			DutySigningRoot: []byte{1, 1, 1, 1},
			Signers:         []types.NodeID{1},
		}
		sig, err := v.signPostConsensusMsg(msg)
		require.NoError(t, err)
		require.NotNil(t, sig.GetSignature())
		require.NotNil(t, sig.GetSigners())

		require.NoError(t, sig.GetSignature().Verify(msg, types.PrimusTestnet, types.PostConsensusSigType, sk1.GetPublicKey().Serialize()))
	})
}
