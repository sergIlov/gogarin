package satellite

import "fmt"

type SpaceCenter struct {
}

func (sc *SpaceCenter) Register(s *Satellite) {
	fmt.Printf("Satellite:\n%+v\n", s.info)

	if len(s.triggers) > 0 {
		fmt.Println("\nTriggers:")
	}
	for i := range s.triggers {
		fmt.Printf("%+v\n", i)
	}
}

type SpaceCenterConfig struct{}

type Satellite struct {
	info     Info
	triggers map[AbilityInfo]Trigger
}

func (s *Satellite) Start(sc *SpaceCenter) {
	sc.Register(s)
}

func (s *Satellite) Stop() {
}

type Trigger func()

func (t Trigger) Validate() {

}

func (s *Satellite) AddTrigger(t Trigger, i AbilityInfo) {
	s.triggers[i] = t
}

func NewSpaceCenter(c SpaceCenterConfig) *SpaceCenter {
	return &SpaceCenter{}
}

func New(i Info) *Satellite {
	return &Satellite{
		info:     i,
		triggers: make(map[AbilityInfo]Trigger),
	}
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
