package qbft

import (
	"bytes"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

func uponPrepare(state State, signedPrepare SignedMessage, prepareMsgContainer, commitMsgContainer MsgContainer) error {
	// TODO - if we receive a prepare before a proposal and return an error we will never process the prepare msg, we still need to add it to the container
	if state.GetProposalAcceptedForCurrentRound() == nil {
		return errors.New("not proposal accepted for prepare")
	}

	if err := validSignedPrepareForHeightRoundAndValue(
		signedPrepare,
		state.GetHeight(),
		state.GetRound(),
		state.GetProposalAcceptedForCurrentRound().GetMessage().GetProposalData().GetData(),
		state.GetConfig().GetNodes(),
	); err != nil {
		return errors.Wrap(err, "invalid prepare msg")
	}

	if !prepareMsgContainer.AddIfDoesntExist(signedPrepare) {
		return nil // uponPrepare was already called
	}

	if !state.GetConfig().HasQuorum(prepareMsgContainer.MessagesForHeightAndRound(state.GetHeight(), state.GetRound())) {
		return nil // no quorum yet
	}

	if didSendCommitForHeightAndRound(state, commitMsgContainer) {
		return nil // already moved to commit stage
	}

	proposedValue := state.GetProposalAcceptedForCurrentRound().GetMessage().GetProposalData().GetData()
	commitMsg := createCommit(state, proposedValue)
	if err := state.GetConfig().GetP2PNetwork().BroadcastSignedMessage(commitMsg); err != nil {
		return errors.Wrap(err, "failed to broadcast commit message")
	}

	state.SetLastPreparedValue(proposedValue)
	state.SetLastPreparedRound(state.GetRound())
	return nil
}

func getRoundChangeJustification(state State, prepareMsgContainer MsgContainer) SignedMessage {
	if state.GetLastPreparedValue() == nil {
		return nil
	}

	prepareMsgs := prepareMsgContainer.MessagesForHeightAndRound(state.GetHeight(), state.GetRound())
	validPrepares := validPreparesForHeightRoundAndDigest(
		prepareMsgs,
		state.GetHeight(),
		state.GetLastPreparedRound(),
		state.GetLastPreparedValue(),
		state.GetConfig().GetNodes(),
	)
	if state.GetConfig().HasQuorum(prepareMsgs) {
		return validPrepares
	}
	return nil
}

// validPreparesForHeightRoundAndDigest returns an aggregated prepare msg for a specific height and round
func validPreparesForHeightRoundAndDigest(
	prepareMessages []SignedMessage,
	height uint64,
	round Round,
	value []byte,
	nodes []*types.Node) SignedMessage {
	var aggregatedPrepareMsg SignedMessage
	for _, signedMsg := range prepareMessages {
		if err := validSignedPrepareForHeightRoundAndValue(signedMsg, height, round, value, nodes); err == nil {
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
	nodes []*types.Node) error {
	if signedPrepare.GetMessage().GetType() != PrepareType {
		return errors.New("prepare msg type is wrong")
	}
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
