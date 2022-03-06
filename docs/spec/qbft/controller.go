package qbft

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

// HistoricalInstanceCapacity represents the upper bound of instances a controller can process messages for as messages are not
// guaranteed to arrive in a timely fashion, we physically limit how far back the controller will process messages for
const HistoricalInstanceCapacity int = 5

type instances [HistoricalInstanceCapacity]*Instance

func (i instances) FindInstance(height uint64) *Instance {
	for _, inst := range i {
		if inst != nil {
			if inst.GetHeight() == height {
				return inst
			}
		}
	}
	return nil
}

// Controller is a QBFT coordinator responsible for starting and following the entire life cycle of multiple QBFT instances
type Controller struct {
	Identifier []byte
	Height     uint64 // incremental Height for instances
	// storedInstances stores the last HistoricalInstanceCapacity in an array for message processing purposes.
	storedInstances instances
	signer          types.SSVSigner
	valueCheck      proposedValueCheck
	storage         Storage
	network         Network
}

// StartNewInstance will start a new QBFT instance, if can't will return error
func (c *Controller) StartNewInstance(value []byte) error {
	if err := c.canStartInstance(value); err != nil {
		return errors.Wrap(err, "can't start new QBFT instance")
	}

	c.bumpHeight()
	newInstance := c.addAndStoreNewInstance()
	newInstance.Start(value, c.Height)

	return nil
}

// ProcessMsg processes a new msg, returns true if Decided, non nil byte slice if Decided (Decided value) and error
// Decided returns just once per instance as true, following messages (for example additional commit msgs) will not return Decided true
func (c *Controller) ProcessMsg(msg *SignedMessage) (bool, []byte, error) {
	if !bytes.Equal(c.Identifier, msg.Message.Identifier) {
		return false, nil, errors.New(fmt.Sprintf("message doesn't belong to Identifier %x", c.Identifier))
	}

	inst := c.InstanceForHeight(msg.Message.Height)
	if inst == nil {
		return false, nil, errors.New(fmt.Sprintf("instance for Height %d,  Identifier %x not found", msg.Message.Height, c.Identifier))
	}

	prevDecided, _ := inst.IsDecided()
	decided, decidedValue, aggregatedCommit, err := inst.ProcessMsg(msg)
	if err != nil {
		return false, nil, errors.Wrap(err, "could not process msg")
	}

	// save the highest Decided
	if decided && inst.GetHeight() == c.Height { // It's the highest instance
		if err := c.storage.SaveHighestDecided(aggregatedCommit); err != nil {
			// LOG
		}
	}

	// if previously Decided we do not return Decided true again
	if prevDecided {
		return false, nil, err
	}

	// Broadcast Decided msg
	if err := c.network.BroadcastDecided(aggregatedCommit); err != nil {
		//TODO We do not return error here, just Log broadcasting error.
		return decided, decidedValue, nil
	}

	return decided, decidedValue, nil
}

func (c *Controller) InstanceForHeight(height uint64) *Instance {
	return c.storedInstances.FindInstance(height)
}

func (c *Controller) bumpHeight() {
	c.Height++
}

// GetIdentifier returns QBFT Identifier, used to identify messages
func (c *Controller) GetIdentifier() []byte {
	return c.Identifier
}

// addAndStoreNewInstance returns creates a new QBFT instance, stores it in an array and returns it
func (c *Controller) addAndStoreNewInstance() Instance {
	panic("implement")
}

func (c *Controller) canStartInstance(value []byte) error {
	if c.Height != 0 {
		// check prev instance if prev instance is not the first instance
		inst := c.storedInstances.FindInstance(c.Height)
		if inst == nil {
			return errors.New("could not find previous instance")
		}
		if decided, _ := inst.IsDecided(); !decided {
			return errors.New("previous instance hasn't Decided")
		}
	}

	// check value
	if err := c.valueCheck(value); err != nil {
		return errors.Wrap(err, "value invalid")
	}

	// TODO - complete more checks
	return nil
}

// Encode implementation
func (c *Controller) Encode() ([]byte, error) {
	m := make(map[string]interface{})

	m["id"] = c.Identifier
	m["Height"] = c.Height

	instances := make([][]byte, 0)
	for _, i := range c.storedInstances {
		byts, err := i.Encode()
		if err != nil {
			return nil, errors.Wrap(err, "could not encode instance")
		}
		instances = append(instances, byts)
	}
	m["instances"] = instances

	return json.Marshal(m)
}

// Decode implementation
func (c *Controller) Decode(data []byte) error {
	panic("implement")
}
