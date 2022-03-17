package spectest

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/ssv"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

var AllTests = []*SpecTest{
	happyFullFlow(),
}

type SpecTest struct {
	Name                    string
	DutyRunner              *ssv.DutyRunner
	Messages                []*types.SSVMessage
	PostDutyRunnerStateRoot string
	ExpectedError           error
}

var baseRunner = func() *ssv.DutyRunner {
	return ssv.NewDutyRunner(
		beacon.RoleTypeAttester,
		testingutils.TestingValidatorPubKey[:],
		nil,
		ssv.NewTestingQBFTController([]byte{1, 2, 3, 4}),
		ssv.NewTestingStorage(),
	)
}
var BaseValidator = func() *ssv.Validator {
	ret := ssv.NewValidator(
		ssv.NewTestingNetwork(),
		ssv.NewTestingBeaconNode(),
		ssv.NewTestingStorage(),
		nil,
		ssv.NewTestingKeyManager(),
		func(data []byte) error {
			return nil
		},
	)
	return ret
}
