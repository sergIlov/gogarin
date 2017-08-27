package satellite

import "fmt"

type Interface interface {
	Info() Info
	Triggers() []Trigger
	Actions() []Action
	Filters() []Filter
	Splitters() []Splitter
	Modifiers() []Modifier
}

type Info struct {
	Name        string
	Version     string
	Description string
}

type Base struct {
}

func (s *Base) Info() Info {
	var i Info
	return i
}

func (s *Base) Triggers() []Trigger {
	var ts []Trigger
	return ts
}

func (s *Base) Actions() []Action {
	var as []Action
	return as
}

func (s *Base) Filters() []Filter {
	var fs []Filter
	return fs
}

func (s *Base) Splitters() []Splitter {
	var ss []Splitter
	return ss
}

func (s *Base) Modifiers() []Modifier {
	var ms []Modifier
	return ms
}

type Message struct {
	Data map[string]interface{}
}

type Ability interface {
	Info() AbilityInfo
}

type AbilityInfo struct {
	Name        string
	Pull        bool
	Push        bool
	Description string
}

// Trigger is a source of messages.
type Trigger interface {
	Ability
	// Messages returns a list of new messages.
	// This is the place where the message processing begins.
	Messages() ([]*Message, error)
}

// Action does some work to achieve Mission objectives.
type Action interface {
	Ability
	// Handle handles each message.
	Handle([]*Message) error
}

// Filter filters messages in the pipeline using predefined criteria.
type Filter interface {
	Ability
	// Select returns only those messages that satisfy predefined criteria.
	Select([]*Message) ([]*Message, error)
}

// Splitter creates new messages from the attributes of the message.
// For example:
//		func Split(ms []*Message) ([]*Message, error) {
//			var result []*Message
//	 		for _, client := range Message.Data.Clients {
//				result = append(result, messages.New(data: client.Data))
//	 		}
//			return result, nil
//		}
type Splitter interface {
	Ability
	// Split can return one or more messages based on the attributes of the message.
	Split([]*Message) ([]*Message, error)
}

// Modifier can add, remove, and delete attributes of a message using predefined configurations.
// In some cases a modifier can completely replace the attributes of a message.
type Modifier interface {
	Ability
	// Modify returns the list of modified messages.
	Modify([]*Message) ([]*Message, error)
}

func Register(s Interface) error {
	fmt.Printf("Satellite:\n%+v\n", s.Info())

	if len(s.Triggers()) > 0 {
		fmt.Println("\nTriggers:")
	}
	for _, t := range s.Triggers() {
		fmt.Printf("%+v\n", t.Info())
	}

	if len(s.Filters()) > 0 {
		fmt.Println("\nFilters:")
	}
	for _, f := range s.Filters() {
		fmt.Printf("%+v\n", f.Info())
	}

	if len(s.Modifiers()) > 0 {
		fmt.Println("\nModifiers:")
	}
	for _, m := range s.Modifiers() {
		fmt.Printf("%+v\n", m.Info())
	}

	if len(s.Splitters()) > 0 {
		fmt.Println("\nSplitters:")
	}
	for _, s := range s.Splitters() {
		fmt.Printf("%+v\n", s.Info())
	}

	if len(s.Actions()) > 0 {
		fmt.Println("\nActions:")
	}
	for _, a := range s.Actions() {
		fmt.Printf("%+v\n", a.Info())
	}

	return nil
}
