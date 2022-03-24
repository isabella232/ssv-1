package postconsensus

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// MultipleSigners tests >1 SignedPostConsensusMessage Signers
func MultipleSigners() *tests.SpecTest {
	dr := testingutils.DecidedRunner()

	noSignerMsg := testingutils.PostConsensusAttestationMsg(testingutils.TestingSK1, 1, qbft.FirstHeight)
	noSignerMsg.Signers = []types.OperatorID{1, 2}
	msgs := []*types.SSVMessage{
		testingutils.SSVMsg(nil, noSignerMsg),
	}

	return &tests.SpecTest{
		Name:                    ">1 SignedPostConsensusMessage Signers",
		DutyRunner:              dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "efacc0e9d5b27b8eab613b8d1ef33738f165edbc6be218dda66a92db0b5a71d2",
		ExpectedError:           "partial sig invalid: SignedPostConsensusMessage allows 1 signer",
	}
}
