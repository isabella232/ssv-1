package ssv

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateMsg(t *testing.T) {
	v := &Validator{
		share: &Share{
			pubKey: []byte{1, 2, 3, 4},
		},
		dutyRunners: DutyRunners{
			beacon.RoleTypeAttester: &DutyRunner{},
		},
	}
	t.Run("valid", func(t *testing.T) {
		msg := types.SSVMessage{
			MsgID:   types.MessageIDForValidatorPKAndRole(types.ValidatorPK{1, 2, 3, 4}, beacon.RoleTypeAttester),
			MsgType: types.SSVConsensusMsgType,
			Data:    []byte{1, 2, 3},
		}
		require.NoError(t, v.validateMessage(msg))
	})

	t.Run("wrong msg id", func(t *testing.T) {
		msg := types.SSVMessage{
			MsgID: types.MessageIDForValidatorPKAndRole(types.ValidatorPK{1, 2, 3, 3}, beacon.RoleTypeAttester),
		}
		require.EqualError(t, v.validateMessage(msg), "msg ID doesn't match validator ID")
	})

	t.Run("no duty runner", func(t *testing.T) {
		msg := types.SSVMessage{
			MsgID: types.MessageIDForValidatorPKAndRole(types.ValidatorPK{1, 2, 3, 4}, beacon.RoleTypeAggregator),
		}
		require.EqualError(t, v.validateMessage(msg), "could not get duty runner for msg ID 0102030402000000")
	})

	t.Run("unknown type", func(t *testing.T) {
		msg := types.SSVMessage{
			MsgID:   types.MessageIDForValidatorPKAndRole(types.ValidatorPK{1, 2, 3, 4}, beacon.RoleTypeAttester),
			MsgType: 5,
		}
		require.EqualError(t, v.validateMessage(msg), "msg type not supported")
	})

	t.Run("data nil", func(t *testing.T) {
		msg := types.SSVMessage{
			MsgID:   types.MessageIDForValidatorPKAndRole(types.ValidatorPK{1, 2, 3, 4}, beacon.RoleTypeAttester),
			MsgType: 0,
			Data:    []byte{},
		}
		require.EqualError(t, v.validateMessage(msg), "msg data is invalid")
	})
}
