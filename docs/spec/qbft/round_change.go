package qbft

import (
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

func uponRoundChange(state State, signedRoundChange *SignedMessage, roundChangeMsgContainer MsgContainer, valCheck types.ValueCheck) error {
	// TODO - Roberto comment: could happen we received a round change before we switched the round and this msg will be rejected (lost)
	if err := validRoundChange(state, signedRoundChange, state.GetHeight(), state.GetRound()); err != nil {
		return errors.Wrap(err, "round change msg invalid")
	}
	if !roundChangeMsgContainer.AddIfDoesntExist(signedRoundChange) {
		return nil // uponCommit was already called
	}

	if highestJustifiedRoundChangeMsg := hasReceivedProposalJustification(state, signedRoundChange, roundChangeMsgContainer, valCheck); highestJustifiedRoundChangeMsg != nil {
		// check if this node is the proposer
		if proposer(state, highestJustifiedRoundChangeMsg.Message.Round) != state.GetConfig().GetID() {
			return nil
		}

		proposal := createProposal(
			state,
			highestJustifiedRoundChangeMsg.Message.GetRoundChangeData().GetNextProposalData(),
			roundChangeMsgContainer.MessagesForHeightAndRound(state.GetHeight(), state.GetRound()), // TODO - might be optimized to include only necessary quorum
			highestJustifiedRoundChangeMsg.Message.GetRoundChangeData().GetRoundChangeJustification(),
		)
		if err := state.GetConfig().GetP2PNetwork().BroadcastSignedMessage(proposal); err != nil {
			return errors.Wrap(err, "failed to broadcast proposal message")
		}
	} else if partialQuorum, rcs := hasReceivedPartialQuorum(state, roundChangeMsgContainer); partialQuorum {
		newRound := minRound(rcs)

		state.SetRound(newRound)
		state.SetProposalAcceptedForCurrentRound(nil)

		roundChange := createRoundChange(state, newRound)
		if err := state.GetConfig().GetP2PNetwork().BroadcastSignedMessage(roundChange); err != nil {
			return errors.Wrap(err, "failed to broadcast round change message")
		}
	}
	return nil
}

func hasReceivedPartialQuorum(state State, roundChangeMsgContainer MsgContainer) (bool, []*SignedMessage) {
	all := roundChangeMsgContainer.AllMessagedForHeight(state.GetHeight())

	rc := make([]*SignedMessage, 0)
	for _, msg := range all {
		if msg.Message.Round > state.GetRound() {
			rc = append(rc, msg)
		}
	}

	return state.GetConfig().HasPartialQuorum(rc), rc
}

func hasReceivedProposalJustification(
	state State,
	signedRoundChange *SignedMessage,
	roundChangeMsgContainer MsgContainer,
	valCheck types.ValueCheck,
) *SignedMessage {
	roundChanges := roundChangeMsgContainer.MessagesForHeightAndRound(state.GetHeight(), state.GetRound())

	// TODO - optimization, if no round change quorum can return false

	// Important!
	// We iterate on all round chance msgs for liveliness in case the last round change msg is malicious.
	for _, msg := range roundChanges {
		prepares := msg.Message.GetRoundChangeData().GetRoundChangeJustification()
		if isReceivedProposalJustification(
			state,
			roundChanges,
			prepares,
			signedRoundChange.Message.Round,
			msg.Message.GetRoundChangeData().GetNextProposalData(),
			valCheck,
		) != nil {
			return msg
		}
	}
	return nil
}

// isReceivedProposalJustification - returns nil if we have a quorum of round change msgs and highest justified value
func isReceivedProposalJustification(
	state State,
	roundChanges, prepares []*SignedMessage,
	newRound Round,
	value []byte,
	valCheck types.ValueCheck,
) error {
	if err := isProposalJustification(
		state,
		roundChanges,
		prepares,
		state.GetHeight(),
		newRound,
		value,
		valCheck,
		proposer(state, newRound),
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

func validRoundChange(state State, signedMsg *SignedMessage, height uint64, round Round) error {
	if signedMsg.Message.MsgType != RoundChangeType {
		return errors.New("round change msg type is wrong")
	}
	if signedMsg.Message.Height != height {
		return errors.New("round change height is wrong")
	}
	if signedMsg.Message.Round != round {
		return errors.New("round change round is wrong")
	}
	if !signedMsg.IsValidSignature(state.GetConfig().GetNodes()) {
		return errors.New("round change msg signature invalid")
	}
	if signedMsg.Message.GetRoundChangeData().GetPreparedRound() == NoRound &&
		signedMsg.Message.GetRoundChangeData().GetPreparedValue() == nil {
		return nil
	} else if signedMsg.Message.GetRoundChangeData().GetPreparedRound() != NoRound &&
		signedMsg.Message.GetRoundChangeData().GetPreparedValue() != nil {

		// TODO Roberto comment: we should add a validation for justification data (sigs and so on)

		if signedMsg.Message.GetRoundChangeData().GetPreparedRound() < round {
			return nil
		}
		return errors.New("prepared round >= round")
	}
	return errors.New("round change prepare round & value are wrong")
}

// highestPrepared returns a round change message with the highest prepared round, returns nil if none found
func highestPrepared(roundChanges []*SignedMessage) *SignedMessage {
	var ret *SignedMessage
	for _, rc := range roundChanges {
		if rc.Message.GetRoundChangeData().GetPreparedRound() == NoRound &&
			rc.Message.GetRoundChangeData().GetPreparedValue() == nil {
			continue
		}

		if ret == nil {
			ret = rc
		} else if ret.Message.GetRoundChangeData().GetPreparedRound() < rc.Message.GetRoundChangeData().GetPreparedRound() {
			ret = rc
		}
	}
	return ret
}

func minRound(roundChangeMsgs []*SignedMessage) Round {
	panic("implement")
}

func createRoundChange(state State, newRound Round) *SignedMessage {
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
