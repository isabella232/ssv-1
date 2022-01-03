package qbft

import (
	"bytes"
	"github.com/pkg/errors"
)

func getRoundChangeJustification(state State, prepareMsgContainer MsgContainer) SignedMessage {
	if state.GetLastPreparedValue() == nil {
		return nil
	}

	validPrepares := validPreparesForHeightRoundAndDigest(
		prepareMsgContainer,
		state.GetHeight(),
		state.GetLastPreparedRound(),
		state.GetInstanceIdentifier(),
		state.GetConfig().GetNodes(),
	)
	if validPrepares.FullQuorumReached(state.GetConfig().GetNodes()) {
		return validPrepares
	}
	return nil
}

// validPreparesForHeightRoundAndDigest returns an aggregated prepare msg for a specific height and round
func validPreparesForHeightRoundAndDigest(
	prepareMsgContainer MsgContainer,
	height uint64,
	round Round,
	identifier []byte,
	nodes []Node) SignedMessage {
	iterator := prepareMsgContainer.Iterator()
	var aggregatedPrepareMsg SignedMessage
	for signedMsg := iterator.Next(); signedMsg != nil; {
		if err := validSignedPrepareForHeightRoundAndDigest(signedMsg, height, round, identifier, nodes); err == nil {
			if aggregatedPrepareMsg == nil {
				aggregatedPrepareMsg = signedMsg
			} else {
				aggregatedPrepareMsg.Aggregate(signedMsg)
			}
		}
	}
	return aggregatedPrepareMsg
}

func validSignedPrepareForHeightRoundAndDigest(
	signedMsg SignedMessage,
	height uint64,
	round Round,
	identifier []byte,
	nodes []Node) error {
	if signedMsg.GetMessage().GetHeight() != height {
		return errors.New("msg height wrong")
	}
	if signedMsg.GetMessage().GetRound() != round {
		return errors.New("msg round wrong")
	}
	if bytes.Compare(signedMsg.GetMessage().GetInstanceIdentifier(), identifier) != 0 {
		return errors.New("msg identifier wrong")
	}
	if !signedMsg.IsValidSignature(nodes) {
		return errors.New("msg signature invalid")
	}
	return nil
}
