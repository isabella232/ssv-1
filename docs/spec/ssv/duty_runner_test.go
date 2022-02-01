package ssv

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDutyExecutionState_AddPartialSig(t *testing.T) {
	t.Run("add to empty", func(t *testing.T) {
		s := NewTestingDutyExecutionState()
		s.AddPartialSig(&testingPostConsensusSigMessage{
			signerID: 1,
		})

		require.Len(t, s.collectedPartialSigs, 1)
	})

	t.Run("add multiple", func(t *testing.T) {
		s := NewTestingDutyExecutionState()
		s.AddPartialSig(&testingPostConsensusSigMessage{
			signerID: 1,
		})
		s.AddPartialSig(&testingPostConsensusSigMessage{
			signerID: 2,
		})
		s.AddPartialSig(&testingPostConsensusSigMessage{
			signerID: 3,
		})

		require.Len(t, s.collectedPartialSigs, 3)
	})

	t.Run("add duplicate", func(t *testing.T) {
		s := NewTestingDutyExecutionState()
		s.AddPartialSig(&testingPostConsensusSigMessage{
			signerID: 1,
		})
		s.AddPartialSig(&testingPostConsensusSigMessage{
			signerID: 1,
		})
		s.AddPartialSig(&testingPostConsensusSigMessage{
			signerID: 3,
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
			collectedPartialSigs: make(map[qbft.NodeID][]byte),
			runningInstance:      inst,
			quorumCount:          3,
			decidedValue: &consensusData{
				Duty:            duty,
				AttestationData: nil,
			},
		}
		dr.dutyExecutionState.AddPartialSig(&testingPostConsensusSigMessage{signerID: 1, sig: []byte{1, 2, 3, 4}})
		dr.dutyExecutionState.AddPartialSig(&testingPostConsensusSigMessage{signerID: 2, sig: []byte{1, 2, 3, 4}})
		dr.dutyExecutionState.AddPartialSig(&testingPostConsensusSigMessage{signerID: 3, sig: []byte{1, 2, 3, 4}})
		err := dr.CanStartNewDuty(&beacon.Duty{
			Type:   beacon.RoleTypeAttester,
			Slot:   12,
			PubKey: testingValidatorPK,
		})
		require.NoError(t, err)
	})
}
