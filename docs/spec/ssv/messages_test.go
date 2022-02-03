package ssv

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMessageIDForValidatorPKAndRole(t *testing.T) {
	require.EqualValues(t, []byte{1, 2, 3, 4, 1, 0, 0, 0}, types.MessageIDForValidatorPKAndRole([]byte{1, 2, 3, 4}, beacon.RoleTypeAttester))
}
