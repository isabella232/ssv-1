package qbft

import (
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

func uponProposal(state State, config Config, signedProposal *SignedMessage, proposeMsgContainer MsgContainer) error {
	valCheck := config.GetValueCheck()
	if err := isValidProposal(state, config, signedProposal, valCheck, state.Share.Committee); err != nil {
		return errors.New("proposal invalid")
	}
	if !proposeMsgContainer.AddIfDoesntExist(signedProposal) {
		return nil // uponProposal was already called
	}

	newRound := signedProposal.Message.Round

	// set state to new round and proposal accepted
	state.ProposalAcceptedForCurrentRound = signedProposal
	if signedProposal.Message.Round > state.Round {
		config.GetTimer().TimeoutForRound(signedProposal.Message.Round)
	}
	state.Round = newRound

	prepare := createPrepare(state, newRound, signedProposal.Message.GetProposalData().GetData())
	if err := config.GetNetwork().Broadcast(prepare); err != nil {
		return errors.Wrap(err, "failed to broadcast prepare message")
	}

	return nil
}

func isValidProposal(
	state State,
	config Config,
	signedProposal *SignedMessage,
	valCheck proposedValueCheck,
	operators []*types.Operator,
) error {
	if signedProposal.Message.MsgType != ProposalMsgType {
		return errors.New("msg type is not proposal")
	}
	if signedProposal.Message.Height != state.Height {
		return errors.New("proposal Height is wrong")
	}
	if len(signedProposal.GetSigners()) != 1 {
		return errors.New("proposal msg allows 1 signer")
	}
	if err := signedProposal.Signature.VerifyByOperators(signedProposal, config.GetSignatureDomainType(), types.QBFTSigType, operators); err != nil {
		return errors.Wrap(err, "proposal msg signature invalid")
	}
	if !signedProposal.MatchedSigners([]types.OperatorID{proposer(state, signedProposal.Message.Round)}) {
		return errors.New("proposal leader invalid")
	}
	if err := isProposalJustification(
		state,
		config,
		signedProposal.Message.GetProposalData().GetRoundChangeJustification(),
		signedProposal.Message.GetProposalData().GetPrepareJustification(),
		state.Height,
		signedProposal.Message.Round,
		signedProposal.Message.GetProposalData().GetData(),
		valCheck,
		signedProposal.Signers[0], // already verified sig so we know there is 1 signer
	); err != nil {
		return errors.Wrap(err, "proposal not justified")
	}

	if (state.ProposalAcceptedForCurrentRound == nil && signedProposal.Message.Round == state.Round) ||
		(state.ProposalAcceptedForCurrentRound != nil && signedProposal.Message.Round > state.Round) {
		return nil
	}
	return errors.New("proposal is not valid with current state")
}

// isProposalJustification returns nil if the signed proposal msg is justified
func isProposalJustification(
	state State,
	config Config,
	roundChangeMsgs []*SignedMessage,
	prepareMsgs []*SignedMessage,
	height uint64,
	round Round,
	value []byte,
	valCheck proposedValueCheck,
	roundLeader types.OperatorID,
) error {
	if err := valCheck(value); err != nil {
		return errors.Wrap(err, "proposal value invalid")
	}

	if round == FirstRound {
		if proposer(state, round) != roundLeader {
			return errors.New("round leader is wrong")
		}
		return nil
	} else {
		if !state.Share.HasQuorum(len(roundChangeMsgs)) {
			return errors.New("change round has not quorum")
		}

		for _, rc := range roundChangeMsgs {
			if err := validRoundChange(state, config, rc, height, round); err != nil {
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
			if !state.Share.HasQuorum(len(prepareMsgs)) {
				return errors.New("change round has not quorum")
			}

			rcm := highestPrepared(roundChangeMsgs)
			if rcm == nil {
				return errors.New("no highest prepared")
			}

			for _, pm := range prepareMsgs {
				if err := validSignedPrepareForHeightRoundAndValue(
					state,
					config,
					pm,
					height,
					rcm.Message.GetRoundChangeData().GetPreparedRound(),
					rcm.Message.GetRoundChangeData().GetPreparedValue(),
					state.Share.Committee,
				); err != nil {
					return errors.New("signed prepare not valid")
				}
			}
			return nil
		}
	}
}

func proposer(state State, round Round) types.OperatorID {
	panic("implement")
}

func createProposal(state State, config Config, value []byte, roundChanged, prepares []*SignedMessage) (*SignedMessage, error) {
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
	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSigType, config.GetSigningPubKey())
	if err != nil {
		return nil, errors.Wrap(err, "failed signing proposal msg")
	}

	signedMsg := &SignedMessage{
		Signature: sig,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Message:   msg,
	}
	return signedMsg, nil
}
