package ssv_test

import (
	"github.com/bloxapp/ssv/docs/spec/ssv"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

var testConsensusData = &types.ConsensusData{
	Duty:            testingutils.TestingDuty,
	AttestationData: testingutils.TestingAttestationData,
}
var TestConsensusDataByts, _ = testConsensusData.Encode()

func NewTestingDutyExecutionState() *ssv.DutyExecutionState {
	return &ssv.DutyExecutionState{
		CollectedPartialSigs: make(map[types.OperatorID][]byte),
	}
}
