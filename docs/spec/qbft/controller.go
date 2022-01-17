package qbft

import (
	"bytes"
	"fmt"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

type instances [HistoricalInstanceSize]*Instance

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

// HistoricalInstanceSize represents the upper bound of instances a controller can process messages for as messages are not
// guaranteed to arrive in a timely fashion, we physically limit how far back the controller will process messages for
const HistoricalInstanceSize int = 5

// Controller is a QBFT coordinator responsible for starting and following the entire life cycle of multiple QBFT instances
type Controller struct {
	identifier []byte
	// storedInstances stores the last HistoricalInstanceSize in an array for message processing purposes.
	storedInstances instances
	valCheck        types.ValueCheck
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

// ProcessMsg processes a new msg, returns true if decided, non nil byte slice if decided (decided value) and error
// decided returns just once per instance as true, following messages (for example additional commit msgs) will not return decided true
func (c *Controller) ProcessMsg(msg SignedMessage) (bool, []byte, error) {
	if !bytes.Equal(c.GetIdentifier(), msg.GetMessage().GetInstanceIdentifier()) {
		return false, nil, errors.New(fmt.Sprintf("message doesn't belong to identifier %x", c.GetIdentifier()))
	}

	inst := c.storedInstances.FindInstance(msg.GetMessage().GetHeight())
	if inst == nil {
		return false, nil, errors.New(fmt.Sprintf("instance for height %d,  identifier %x not found", msg.GetMessage().GetHeight(), c.GetIdentifier()))
	}

	prevDecided, _ := inst.IsDecided()
	if prevDecided {
		if _, _, err := inst.ProcessMsg(msg); err != nil {
			return false, nil, err
		}
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
	// check prev instance
	inst := c.storedInstances.FindInstance(c.Height())
	if inst == nil {
		return errors.New("could not find previous instance")
	}
	if decided, _ := inst.IsDecided(); !decided {
		return errors.New("previous instance hasn't decided")
	}

	// check value
	if err := c.valCheck.Check(value); err != nil {
		return errors.Wrap(err, "value invalid")
	}
	panic("implement")
}
