package ssv

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/types"
)

type DutyRunnerState struct {
	BeaconRoleType beacon.RoleType
	Share          *types.Share
	// DutyExecutionState holds all relevant params for a full duty execution (consensus & post consensus)
	DutyExecutionState *DutyExecutionState
	QBFTController     qbft.IController
}
