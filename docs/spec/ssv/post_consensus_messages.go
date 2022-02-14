package ssv

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

func (v *Validator) processPostConsensusSig(dutyRunner *DutyRunner, sigMsg *PostConsensusSigMessage) error {
	if err := v.validatePostConsensusPartialSig(dutyRunner, sigMsg); err != nil {
		return errors.Wrap(err, "partial sig invalid")
	}

	postCons := dutyRunner.PostConsensusStateForHeight(sigMsg.GetHeight())
	if postCons == nil {
		return errors.New("PostConsensusSigMessage height doesn't match duty runner's height'")
	}
	postCons.AddPartialSig(sigMsg)

	if !postCons.HasPostConsensusSigQuorum() {
		return nil
	}

	switch dutyRunner.beaconRoleType {
	case beacon.RoleTypeAttester:
		att, err := postCons.ReconstructAttestationSig()
		if err != nil {
			return errors.Wrap(err, "could not reconstruct post consensus sig")
		}
		if err := v.beacon.SubmitAttestation(att); err != nil {
			return errors.Wrap(err, "could not submit to beacon chain reconstructed attestation")
		}
	default:
		return errors.Errorf("unknown duty post consensus sig %s", dutyRunner.beaconRoleType.String())
	}
	return nil
}

func (v *Validator) validatePostConsensusPartialSig(dutyRunner *DutyRunner, sigMsg *PostConsensusSigMessage) error {
	if err := sigMsg.signature.VerifyByNodes(sigMsg, v.share.domainType, types.PostConsensusSigType, v.share.GetQBFTCommittee()); err != nil {
		return errors.Wrap(err, "failed to verify signature")
	}
	return nil
}
