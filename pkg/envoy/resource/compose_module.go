package resource

import (
	"fmt"

	"github.com/cortezaproject/corteza-server/compose/types"
)

const (
	COMPOSE_MODULE_RESOURCE_TYPE = "composeModule"
)

type (
	composeModule struct {
		*base
		Res *types.Module

		// Might keep track of related NS
	}
)

func ComposeModule(res *types.Module) *composeModule {
	r := &composeModule{base: &base{}}
	r.SetResourceType(COMPOSE_MODULE_RESOURCE_TYPE)
	r.Res = res

	r.AddIdentifier(identifiers(res.Handle, res.Name, res.ID)...)

	// Field deps.
	for _, f := range res.Fields {
		switch f.Kind {
		case "Record":
			refM := f.Options.String("module")
			if refM != "" && refM != "0" {
				r.AddRef(COMPOSE_MODULE_RESOURCE_TYPE, refM)
			}
		}
	}

	return r
}

func (m *composeModule) SearchQuery() types.ModuleFilter {
	f := types.ModuleFilter{
		Handle: m.Res.Handle,
		Name:   m.Res.Name,
	}

	if m.Res.ID > 0 {
		f.Query = fmt.Sprintf("moduleID=%d", m.Res.ID)
	}

	return f
}
