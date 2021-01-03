package service

import (
	"github.com/cortezaproject/corteza-server/automation/types/fn"
	"sync"
)

var (
	fnRegistry = make(map[string]*fn.Function)
	fnRegMutex sync.RWMutex
)

func RegisterFn(fn *fn.Function) {
	defer fnRegMutex.Unlock()
	fnRegMutex.Lock()
	fnRegistry[fn.Ref] = fn
}

func UnregisterFn(name string) {
	fnRegMutex.Lock()
	delete(fnRegistry, name)
	fnRegMutex.Unlock()
}

func RegisteredFn(name string) *fn.Function {
	defer fnRegMutex.RUnlock()
	fnRegMutex.RLock()
	return fnRegistry[name]
}
