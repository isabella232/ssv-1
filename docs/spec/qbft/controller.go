package qbft

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
)

type instances [HistoricalInstanceCapacity]*Instance

func (i instances) FindInstance(height uint64) *Instance {
	for _, inst := range i {
		if inst != nil {
			if inst.state.GetHeight() == height {
				return inst
			}
		}
	}
	return nil
}

// HistoricalInstanceCapacity represents the upper bound of instances a controller can process messages for as messages are not
// guaranteed to arrive in a timely fashion, we physically limit how far back the controller will process messages for
const HistoricalInstanceCapacity int = 5

// Controller is a QBFT coordinator responsible for starting and following the entire life cycle of multiple QBFT instances
type Controller struct {
	identifier []byte
	// storedInstances stores the last HistoricalInstanceCapacity in an array for message processing purposes.
	storedInstances instances
}

// StartNewInstance will start a new QBFT instance, if can't will return error
func (c *Controller) StartNewInstance(value []byte) error {
	if err := c.canStartInstance(value); err != nil {
		return errors.Wrap(err, "can't start new QBFT instance")
	}

	newInstance := c.addAndStoreNewInstance()
	newInstance.Start(value)

	return nil
}

// ProcessMsg processes a new msg
func (c *Controller) ProcessMsg(msg SignedMessage) error {
	if !bytes.Equal(c.GetIdentifier(), msg.GetMessage().GetInstanceIdentifier()) {
		return errors.New(fmt.Sprintf("message doesn't belong to identifier %x", c.GetIdentifier()))
	}

	inst := c.storedInstances.FindInstance(msg.GetMessage().GetHeight())
	if inst == nil {
		return errors.New(fmt.Sprintf("instance for height %d,  identifier %x not found", msg.GetMessage().GetHeight(), c.GetIdentifier()))
	}

	return inst.ProcessMsg(msg)
}

// Height returns the current running instance height or, if not started, the last decided height
func (c *Controller) Height() uint64 {
	panic("implement")
}

// GetIdentifier returns QBFT identifier, used to identify messages
func (c *Controller) GetIdentifier() []byte {
	return c.identifier
}

// addAndStoreNewInstance returns creates a new QBFT instance, stores it in an array and returns it
func (c *Controller) addAndStoreNewInstance() Instance {
	panic("implement")
}

func (c *Controller) canStartInstance(value []byte) error {
	panic("implement")
}
