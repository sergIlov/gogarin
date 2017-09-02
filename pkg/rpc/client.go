package rpc

import (
	"errors"
	"time"
)

var ErrInvalidResponse = errors.New("invalid response")
var ErrTimeout = errors.New("timeout is reached")

type Client interface {
	Send(topic string, data interface{}, timeout time.Duration) (result interface{}, err error)
	Receive(topic string, timeout time.Duration) (replyTo string, data interface{}, err error)
	Respond(topic string, data interface{}) error
}
