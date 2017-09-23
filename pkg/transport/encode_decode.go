package transport

import (
	"context"
)

// DecodeRequestFunc extracts a user-domain request object from a request.
// It's designed to be used in Space center, for server-side endpoints.
type DecodeRequestFunc func(context.Context, interface{}) (request interface{}, err error)

// EncodeRequestFunc encodes the passed request object into the request object.
// It's designed to be used in Satellites, for client-side endpoints.
type EncodeRequestFunc func(context.Context, interface{}) (request interface{}, err error)

// EncodeResponseFunc encodes the passed response object to the response message.
// It's designed to be used in Space center, for server-side endpoints.
type EncodeResponseFunc func(context.Context, interface{}) (response interface{}, err error)

// DecodeResponseFunc extracts a user-domain response object from a response object.
// It's designed to be used in Satellites, for client-side
type DecodeResponseFunc func(context.Context, interface{}) (response interface{}, err error)
