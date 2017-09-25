package satellite

import (
	"context"

	"github.com/antonkuzmenko/gogarin/pkg/transport"
)

type Config struct {
	Transport TransportConfig
	Logger    string `default:"json"`
}

func New(r transport.Connection, i Info) *Satellite {
	return &Satellite{
		conn: r,
		Info: i,
	}
}

type Satellite struct {
	conn     transport.Connection
	Info     Info
	Triggers []Trigger
}

func (s *Satellite) AddTrigger(t Trigger) {
	s.Triggers = append(s.Triggers, t)
}

func (s *Satellite) Start(c Config) error {
	registerEndpoint := makeRegisterEndpoint(c, s.conn)
	_, err := registerEndpoint(context.TODO(), s.Info)
	return err
}

func (s *Satellite) Stop() {
}

type Info struct {
	Name        string
	Version     string
	Description string
}

type Trigger struct {
	Call      func()
	Info      AbilityInfo
	Config    interface{}
	Validator func(config interface{})
}

type AbilityInfo struct {
	Name        string
	Description string
}
