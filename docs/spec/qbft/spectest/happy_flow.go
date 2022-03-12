package spectest

import "github.com/bloxapp/ssv/docs/spec/qbft"

func happyFlow() *SpecTest {
	pre := baseInstance()
	post := baseInstance()
	msgs := []*qbft.SignedMessage{
		signMsg(&qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       proposalDataByts([]byte{1, 2, 3, 4}, nil, nil),
		}),
	}
	return &SpecTest{
		Name:     "happy flow",
		Pre:      pre,
		Post:     post,
		Messages: msgs,
	}
}
