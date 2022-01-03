package qbft

type Iterator interface {
	Next() SignedMessage
}

type MsgContainer interface {
	Iterator() Iterator
}
