package spectest

import (
	"encoding/hex"
	"encoding/json"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestAll(t *testing.T) {
	for _, test := range AllTests {
		t.Run(test.Name, func(t *testing.T) {
			runTest(t, test)
		})
	}
}

func TestJson(t *testing.T) {
	basedir, _ := os.Getwd()
	path := filepath.Join(basedir, "generate")
	fileName := "tests.json"
	tests := map[string]*SpecTest{}
	byteValue, err := ioutil.ReadFile(path + "/" + fileName)
	if err != nil {
		panic(err.Error())
	}

	if err := json.Unmarshal(byteValue, &tests); err != nil {
		panic(err.Error())
	}

	for _, test := range tests {
		byts, _ := test.Pre.Encode()

		// a little trick we do to instantiate all the internal instance params
		pre := qbft.NewInstance(testingutils.TestingConfig, test.Pre.State.Share, test.Pre.State.ID)
		pre.Decode(byts)
		test.Pre = pre
		runTest(t, test)
	}
}

func runTest(t *testing.T, test *SpecTest) {
	for _, msg := range test.Messages {
		_, _, _, err := test.Pre.ProcessMsg(msg)
		require.NoError(t, err)
	}

	postRoot, err := test.Pre.State.GetRoot()
	require.NoError(t, err)

	require.EqualValues(t, test.PostRoot, hex.EncodeToString(postRoot), "post root not valid")
}
