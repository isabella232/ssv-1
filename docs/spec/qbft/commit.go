package qbft

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
