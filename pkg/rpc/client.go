package rpc

import (
	"errors"
	"io"
	"time"
)

var ErrInvalidResponse = errors.New("invalid response")
var ErrTimeout = errors.New("timeout is reached")

type Client interface {
	Respond(topic string, data io.Reader) error
	Send(topic string, data io.Reader, timeout time.Duration) (result io.Reader, err error)
	Receive(topic string, timeout time.Duration) (replyTo string, data io.Reader, err error)
}
