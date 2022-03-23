package postconsensus

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// ValidMessage tests a full valid SignedPostConsensusMessage
func ValidMessage() *tests.SpecTest {
	dr := testingutils.DecidedRunner()

	msgs := []*types.SSVMessage{
		testingutils.SSVMsg(nil, testingutils.PostConsensusAttestationMsg(testingutils.TestingSK1, 1, qbft.FirstHeight)),
	}

	return &tests.SpecTest{
		Name:                    "valid SignedPostConsensusMessage",
		DutyRunner:              dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "7926d81e439abfb30e42d08aa6511008e037a87405e99dbc7000b82dd8cf3ce7",
	}
}
