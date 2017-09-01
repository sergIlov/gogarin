package redis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/antonkuzmenko/gogarin/pkg/rpc"
	"github.com/garyburd/redigo/redis"
	"github.com/oklog/ulid"
)

func New(c Config) rpc.Client {
	pool := &redis.Pool{
		MaxActive:   c.MaxActiveConnections,
		MaxIdle:     c.MaxIdleConnections,
		IdleTimeout: time.Duration(c.ConnectionIdleTimeoutInSeconds) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", c.Address)
		},
	}

	return &redisClient{pool: pool}
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
	Data    []byte `json:"data"`
}

type redisClient struct {
	pool *redis.Pool
}

func (r *redisClient) Send(topic string, data io.Reader, timeout time.Duration) (result io.Reader, err error) {
	con := r.pool.Get()
	defer con.Close()
	emptyReader := bytes.NewReader(nil)

	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id, err := ulid.New(ulid.Timestamp(t), entropy)
	if err != nil {
		return emptyReader, err
	}
	replyTo := topic + ":reply:" + fmt.Sprintf("%s", id)

	err = send(con, topic, replyTo, data)
	if err != nil {
		return emptyReader, err
	}

	_, result, err = receive(con, replyTo, timeout)

	return result, err
}

func (r *redisClient) Respond(replyTo string, data io.Reader) error {
	con := r.pool.Get()
	defer con.Close()
	return send(con, replyTo, "", data)
}

func (r *redisClient) Receive(topic string, timeout time.Duration) (replyTo string, data io.Reader, err error) {
	con := r.pool.Get()
	defer con.Close()
	return receive(con, topic, timeout)
}

func send(con redis.Conn, topic, replyTo string, data io.Reader) error {
	const command = "LPUSH"

	if con.Err() != nil {
		return con.Err()
	}

	d, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}

	m := message{ReplyTo: replyTo, Data: d}
	mbyte, err := json.Marshal(m)
	if err != nil {
		return err
	}

	_, err = con.Do(command, topic, mbyte)
	return err
}

func receive(con redis.Conn, topic string, timeout time.Duration) (replyTo string, data io.Reader, err error) {
	const command = "BRPOP"
	emptyReader := bytes.NewReader(nil)

	if con.Err() != nil {
		return "", emptyReader, con.Err()
	}

	res, err := con.Do(command, topic, int(timeout.Seconds()))
	if err == redis.ErrNil {
		return "", emptyReader, rpc.ErrTimeout
	}
	if err != nil {
		return "", emptyReader, err
	}
	bts, err := redis.ByteSlices(res, err)

	if err != nil {
		return "", emptyReader, err
	}

	if len(bts) != 2 {
		return "", emptyReader, rpc.ErrInvalidResponse
	}

	var msg message
	err = json.Unmarshal(bts[1], &msg)
	if err != nil {
		return "", emptyReader, err
	}

	return msg.ReplyTo, bytes.NewReader(msg.Data), err
}
