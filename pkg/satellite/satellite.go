package satellite

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/antonkuzmenko/gogarin/pkg/transport"
)

type Config struct {
	RPC RPCConfig
}

func New(r transport.Connection, i Info) *Satellite {
	return &Satellite{
		client: r,
		Info:   i,
	}
}

type Satellite struct {
	client   transport.Connection
	Info     Info
	Triggers []Trigger
}

func (s *Satellite) AddTrigger(t Trigger) {
	s.Triggers = append(s.Triggers, t)
}

func (s *Satellite) Start() error {
	d, err := json.Marshal(s.Info)
	if err != nil {
		return err
	}
	res, err := s.client.Send("satellite.register", d, time.Duration(1)*time.Second)
	if err != nil {
		return err
	}
	r := res.([]byte)
	var i Info
	err = json.Unmarshal(r, &i)
	if err != nil {
		return err
	}
	fmt.Printf("%+v", i)
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

type AbilityInfo struct {
	Name        string
	Description string
}

const (
	STRING = "string"
	OBJECT = "object"
)

type AbilityFields map[string]*AbilityField

type AbilityField struct {
	Name        string
	Type        string
	Description string
	Fields      AbilityFields
}