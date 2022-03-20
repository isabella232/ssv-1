package testingutils

import (
	"github.com/bloxapp/ssv/utils/threshold"
	"github.com/herumi/bls-eth-go-binary/bls"
)

var TestingSK1 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("19153f7f7f004a839750f488fe85984c61e8d8342927e5fa027b72b8a5b23bda")
	return ret
}()
var TestingSK2 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("272c5e6350563fee5d8d7f0fb316ec94da8b352c1319fbe1edfc417dffe37ed8")
	return ret
}()
var TestingSK3 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("099ab2c89c0f882622288cc88a622c5591573b9f736e6be702097da249b45bac")
	return ret
}()
var TestingSK4 = func() *bls.SecretKey {
	threshold.Init()
	ret := &bls.SecretKey{}
	ret.DeserializeHexStr("344de4028bc9a073185bf5bb8e092f93da0a8f914a2392083ea327248324d257")
	return ret
}()
