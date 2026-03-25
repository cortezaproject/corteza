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
		functions: make([]types.ConstructFunction, 0, 2),
		triggers:  make([]types.ConstructTrigger, 0, 2),
	}
}

func (r *constructRegistry) AddFunctions(ff ...types.ConstructFunction) {
	r.mux.Lock()
	defer r.mux.Unlock()

	for _, fn := range ff {
		replaced := false
		for i, existing := range r.functions {
			if existing.Ref == fn.Ref {
				r.functions[i] = fn
				replaced = true
				break
			}
		}
		if !replaced {
			r.functions = append(r.functions, fn)
		}
	}
}

func (r *constructRegistry) AddTriggers(tt ...types.ConstructTrigger) {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.triggers = append(r.triggers, tt...)
}

func (r *constructRegistry) Function(ref string) (out types.ConstructFunction, ok bool) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	for _, fn := range r.functions {
		if fn.Ref == ref {
			return fn, true
		}
	}

	ok = false
	return
}

func (r *constructRegistry) Functions() []types.ConstructFunction {
	r.mux.RLock()
	defer r.mux.RUnlock()

	return r.functions
}

func (r *constructRegistry) Triggers() []types.ConstructTrigger {
	r.mux.RLock()
	defer r.mux.RUnlock()

	return r.triggers
}
