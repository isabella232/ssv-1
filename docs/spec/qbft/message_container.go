package qbft

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
)

type MsgContainer struct {
	Msgs map[Round][]*SignedMessage
}

// AllMessaged returns all messages
func (c *MsgContainer) AllMessaged() []*SignedMessage {
	ret := make([]*SignedMessage, 0)
	for _, roundMsgs := range c.Msgs {
		for _, msg := range roundMsgs {
			ret = append(ret, msg)
		}
	}
	return ret
}

// MessagesForRound returns all msgs for Height and round, empty slice otherwise
func (c *MsgContainer) MessagesForRound(round Round) []*SignedMessage {
	if c.Msgs[round] != nil {
		return c.Msgs[round]
	}
	return make([]*SignedMessage, 0)
}

// AddIfDoesntExist will add a msg with the following unique params: (round, signers)
// returns true if added
func (c *MsgContainer) AddIfDoesntExist(msg *SignedMessage) (bool, error) {
	if c.Msgs[msg.Message.Round] == nil {
		c.Msgs[msg.Message.Round] = make([]*SignedMessage, 0)
	}

	r, err := msg.GetRoot()
	if err != nil {
		return false, errors.Wrap(err, "could not get signed msg root")
	}

	for _, existingMsg := range c.Msgs[msg.Message.Round] {
		toMatchRoot, err := existingMsg.GetRoot()
		if err != nil {
			return false, errors.Wrap(err, "could not get existing signed msg root")
		}
		if bytes.Equal(r, toMatchRoot) && existingMsg.MatchedSigners(msg.Signers) {
			// TODO - what about 2 msgs with common signers but not identical?
			return false, nil
		}
	}

	// add msg
	c.Msgs[msg.Message.Round] = append(c.Msgs[msg.Message.Round], msg)
	return true, nil
}

// Encode returns the encoded struct in bytes or error
func (c *MsgContainer) Encode() ([]byte, error) {
	return json.Marshal(c)
}

// Decode returns error if decoding failed
func (c *MsgContainer) Decode(data []byte) error {
	return json.Unmarshal(data, &c)
}
