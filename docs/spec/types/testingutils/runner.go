package testingutils

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/ssv"
)

var BaseRunner = func() *ssv.DutyRunner {
	return ssv.NewDutyRunner(
		beacon.RoleTypeAttester,
		TestingShare,
		NewTestingQBFTController([]byte{1, 2, 3, 4}),
		NewTestingStorage(),
	)
}
