package qbft

import (
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

func uponRoundChange(state State, signedRoundChange SignedMessage, roundChangeMsgContainer MsgContainer, valCheck types.ValueCheck) error {
	if err := validRoundChange(state, signedRoundChange, state.GetHeight(), state.GetRound()); err != nil {
		return errors.Wrap(err, "round change msg invalid")
	}
	if !roundChangeMsgContainer.AddIfDoesntExist(signedRoundChange) {
		return nil // uponCommit was already called
	}

	if hasReceivedProposalJustification(state, signedRoundChange, roundChangeMsgContainer, valCheck) {
		var value []byte
		if state.GetLastPreparedValue() != nil {
			value = state.GetLastPreparedValue()
		} else {
			// TODO - set to start value
		}

		proposal := createProposal(state, value)
		if err := state.GetConfig().GetNetwork().BroadcastSignedMessage(proposal); err != nil {
			return errors.Wrap(err, "failed to broadcast proposal message")
		}

		state.SetRound(111) // TODO - why do we set round? and if so, what is the value of newRound from the spec?
		state.SetProposalAcceptedForCurrentRound(nil)
	} else if hasReceivedPartialQuorum(state, roundChangeMsgContainer) {
		newRound := minRound(roundChangeMsgContainer.MessagesForHeightAndRound(signedRoundChange.GetMessage().GetHeight(), signedRoundChange.GetMessage().GetRound()))

		roundChange := createRoundChange(state, newRound)
		if err := state.GetConfig().GetNetwork().BroadcastSignedMessage(roundChange); err != nil {
			return errors.Wrap(err, "failed to broadcast round change message")
		}

		state.SetRound(newRound) // TODO - why do we set round?
		state.SetProposalAcceptedForCurrentRound(nil)
	}
	return nil
}

func hasReceivedPartialQuorum(state State, roundChangeMsgContainer MsgContainer) bool {
	rc := roundChangeMsgContainer.MessagesForHeightAndRound(state.GetHeight(), state.GetRound())
	return state.GetConfig().HasPartialQuorum(rc)
}

func hasReceivedProposalJustification(
	state State,
	signedRoundChange SignedMessage,
	roundChangeMsgContainer MsgContainer,
	valCheck types.ValueCheck,
) bool {
	roundChanges := roundChangeMsgContainer.MessagesForHeightAndRound(state.GetHeight(), state.GetRound())
	prepares := signedRoundChange.GetMessage().GetRoundChangeData().GetRoundChangeJustification()
	return isReceivedProposalJustification(
		state,
		roundChanges,
		prepares,
		signedRoundChange.GetMessage().GetRound(),
		signedRoundChange.GetMessage().GetRoundChangeData().GetNextProposalData(),
		valCheck,
	) != nil
}

func isReceivedProposalJustification(
	state State,
	roundChanges, prepares []SignedMessage,
	newRound Round,
	value []byte,
	valCheck types.ValueCheck,
) error {
	/**
			&& roundChanges <= receivedRoundChanges(current)
	        && prepares <= receivedPrepares(current)
	TODO - not sure what does this check?
	*/

	if err := isProposalJustification(
		state,
		roundChanges,
		prepares,
		state.GetHeight(),
		state.GetRound(),
		value,
		valCheck,
		state.GetConfig().GetID(), // checks if this node is the leader
	); err != nil {
		return errors.Wrap(err, "round change ")
	}

	noPrevProposal := state.GetProposalAcceptedForCurrentRound() == nil && state.GetRound() == newRound
	prevProposal := state.GetProposalAcceptedForCurrentRound() != nil && newRound > state.GetRound()

	if !noPrevProposal && !prevProposal {
		return errors.New("prev proposal and new round mismatch")
	}
	return nil
}

func validRoundChange(state State, signedMsg SignedMessage, height uint64, round Round) error {
	if signedMsg.GetMessage().GetType() != RoundChangeType {
		return errors.New("round change msg type is wrong")
	}
	if signedMsg.GetMessage().GetHeight() != height {
		return errors.New("round change height is wrong")
	}
	if signedMsg.GetMessage().GetRound() != round {
		return errors.New("round change round is wrong")
	}
	if !signedMsg.IsValidSignature(state.GetConfig().GetNodes()) {
		return errors.New("round change msg signature invalid")
	}

	if signedMsg.GetMessage().GetRoundChangeData().GetPreparedRound() == NoRound &&
		signedMsg.GetMessage().GetRoundChangeData().GetPreparedValue() == nil {
		return nil
	} else if signedMsg.GetMessage().GetRoundChangeData().GetPreparedRound() != NoRound &&
		signedMsg.GetMessage().GetRoundChangeData().GetPreparedValue() != nil {
		if signedMsg.GetMessage().GetRoundChangeData().GetPreparedRound() < round {
			return nil
		}
		return errors.New("prepared round >= round")
	}
	return errors.New("round change prepare round & value are wrong")
}

// highestPrepared returns a round change message with the highest prepared round, returns nil if none found
func highestPrepared(roundChanges []SignedMessage) SignedMessage {
	var ret SignedMessage
	for _, rc := range roundChanges {
		if rc.GetMessage().GetRoundChangeData().GetPreparedRound() == NoRound &&
			rc.GetMessage().GetRoundChangeData().GetPreparedValue() == nil {
			continue
		}

		if ret == nil {
			ret = rc
		} else if ret.GetMessage().GetRoundChangeData().GetPreparedRound() < rc.GetMessage().GetRoundChangeData().GetPreparedRound() {
			ret = rc
		}
	}
	return ret
}

func minRound(roundChangeMsgs []SignedMessage) Round {
	panic("implement")
}

func createRoundChange(state State, newRound Round) SignedMessage {
	/**
	RoundChange(
	           signRoundChange(
	               UnsignedRoundChange(
	                   |current.blockchain|,
	                   newRound,
	                   digestOptionalBlock(current.lastPreparedBlock),
	                   current.lastPreparedRound),
	           current.id),
	           current.lastPreparedBlock,
	           getRoundChangeJustification(current)
	       )
	*/
	panic("implement")
}
