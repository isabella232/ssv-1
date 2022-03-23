package spectest

import (
	tests2 "github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests/postconsensus"
)

var AllTests = []*tests2.SpecTest{
	tests2.HappyFullFlow(),
	postconsensus.ValidMessage(),
}
