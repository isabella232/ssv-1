package qbft

import (
	"bytes"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

// uponCommit returns true if a quorum of commit messages was received.
func uponCommit(state State, config Config, signedCommit *SignedMessage, commitMsgContainer MsgContainer) (bool, []byte, *SignedMessage, error) {
	if state.ProposalAcceptedForCurrentRound == nil {
		return false, nil, nil, errors.New("did not receive proposal for this round")
	}

	if err := validateCommit(
		state,
		config,
		signedCommit,
		state.Height,
		state.Round,
		state.ProposalAcceptedForCurrentRound,
		state.Share.GetQBFTCommittee(),
	); err != nil {
		return false, nil, nil, errors.Wrap(err, "commit msg invalid")
	}
	if !commitMsgContainer.AddIfDoesntExist(signedCommit) {
		return false, nil, nil, nil // uponCommit was already called
	}

	value := signedCommit.Message.GetCommitData().GetData()
	if quorum, commitMsgs := commitQuorumForValue(state, config, commitMsgContainer, value); quorum {
		agg, err := aggregateCommitMsgs(commitMsgs)
		if err != nil {
			return false, nil, nil, errors.Wrap(err, "could not aggregate commit msgs")
		}
		return true, value, agg, nil
	}
	return false, nil, nil, nil
}

func commitQuorumForValue(state State, config Config, commitMsgContainer MsgContainer, value []byte) (bool, []*SignedMessage) {
	commitMsgs := commitMsgContainer.MessagesForHeightAndRound(state.Height, state.Round)
	valueFiltered := make([]*SignedMessage, 0)
	for _, msg := range commitMsgs {
		if bytes.Equal(msg.Message.GetCommitData().GetData(), value) {
			valueFiltered = append(valueFiltered, msg)
		}
	}

	return config.HasQuorum(valueFiltered), valueFiltered
}

func aggregateCommitMsgs(msgs []*SignedMessage) (*SignedMessage, error) {
	if len(msgs) == 0 {
		return nil, errors.New("can't aggregate zero commit msgs")
	}

	var ret *SignedMessage
	for _, m := range msgs {
		if ret == nil {
			ret = m.DeepCopy()
		} else {
			if err := ret.Aggregate(m); err != nil {
				return nil, errors.Wrap(err, "could not aggregate commit msg")
			}
		}
	}
	return ret, nil
}

// didSendCommitForHeightAndRound returns true if sent commit msg for specific height and round
func didSendCommitForHeightAndRound(state State, commitMsgContainer MsgContainer) bool {
	/**
	!exists m :: && m in current.messagesReceived
	                            && m.Commit?
	                            && var uPayload := m.commitPayload.unsignedPayload;
	                            && uPayload.height == |current.blockchain|
	                            && uPayload.round == current.round
	                            && recoverSignedCommitAuthor(m.commitPayload) == current.id
	*/

	panic("implement")
}

func createCommit(state State, value []byte) *SignedMessage {
	/**
	Commit(
	                    signCommit(
	                        UnsignedCommit(
	                            |current.blockchain|,
	                            current.round,
	                            signHash(hashBlockForCommitSeal(proposedBlock), current.id),
	                            digest(proposedBlock)),
	                            current.id
	                        )
	                    );
	*/
	panic("implement")
}

func validateCommit(
	state State,
	config Config,
	signedCommit *SignedMessage,
	height uint64,
	round Round,
	proposedMsg *SignedMessage,
	operators []*types.Operator,
) error {
	if signedCommit.Message.MsgType != CommitMsgType {
		return errors.New("commit msg type is wrong")
	}
	if signedCommit.Message.Height != height {
		return errors.New("commit height is wrong")
	}
	if signedCommit.Message.Round != round {
		return errors.New("commit round is wrong")
	}
	if !bytes.Equal(proposedMsg.Message.GetCommitData().GetData(), signedCommit.Message.GetCommitData().GetData()) {
		return errors.New("proposed data different than commit msg data")
	}

	if len(signedCommit.GetSigners()) != 1 {
		return errors.New("commit msg allows 1 signer")
	}
	// TODO how to process decided msgs with multiple signer?

	if err := signedCommit.Signature.VerifyByOperators(signedCommit, config.GetSignatureDomainType(), types.QBFTSigType, operators); err != nil {
		return errors.Wrap(err, "commit msg signature invalid")
	}
	return nil
}
