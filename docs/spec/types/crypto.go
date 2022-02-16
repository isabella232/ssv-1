package types

import (
	"crypto/sha256"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

type SignatureDomain []byte
type Signature []byte

func (s Signature) VerifyByNodes(data MessageSignature, domain DomainType, sigType SignatureType, nodes []*Node) error {
	pks := make([][]byte, 0)
	for _, id := range data.GetSigners() {
		found := false
		for _, n := range nodes {
			if id == n.GetID() {
				pks = append(pks, n.GetPublicKey())
				found = true
			}
		}
		if !found {
			return errors.New("signer not found in nodes")
		}
	}
	return s.VerifyMultiPubKey(data, domain, sigType, pks)
}

func (s Signature) VerifyMultiPubKey(data MessageRoot, domain DomainType, sigType SignatureType, pks [][]byte) error {
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

	return s.Verify(data, domain, sigType, aggPK.Serialize())
}

func (s Signature) Verify(data MessageRoot, domain DomainType, sigType SignatureType, pkByts []byte) error {
	sign := &bls.Sign{}
	if err := sign.Deserialize(s); err != nil {
		return errors.Wrap(err, "failed to deserialize signature")
	}

	pk := &bls.PublicKey{}
	if err := pk.Deserialize(pkByts); err != nil {
		return errors.Wrap(err, "failed to deserialize public key")
	}

	computedRoot, err := ComputeSigningRoot(data, ComputeSignatureDomain(domain, sigType))
	if err != nil {
		return errors.Wrap(err, "could not compute signing root")
	}
	if res := sign.VerifyByte(pk, computedRoot); !res {
		return errors.New("failed to verify signature")
	}
	return nil
}

func ComputeSigningRoot(data MessageRoot, domain SignatureDomain) ([]byte, error) {
	dataRoot, err := data.GetRoot()
	if err != nil {
		return nil, errors.Wrap(err, "could not get root from MessageRoot")
	}

	ret := sha256.Sum256(append(dataRoot, domain...))
	return ret[:], nil
}

func ComputeSignatureDomain(domain DomainType, sigType SignatureType) SignatureDomain {
	return SignatureDomain(append(domain, sigType...))
}
