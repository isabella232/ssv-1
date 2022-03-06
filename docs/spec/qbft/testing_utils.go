package qbft

import (
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/utils/threshold"
	"github.com/herumi/bls-eth-go-binary/bls"
)

var testingMessage = &Message{
	MsgType:    ProposalMsgType,
	Height:     1,
	Round:      1,
	Identifier: []byte{1, 2, 3, 4},
	Data:       []byte{1, 2, 3, 4},
}
var TestingSignedMsg = func() *SignedMessage {
	id := types.OperatorID(1)
	domain := types.PrimusTestnet
	sigType := types.QBFTSigType

	r, _ := types.ComputeSigningRoot(testingMessage, types.ComputeSignatureDomain(domain, sigType))
	sig := TestingSK.SignByte(r)

	return &SignedMessage{
		Message:   testingMessage,
		Signers:   []types.OperatorID{id},
		Signature: sig.Serialize(),
	}
}()
var TestingSK = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.SetByCSPRNG()
	return ret
}()
