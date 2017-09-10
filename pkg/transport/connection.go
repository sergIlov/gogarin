package transport

import (
	"errors"
	"time"
)

// ErrInvalidResponse indicates corrupted response from a message broker.
var ErrInvalidResponse = errors.New("invalid response")

// ErrTimeout is returned when Send or Receive timeouts.
var ErrTimeout = errors.New("timeout is reached")

// Connection is a RPC over a message broker.
type Connection interface {
	// Send sends data to the space center from a satellite.
	Send(topic string, data interface{}, timeout time.Duration) (result interface{}, err error)

	// Receive receives a message from a satellite.
	Receive(topic string, timeout time.Duration) (replyTo string, data interface{}, err error)

	// Respond responds to Send called from a satellite.
	Respond(topic string, data interface{}) error
}
