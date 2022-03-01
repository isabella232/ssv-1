package types

type MessageEncoder interface {
	// Encode returns a msg encoded bytes or error
	Encode() ([]byte, error)
	// Decode returns error if decoding failed
	Decode(data []byte) error
}

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
