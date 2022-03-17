package spectest

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
)

var AllTests = []*SpecTest{
	happyFullFlow(),
}

type SpecTest struct {
	Name     string
	Pre      *qbft.Instance
	PostRoot string
	Messages []*qbft.SignedMessage
	err      error
}
