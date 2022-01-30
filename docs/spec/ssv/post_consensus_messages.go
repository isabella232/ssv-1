package ssv

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/network"
	"github.com/pkg/errors"
)

func (v *Validator) processPostConsensusSig(dutyRunner *DutyRunner, sigMsg PostConsensusSigMessage) error {
	postCons := dutyRunner.PostConsensusStateForHeight(sigMsg.GetHeight())
	if postCons == nil {
		return errors.New("PostConsensusSigMessage height doesn't match duty runner's height'")
	}

	if postCons.collectedPartialSigs[sigMsg.GetSignerID()] == nil {
		if err := v.verifyPostConsensusPartialSig(dutyRunner, sigMsg); err != nil {
			return errors.Wrap(err, "partial sig invalid")
		}
		postCons.collectedPartialSigs[sigMsg.GetSignerID()] = sigMsg.GetSig()
	}

	if !postCons.HasPostConsensusSigQuorum() {
		return nil
	}

	switch dutyRunner.runningDuty.Type {
	case beacon.RoleTypeAttester:
		att, err := postCons.ReconstructAttestationSig()
		if err != nil {
			return errors.Wrap(err, "could not reconstruct post consensus sig")
		}
		if err := v.beacon.SubmitAttestation(att); err != nil {
			return errors.Wrap(err, "could not submit to beacon chain reconstructed attestation")
		}
	default:
		return errors.Errorf("unknown duty post consensus sig %s", dutyRunner.runningDuty.Type.String())
	}
	return nil
}

func (v *Validator) verifyPostConsensusPartialSig(dutyRunner *DutyRunner, sigMsg PostConsensusSigMessage) error {
	panic("implement")
}

func (v *Validator) createPartialSigMsg(sig []byte) network.Message {
	panic("implement")
}
