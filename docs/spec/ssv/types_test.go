package ssv

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidatorID_MessageIDBelongs(t *testing.T) {
	t.Run("msg id belongs", func(t *testing.T) {
		msgID := []byte{1, 2, 3, 4, 1, 0, 0, 0}
		valID := ValidatorID{1, 2, 3, 4}
		require.True(t, valID.MessageIDBelongs(msgID))
	})

	t.Run("msg id doesn't belong", func(t *testing.T) {
		msgID := []byte{1, 2, 3, 4, 1, 0, 0, 0}
		valID := ValidatorID{1, 2, 3, 3}
		require.False(t, valID.MessageIDBelongs(msgID))
	})

	t.Run("msg id doesn't belong", func(t *testing.T) {
		msgID := []byte{1, 2, 3, 4, 1, 0, 0, 0}
		valID := ValidatorID{1, 2, 3, 4, 4}
		require.False(t, valID.MessageIDBelongs(msgID))
	})
}
