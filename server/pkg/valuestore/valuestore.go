package valuestore

import (
	"strings"
)

type (
	store struct {
		env map[string]any
	}
)

var (
	gStore *store
)

func SetGlobal(s *store) {
	gStore = s
}

func New() *store {
	return &store{}
}

func Global() *store {
	return gStore
}

func EnvGetter() func(string) any {
	return gStore.Env
}

func (s *store) SetEnv(env map[string]any) {
	if s.env != nil {
		panic("cannot redefine environment variables")
	}

	s.env = env
}

func (s *store) Env(k string) (v any) {
	if s.env == nil {
		panic("valuestore env not initialized")
	}
	return s.env[strings.ToLower(k)]
}
