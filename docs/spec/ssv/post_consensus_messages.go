package ssv

import (
	"bytes"
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

func (v *Validator) processPostConsensusSig(dutyRunner *DutyRunner, signedMsg *SignedPostConsensusMessage) error {
	postCons := dutyRunner.PostConsensusStateForHeight(signedMsg.message.Height)
	if postCons == nil {
		return errors.New("PostConsensusMessage Height doesn't match duty runner's Height'")
	}

	if err := v.validatePostConsensusPartialSig(postCons, signedMsg); err != nil {
		return errors.Wrap(err, "partial sig invalid")
	}

	postCons.AddPartialSig(signedMsg.message)

	if !postCons.HasPostConsensusSigQuorum() {
		return nil
	}

	// if finished, no need to proceed with reconstructing the DutySignature
	if postCons.IsFinished() {
		return nil
	}
	postCons.SetFinished()

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

func (v *Validator) validatePostConsensusPartialSig(executionState *dutyExecutionState, SignedMsg *SignedPostConsensusMessage) error {
	if err := SignedMsg.GetSignature().VerifyByNodes(SignedMsg, v.share.domainType, types.PostConsensusSigType, v.share.GetQBFTCommittee()); err != nil {
		return errors.Wrap(err, "failed to verify DutySignature")
	}

	// validate signing root equal to decided
	if !bytes.Equal(executionState.postConsensusSigRoot, SignedMsg.message.DutySigningRoot) {
		return errors.New("pos consensus message signing root is wrong")
	}

	// TODO verify actual sig with signing root

	return nil
}

func (v *Validator) signPostConsensusMsg(msg *PostConsensusMessage) (*SignedPostConsensusMessage, error) {
	signature, err := v.signer.SignRoot(msg, types.PostConsensusSigType, v.share.pubKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not compute PostConsensusMessage root")
	}

	return &SignedPostConsensusMessage{
		message:   msg,
		signature: signature,
		signers:   []types.NodeID{v.share.GetNodeID()},
	}, nil
}
