package qbft

import (
	"bytes"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

func uponPrepare(state State, signedPrepare *SignedMessage, prepareMsgContainer, commitMsgContainer MsgContainer) error {
	// TODO - if we receive a prepare before a proposal and return an error we will never process the prepare msg, we still need to add it to the container
	if state.GetProposalAcceptedForCurrentRound() == nil {
		return errors.New("not proposal accepted for prepare")
	}

	if err := validSignedPrepareForHeightRoundAndValue(
		state,
		signedPrepare,
		state.GetHeight(),
		state.GetRound(),
		state.GetProposalAcceptedForCurrentRound().Message.GetProposalData().GetData(),
		state.GetConfig().GetOperators(),
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

	proposedValue := state.GetProposalAcceptedForCurrentRound().Message.GetProposalData().GetData()

	state.SetLastPreparedValue(proposedValue)
	state.SetLastPreparedRound(state.GetRound())

	commitMsg := createCommit(state, proposedValue)
	if err := state.GetConfig().GetNetwork().Broadcast(commitMsg); err != nil {
		return errors.Wrap(err, "failed to broadcast commit message")
	}

	return nil
}

func getRoundChangeJustification(state State, prepareMsgContainer MsgContainer) *SignedMessage {
	if state.GetLastPreparedValue() == nil {
		return nil
	}

	prepareMsgs := prepareMsgContainer.MessagesForHeightAndRound(state.GetHeight(), state.GetLastPreparedRound())
	validPrepares := validPreparesForHeightRoundAndDigest(
		state,
		prepareMsgs,
		state.GetHeight(),
		state.GetLastPreparedRound(),
		state.GetLastPreparedValue(),
		state.GetConfig().GetOperators(),
	)
	if state.GetConfig().HasQuorum(prepareMsgs) {
		return validPrepares
	}
	return nil
}

// validPreparesForHeightRoundAndDigest returns an aggregated prepare msg for a specific height and round
func validPreparesForHeightRoundAndDigest(
	state State,
	prepareMessages []*SignedMessage,
	height uint64,
	round Round,
	value []byte,
	operators []*types.Operator) *SignedMessage {
	var aggregatedPrepareMsg *SignedMessage
	for _, signedMsg := range prepareMessages {
		if err := validSignedPrepareForHeightRoundAndValue(state, signedMsg, height, round, value, operators); err == nil {
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
	state State,
	signedPrepare *SignedMessage,
	height uint64,
	round Round,
	value []byte,
	operators []*types.Operator) error {
	if signedPrepare.Message.MsgType != PrepareMsgType {
		return errors.New("prepare msg type is wrong")
	}
	if signedPrepare.Message.Height != height {
		return errors.New("msg height wrong")
	}
	if signedPrepare.Message.Round != round {
		return errors.New("msg round wrong")
	}
	if bytes.Compare(signedPrepare.Message.GetPrepareData().GetData(), value) != 0 {
		return errors.New("msg identifier wrong")
	}

	if err := signedPrepare.Signature.VerifyByOperators(signedPrepare, state.GetConfig().GetSignatureDomainType(), types.QBFTSigType, operators); err != nil {
		return errors.Wrap(err, "prepare msg signature invalid")
	}
	return nil
}

func createPrepare(state State, newRound Round, value []byte) *SignedMessage {
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
