package testingutils

import (
	"encoding/hex"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/go-bitfield"
)

type testingKeyManager struct {
	keys   map[string]*bls.SecretKey
	domain types.DomainType
}

func NewTestingKeyManager() types.KeyManager {
	ret := &testingKeyManager{
		keys:   map[string]*bls.SecretKey{},
		domain: types.PrimusTestnet,
	}
	ret.AddShare(TestingSK1)
	ret.AddShare(TestingSK2)
	ret.AddShare(TestingSK3)
	ret.AddShare(TestingSK4)
	ret.AddShare(TestingWrongSK)
	return ret
}

// SignAttestation signs the given attestation
func (km *testingKeyManager) SignAttestation(data *spec.AttestationData, duty *beacon.Duty, pk []byte) (*spec.Attestation, []byte, error) {
	if k, found := km.keys[hex.EncodeToString(pk)]; found {
		sig := k.SignByte(TestingAttestationRoot)
		blsSig := spec.BLSSignature{}
		copy(blsSig[:], sig.Serialize())

		aggregationBitfield := bitfield.NewBitlist(duty.CommitteeLength)
		aggregationBitfield.SetBitAt(duty.ValidatorCommitteeIndex, true)

		return &spec.Attestation{
			AggregationBits: aggregationBitfield,
			Data:            data,
			Signature:       blsSig,
		}, TestingAttestationRoot, nil
	}
	return nil, nil, errors.New("pk not found")
}

// IsAttestationSlashable returns error if attestation is slashable
func (km *testingKeyManager) IsAttestationSlashable(data *spec.AttestationData) error {
	return nil
}

func (km *testingKeyManager) SignRoot(data types.Root, sigType types.SignatureType, pk []byte) (types.Signature, error) {
	if k, found := km.keys[hex.EncodeToString(pk)]; found {
		computedRoot, err := types.ComputeSigningRoot(data, types.ComputeSignatureDomain(km.domain, sigType))
		if err != nil {
			return nil, errors.Wrap(err, "could not sign root")
		}

		return k.SignByte(computedRoot).Serialize(), nil
	}
	return nil, errors.New("pk not found")
}

func (km *testingKeyManager) AddShare(shareKey *bls.SecretKey) error {
	km.keys[hex.EncodeToString(shareKey.GetPublicKey().Serialize())] = shareKey
	return nil
}
