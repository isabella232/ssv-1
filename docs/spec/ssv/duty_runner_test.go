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
			Signers: []types.NodeID{1},
		})

		require.Len(t, s.collectedPartialSigs, 1)
	})

	t.Run("add multiple", func(t *testing.T) {
		s := newTestingDutyExecutionState()
		s.AddPartialSig(&PostConsensusMessage{
			Signers: []types.NodeID{1},
		})
		s.AddPartialSig(&PostConsensusMessage{
			Signers: []types.NodeID{2},
		})
		s.AddPartialSig(&PostConsensusMessage{
			Signers: []types.NodeID{3},
		})

		require.Len(t, s.collectedPartialSigs, 3)
	})

	t.Run("add duplicate", func(t *testing.T) {
		s := newTestingDutyExecutionState()
		s.AddPartialSig(&PostConsensusMessage{
			Signers: []types.NodeID{1},
		})
		s.AddPartialSig(&PostConsensusMessage{
			Signers: []types.NodeID{1},
		})
		s.AddPartialSig(&PostConsensusMessage{
			Signers: []types.NodeID{3},
		})

		require.Len(t, s.collectedPartialSigs, 2)
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
		inst.decided = false
		dr.dutyExecutionState = &dutyExecutionState{
			runningInstance: inst,
		}
		err := dr.CanStartNewDuty(&beacon.Duty{
			Type:   beacon.RoleTypeAttester,
			PubKey: testingValidatorPK,
		})
		require.EqualError(t, err, "consensus on duty is running")
	})

	t.Run("decided but still collecting sigs", func(t *testing.T) {
		dr := newTestingDutyRunner()
		duty := &beacon.Duty{
			Type:   beacon.RoleTypeAttester,
			Slot:   12,
			PubKey: testingValidatorPK,
		}
		inst := newTestingQBFTInstance()
		inst.decided = true
		dr.dutyExecutionState = &dutyExecutionState{
			runningInstance: inst,
			quorumCount:     3,
			decidedValue: &consensusData{
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

	t.Run("decided, not collected enough sigs but passed PostConsensusSigCollectionSlotTimeout slots", func(t *testing.T) {
		dr := newTestingDutyRunner()
		duty := &beacon.Duty{
			Type:   beacon.RoleTypeAttester,
			Slot:   12,
			PubKey: testingValidatorPK,
		}
		inst := newTestingQBFTInstance()
		inst.decided = true
		dr.dutyExecutionState = &dutyExecutionState{
			runningInstance: inst,
			quorumCount:     3,
			decidedValue: &consensusData{
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

	t.Run("decided, not collected enough sigs but passed > PostConsensusSigCollectionSlotTimeout slots", func(t *testing.T) {
		dr := newTestingDutyRunner()
		duty := &beacon.Duty{
			Type:   beacon.RoleTypeAttester,
			Slot:   12,
			PubKey: testingValidatorPK,
		}
		inst := newTestingQBFTInstance()
		inst.decided = true
		dr.dutyExecutionState = &dutyExecutionState{
			runningInstance: inst,
			quorumCount:     3,
			decidedValue: &consensusData{
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

	t.Run("decided, collected enough sigs", func(t *testing.T) {
		dr := newTestingDutyRunner()
		duty := &beacon.Duty{
			Type:   beacon.RoleTypeAttester,
			Slot:   12,
			PubKey: testingValidatorPK,
		}
		inst := newTestingQBFTInstance()
		inst.decided = true
		dr.dutyExecutionState = &dutyExecutionState{
			collectedPartialSigs: make(map[types.NodeID][]byte),
			runningInstance:      inst,
			quorumCount:          3,
			decidedValue: &consensusData{
				Duty:            duty,
				AttestationData: nil,
			},
		}
		dr.dutyExecutionState.AddPartialSig(&PostConsensusMessage{Signers: []types.NodeID{1}, DutySignature: []byte{1, 2, 3, 4}})
		dr.dutyExecutionState.AddPartialSig(&PostConsensusMessage{Signers: []types.NodeID{2}, DutySignature: []byte{1, 2, 3, 4}})
		dr.dutyExecutionState.AddPartialSig(&PostConsensusMessage{Signers: []types.NodeID{3}, DutySignature: []byte{1, 2, 3, 4}})
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
		require.NotNil(t, dr.dutyExecutionState)
		require.EqualValues(t, 1, dr.dutyExecutionState.height)
		require.NotNil(t, dr.dutyExecutionState.runningInstance)
		require.EqualValues(t, 3, dr.dutyExecutionState.quorumCount)
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
		dr.dutyExecutionState = &dutyExecutionState{
			collectedPartialSigs: make(map[types.NodeID][]byte),
			quorumCount:          3,
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
		require.NotNil(t, dr.dutyExecutionState.decidedValue)
		require.NotNil(t, dr.dutyExecutionState.signedAttestation)
		require.NotNil(t, dr.dutyExecutionState.postConsensusSigRoot)
		require.NotNil(t, dr.dutyExecutionState.collectedPartialSigs)
	})
}
