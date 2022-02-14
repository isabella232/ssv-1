package types

type MessageEncoder interface {
	// Encode returns a msg encoded bytes or error
	Encode() ([]byte, error)
	// Decode returns error if decoding failed
	Decode(data []byte) error
}

type MessageRoot interface {
	// GetRoot returns the root used for signing and verification
	GetRoot() []byte
}

type MessageSignature interface {
	MessageRoot
	GetSignature() Signature
	GetSigners() []NodeID
	// MatchedSigners returns true if the provided signer ids are equal to GetSignerIds() without order significance
	MatchedSigners(ids []NodeID) bool
	// Aggregate will aggregate the signed message if possible (unique signers, same digest, valid)
	Aggregate(signedMsg MessageSignature) error
}
