package corredor

import (
	"github.com/cortezaproject/corteza-server/pkg/automation"
)

type (
	Runnable interface {
		IsAsync() bool
		GetName() string
		GetSource() string
		GetTimeout() uint32
	}
)

func FromScript(s *automation.Script) *Script {
	return &Script{
		Source:  s.Source,
		Name:    s.Name,
		Timeout: uint32(s.Timeout),
		Async:   s.Async,
	}
}
