package eth1

import (
	"github.com/prysmaticlabs/prysm/async/event"
	"math/big"
	"time"
)

// ClientMock implements eth1.Client interface
type ClientMock struct {
	Feed *event.Feed

	SyncTimeout  time.Duration
	SyncResponse error
}

// EventsFeed returns the contract events feed
func (ec *ClientMock) EventsFeed() *event.Feed {
	return ec.Feed
}

// Subscribe subscribes to contract events
func (ec *ClientMock) Subscribe(handler EventHandler) error {
	return nil
}

// Start mocking client init
func (ec *ClientMock) Start() error {
	return nil
}

// Sync mocking events sync
func (ec *ClientMock) Sync(fromBlock *big.Int, handler EventHandler) error {
	<-time.After(ec.SyncTimeout)
	return ec.SyncResponse
}
