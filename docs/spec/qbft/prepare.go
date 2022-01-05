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
	if state.GetConfig().HasQuorum([]SignedMessage{validPrepares}) {
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
		if err := validSignedPrepareForHeightRoundAndValue(signedMsg, height, round, identifier, nodes); err == nil {
			if aggregatedPrepareMsg == nil {
				aggregatedPrepareMsg = signedMsg
			} else {
				aggregatedPrepareMsg.Aggregate(signedMsg)
			}
		}
	}
	return aggregatedPrepareMsg
}

// validSignedPrepareForHeightRoundAndValue known in dafny spec as validSignedPrepareForHeightRoundAndDigest
// https://entethalliance.github.io/client-spec/qbft_spec.html#dfn-qbftspecification
func validSignedPrepareForHeightRoundAndValue(
	signedPrepare SignedMessage,
	height uint64,
	round Round,
	value []byte,
	nodes []Node) error {
	if signedPrepare.GetMessage().GetHeight() != height {
		return errors.New("msg height wrong")
	}
	if signedPrepare.GetMessage().GetRound() != round {
		return errors.New("msg round wrong")
	}
	if bytes.Compare(signedPrepare.GetMessage().GetPrepareData().GetData(), value) != 0 {
		return errors.New("msg identifier wrong")
	}
	if !signedPrepare.IsValidSignature(nodes) {
		return errors.New("prepare msg signature invalid")
	}
	return nil
}
