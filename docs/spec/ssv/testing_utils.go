package ssv

import "github.com/bloxapp/ssv/docs/spec/qbft"

type testingPostConsensusSigMessage struct {
	height   uint64
	signerID qbft.NodeID
	sig      []byte
}

func (tpcsm *testingPostConsensusSigMessage) GetHeight() uint64 {
	return tpcsm.height
}

func (tpcsm *testingPostConsensusSigMessage) GetSignerID() qbft.NodeID {
	return tpcsm.signerID
}

func (tpcsm *testingPostConsensusSigMessage) GetSig() []byte {
	return tpcsm.sig
}

func NewTestingDutyExecutionState() *dutyExecutionState {
	return &dutyExecutionState{
		collectedPartialSigs: make(map[qbft.NodeID][]byte),
	}
}
