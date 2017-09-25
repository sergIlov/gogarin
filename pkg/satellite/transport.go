package satellite

import (
	"os"

	"github.com/antonkuzmenko/gogarin/pkg/transport"
	"github.com/antonkuzmenko/gogarin/pkg/transport/redis"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// TransportConfig contains configurations for RPC transport.
type TransportConfig struct {
	Adapter              string `required:"true"`
	Redis                redis.Config
	RegisterTimeoutInSec int `default:"10000"`
}

const (
	redisTransport = "redis"
)

// NewConnection creates new transport.Connection.
func NewConnection(c Config, logger log.Logger) transport.Connection {
	if c.Transport.Adapter == redisTransport {
		return redis.New(c.Transport.Redis)
	}

	level.Error(logger).Log("err", "invalid Transport.Adapter", "adapter", c.Transport.Adapter)
	os.Exit(1)
	return nil
}
