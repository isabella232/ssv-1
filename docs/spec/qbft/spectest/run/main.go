package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/qbft/spectest"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	basedir, _ := os.Getwd()
	path := filepath.Join(basedir, "docs", "spec", "qbft", "spectest", "run")
	fileName := "tests.json"
	fmt.Printf("reading %s/%s\n", path, fileName)

	tests := map[string]*spectest.SpecTest{}
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
		pre := qbft.NewInstance(spectest.TestingConfig)
		pre.Decode(byts)
		test.Pre = pre
		runTest(test)
	}
}

func runTest(test *spectest.SpecTest) {
	for _, msg := range test.Messages {
		_, _, _, err := test.Pre.ProcessMsg(msg)
		if err != nil {
			panic(err.Error())
		}
	}

	postRoot, err := test.Pre.State.GetRoot()
	if err != nil {
		panic(err.Error())
	}

	if test.PostRoot != hex.EncodeToString(postRoot) {
		panic("post state root is wrong")
	}
}
