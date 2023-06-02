package automation

import (
	"github.com/cortezaproject/corteza/server/automation/types"
)

type (
	AutomationRegistry interface {
		AddFunctions(ff ...*types.Function)
		// AddTypes(tt ...expr.Type)
	}
)
