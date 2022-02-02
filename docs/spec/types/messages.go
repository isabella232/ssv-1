package types

type MessageEncoder interface {
	// Encode returns a msg encoded bytes or error
	Encode() ([]byte, error)
	// Decode returns error if decoding failed
	Decode(data []byte) error
}

type MessageDigest interface {
	// Digest returns a digest of the msg including any msg specific data, used for signing and verification
	Digest() []byte
}
