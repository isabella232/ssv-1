package spectest

import (
	"github.com/bloxapp/ssv/docs/spec/ssv"
	"github.com/bloxapp/ssv/docs/spec/types"
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
