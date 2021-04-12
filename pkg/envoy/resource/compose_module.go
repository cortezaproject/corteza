package resource

import (
	"fmt"
	"strconv"

	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	ComposeModule struct {
		*base
		Res *types.Module

		// Might keep track of related NS
		RefNs   *Ref
		RefMods RefSet
	}
)

func NewComposeModule(res *types.Module, nsRef string) *ComposeModule {
	r := &ComposeModule{
		base:    &base{},
		RefMods: make(RefSet, 0, len(res.Fields)),
	}
	r.SetResourceType(COMPOSE_MODULE_RESOURCE_TYPE)
	r.Res = res

	r.AddIdentifier(identifiers(res.Handle, res.Name, res.ID)...)

	r.RefNs = r.AddRef(COMPOSE_NAMESPACE_RESOURCE_TYPE, nsRef)

	// Field deps
	for _, f := range res.Fields {
		switch f.Kind {
		case "Record":
			refMod := f.Options.String("module")
			if refMod == "" {
				refMod = f.Options.String("moduleID")
			}
			if refMod != "" && refMod != "0" {
				r.RefMods = append(r.RefMods, r.AddRef(COMPOSE_MODULE_RESOURCE_TYPE, refMod).Constraint(r.RefNs))
			}
		}
	}

	// Initial timestamps
	r.SetTimestamps(MakeTimestampsCUDA(&res.CreatedAt, res.UpdatedAt, res.DeletedAt, nil))

	return r
}

func (r *ComposeModule) SysID() uint64 {
	return r.Res.ID
}

func (r *ComposeModule) Ref() string {
	return firstOkString(r.Res.Handle, r.Res.Name, strconv.FormatUint(r.Res.ID, 10))
}

// FindComposeModule looks for the module in the resource set
func FindComposeModule(rr InterfaceSet, ii Identifiers) (ns *types.Module) {
	var modRes *ComposeModule

	rr.Walk(func(r Interface) error {
		mr, ok := r.(*ComposeModule)
		if !ok {
			return nil
		}

		if mr.Identifiers().HasAny(ii) {
			modRes = mr
		}
		return nil
	})

	// Found it
	if modRes != nil {
		return modRes.Res
	}
	return nil
}

func FindComposeModuleResource(rr InterfaceSet, ii Identifiers) (mod *ComposeModule) {
	rr.Walk(func(r Interface) error {
		mr, ok := r.(*ComposeModule)
		if !ok {
			return nil
		}

		if mr.Identifiers().HasAny(ii) {
			mod = mr
		}
		return nil
	})

	return mod
}

func ComposeModuleErrUnresolved(ii Identifiers) error {
	return fmt.Errorf("compose module unresolved %v", ii.StringSlice())
}
