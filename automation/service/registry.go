package service

import (
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"sort"
	"sync"
)

type (
	registry struct {
		lock      *sync.RWMutex
		functions map[string]*types.Function
		types     map[string]expr.Type
	}

	handler interface {
		Functions() []*types.Function
	}
)

var (
	defaultRegistry = initRegistry()
)

func Registry() *registry {
	return defaultRegistry
}

func initRegistry() *registry {
	return &registry{
		lock:      &sync.RWMutex{},
		functions: make(map[string]*types.Function),
		types:     make(map[string]expr.Type),
	}
}

func (r *registry) AddFunctions(ff ...*types.Function) {
	defer r.lock.Unlock()
	r.lock.Lock()
	for _, fn := range ff {
		r.functions[fn.Ref] = fn
	}
}

func (r *registry) AddHandlers(hh ...handler) {
	for _, h := range hh {
		r.AddFunctions(h.Functions()...)
	}
}

func (r registry) Function(ref string) *types.Function {
	defer r.lock.RUnlock()
	r.lock.RLock()
	return r.functions[ref]
}

func (r registry) Functions() []*types.Function {
	var (
		rr = make([]string, 0, len(r.functions))
		ff = make([]*types.Function, 0, len(r.functions))
	)

	for ref := range r.functions {
		rr = append(rr, ref)
	}

	sort.Strings(rr)

	for _, ref := range rr {
		ff = append(ff, r.functions[ref])
	}

	return ff
}

func (r *registry) AddTypes(tt ...expr.Type) {
	defer r.lock.Unlock()
	r.lock.Lock()
	for _, t := range tt {
		r.types[t.Type()] = t
	}
}

func (r *registry) Type(ref string) expr.Type {
	defer r.lock.RUnlock()
	r.lock.RLock()
	return r.types[ref]
}

func (r *registry) Types() []string {
	var (
		rr = make([]string, 0, len(r.types))
	)

	for ref := range r.types {
		rr = append(rr, ref)
	}

	sort.Strings(rr)

	return rr
}
