package satellite

import (
	"errors"
	"fmt"

	"github.com/antonkuzmenko/gogarin/pkg/rpc"
	"github.com/antonkuzmenko/gogarin/pkg/rpc/redis"
)

type RPCConfig struct {
	Adapter string `required:"true"`
	Redis   redis.Config
}

const (
	redisRPC = "redis"
)

func NewRPC(c RPCConfig) (rpc.Client, error) {
	if c.Adapter == redisRPC {
		return redis.New(c.Redis), nil
	}

	return nil, errors.New(fmt.Sprintf("RPC client for %q is not found", c.Adapter))
}
