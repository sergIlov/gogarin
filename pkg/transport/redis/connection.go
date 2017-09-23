package redis

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/antonkuzmenko/gogarin/pkg/transport"
	"github.com/garyburd/redigo/redis"
)

// New creates a connection pool that implements transport.Connection.
func New(c Config) transport.Connection {
	pool := &redis.Pool{
		MaxActive:   c.MaxActiveConnections,
		MaxIdle:     c.MaxIdleConnections,
		IdleTimeout: time.Duration(c.ConnectionIdleTimeoutInMs) * time.Millisecond,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				c.Address,
				redis.DialConnectTimeout(time.Duration(c.ConnectTimeoutInMs)*time.Millisecond),
				redis.DialReadTimeout(time.Duration(c.ReadTimeoutInMs)*time.Millisecond),
				redis.DialWriteTimeout(time.Duration(c.WriteTimeoutInMs)*time.Millisecond),
			)
		},
	}

	return &Connection{pool: pool}
}

// Config for redis.Pool.
type Config struct {
	// Address specifies the location of the redis sever and is used when dialing a Connection.
	Address string `default:"localhost:6379"`

	// DB specifies the database to select when dialing a Connection.
	DB int `default:"0"`

	// MaxIdleConnections is a maximum number of idle connections in the pool.
	MaxIdleConnections int `default:"50"`

	// MaxActiveConnections is a maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	MaxActiveConnections int `default:"1000"`

	// ConnectionIdleTimeoutInMs close connections after remaining idle for this duration.
	// If the value is zero, then idle connections are not closed. Applications should set
	// the timeout to a value less than the server's timeout.
	// The default ConnectionIdleTimeoutInMs is 300000ms/300s/5m.
	ConnectionIdleTimeoutInMs int `default:"300000"`

	// ConnectTimeoutInMs specifies the timeout for connecting to the Redis server.
	// The default ConnectTimeoutInMs is 10000ms/10s.
	ConnectTimeoutInMs int `default:"10000"`

	// ReadTimeoutInMs specifies the timeout for reading a single command reply.
	// The default ReadTimeoutInMs is 10000ms/10s.
	ReadTimeoutInMs int `default:"10000"`

	// WriteTimeoutInMs specifies the timeout for writing a single command.
	// The default WriteTimeoutInMs is 10000ms/10s.
	WriteTimeoutInMs int `default:"10000"`
}

type message struct {
	ReplyTopic string      `json:"reply_topic"`
	Data       interface{} `json:"data"`
}

type Connection struct {
	pool *redis.Pool
}

func (r *Connection) Send(topic, replyTopic string, data interface{}) (err error) {
	const command = "LPUSH"

	con := r.pool.Get()
	defer con.Close() // nolint: errcheck
	if con.Err() != nil {
		return con.Err()
	}

	m := message{ReplyTopic: replyTopic, Data: data}
	msg, err := json.Marshal(m)
	if err != nil {
		return err
	}

	_, err = con.Do(command, topic, msg)
	return err
}

func (r *Connection) Receive(topic string, timeout time.Duration) (replyTopic string, data interface{}, err error) {
	const command = "BRPOP"

	con := r.pool.Get()
	defer con.Close() // nolint: errcheck
	if con.Err() != nil {
		return "", nil, con.Err()
	}

	res, err := con.Do(command, topic, int(timeout.Seconds()))
	if err != nil {
		return "", nil, err
	}

	result, err := redis.ByteSlices(res, err)
	if err == redis.ErrNil {
		return "", nil, transport.ErrTimeout
	}
	if err != nil {
		return "", nil, err
	}

	if len(result) != 2 {
		return "", nil, transport.ErrInvalidResponse
	}

	var msg message
	err = json.Unmarshal(result[1], &msg)
	if err != nil {
		return "", nil, err
	}

	data, err = base64.StdEncoding.DecodeString(msg.Data.(string))
	if err != nil {
		return "", nil, err
	}

	return msg.ReplyTopic, data, err
}
