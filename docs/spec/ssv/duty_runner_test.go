package ssv

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDutyExecutionState_AddPartialSig(t *testing.T) {
	t.Run("add to empty", func(t *testing.T) {
		s := newTestingDutyExecutionState()
		s.AddPartialSig(&PostConsensusMessage{
			Signers: []types.OperatorID{1},
		})

		require.Len(t, s.CollectedPartialSigs, 1)
	})

	t.Run("add multiple", func(t *testing.T) {
		s := newTestingDutyExecutionState()
		s.AddPartialSig(&PostConsensusMessage{
			Signers: []types.OperatorID{1},
		})
		s.AddPartialSig(&PostConsensusMessage{
			Signers: []types.OperatorID{2},
		})
		s.AddPartialSig(&PostConsensusMessage{
			Signers: []types.OperatorID{3},
		})

		require.Len(t, s.CollectedPartialSigs, 3)
	})

	t.Run("add duplicate", func(t *testing.T) {
		s := newTestingDutyExecutionState()
		s.AddPartialSig(&PostConsensusMessage{
			Signers: []types.OperatorID{1},
		})
		s.AddPartialSig(&PostConsensusMessage{
			Signers: []types.OperatorID{1},
		})
		s.AddPartialSig(&PostConsensusMessage{
			Signers: []types.OperatorID{3},
		})

		require.Len(t, s.CollectedPartialSigs, 2)
	})
}

func TestDutyRunner_CanStartNewDuty(t *testing.T) {
	t.Run("no prev start", func(t *testing.T) {
		dr := newTestingDutyRunner()
		err := dr.CanStartNewDuty(&beacon.Duty{
			Type: beacon.RoleTypeAttester,
		})
		require.NoError(t, err)
	})

	t.Run("running instance", func(t *testing.T) {
		dr := newTestingDutyRunner()
		inst := newTestingQBFTInstance()
		inst.Decided = false
		dr.State.DutyExecutionState = &DutyExecutionState{
			RunningInstance: inst,
		}
		err := dr.CanStartNewDuty(&beacon.Duty{
			Type:   beacon.RoleTypeAttester,
			PubKey: testingValidatorPK,
		})
		require.EqualError(t, err, "consensus on duty is running")
	})

	t.Run("Decided but still collecting sigs", func(t *testing.T) {
		dr := newTestingDutyRunner()
		duty := &beacon.Duty{
			Type:   beacon.RoleTypeAttester,
			Slot:   12,
			PubKey: testingValidatorPK,
		}
		inst := newTestingQBFTInstance()
		inst.Decided = true
		dr.State.DutyExecutionState = &DutyExecutionState{
			RunningInstance: inst,
			Quorum:          3,
			DecidedValue: &consensusData{
				Duty:            duty,
				AttestationData: nil,
			},
		}
		err := dr.CanStartNewDuty(&beacon.Duty{
			Type:   beacon.RoleTypeAttester,
			Slot:   12,
			PubKey: testingValidatorPK,
		})
		require.EqualError(t, err, "post consensus sig collection is running")
	})

	t.Run("Decided, not collected enough sigs but passed PostConsensusSigCollectionSlotTimeout slots", func(t *testing.T) {
		dr := newTestingDutyRunner()
		duty := &beacon.Duty{
			Type:   beacon.RoleTypeAttester,
			Slot:   12,
			PubKey: testingValidatorPK,
		}
		inst := newTestingQBFTInstance()
		inst.Decided = true
		dr.State.DutyExecutionState = &DutyExecutionState{
			RunningInstance: inst,
			Quorum:          3,
			DecidedValue: &consensusData{
				Duty:            duty,
				AttestationData: nil,
			},
		}
		err := dr.CanStartNewDuty(&beacon.Duty{
			Type:   beacon.RoleTypeAttester,
			Slot:   12 + PostConsensusSigCollectionSlotTimeout,
			PubKey: testingValidatorPK,
		})
		require.EqualError(t, err, "post consensus sig collection is running")
	})

	t.Run("Decided, not collected enough sigs but passed > PostConsensusSigCollectionSlotTimeout slots", func(t *testing.T) {
		dr := newTestingDutyRunner()
		duty := &beacon.Duty{
			Type:   beacon.RoleTypeAttester,
			Slot:   12,
			PubKey: testingValidatorPK,
		}
		inst := newTestingQBFTInstance()
		inst.Decided = true
		dr.State.DutyExecutionState = &DutyExecutionState{
			RunningInstance: inst,
			Quorum:          3,
			DecidedValue: &consensusData{
				Duty:            duty,
				AttestationData: nil,
			},
		}
		err := dr.CanStartNewDuty(&beacon.Duty{
			Type:   beacon.RoleTypeAttester,
			Slot:   12 + PostConsensusSigCollectionSlotTimeout + 1,
			PubKey: testingValidatorPK,
		})
		require.NoError(t, err)
	})

	t.Run("Decided, collected enough sigs", func(t *testing.T) {
		dr := newTestingDutyRunner()
		duty := &beacon.Duty{
			Type:   beacon.RoleTypeAttester,
			Slot:   12,
			PubKey: testingValidatorPK,
		}
		inst := newTestingQBFTInstance()
		inst.Decided = true
		dr.State.DutyExecutionState = &DutyExecutionState{
			CollectedPartialSigs: make(map[types.OperatorID][]byte),
			RunningInstance:      inst,
			Quorum:               3,
			DecidedValue: &consensusData{
				Duty:            duty,
				AttestationData: nil,
			},
		}
		dr.State.DutyExecutionState.AddPartialSig(&PostConsensusMessage{Signers: []types.OperatorID{1}, DutySignature: []byte{1, 2, 3, 4}})
		dr.State.DutyExecutionState.AddPartialSig(&PostConsensusMessage{Signers: []types.OperatorID{2}, DutySignature: []byte{1, 2, 3, 4}})
		dr.State.DutyExecutionState.AddPartialSig(&PostConsensusMessage{Signers: []types.OperatorID{3}, DutySignature: []byte{1, 2, 3, 4}})
		err := dr.CanStartNewDuty(&beacon.Duty{
			Type:   beacon.RoleTypeAttester,
			Slot:   12,
			PubKey: testingValidatorPK,
		})
		require.NoError(t, err)
	})
}

