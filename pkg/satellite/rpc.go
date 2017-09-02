package satellite

import (
	"errors"
	"fmt"

	"github.com/antonkuzmenko/gogarin/pkg/transport"
	"github.com/antonkuzmenko/gogarin/pkg/transport/redis"
)

type RPCConfig struct {
	Adapter string `required:"true"`
	Redis   redis.Config
}

const (
	redisRPC = "redis"
)

func NewRPC(c RPCConfig) (transport.Connection, error) {
	if c.Adapter == redisRPC {
		return redis.New(c.Redis), nil
	}

	return nil, errors.New(fmt.Sprintf("RPC client for %q is not found", c.Adapter))
}
