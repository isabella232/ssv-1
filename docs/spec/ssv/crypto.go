package ssv

import (
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
)

func verifySig(sig *bls.Sign, signers []types.NodeID, nodes []*types.Node, root []byte) bool {
	pks := make([]bls.PublicKey, 0)
	for _, signer := range signers {
		for _, n := range nodes {
			if n.GetID() == signer {
				pk := bls.PublicKey{}
				if err := pk.Deserialize(n.GetPublicKey()); err == nil {
					pks = append(pks, pk)
				}
			}
		}
	}

	return sig.FastAggregateVerify(pks, root)
}
