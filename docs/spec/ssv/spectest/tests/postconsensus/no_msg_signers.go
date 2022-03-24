package postconsensus

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// NoMessageSigners tests an empty PostConsensusMessage Signers
func NoMessageSigners() *tests.SpecTest {
	dr := testingutils.DecidedRunner()

	msgs := []*types.SSVMessage{
		testingutils.SSVMsg(nil, testingutils.PostConsensusAttestationMsgWithMsgSigners(testingutils.TestingSK1, 1, qbft.FirstHeight)),
	}

	return &tests.SpecTest{
		Name:                    "NoSigners PostConsensusMessage",
		DutyRunner:              dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "efacc0e9d5b27b8eab613b8d1ef33738f165edbc6be218dda66a92db0b5a71d2",
		ExpectedError:           "partial sig invalid: could not verify beacon partial Signature: PostConsensusMessage allows 1 signer",
	}
}
