package types

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMessageIDForValidatorPKAndRole(t *testing.T) {
	require.EqualValues(t, []byte{1, 2, 3, 4, 1, 0, 0, 0}, MessageIDForValidatorPKAndRole([]byte{1, 2, 3, 4}, beacon.RoleTypeAttester))
}

func TestMessageID_GetRoleType(t *testing.T) {
	t.Run("attester", func(t *testing.T) {
		msgID := MessageIDForValidatorPKAndRole([]byte{1, 2, 3, 4}, beacon.RoleTypeAttester)
		require.EqualValues(t, beacon.RoleTypeAttester, msgID.GetRoleType())
	})

	t.Run("proposer", func(t *testing.T) {
		msgID := MessageIDForValidatorPKAndRole([]byte{1, 2, 3, 4}, beacon.RoleTypeProposer)
		require.EqualValues(t, beacon.RoleTypeProposer, msgID.GetRoleType())
	})

	t.Run("long pk", func(t *testing.T) {
		msgID := MessageIDForValidatorPKAndRole([]byte{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4}, beacon.RoleTypeProposer)
		require.EqualValues(t, beacon.RoleTypeProposer, msgID.GetRoleType())
	})
}
