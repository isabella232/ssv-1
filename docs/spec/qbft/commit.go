package qbft

import (
	"bytes"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

// uponCommit returns true if a quorum of commit messages was received.
func uponCommit(state State, signedCommit SignedMessage, commitMsgContainer MsgContainer) (bool, []byte, error) {
	if err := validateCommit(
		signedCommit,
		state.GetHeight(),
		state.GetRound(),
		state.GetConfig().GetNodes(),
	); err != nil {
		return false, nil, errors.Wrap(err, "commit msg invalid")
	}
	if !commitMsgContainer.AddIfDoesntExist(signedCommit) {
		return false, nil, nil // uponCommit was already called
	}

	value := signedCommit.GetMessage().GetCommitData().GetData()
	if commitQuorumForValue(state, commitMsgContainer, value) {
		return true, value, nil
	}
	return false, nil, nil
}

func commitQuorumForValue(state State, commitMsgContainer MsgContainer, value []byte) bool {
	commitMsgs := commitMsgContainer.MessagesForHeightAndRound(state.GetHeight(), state.GetRound())
	valueFiltered := make([]SignedMessage, 0)
	for _, msg := range commitMsgs {
		if bytes.Equal(msg.GetMessage().GetCommitData().GetData(), value) {
			valueFiltered = append(valueFiltered, msg)
		}
	}

	return state.GetConfig().HasQuorum(valueFiltered)
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

func createCommit(state State, value []byte) SignedMessage {
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
	signedCommit SignedMessage,
	height uint64,
	round Round,
	nodes []types.Node,
) error {
	if signedCommit.GetMessage().GetType() != CommitType {
		return errors.New("commit msg type is wrong")
	}
	if signedCommit.GetMessage().GetHeight() != height {
		return errors.New("commit height is wrong")
	}
	if signedCommit.GetMessage().GetRound() != round {
		return errors.New("commit round is wrong")
	}
	if !signedCommit.IsValidSignature(nodes) {
		return errors.New("commit msg signature invalid")
	}
	return nil
}
