package types

// ValidatorPK is an eth2 validator public key
type ValidatorPK []byte

type Root interface {
	// GetRoot returns the root used for signing and verification
	GetRoot() ([]byte, error)
}

type MessageSignature interface {
	Root
	GetSignature() Signature
	GetSigners() []OperatorID
	// MatchedSigners returns true if the provided signer ids are equal to GetSignerIds() without order significance
	MatchedSigners(ids []OperatorID) bool
	// Aggregate will aggregate the signed message if possible (unique signers, same digest, valid)
	Aggregate(signedMsg MessageSignature) error
}
