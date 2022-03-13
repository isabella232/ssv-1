package main

import (
	"encoding/json"
	"fmt"
	"github.com/bloxapp/ssv/docs/spec/qbft/spectest"
)

func main() {
	for _, t := range spectest.AllTests {
		printTest(t)
	}
}

func printTest(t *spectest.SpecTest) {
	byts, _ := json.Marshal(t)
	fmt.Printf("%s: %s\n", t.Name, string(byts))
}
