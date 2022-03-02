package types

type Encoder interface {
	// Encode returns a msg encoded bytes or error
	Encode() ([]byte, error)
	// Decode returns error if decoding failed
	Decode(data []byte) error
}
