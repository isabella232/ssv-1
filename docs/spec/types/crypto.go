package types

import (
	"crypto/sha256"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

type SignatureDomain []byte
type Signature []byte

func (s Signature) VerifyMultiPubKey(root []byte, domain DomainType, sigType SignatureType, pks [][]byte) error {
	var aggPK *bls.PublicKey
	for _, pkByts := range pks {
		pk := &bls.PublicKey{}
		if err := pk.Deserialize(pkByts); err != nil {
			return errors.Wrap(err, "failed to deserialize public key")
		}

		if aggPK == nil {
			aggPK = pk
		} else {
			aggPK.Add(pk)
		}
	}

	if aggPK == nil {
		return errors.New("no public keys found")
	}

	return s.Verify(root, domain, sigType, aggPK.Serialize())
}

func (s Signature) Verify(root []byte, domain DomainType, sigType SignatureType, pkByts []byte) error {
	sign := &bls.Sign{}
	if err := sign.Deserialize(s); err != nil {
		return errors.Wrap(err, "failed to deserialize signature")
	}

	pk := &bls.PublicKey{}
	if err := pk.Deserialize(pkByts); err != nil {
		return errors.Wrap(err, "failed to deserialize public key")
	}

	computedRoot := ComputeSigningRoot(root, ComputeSignatureDomain(domain, sigType))
	if res := sign.VerifyByte(pk, computedRoot); !res {
		return errors.New("failed to verify signature")
	}
	return nil
}

func ComputeSigningRoot(root []byte, domain SignatureDomain) []byte {
	ret := sha256.Sum256(append(root, domain...))
	return ret[:]
}

func ComputeSignatureDomain(domain DomainType, sigType SignatureType) SignatureDomain {
	return SignatureDomain(append(domain, sigType...))
}
