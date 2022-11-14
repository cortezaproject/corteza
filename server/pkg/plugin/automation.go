package plugin

// Collection of boot-lifecycle related functions
// that exec plugin functions

import (
	"github.com/cortezaproject/corteza-server/automation/types"
	sdk "github.com/cortezaproject/corteza-server/sdk/plugin"
)

type (
	automationRegistry interface {
		AddFunctions(ff ...*types.Function)
		// AddTypes(tt ...expr.Type)
	}
)

func (pp Set) RegisterAutomation(r automationRegistry) error {
	for _, p := range pp {
		d, is := p.def.(sdk.AutomationFunctionsProvider)
		if !is {
			continue
		}

		r.AddFunctions(d.AutomationFunctions()...)
	}

	return nil
}
