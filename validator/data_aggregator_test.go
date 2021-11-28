package validator

import (
	"fmt"
	"github.com/bloxapp/ssv/utils/threshold"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsAggregator(t *testing.T) {
	threshold.Init()
	sk := &bls.SecretKey{}
	sk.SetByCSPRNG()

	fmt.Printf("sk: %s\n", sk.GetHexString())

	cnt := 0
	modulo := 10
	for i := 1; i <= 100; i++ {
		res, err := isAggregator(sk, uint64(i), uint64(modulo))
		require.NoError(t, err)
		fmt.Printf("aggregator val %d: %t\n", i, res)
		if res {
			cnt++
		}
	}
	fmt.Printf("aggregator in %d validators\n", cnt)
}

func TestIsAggregator2(t *testing.T) {
	threshold.Init()

	sum := 0
	max := 100
	modulo := 10
	for j := 0; j < max; j++ {
		sk := &bls.SecretKey{}
		sk.SetByCSPRNG()
		cnt := 0
		for i := 1; i <= 50; i++ {
			res, err := isAggregator(sk, uint64(i), uint64(modulo))
			require.NoError(t, err)
			if res {
				cnt++
			}
		}
		sum += cnt
		fmt.Printf("aggregator in %d validators\n", cnt)
	}
	fmt.Printf("aggregator avg: %f\n", float64(sum)/float64(max))
}
