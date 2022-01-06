package qbft

import "github.com/pkg/errors"

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
	if err := state.GetConfig().GetNetwork().BroadcastSignedMessage(prepare); err != nil {
		return errors.Wrap(err, "failed to broadcast prepare message")
	}

	state.SetProposalAcceptedForCurrentRound(signedProposal)
	if signedProposal.GetMessage().GetRound() > state.GetRound() {
		state.GetConfig().GetTimer().TimeoutForRound(signedProposal.GetMessage().GetRound())
	}
	state.SetRound(newRound)

	return nil
}

func isValidProposal(state State, signedProposal SignedMessage, valCheck ValueCheck) error {
	if signedProposal.GetMessage().GetType() != ProposalType {
		return errors.New("msg type is not proposal")
	}
	if signedProposal.GetMessage().GetHeight() != state.GetHeight() {
		return errors.New("proposal height is wrong")
	}
	if err := isProposalJustification(
		state,
		signedProposal,
		state.GetHeight(),
		signedProposal.GetMessage().GetRound(),
		proposer(state),
		valCheck,
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
// roundChangeMsgContainer MUST contain only relevant msgs (identifier, height and round)
// roundChangeMsgContainer MUST contain only relevant msgs (identifier, height and round)
func isProposalJustification(
	state State,
	signedProposal SignedMessage,
	height uint64,
	round Round,
	leaderID NodeID,
	valCheck ValueCheck,
) error {
	roundChangeMsgs := signedProposal.GetMessage().GetProposalData().GetRoundChangeJustification()
	prepareMsgs := signedProposal.GetMessage().GetProposalData().GetPrepareJustification()

	if !signedProposal.MatchedSigners([]NodeID{leaderID}) {
		return errors.New("proposal leader invalid")
	}
	if err := valCheck.Check(signedProposal.GetMessage().GetProposalData().GetData()); err != nil {
		return errors.Wrap(err, "proposal value invalid")
	}
	if signedProposal.GetMessage().GetRound() != round {
		return errors.New("proposal round is wrong")
	}
	if signedProposal.GetMessage().GetHeight() != height {
		return errors.New("proposal height is wrong")
	}

	if round == FirstRound {
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

func proposer(state State) NodeID {
	panic("implement")
}

func createPrepare(state State, newRound Round, value []byte) SignedMessage {
	/**
	Prepare(
	                    signPrepare(
	                        UnsignedPrepare(
	                            |current.blockchain|,
	                            newRound,
	                            digest(m.proposedBlock)),
	                        current.id
	                        )
	                );
	*/
	panic("implement")
}
