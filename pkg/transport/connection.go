package transport

import (
	"errors"
	"time"
)

const NoReply = ""

// ErrInvalidResponse indicates corrupted response from a message broker.
var ErrInvalidResponse = errors.New("invalid response")

// ErrTimeout is returned when Send or Receive timeouts.
var ErrTimeout = errors.New("timeout is reached")

// Connection is a RPC over a message broker.
type Connection interface {
	// Send sends data to the topic.
	Send(topic, replyTopic string, data interface{}) error

	// Receive receives data from the topic.
	Receive(topic string, timeout time.Duration) (replyTopic string, data interface{}, err error)
}
