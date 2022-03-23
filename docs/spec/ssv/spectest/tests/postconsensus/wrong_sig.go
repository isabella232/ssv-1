package postconsensus

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// InvaliSignature tests an invalid SignedPostConsensusMessage sig
func InvaliSignature() *tests.SpecTest {
	dr := testingutils.DecidedRunner()

	msgs := []*types.SSVMessage{
		testingutils.SSVMsg(nil, testingutils.PostConsensusAttestationMsg(testingutils.TestingSK1, 2, qbft.FirstHeight)),
	}

	return &tests.SpecTest{
		Name:                    "Invalid SignedPostConsensusMessage signature",
		DutyRunner:              dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "efacc0e9d5b27b8eab613b8d1ef33738f165edbc6be218dda66a92db0b5a71d2",
		ExpectedError:           "partial sig invalid: failed to verify DutySignature: failed to verify signature",
	}
}
