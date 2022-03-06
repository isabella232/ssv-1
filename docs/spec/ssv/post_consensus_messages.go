package ssv

import (
	"bytes"
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

func (v *Validator) processPostConsensusSig(dutyRunner *DutyRunner, signedMsg *SignedPostConsensusMessage) error {
	postCons := dutyRunner.PostConsensusStateForHeight(signedMsg.message.Height)
	if postCons == nil {
		return errors.New("PostConsensusMessage Height doesn't match duty runner's Height'")
	}

	if err := v.validatePostConsensusMsg(postCons, signedMsg); err != nil {
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

	switch dutyRunner.State.BeaconRoleType {
	case beacon.RoleTypeAttester:
		att, err := postCons.ReconstructAttestationSig()
		if err != nil {
			return errors.Wrap(err, "could not reconstruct post consensus sig")
		}
		if err := v.beacon.SubmitAttestation(att); err != nil {
			return errors.Wrap(err, "could not submit to beacon chain reconstructed attestation")
		}
	default:
		return errors.Errorf("unknown duty post consensus sig %s", dutyRunner.State.BeaconRoleType.String())
	}
	return nil
}

func (v *Validator) validatePostConsensusMsg(executionState *DutyExecutionState, SignedMsg *SignedPostConsensusMessage) error {
	if len(SignedMsg.GetSigners()) != 1 {
		return errors.New("SignedPostConsensusMessage allows 1 signer")
	}

	if err := SignedMsg.GetSignature().VerifyByOperators(SignedMsg, v.share.DomainType, types.PostConsensusSigType, v.share.Committee); err != nil {
		return errors.Wrap(err, "failed to verify DutySignature")
	}

	// validate signing root equal to Decided
	if !bytes.Equal(executionState.PostConsensusSigRoot, SignedMsg.message.DutySigningRoot) {
		return errors.New("post consensus message signing root is wrong")
	}

	if len(SignedMsg.message.Signers) != 1 {
		return errors.New("PostConsensusMessage allows 1 signer")
	}

	if err := v.verifyBeaconPartialSignature(SignedMsg.message); err != nil {
		return errors.Wrap(err, "could not verify beacon partial signature")
	}

	return nil
}

func (v *Validator) verifyBeaconPartialSignature(msg *PostConsensusMessage) error {
	signer := msg.Signers[0]
	signature := msg.DutySignature
	root := msg.DutySigningRoot

	for _, n := range v.share.Committee {
		if n.GetID() == signer {
			pk := &bls.PublicKey{}
			if err := pk.Deserialize(n.GetPublicKey()); err != nil {
				return errors.Wrap(err, "could not deserialized pk")
			}
			sig := &bls.Sign{}
			if err := sig.Deserialize(signature); err != nil {
				return errors.Wrap(err, "could not deserialized signature")
			}

			// protect nil root
			root = ensureRoot(root)
			// verify
			if !sig.VerifyByte(pk, root) {
				return errors.Errorf("could not verify signature from iBFT member %d", signer)
			}
			return nil
		}
	}
	return errors.New("beacon partial signature signer not found")
}

func (v *Validator) signPostConsensusMsg(msg *PostConsensusMessage) (*SignedPostConsensusMessage, error) {
	signature, err := v.signer.SignRoot(msg, types.PostConsensusSigType, v.share.PubKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not compute PostConsensusMessage root")
	}

	return &SignedPostConsensusMessage{
		message:   msg,
		signature: signature,
		signers:   []types.OperatorID{v.share.OperatorID},
	}, nil
}
