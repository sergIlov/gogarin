package redis

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

// Server wraps an endpoint and implements a Handler.
type Server struct {
	e            endpoint.Endpoint
	dec          DecodeRequestFunc
	enc          EncodeResponseFunc
	errorEncoder ErrorEncoder
	logger       log.Logger
}

// NewServer constructs a new server, which implements transport.Handler and wraps
// the provided endpoint.
func NewServer(
	e endpoint.Endpoint,
	dec DecodeRequestFunc,
	enc EncodeResponseFunc,
	options ...ServerOption,
) *Server {
	s := &Server{
		e:            e,
		dec:          dec,
		enc:          enc,
		errorEncoder: DefaultErrorEncoder,
		logger:       log.NewNopLogger(),
	}
	for _, option := range options {
		option(s)
	}
	return s
}

// ServerOption sets an optional parameter for servers.
type ServerOption func(*Server)

// ServerErrorEncoder is used to encode errors whenever they're
// encountered in the processing of a request. Clients can use this
// to provide custom // error formatting and response codes.
// By default, errors will be written with the DefaultErrorEncoder.
func ServerErrorEncoder(ee ErrorEncoder) ServerOption {
	return func(s *Server) { s.errorEncoder = ee }
}

// ServerLogger is used to log non-terminal errors. By default, no errors
// are logged. This is intended as a diagnostic measure.
func ServerLogger(logger log.Logger) ServerOption {
	return func(s *Server) { s.logger = logger }
}

func (s Server) ServeRPC(ctx context.Context, req interface{}) interface{} {
	request, err := s.dec(ctx, req)
	if err != nil {
		s.logger.Log("err", err)
		return s.errorEncoder(ctx, err)
	}

	response, err := s.e(ctx, request)
	if err != nil {
		s.logger.Log("err", err)
		return s.errorEncoder(ctx, err)
	}

	res, err := s.enc(ctx, response)
	if err != nil {
		s.logger.Log("err", err)
		return s.errorEncoder(ctx, err)
	}

	return res
}

// ErrorEncoder is responsible for encoding an error.
type ErrorEncoder func(context.Context, error) interface{}

// DefaultErrorEncoder encodes the error to the JSON.
func DefaultErrorEncoder(_ context.Context, err error) interface{} {
	type response struct {
		Ok     bool
		Errors []error
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(response{Ok: false, Errors: []error{err}})
	if err != nil {
		buf = bytes.Buffer{}
		res := response{Ok: false, Errors: []error{errors.New("Could not encode an error")}}
		_ = json.NewEncoder(&buf).Encode(res)
	}

	return buf
}
