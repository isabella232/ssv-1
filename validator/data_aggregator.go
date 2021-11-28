package validator

import (
	"encoding/binary"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

func indexSig(sk *bls.SecretKey, index uint64) ([]byte, error) {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, index)
	return sk.SignByte(data).Serialize(), nil
}

func isAggregator(sk *bls.SecretKey, index uint64, modulo uint64) (bool, error) {
	sig, err := indexSig(sk, index)
	if err != nil {
		return false, errors.Wrap(err, "could not sign index")
	}

	base := binary.LittleEndian.Uint64(sig[0:8])
	return base%modulo == 0, nil
}
