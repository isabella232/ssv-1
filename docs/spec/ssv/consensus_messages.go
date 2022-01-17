package ssv

import (
	"fmt"
	"github.com/bloxapp/ssv/docs/spec/network"
	"github.com/pkg/errors"
)

func (v *Validator) processConsensusMsg(msg network.Message) error {
	signedMsg, err := msg.GetSignedMessage()
	if err != nil {
		return errors.Wrap(err, "could not get SignedMessage from network message")
	}

	dutyRunner := v.dutyRunners.DutyRunnerForMsgID(msg.GetID())
	if dutyRunner == nil {
		return errors.Errorf("duty runner not found for msg id: %s", msg.GetID())
	}

	decided, decidedValue, err := dutyRunner.qbftController.ProcessMsg(signedMsg)

	if !decided {
		return nil
	}

	// check decided value
	if err := v.valCheck.Check(decidedValue); err != nil {
		return errors.Wrap(err, "post consensus value check fail")
	}

	partialSig, err := v.partialSignDecidedValue(decidedValue)
	if err != nil {
		return errors.Wrap(err, "could not partially sign decided value")
	}
}

// partialSignDecidedValue takes a slice of bytes after consensus, parses it into an eth2 object and signs it
func (v *Validator) partialSignDecidedValue(decidedValue []byte) ([]byte, error) {
	panic("implement")
}
