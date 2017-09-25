package satellite

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/antonkuzmenko/gogarin/pkg/transport"
	"github.com/go-kit/kit/endpoint"
)

func makeRegisterEndpoint(config Config, conn transport.Connection) endpoint.Endpoint {
	return transport.NewClient(
		conn,
		"satellite.register",
		time.Duration(config.Transport.RegisterTimeoutInSec)*time.Second,
		encodeJSONRequest,
		decodeJSONRegisterResponse,
	).Endpoint()
}

func decodeJSONRegisterResponse(_ context.Context, data interface{}) (response interface{}, err error) {
	r := data.([]byte)
	var i Info
	err = json.Unmarshal(r, &i)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func encodeJSONRequest(_ context.Context, data interface{}) (request interface{}, err error) {
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
