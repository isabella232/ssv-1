package qbft

import (
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

func uponProposal(state State, signedProposal SignedMessage, proposeMsgContainer MsgContainer) error {
	valCheck := state.GetConfig().GetValueCheck()
	if err := isValidProposal(state, signedProposal, valCheck); err != nil {
		return errors.New("proposal invalid")
	}
	if !proposeMsgContainer.AddIfDoesntExist(signedProposal) {
		return nil // uponProposal was already called
	}

	newRound := signedProposal.GetMessage().GetRound()

	prepare := createPrepare(state, newRound, signedProposal.GetMessage().GetProposalData().GetData())
	if err := state.GetConfig().GetP2PNetwork().BroadcastSignedMessage(prepare); err != nil {
		return errors.Wrap(err, "failed to broadcast prepare message")
	}

	state.SetProposalAcceptedForCurrentRound(signedProposal)
	if signedProposal.GetMessage().GetRound() > state.GetRound() {
		state.GetConfig().GetTimer().TimeoutForRound(signedProposal.GetMessage().GetRound())
	}
	state.SetRound(newRound)

	return nil
}

func isValidProposal(state State, signedProposal SignedMessage, valCheck types.ValueCheck) error {
	if signedProposal.GetMessage().GetType() != ProposalType {
		return errors.New("msg type is not proposal")
	}
	if signedProposal.GetMessage().GetHeight() != state.GetHeight() {
		return errors.New("proposal height is wrong")
	}
	// TODO - where do we check signedProposal sig
	if !signedProposal.MatchedSigners([]types.NodeID{proposer(state, signedProposal.GetMessage().GetRound())}) {
		return errors.New("proposal leader invalid")
	}
	if err := isProposalJustification(
		state,
		signedProposal.GetMessage().GetProposalData().GetRoundChangeJustification(),
		signedProposal.GetMessage().GetProposalData().GetPrepareJustification(),
		state.GetHeight(),
		signedProposal.GetMessage().GetRound(),
		signedProposal.GetMessage().GetProposalData().GetData(),
		valCheck,
		signedProposal.GetSigners()[0], // already verified sig so we know there is 1 signer
	); err != nil {
		return errors.Wrap(err, "proposal not justified")
	}

	if (state.GetProposalAcceptedForCurrentRound() == nil && signedProposal.GetMessage().GetRound() == state.GetRound()) ||
		(state.GetProposalAcceptedForCurrentRound() != nil && signedProposal.GetMessage().GetRound() > state.GetRound()) {
		return nil
	}
	return errors.New("proposal is not valid with current state")
}

// isProposalJustification returns nil if the signed proposal msg is justified
func isProposalJustification(
	state State,
	roundChangeMsgs []SignedMessage,
	prepareMsgs []SignedMessage,
	height uint64,
	round Round,
	value []byte,
	valCheck types.ValueCheck,
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
				if rc.GetMessage().GetRoundChangeData().GetPreparedRound() != NoRound &&
					rc.GetMessage().GetRoundChangeData().GetPreparedValue() != nil {
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
					pm,
					height,
					rcm.GetMessage().GetRoundChangeData().GetPreparedRound(),
					rcm.GetMessage().GetRoundChangeData().GetPreparedValue(),
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

func createProposal(state State, value []byte) SignedMessage {
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
	panic("implementation")
}
