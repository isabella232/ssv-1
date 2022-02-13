package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type testSigningRoot struct {
	root []byte
}

func (r *testSigningRoot) GetRoot() []byte {
	return r.root
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
