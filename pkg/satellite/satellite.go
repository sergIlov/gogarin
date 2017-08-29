package satellite

func New(i Info) *Satellite {
	return &Satellite{
		Info:     i,
		Triggers: make(map[AbilityInfo]Trigger),
	}
}

type Satellite struct {
	Info     Info
	Triggers map[AbilityInfo]Trigger
}

func (s *Satellite) AddTrigger(t Trigger, i AbilityInfo) {
	s.Triggers[i] = t
}

func (s *Satellite) Start(c Connector) {
	c.Register(s)
}

func (s *Satellite) Stop() {
}

type Trigger func()

type Connector interface {
	Register(*Satellite) error
}

type Info struct {
	Name        string
	Version     string
	Description string
}

type AbilityInfo struct {
	Name        string
	Description string
}

//
//type Ability interface {
//	ConfigScheme() (AbilityConfigScheme, error)
//	Validate(AbilitySettings) error
//	Join(AbilitySettings) error
//	Leave(AbilitySettings) error
//}
//
