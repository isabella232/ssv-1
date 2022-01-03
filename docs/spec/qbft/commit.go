package qbft

// didSendCommitForHeightAndRound returns true if sent commit msg for specific height and round
func didSendCommitForHeightAndRound() bool {
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
