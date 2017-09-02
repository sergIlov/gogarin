package redis

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/antonkuzmenko/gogarin/pkg/transport"
	"github.com/garyburd/redigo/redis"
	"github.com/oklog/ulid"
)

func New(c Config) transport.Connection {
	pool := &redis.Pool{
		MaxActive:   c.MaxActiveConnections,
		MaxIdle:     c.MaxIdleConnections,
		IdleTimeout: time.Duration(c.ConnectionIdleTimeoutInSeconds) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", c.Address)
		},
	}

	return &connection{pool: pool}
}

type Config struct {
	Address string `default:"localhost:6379"`

	Db int `default:"0"`

	// Maximum number of idle connections in the pool.
	MaxIdleConnections int `default:"2"`

	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	MaxActiveConnections int `default:"2"`

	// Close connections after remaining idle for this duration. If the value
	// is zero, then idle connections are not closed. Applications should set
	// the timeout to a value less than the server's timeout.
	ConnectionIdleTimeoutInSeconds int `default:"120"`

	ReadTimeoutInSeconds int `default:"10"`

	WriteTimeoutInSeconds int `default:"10"`
}

type message struct {
	ReplyTo string `json:"reply_to"`
	Data    interface{} `json:"data"`
}

type connection struct {
	pool *redis.Pool
}

func (r *connection) Send(topic string, data interface{}, timeout time.Duration) (result interface{}, err error) {
	con := r.pool.Get()
	defer con.Close()

	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id, err := ulid.New(ulid.Timestamp(t), entropy)
	if err != nil {
		return nil, err
	}
	replyTo := topic + ":reply:" + fmt.Sprintf("%s", id)

	err = send(con, topic, replyTo, data)
	if err != nil {
		return nil, err
	}

	_, result, err = receive(con, replyTo, timeout)

	return result, err
}

func (r *connection) Receive(topic string, timeout time.Duration) (replyTo string, data interface{}, err error) {
	con := r.pool.Get()
	defer con.Close()
	return receive(con, topic, timeout)
}

func (r *connection) Respond(topic string, data interface{}) error {
	con := r.pool.Get()
	defer con.Close()
	return send(con, topic, "", data)
}

func send(con redis.Conn, topic, replyTo string, data interface{}) error {
	const command = "LPUSH"

	if con.Err() != nil {
		return con.Err()
	}

	m := message{ReplyTo: replyTo, Data: data}
	mbyte, err := json.Marshal(m)
	if err != nil {
		return err
	}

	_, err = con.Do(command, topic, mbyte)
	return err
}

func receive(con redis.Conn, topic string, timeout time.Duration) (replyTo string, data interface{}, err error) {
	const command = "BRPOP"

	if con.Err() != nil {
		return "", nil, con.Err()
	}

	res, err := con.Do(command, topic, int(timeout.Seconds()))
	if err == redis.ErrNil {
		return "", nil, transport.ErrTimeout
	}
	if err != nil {
		return "", nil, err
	}
	bts, err := redis.ByteSlices(res, err)

	if err != nil {
		return "", nil, err
	}

	if len(bts) != 2 {
		return "", nil, transport.ErrInvalidResponse
	}

	var msg message
	err = json.Unmarshal(bts[1], &msg)
	if err != nil {
		return "", nil, err
	}

	return msg.ReplyTo, msg.Data, err
}
