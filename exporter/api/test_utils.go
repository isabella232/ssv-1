package api

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/url"
	"sync"
)

// WSClient represents a client connection to be used in tests
type WSClient struct {
	ctx    context.Context
	logger *zap.Logger
	mut    sync.Mutex
	msgs   []Message

	out chan Message
}

// NewWSClient creates a new instance of ws client
func NewWSClient(ctx context.Context, logger *zap.Logger) *WSClient {
	return &WSClient{
		ctx:    ctx,
		logger: logger,
		mut:    sync.Mutex{},
		msgs:   make([]Message, 0),
		out:    make(chan Message),
	}
}

// StartStream initiates stream
func (client *WSClient) StartStream(addr, path string) error {
	u := url.URL{Scheme: "ws", Host: addr, Path: path}
	client.logger.Debug("connecting to server", zap.String("addr", u.String()))
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return errors.Wrap(err, "dial error")
	}
	defer func() {
		_ = c.Close()
	}()

	for {
		if client.ctx.Err() != nil {
			return nil
		}
		_, raw, err := c.ReadMessage()
		if err != nil {
			return errors.Wrap(err, "read error")
		}
		var msg Message
		if err := json.Unmarshal(raw, &msg); err != nil {
			client.logger.Error("failed to parse message", zap.Error(err))
			continue
		}
		client.mut.Lock()
		client.msgs = append(client.msgs, msg)
		client.mut.Unlock()
	}
}

// StartQuery initiates query requests
func (client *WSClient) StartQuery(addr, path string) error {
	u := url.URL{Scheme: "ws", Host: addr, Path: path}
	client.logger.Debug("connecting to server", zap.String("addr", u.String()))
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return errors.Wrap(err, "dial error")
	}
	defer func() {
		_ = c.Close()
	}()

msgLoop:
	for {
		select {
		case <-client.ctx.Done():
			return nil
		case m := <-client.out:
			if err := c.WriteJSON(&m); err != nil {
				return errors.Wrap(err, "send error")
			}
			_, raw, err := c.ReadMessage()
			if err != nil {
				return errors.Wrap(err, "read error")
			}
			var msg Message
			if err := json.Unmarshal(raw, &msg); err != nil {
				client.logger.Error("failed to parse message", zap.Error(err))
				continue msgLoop
			}
			client.mut.Lock()
			client.msgs = append(client.msgs, msg)
			client.mut.Unlock()
		}
	}
}

// MessageCount returns the count of incoming messages
func (client *WSClient) MessageCount() int {
	client.mut.Lock()
	defer client.mut.Unlock()

	return len(client.msgs)
}
