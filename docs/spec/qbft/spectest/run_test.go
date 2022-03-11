package spectest

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAll(t *testing.T) {
	for _, test := range AllTests {
		t.Run(test.Name, func(t *testing.T) {
			runTest(t, test)
		})
	}
}

func runTest(t *testing.T, test *SpecTest) {
	for _, msg := range test.Messages {
		_, _, _, err := test.Pre.ProcessMsg(msg)
		require.NoError(t, err)
	}
}
