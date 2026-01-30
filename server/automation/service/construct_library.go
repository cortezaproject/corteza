package service

import (
	"sync"

	"github.com/cortezaproject/corteza/server/automation/types"
)

type (
	constructRegistry struct {
		mux       *sync.RWMutex
		triggers  []types.ConstructTrigger
		functions []types.ConstructFunction
	}
)

var (
	defaultConstructRegistry = initConstructRegistry()
)

func ConstructLibrary() *constructRegistry {
	return defaultConstructRegistry
}

func initConstructRegistry() *constructRegistry {
	return &constructRegistry{
		mux:       &sync.RWMutex{},
		functions: make([]types.ConstructFunction, 2),
		triggers:  make([]types.ConstructTrigger, 2),
	}
}

func (r *constructRegistry) Functions() []types.ConstructFunction {
	return r.functions
}

func (r *constructRegistry) Triggers() []types.ConstructTrigger {
	return r.triggers
}
