package qbft

import (
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

func uponProposal(state State, signedProposal *SignedMessage, proposeMsgContainer MsgContainer) error {
	valCheck := state.GetConfig().GetValueCheck()
	if err := isValidProposal(state, signedProposal, valCheck); err != nil {
		return errors.New("proposal invalid")
	}
	if !proposeMsgContainer.AddIfDoesntExist(signedProposal) {
		return nil // uponProposal was already called
	}

	newRound := signedProposal.Message.Round

	// set state to new round and proposal accepted
	state.SetProposalAcceptedForCurrentRound(signedProposal)
	if signedProposal.Message.Round > state.GetRound() {
		state.GetConfig().GetTimer().TimeoutForRound(signedProposal.Message.Round)
	}
	state.SetRound(newRound)

	prepare := createPrepare(state, newRound, signedProposal.Message.GetProposalData().GetData())
	if err := state.GetConfig().GetP2PNetwork().BroadcastSignedMessage(prepare); err != nil {
		return errors.Wrap(err, "failed to broadcast prepare message")
	}

	return nil
}

func isValidProposal(state State, signedProposal *SignedMessage, valCheck types.BeaconDataCheck) error {
	if signedProposal.Message.MsgType != ProposalType {
		return errors.New("msg type is not proposal")
	}
	if signedProposal.Message.Height != state.GetHeight() {
		return errors.New("proposal height is wrong")
	}
	// TODO - Roberto comment: we should check signedProposal sig (added here https://github.com/ConsenSys/qbft-formal-spec-and-verification/blob/main/dafny/spec/L1/node_auxiliary_functions.dfy#L573)
	if !signedProposal.MatchedSigners([]types.NodeID{proposer(state, signedProposal.Message.Round)}) {
		return errors.New("proposal leader invalid")
	}
	if err := isProposalJustification(
		state,
		signedProposal.Message.GetProposalData().GetRoundChangeJustification(),
		signedProposal.Message.GetProposalData().GetPrepareJustification(),
		state.GetHeight(),
		signedProposal.Message.Round,
		signedProposal.Message.GetProposalData().GetData(),
		valCheck,
		signedProposal.Signers[0], // already verified sig so we know there is 1 signer
	); err != nil {
		return errors.Wrap(err, "proposal not justified")
	}

	if (state.GetProposalAcceptedForCurrentRound() == nil && signedProposal.Message.Round == state.GetRound()) ||
		(state.GetProposalAcceptedForCurrentRound() != nil && signedProposal.Message.Round > state.GetRound()) {
		return nil
	}
	return errors.New("proposal is not valid with current state")
}

// isProposalJustification returns nil if the signed proposal msg is justified
func isProposalJustification(
	state State,
	roundChangeMsgs []*SignedMessage,
	prepareMsgs []*SignedMessage,
	height uint64,
	round Round,
	value []byte,
	valCheck types.BeaconDataCheck,
	roundLeader types.NodeID,
) error {
	if err := valCheck.Check(value); err != nil {
		return errors.Wrap(err, "proposal value invalid")
	}

	if round == FirstRound {
		if proposer(state, round) != roundLeader {
			return errors.New("round leader is wrong")
		}
		return nil
	} else {
		if !state.GetConfig().HasQuorum(roundChangeMsgs) {
			return errors.New("change round has not quorum")
		}

		for _, rc := range roundChangeMsgs {
			if err := validRoundChange(state, rc, height, round); err != nil {
				return errors.Wrap(err, "change round msg not valid")
			}
		}

		previouslyPreparedF := func() bool {
			for _, rc := range roundChangeMsgs { // TODO - might be redundant as it's checked in validRoundChange
				if rc.Message.GetRoundChangeData().GetPreparedRound() != NoRound &&
					rc.Message.GetRoundChangeData().GetPreparedValue() != nil {
					return true
				}
			}
			return false
		}

		if !previouslyPreparedF() {
			if proposer(state, round) != roundLeader {
				return errors.New("round leader is wrong")
			}
			return nil
		} else {
			if !state.GetConfig().HasQuorum(prepareMsgs) {
				return errors.New("change round has not quorum")
			}

			rcm := highestPrepared(roundChangeMsgs)
			if rcm == nil {
				return errors.New("no highest prepared")
			}

			for _, pm := range prepareMsgs {
				if err := validSignedPrepareForHeightRoundAndValue(
					state,
					pm,
					height,
					rcm.Message.GetRoundChangeData().GetPreparedRound(),
					rcm.Message.GetRoundChangeData().GetPreparedValue(),
					state.GetConfig().GetNodes(),
				); err != nil {
					return errors.New("signed prepare not valid")
				}
			}
			return nil
		}
	}
}

func proposer(state State, round Round) types.NodeID {
	panic("implement")
}

func createProposal(state State, value []byte, roundChanged, prepares []*SignedMessage) (*SignedMessage, error) {
	/**
	  	Proposal(
	                        signProposal(
	                            UnsignedProposal(
	                                |current.blockchain|,
	                                newRound,
	                                digest(block)),
	                            current.id),
	                        block,
	                        extractSignedRoundChanges(roundChanges),
	                        extractSignedPrepares(prepares));
	*/
	panic("implement")
	msg := &Message{}
	sig, err := state.GetConfig().GetSigner().SignRoot(msg.GetRoot(), types.QBFTSigType, state.GetConfig().GetSigningPubKey())
	if err != nil {
		return nil, errors.Wrap(err, "failed signing proposal msg")
	}

	signedMsg := &SignedMessage{
		Signature: sig,
		Signers:   []types.NodeID{state.GetConfig().GetID()},
		Message:   msg,
	}
	return signedMsg, nil
}
