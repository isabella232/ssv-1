package main

import (
	"encoding/json"
	"fmt"
	"github.com/bloxapp/ssv/docs/spec/qbft/spectest"
)

func main() {
	all := map[string]*spectest.SpecTest{}
	for _, t := range spectest.AllTests {
		all[t.Name] = t
	}

	byts, _ := json.Marshal(all)
	fmt.Printf("\n\n%s\n\n", string(byts))
}
