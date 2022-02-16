package ssv

import (
	"encoding/hex"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidatorID_MessageIDBelongs(t *testing.T) {
	t.Run("msg id belongs", func(t *testing.T) {
		msgID := []byte{1, 2, 3, 4, 1, 0, 0, 0}
		valID := types.ValidatorID{1, 2, 3, 4}
		require.True(t, valID.MessageIDBelongs(msgID))
	})

	t.Run("msg id doesn't belong", func(t *testing.T) {
		msgID := []byte{1, 2, 3, 4, 1, 0, 0, 0}
		valID := types.ValidatorID{1, 2, 3, 3}
		require.False(t, valID.MessageIDBelongs(msgID))
	})

	t.Run("msg id doesn't belong", func(t *testing.T) {
		msgID := []byte{1, 2, 3, 4, 1, 0, 0, 0}
		valID := types.ValidatorID{1, 2, 3, 4, 4}
		require.False(t, valID.MessageIDBelongs(msgID))
	})
}

func TestConsensusData_Marshaling(t *testing.T) {
	expected, _ := hex.DecodeString("7b226174745f64617461223a2230313030303030303030303030303030303330303030303030303030303030303031303230333034303530363037303830393061303130323033303430353036303730383039306130313032303330343035303630373038303930613031303230313030303030303030303030303030303130323033303430353036303730383039306130313032303330343035303630373038303930613031303230333034303530363037303830393061303130323032303030303030303030303030303030313032303330343035303630373038303930613031303230333034303530363037303830393061303130323033303430353036303730383039306130313032222c2264757479223a223762323235343739373036353232336133313263323235303735363234623635373932323361356233313263333232633333326333343263333532633336326333373263333832633339326333313330326333313263333232633333326333343263333532633336326333373263333832633339326333313330326333313263333232633333326333343263333532633336326333373263333832633339326333313330326333313263333232633333326333343263333532633336326333373263333832633339326333313330326333313263333232633333326333343263333532633336326333373263333835643263323235333663366637343232336133313332326332323536363136633639363436313734366637323439366536343635373832323361333132633232343336663664366436393734373436353635343936653634363537383232336133323332326332323433366636643664363937343734363536353463363536653637373436383232336133303263323234333666366436643639373437343635363537333431373435333663366637343232336133333336326332323536363136633639363436313734366637323433366636643664363937343734363536353439366536343635373832323361333133313764227d")

	t.Run("attestation data", func(t *testing.T) {
		c := testConsensusData
		byts, err := c.Encode()
		require.NoError(t, err)
		require.EqualValues(t, expected, byts)
	})

	t.Run("marshal with attestation data", func(t *testing.T) {
		c := &consensusData{}
		require.NoError(t, c.Decode(expected))
		require.EqualValues(t, beacon.RoleTypeAttester, c.Duty.Type)
		require.EqualValues(t, spec.BLSPubKey{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8}, c.Duty.PubKey)
		require.EqualValues(t, spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2}, c.AttestationData.BeaconBlockRoot)
	})
}
