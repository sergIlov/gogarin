package satellite

func New(i Info) *Satellite {
	return &Satellite{
		Info: i,
	}
}

type Satellite struct {
	Info     Info
	Triggers []Trigger
}

func (s *Satellite) AddTrigger(t Trigger) {
	s.Triggers = append(s.Triggers, t)
}

func (s *Satellite) Start(c SpaceCenterConnector) {
	c.Register(s)
}

func (s *Satellite) Stop() {
}

type SpaceCenterConnector interface {
	Register(*Satellite) error
}

type Info struct {
	Name        string
	Version     string
	Description string
}

type Trigger struct {
	Call   func()
	Info   AbilityInfo
	Config interface{}
	Validator func(config interface{})

}

type AbilityInfo struct {
	Name        string
	Description string
}
