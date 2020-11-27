package resource

import (
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	ComposeModule struct {
		*base
		Res *types.Module

		// Might keep track of related NS
		NsRef  *Ref
		ModRef RefSet
	}
)

func NewComposeModule(res *types.Module, nsRef string) *ComposeModule {
	r := &ComposeModule{
		base:   &base{},
		ModRef: make(RefSet, 0, len(res.Fields)),
	}
	r.SetResourceType(COMPOSE_MODULE_RESOURCE_TYPE)
	r.Res = res

	r.AddIdentifier(identifiers(res.Handle, res.Name, res.ID)...)

	r.NsRef = r.AddRef(COMPOSE_NAMESPACE_RESOURCE_TYPE, nsRef)

	// Field deps
	for _, f := range res.Fields {
		switch f.Kind {
		case "Record":
			refM := f.Options.String("module")
			if refM != "" && refM != "0" {
				r.ModRef = append(r.ModRef, r.AddRef(COMPOSE_MODULE_RESOURCE_TYPE, refM))
			}
		}
	}

	return r
}

func (r *ComposeModule) SysID() uint64 {
	return r.Res.ID
}
