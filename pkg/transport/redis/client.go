package redis

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
)

// Client wraps a topic and provides a method that implements endpoint.Endpoint.
type Client struct {
	conn           *Connection
	topic          string
	receiveTimeout time.Duration
	enc            EncodeRequestFunc
	dec            DecodeResponseFunc
	before         []ClientRequestFunc
	after          []ClientResponseFunc
}

// NewClient constructs a usable Client for a single remote method.
func NewClient(
	conn *Connection,
	topic string,
	receiveTimeout time.Duration,
	enc EncodeRequestFunc,
	dec DecodeResponseFunc,
	options ...ClientOption,
) *Client {
	c := &Client{
		conn:           conn,
		topic:          topic,
		receiveTimeout: receiveTimeout,
		enc:            enc,
		dec:            dec,
	}
	for _, option := range options {
		option(c)
	}
	return c
}

// ClientRequestFunc may take information from a RPC request and put it into a
// request context. RequestFuncs are executed before invoking conn.Send.
type ClientRequestFunc func(ctx context.Context, req interface{}) context.Context

// ClientResponseFunc may take information from a RPC request and make the
// response available for consumption. ClientResponseFuncs are executed
// after a request has been made, but prior to it being decoded.
type ClientResponseFunc func(ctx context.Context, res interface{}) context.Context

// ClientOption sets an optional parameter for clients.
type ClientOption func(*Client)

// SetConnection sets the underlying redis.Connection used for requests.
func SetConnection(conn *Connection) ClientOption {
	return func(c *Client) { c.conn = conn }
}

// ClientBefore sets the ClientRequestFuncs that are applied to the outgoing RPC
// request before it's invoked.
func ClientBefore(before ...ClientRequestFunc) ClientOption {
	return func(c *Client) { c.before = append(c.before, before...) }
}

// ClientAfter sets the ClientResponseFuncs applied to the incoming RPC
// request prior to it being decoded. This is useful for obtaining anything off
// of the response and adding onto the context prior to decoding.
func ClientAfter(after ...ClientResponseFunc) ClientOption {
	return func(c *Client) { c.after = append(c.after, after...) }
}

// Endpoint returns a usable endpoint that invokes the remote endpoint.
func (c Client) Endpoint() endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (res interface{}, err error) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		request, err := c.enc(ctx, req)
		if err != nil {
			return nil, err
		}

		for _, f := range c.before {
			ctx = f(ctx, req)
		}

		response, err := c.conn.Send(c.topic, request, c.receiveTimeout)

		if err != nil {
			return nil, err
		}

		for _, f := range c.after {
			ctx = f(ctx, response)
		}

		res, err = c.dec(ctx, response)
		if err != nil {
			return nil, err
		}

		return res, nil
	}
}

// EncodeJSONRequest is an EncodeRequestFunc that serializes the request as a
// JSON object to the Request body. Many JSON-over-RPC services can use it as
// a sensible default.
func EncodeJSONRequest(ctx context.Context, req interface{}) error {
	panic("todo")
}
