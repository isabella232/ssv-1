package ssv

import (
	"crypto/sha256"
	"encoding/json"
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

type DutyRunnerState struct {
	BeaconRoleType beacon.RoleType
	Share          types.Share
	// DutyExecutionState holds all relevant params for a full duty execution (consensus & post consensus)
	DutyExecutionState *DutyExecutionState
	QBFTController     *qbft.Controller
}

// GetRoot returns the root used for signing and verification
func (s *DutyRunnerState) GetRoot() ([]byte, error) {
	marshaledRoot, err := s.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode DutyRunnerState")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret[:], nil
}

// Encode returns the encoded struct in bytes or error
func (s *DutyRunnerState) Encode() ([]byte, error) {
	return json.Marshal(s)
	//m := make(map[string]interface{})
	//
	//m["role_type"] = s.BeaconRoleType
	//
	//byts, err := s.Share.Encode()
	//if err != nil {
	//	return nil, errors.Wrap(err, "could not encode share")
	//}
	//m["share"] = byts
	//
	//byts, err = s.DutyExecutionState.Encode()
	//if err != nil {
	//	return nil, errors.Wrap(err, "could not encode DutyExecutionState")
	//}
	//m["duty_execution_state"] = byts
	//
	//byts, err = s.QBFTController.Encode()
	//if err != nil {
	//	return nil, errors.Wrap(err, "could not encode QBFTController")
	//}
	//m["controller"] = byts
	//
	//return json.Marshal(m)
}

// Decode returns error if decoding failed
func (s *DutyRunnerState) Decode(data []byte) error {
	return json.Unmarshal(data, &s)
}
