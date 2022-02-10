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
	GetSignature() []byte
	GetSigners() []NodeID
	// IsValidSignature returns true if signature is valid (against message and signers)
	IsValidSignature(domain DomainType, nodes []*Node) error
	// MatchedSigners returns true if the provided signer ids are equal to GetSignerIds() without order significance
	MatchedSigners(ids []NodeID) bool
	// Aggregate will aggregate the signed message if possible (unique signers, same digest, valid)
	Aggregate(signedMsg MessageSignature) error
}
