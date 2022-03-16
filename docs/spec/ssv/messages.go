package ssv

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

type PostConsensusMessage struct {
	Height          qbft.Height
	DutySignature   []byte // The beacon chain partial signature for a duty
	DutySigningRoot []byte // the root signed in DutySignature
	Signers         []types.OperatorID
}

// Encode returns a msg encoded bytes or error
func (pcsm *PostConsensusMessage) Encode() ([]byte, error) {
	return json.Marshal(pcsm)
}

// Decode returns error if decoding failed
func (pcsm *PostConsensusMessage) Decode(data []byte) error {
	return json.Unmarshal(data, pcsm)
}

func (pcsm *PostConsensusMessage) GetRoot() ([]byte, error) {
	marshaledRoot, err := pcsm.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode PostConsensusMessage")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret[:], nil
}

type SignedPostConsensusMessage struct {
	message   *PostConsensusMessage
	signature types.Signature
	signers   []types.OperatorID
}

// Encode returns a msg encoded bytes or error
func (spcsm *SignedPostConsensusMessage) Encode() ([]byte, error) {
	d := make(map[string]interface{})

	if spcsm.message != nil {
		data, err := spcsm.message.Encode()
		if err != nil {
			return nil, errors.Wrap(err, "could not encode PostConsensusMessage")
		}
		d["msg"] = hex.EncodeToString(data)
	}
	d["signature"] = hex.EncodeToString(spcsm.signature)
	d["signers"] = spcsm.signers

	return json.Marshal(d)
}

// Decode returns error if decoding failed
func (spcsm *SignedPostConsensusMessage) Decode(data []byte) error {
	d := make(map[string]interface{})
	if err := json.Unmarshal(data, &d); err != nil {
		return errors.Wrap(err, "could not unmarshal SignedPostConsensusMessage")
	}

	if d["msg"] != nil {
		pcm := &PostConsensusMessage{}
		dataByts, err := hex.DecodeString(d["msg"].(string))
		if err != nil {
			return errors.Wrap(err, "could not decode PostConsensusMessage message")
		}
		if err := pcm.Decode(dataByts); err != nil {
			return errors.Wrap(err, "could not unmarshal PostConsensusMessage")
		}
		spcsm.message = pcm
	}

	if d["signature"] != nil {
		dataByts, err := hex.DecodeString(d["signature"].(string))
		if err != nil {
			return errors.Wrap(err, "could not decode SignedPostConsensusMessage signature")
		}
		spcsm.signature = dataByts
	}

	if d["signers"] != nil {
		spcsm.signers = make([]types.OperatorID, 0)
		for _, s := range d["signers"].([]interface{}) {
			spcsm.signers = append(spcsm.signers, types.OperatorID(s.(float64)))
		}
	}

	return nil
}

func (spcsm *SignedPostConsensusMessage) GetSignature() types.Signature {
	return spcsm.signature
}

func (spcsm *SignedPostConsensusMessage) GetSigners() []types.OperatorID {
	return spcsm.signers
}

func (spcsm *SignedPostConsensusMessage) GetRoot() ([]byte, error) {
	return spcsm.message.GetRoot()
}

func (spcsm *SignedPostConsensusMessage) Aggregate(signedMsg types.MessageSignature) error {
	//if !bytes.Equal(spcsm.GetRoot(), signedMsg.GetRoot()) {
	//	return errors.New("can't aggregate msgs with different roots")
	//}
	//
	//// verify no matching Signers
	//for _, signerID := range spcsm.signers {
	//	for _, toMatchID := range signedMsg.GetSigners() {
	//		if signerID == toMatchID {
	//			return errors.New("signer IDs partially/ fully match")
	//		}
	//	}
	//}
	//
	//allSigners := append(spcsm.signers, signedMsg.GetSigners()...)
	//
	//// verify and aggregate
	//sig1, err := blsSig(spcsm.signature)
	//if err != nil {
	//	return errors.Wrap(err, "could not parse DutySignature")
	//}
	//
	//sig2, err := blsSig(signedMsg.GetSignature())
	//if err != nil {
	//	return errors.Wrap(err, "could not parse DutySignature")
	//}
	//
	//sig1.Add(sig2)
	//spcsm.signature = sig1.Serialize()
	//spcsm.signers = allSigners
	//return nil
	panic("implement")
}

// MatchedSigners returns true if the provided signer ids are equal to GetSignerIds() without order significance
func (spcsm *SignedPostConsensusMessage) MatchedSigners(ids []types.OperatorID) bool {
	toMatchCnt := make(map[types.OperatorID]int)
	for _, id := range ids {
		toMatchCnt[id]++
	}

	foundCnt := make(map[types.OperatorID]int)
	for _, id := range spcsm.GetSigners() {
		foundCnt[id]++
	}

	for id, cnt := range toMatchCnt {
		if cnt != foundCnt[id] {
			return false
		}
	}
	return true
}

func blsSig(sig []byte) (*bls.Sign, error) {
	ret := &bls.Sign{}
	if err := ret.Deserialize(sig); err != nil {
		return nil, errors.Wrap(err, "could not covert DutySignature byts to bls.sign")
	}
	return ret, nil
}
