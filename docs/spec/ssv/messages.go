package ssv

import (
	"bytes"
	"encoding/json"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

type PostConsensusSigMessage struct {
	height    uint64
	signature types.Signature
	root      []byte
	signers   []types.NodeID
}

// Encode returns a msg encoded bytes or error
func (pcsm *PostConsensusSigMessage) Encode() ([]byte, error) {
	return json.Marshal(pcsm)
}

// Decode returns error if decoding failed
func (pcsm *PostConsensusSigMessage) Decode(data []byte) error {
	return json.Unmarshal(data, pcsm)
}

func (pcsm *PostConsensusSigMessage) GetSignature() []byte {
	return pcsm.signature
}

func (pcsm *PostConsensusSigMessage) GetSigners() []types.NodeID {
	return pcsm.signers
}

func (pcsm *PostConsensusSigMessage) Aggregate(signedMsg types.MessageSignature) error {
	if !bytes.Equal(pcsm.root, signedMsg.GetRoot()) {
		return errors.New("can't aggregate msgs with different roots")
	}

	// verify no matching signers
	for _, signerID := range pcsm.signers {
		for _, toMatchID := range signedMsg.GetSigners() {
			if signerID == toMatchID {
				return errors.New("signer IDs partially/ fully match")
			}
		}
	}

	allSigners := append(pcsm.signers, signedMsg.GetSigners()...)

	// verify and aggregate
	sig1, err := blsSig(pcsm.signature)
	if err != nil {
		return errors.Wrap(err, "could not parse signature")
	}

	sig2, err := blsSig(signedMsg.GetSignature())
	if err != nil {
		return errors.Wrap(err, "could not parse signature")
	}

	sig1.Add(sig2)
	pcsm.signature = sig1.Serialize()
	pcsm.signers = allSigners
	return nil
}

// MatchedSigners returns true if the provided signer ids are equal to GetSignerIds() without order significance
func (pcsm *PostConsensusSigMessage) MatchedSigners(ids []types.NodeID) bool {
	toMatchCnt := make(map[types.NodeID]int)
	for _, id := range ids {
		toMatchCnt[id]++
	}

	foundCnt := make(map[types.NodeID]int)
	for _, id := range pcsm.GetSigners() {
		foundCnt[id]++
	}

	for id, cnt := range toMatchCnt {
		if cnt != foundCnt[id] {
			return false
		}
	}
	return true
}

func (pcsm *PostConsensusSigMessage) GetHeight() uint64 {
	return pcsm.height
}

func (pcsm *PostConsensusSigMessage) GetRoot() []byte {
	return pcsm.root
}

func blsSig(sig []byte) (*bls.Sign, error) {
	ret := &bls.Sign{}
	if err := ret.Deserialize(sig); err != nil {
		return nil, errors.Wrap(err, "could not covert signature byts to bls.sign")
	}
	return ret, nil
}