func TestDutyRunner_StartNewInstance(t *testing.T) {
	t.Run("value nil", func(t *testing.T) {
		dr := newTestingDutyRunner()
		require.EqualError(t, dr.StartNewInstance(nil), "new instance value nil")
	})

	t.Run("valid start", func(t *testing.T) {
		dr := newTestingDutyRunner()
		require.NoError(t, dr.StartNewInstance([]byte{1, 2, 3, 4}))
		require.NotNil(t, dr.State.DutyExecutionState)
		require.NotNil(t, dr.State.DutyExecutionState.RunningInstance)
		require.EqualValues(t, 3, dr.State.DutyExecutionState.Quorum)
	})
}

func TestDutyRunner_PostConsensusStateForHeight(t *testing.T) {
	t.Run("no return", func(t *testing.T) {
		dr := newTestingDutyRunner()
		require.Nil(t, dr.PostConsensusStateForHeight(10))
	})

	t.Run("returns", func(t *testing.T) {
		dr := newTestingDutyRunner()
		require.NoError(t, dr.StartNewInstance([]byte{1, 2, 3, 4}))
		require.NotNil(t, dr.PostConsensusStateForHeight(1))
	})
}

func TestDutyRunner_DecideRunningInstance(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		dr := newTestingDutyRunner()
		dr.State.DutyExecutionState = &DutyExecutionState{
			CollectedPartialSigs: make(map[types.OperatorID][]byte),
			Quorum:               3,
		}
		decidedValue := &consensusData{
			Duty: &beacon.Duty{
				Type:   beacon.RoleTypeAttester,
				Slot:   12,
				PubKey: testingValidatorPK,
			},
			AttestationData: nil,
		}
		_, err := dr.DecideRunningInstance(decidedValue, &testingKeyManager{})
		require.NoError(t, err)
		require.NotNil(t, dr.State.DutyExecutionState.DecidedValue)
		require.NotNil(t, dr.State.DutyExecutionState.SignedAttestation)
		require.NotNil(t, dr.State.DutyExecutionState.PostConsensusSigRoot)
		require.NotNil(t, dr.State.DutyExecutionState.CollectedPartialSigs)
	})
}
