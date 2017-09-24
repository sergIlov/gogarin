package satellite

import (
	"context"
	"fmt"

	"github.com/antonkuzmenko/gogarin/pkg/transport"
)

type Config struct {
	RPC    RPCConfig
	Logger string `default:"json"`
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
	Actions  []Action
}

func (s *Satellite) AddTrigger(t Trigger) {
	s.Triggers = append(s.Triggers, t)
}

func (s *Satellite) AddAction(a Action) {
	s.Actions = append(s.Actions, a)
}

func (s *Satellite) Start(c Config) error {
	registerEndpoint := createRegisterEndpoint(c, s.conn)
	res, err := registerEndpoint(context.TODO(), s.Info)
	if err != nil {
		return err
	}

	info := res.(Info)
	fmt.Printf("%+v", info)
	return nil
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

type Action struct {
	Call      func()
	Info      AbilityInfo
	Config    interface{}
	Validator func(config interface{})
}

type AbilityInfo struct {
	Name        string
	Description string
}
