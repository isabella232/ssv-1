package types

import (
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

type testSigningRoot struct {
	root      []byte
	Signature []byte
	signers   []NodeID
}

func (r *testSigningRoot) GetRoot() []byte {
	return r.root
}

func (r *testSigningRoot) GetSignature() []byte {
	return r.Signature
}

func (r *testSigningRoot) GetSigners() []NodeID {
	return r.signers
}

// IsValidSignature returns true if signature is valid (against message and signers)
func (r *testSigningRoot) IsValidSignature(domain DomainType, nodes []*Node) error {
	panic("fail")
}

// MatchedSigners returns true if the provided signer ids are equal to GetSignerIds() without order significance
func (r *testSigningRoot) MatchedSigners(ids []NodeID) bool {
	panic("fail")
}

// Aggregate will aggregate the signed message if possible (unique signers, same digest, valid)
func (r *testSigningRoot) Aggregate(signedMsg MessageSignature) error {
	panic("fail")
}

func TestComputeSigningRoot(t *testing.T) {
	t.Run("", func(t *testing.T) {
		root := &testSigningRoot{root: []byte{1, 2, 3, 4}}
		domain := PrimusTestnet
		sigType := QBFTSigType
		require.EqualValues(t, []byte{0x8e, 0x9e, 0xa8, 0x82, 0x0, 0x46, 0xb7, 0x5d, 0xe9, 0x0, 0xb5, 0xdc, 0x1c, 0xb, 0xa5, 0x82, 0xf7, 0xc6, 0x79, 0xc7, 0x3d, 0x20, 0xf, 0x95, 0x81, 0x23, 0xa5, 0xbc, 0x2f, 0x2c, 0xd8, 0x3e}, ComputeSigningRoot(root, ComputeSignatureDomain(domain, sigType)))
	})
}

func TestComputeSignatureDomain(t *testing.T) {
	require.EqualValues(t, []byte{1, 2, 3, 4, 1, 2, 3, 4}, ComputeSignatureDomain([]byte{1, 2, 3, 4}, []byte{1, 2, 3, 4}))
}

func TestSignature_Verify(t *testing.T) {
	msgRoot := &testSigningRoot{root: []byte{1, 2, 3, 4}}
	domain := PrimusTestnet
	sigType := QBFTSigType

	computedRoot := ComputeSigningRoot(msgRoot, ComputeSignatureDomain(domain, sigType))

	sk := &bls.SecretKey{}
	sk.SetByCSPRNG()
	sig := sk.SignByte(computedRoot)

	t.Run("valid sig", func(t *testing.T) {
		require.NoError(t, Signature(sig.Serialize()).Verify(msgRoot, domain, sigType, sk.GetPublicKey().Serialize()))
	})

	t.Run("invalid sig", func(t *testing.T) {
		wrongSK := &bls.SecretKey{}
		wrongSK.SetByCSPRNG()
		require.EqualError(t, Signature(sig.Serialize()).Verify(msgRoot, domain, sigType, wrongSK.GetPublicKey().Serialize()), "failed to verify signature")
	})

	t.Run("invalid sig", func(t *testing.T) {
		require.EqualError(t, Signature([]byte{1, 2, 3, 4}).Verify(msgRoot, domain, sigType, sk.GetPublicKey().Serialize()), "failed to deserialize signature: err blsSignatureDeserialize 01020304")
	})

	t.Run("invalid pk", func(t *testing.T) {
		require.EqualError(t, Signature(sig.Serialize()).Verify(msgRoot, domain, sigType, []byte{1, 2, 3, 4}), "failed to deserialize public key: err blsPublicKeyDeserialize 01020304")
	})
}

func TestSignature_VerifyMultiPubKey(t *testing.T) {
	msgRoot := &testSigningRoot{root: []byte{1, 2, 3, 4}}
	domain := PrimusTestnet
	sigType := QBFTSigType

	computedRoot := ComputeSigningRoot(msgRoot, ComputeSignatureDomain(domain, sigType))

	sk1 := &bls.SecretKey{}
	sk1.SetByCSPRNG()

	sk2 := &bls.SecretKey{}
	sk2.SetByCSPRNG()

	sk3 := &bls.SecretKey{}
	sk3.SetByCSPRNG()

	t.Run("single signer", func(t *testing.T) {
		require.NoError(t, Signature(sk1.SignByte(computedRoot).Serialize()).VerifyMultiPubKey(msgRoot, domain, sigType, [][]byte{sk1.GetPublicKey().Serialize()}))
	})

	t.Run("2 signers", func(t *testing.T) {
		agg := sk1.SignByte(computedRoot)
		agg.Add(sk2.SignByte(computedRoot))

		pks := [][]byte{
			sk1.GetPublicKey().Serialize(),
			sk2.GetPublicKey().Serialize(),
		}

		require.NoError(t, Signature(agg.Serialize()).VerifyMultiPubKey(msgRoot, domain, sigType, pks))
	})

	t.Run("2 wrong signers", func(t *testing.T) {
		agg := sk1.SignByte(computedRoot)
		agg.Add(sk2.SignByte(computedRoot))

		pks := [][]byte{
			sk1.GetPublicKey().Serialize(),
			sk3.GetPublicKey().Serialize(),
		}

		require.EqualError(t, Signature(agg.Serialize()).VerifyMultiPubKey(msgRoot, domain, sigType, pks), "failed to verify signature")
	})

	t.Run("3 signers", func(t *testing.T) {
		agg := sk1.SignByte(computedRoot)
		agg.Add(sk2.SignByte(computedRoot))
		agg.Add(sk3.SignByte(computedRoot))

		pks := [][]byte{
			sk1.GetPublicKey().Serialize(),
			sk2.GetPublicKey().Serialize(),
			sk3.GetPublicKey().Serialize(),
		}

		require.NoError(t, Signature(agg.Serialize()).VerifyMultiPubKey(msgRoot, domain, sigType, pks))
	})

	t.Run("3 signers, 2 pks", func(t *testing.T) {
		agg := sk1.SignByte(computedRoot)
		agg.Add(sk2.SignByte(computedRoot))
		agg.Add(sk3.SignByte(computedRoot))

		pks := [][]byte{
			sk1.GetPublicKey().Serialize(),
			sk2.GetPublicKey().Serialize(),
		}

		require.EqualError(t, Signature(agg.Serialize()).VerifyMultiPubKey(msgRoot, domain, sigType, pks), "failed to verify signature")
	})
}

func TestSignature_VerifyByNodes(t *testing.T) {
	msgRoot := &testSigningRoot{root: []byte{1, 2, 3, 4}}
	domain := PrimusTestnet
	sigType := QBFTSigType

	computedRoot := ComputeSigningRoot(msgRoot, ComputeSignatureDomain(domain, sigType))

	sk1 := &bls.SecretKey{}
	sk1.SetByCSPRNG()

	sk2 := &bls.SecretKey{}
	sk2.SetByCSPRNG()

	sk3 := &bls.SecretKey{}
	sk3.SetByCSPRNG()

	t.Run("valid sig", func(t *testing.T) {
		nodes := []*Node{
			{
				NodeID: 1,
				PubKey: sk1.GetPublicKey().Serialize(),
			},
			{
				NodeID: 2,
				PubKey: sk2.GetPublicKey().Serialize(),
			},
			{
				NodeID: 3,
				PubKey: sk3.GetPublicKey().Serialize(),
			},
		}

		agg := sk1.SignByte(computedRoot)
		agg.Add(sk2.SignByte(computedRoot))
		agg.Add(sk3.SignByte(computedRoot))

		msgRoot.Signature = agg.Serialize()
		msgRoot.signers = []NodeID{1, 2, 3}

		require.NoError(t, Signature(agg.Serialize()).VerifyByNodes(msgRoot, domain, sigType, nodes))
	})

	t.Run("missing id", func(t *testing.T) {
		nodes := []*Node{
			{
				NodeID: 1,
				PubKey: sk1.GetPublicKey().Serialize(),
			},
			{
				NodeID: 2,
				PubKey: sk2.GetPublicKey().Serialize(),
			},
		}

		agg := sk1.SignByte(computedRoot)
		agg.Add(sk2.SignByte(computedRoot))
		agg.Add(sk3.SignByte(computedRoot))

		msgRoot.Signature = agg.Serialize()
		msgRoot.signers = []NodeID{1, 2, 3}

		require.EqualError(t, Signature(agg.Serialize()).VerifyByNodes(msgRoot, domain, sigType, nodes), "signer not found in nodes")
	})
}
