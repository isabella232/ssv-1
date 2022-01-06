package qbft

import "github.com/pkg/errors"

func createRoundChange(state State) SignedMessage {
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
