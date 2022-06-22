package check

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type State string

const (
	StateOK      State = "ok"
	StateFailed  State = "failed"
	StateWarning State = "warning"
	StateUnable  State = "unable to run"
)

type Check interface {
	Run()
	Export() GenericCheck

	ValidateParameters(requiredKeywords []string) error

	SetParameter(key, value string)
	SetHelp(text string)
	SetName(name string)
	SetIdentifier(identifier string)

	StartTiming()
	FinishTiming()

	Logger() *log.Logger
}

type GenericCheck struct {
	Name       string
	Help       string
	State      State
	Error      error
	Identifier string
	Parameters map[string]string

	Duration  time.Duration
	startTime time.Time
}

func (generic *GenericCheck) SetName(name string) {
	generic.Name = name
}
func (generic *GenericCheck) SetIdentifier(identifier string) {
	generic.Identifier = identifier
}

func (generic *GenericCheck) SetHelp(text string) {
	generic.Help = text
}

func (generic *GenericCheck) StartTiming() {
	generic.startTime = time.Now()
}

func (generic *GenericCheck) FinishTiming() {
	generic.Duration = time.Since(generic.startTime)
}

func (generic *GenericCheck) SetParameter(key, value string) {
	if generic.Parameters == nil {
		generic.Parameters = make(map[string]string)
	}
	generic.Parameters[key] = value
}

func (generic *GenericCheck) Export() GenericCheck {
	return *generic
}

func (generic *GenericCheck) Logger() *log.Logger {
	return log.New(os.Stderr, "", log.Lmsgprefix)
}

func (generic *GenericCheck) ValidateParameters(requiredKeywords []string) error {
	for _, keyword := range requiredKeywords {
		if value, ok := generic.Parameters[keyword]; value == "" || !ok {
			err := fmt.Errorf("keyword %+q not set", keyword)
			generic.State = StateUnable
			generic.Error = err
			return err
		}
		continue
	}
	return nil
}

type Registry struct {
	Checks map[string]func() Check

	mux sync.Mutex
}

func (r *Registry) Register(name string, factory func() Check) {
	if r.Checks == nil {
		r.Checks = make(map[string]func() Check)
	}
	r.mux.Lock()
	r.Checks[name] = factory
	r.mux.Unlock()
}

var RegistryInstance Registry
