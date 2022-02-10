package types

import (
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/beacon"
	"github.com/herumi/bls-eth-go-binary/bls"
)

type SignatureType byte

var (
	QBFTSigType          = byte(1)
	PostConsensusSigType = byte(2)
)

type BeaconSigner interface {
	// SignAttestation signs the given attestation
	SignAttestation(data *spec.AttestationData, duty *beacon.Duty, pk []byte) (*spec.Attestation, []byte, error)
}

type Signer interface {
	SignRoot(root, sigType, pk []byte) ([]byte, error)
}

// KeyManager is an interface responsible for all key manager functions
type KeyManager interface {
	BeaconSigner
	Signer
	// AddShare saves a share key
	AddShare(shareKey *bls.SecretKey) error
}

// SSVKeyManager implements the KeyManager interface with all of its funcs
type SSVKeyManager struct {
	keys map[string]*bls.SecretKey // holds pub keys as key and secret key as value
}

// SignAttestation signs the given attestation
func (s *SSVKeyManager) SignAttestation(data *spec.AttestationData, duty *beacon.Duty, pk []byte) (*spec.Attestation, []byte, error) {
	panic("implement from beacon ")
}

func (s *SSVKeyManager) SignRoot(root, sigType, pk []byte) ([]byte, error) {

}

// AddShare saves a share key
func (s *SSVKeyManager) AddShare(sk *bls.SecretKey) error {
	s.keys[sk.GetPublicKey().GetHexString()] = sk
	return nil
}
