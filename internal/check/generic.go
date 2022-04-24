package check

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type State string

type Check interface {
	Run()
	Export() GenericCheck

	ValidateParameters(requiredKeywords []string) error

	SetParameter(key, value string)
	SetHelp(text string)
	SetName(name string)

	Logger() *log.Logger
}

type GenericCheck struct {
	Name       string
	Help       string
	State      State
	Error      error
	Parameters map[string]string
}

func (generic *GenericCheck) SetName(name string) {
	generic.Name = name
}

func (generic *GenericCheck) SetHelp(text string) {
	generic.Help = text
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

const (
	StateOK      State = "ok"
	StateFailed  State = "failed"
	StateWarning State = "warning"
	StateUnable  State = "unable to run"
)
