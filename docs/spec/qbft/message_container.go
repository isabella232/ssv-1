package qbft

type MsgContainer interface {
	// MessagesForHeightAndRound returns all msgs for height and round, empty slice otherwise
	MessagesForHeightAndRound(height uint64, round Round) []SignedMessage
	// AddIfDoesntExist will add the msg if it doesn't exist in the container, returns true if added
	AddIfDoesntExist(msg SignedMessage) bool
}
