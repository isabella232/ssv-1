package testingutils

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
)

var TestConsensusData = &types.ConsensusData{
	Duty:            TestingDuty,
	AttestationData: TestingAttestationData,
}
var TestConsensusDataByts, _ = TestConsensusData.Encode()

var SSVMsg = func(qbftMsg *qbft.SignedMessage, postMsg *ssv.SignedPostConsensusMessage) *types.SSVMessage {
	var msgType types.MsgType
	var data []byte
	if qbftMsg != nil {
		msgType = types.SSVConsensusMsgType
		data, _ = qbftMsg.Encode()
	} else if postMsg != nil {
		msgType = types.SSVPostConsensusMsgType
		data, _ = postMsg.Encode()
	} else {
		panic("msg type undefined")
	}

	msgID := types.MessageIDForValidatorPKAndRole(TestingValidatorPubKey[:], beacon.RoleTypeAttester)

	return &types.SSVMessage{
		MsgType: msgType,
		MsgID:   msgID,
		Data:    data,
	}
}

var PostConsensusAttestationMsgWithWrongRoot = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *ssv.SignedPostConsensusMessage {
	return postConsensusAttestationMsg(sk, id, height, true)
}

var PostConsensusAttestationMsg = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *ssv.SignedPostConsensusMessage {
	return postConsensusAttestationMsg(sk, id, height, false)
}

var postConsensusAttestationMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	height qbft.Height,
	wrongRoot bool,
) *ssv.SignedPostConsensusMessage {
	signer := NewTestingKeyManager()
	signedAtt, root, _ := signer.SignAttestation(TestingAttestationData, TestingDuty, sk.GetPublicKey().Serialize())

	if wrongRoot {
		root = []byte{1, 2, 3, 4}
	}

	postConsensusMsg := &ssv.PostConsensusMessage{
		Height:          height,
		DutySignature:   signedAtt.Signature[:],
		DutySigningRoot: root,
		Signers:         []types.OperatorID{id},
	}

	sig, _ := signer.SignRoot(postConsensusMsg, types.PostConsensusSigType, sk.GetPublicKey().Serialize())
	return &ssv.SignedPostConsensusMessage{
		Message:   postConsensusMsg,
		Signature: sig,
		Signers:   []types.OperatorID{id},
	}
}
