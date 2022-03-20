package testingutils

import "github.com/bloxapp/ssv/docs/spec/ssv"

var BaseValidator = func() *ssv.Validator {
	ret := ssv.NewValidator(
		NewTestingNetwork(),
		NewTestingBeaconNode(),
		NewTestingStorage(),
		TestingShare,
		NewTestingKeyManager(),
		func(data []byte) error {
			return nil
		},
	)
	return ret
}
