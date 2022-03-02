package types

import (
	"encoding/json"
)

// Share holds all info about the QBFT/ SSV committee for msg signing and verification
type Share struct {
	operatorID OperatorID
	pubKey     ValidatorPK
	committee  []*Operator
	quorum     uint64
	domainType DomainType
}

// GetOperatorID returns the node ID for this share
func (share *Share) GetOperatorID() OperatorID {
	return share.operatorID
}

// GetValidatorPubKey returns the validator public key to which the share is associated with
func (share *Share) GetValidatorPubKey() ValidatorPK {
	return share.pubKey
}

func (share *Share) GetQBFTCommittee() []*Operator {
	return share.committee
}

func (share *Share) GetQuorumCount() uint64 {
	return share.quorum
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

func (share *Share) GetDomainType() DomainType {
	return share.domainType
}

func (share *Share) Encode() ([]byte, error) {
	return json.Marshal(share)
}

func (share *Share) Decode(data []byte) error {
	return json.Unmarshal(data, share)
}
