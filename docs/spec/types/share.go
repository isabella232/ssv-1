package types

import (
	"encoding/json"
)

// Share holds all info about the QBFT/ SSV Committee for msg signing and verification
type Share struct {
	OperatorID    OperatorID
	PubKey        ValidatorPK
	Committee     []*Operator
	Quorum        uint64
	PartialQuorum uint64
	DomainType    DomainType
}

// HasQuorum returns true if at least 2f+1 items are present (cnt is the number of items). It assumes nothing about those items, not their type or structure
// https://github.com/ConsenSys/qbft-formal-spec-and-verification/blob/main/dafny/spec/L1/node_auxiliary_functions.dfy#L259
func (share *Share) HasQuorum(cnt int) bool {
	panic("implement")
}

// HasPartialQuorum returns true if at least f+1 items present (cnt is the number of items). It assumes nothing about those items, not their type or structure.
// https://github.com/ConsenSys/qbft-formal-spec-and-verification/blob/main/dafny/spec/L1/node_auxiliary_functions.dfy#L244
func (share *Share) HasPartialQuorum(cnt int) bool {
	panic("implement")
}

func (share *Share) Encode() ([]byte, error) {
	return json.Marshal(share)
}

func (share *Share) Decode(data []byte) error {
	return json.Unmarshal(data, share)
}
