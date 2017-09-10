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

type AbilityFieldType struct {
	Name        string
	Description string
}

func NewAbilityFieldType(name string, description string) AbilityFieldType {
	return AbilityFieldType{name, description}
}

var (
	Boolean    = NewAbilityFieldType("boolean", "")
	Integer    = NewAbilityFieldType("integer", "")
	Float      = NewAbilityFieldType("Float", "")
	String     = NewAbilityFieldType("string", "")
	Date       = NewAbilityFieldType("date", "")
	Datetime   = NewAbilityFieldType("datetime", "")
	Object     = NewAbilityFieldType("object", "")
	Collection = NewAbilityFieldType("collection", "")
)

type AbilityFields map[string]*AbilityField

type AbilityField struct {
	Name        string
	Type        AbilityFieldType
	Description string
	Fields      AbilityFields
}
