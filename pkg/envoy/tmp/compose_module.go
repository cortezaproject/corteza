package tmp

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	composeModule struct {
		es *encoderState
		s  store.Storer
	}
)

func NewComposeModule(es *encoderState, s store.Storer) *composeModule {
	return &composeModule{
		es: es,
		s:  s,
	}
}

func (p *composeModule) Process(ctx context.Context, state *envoy.ExecState) error {
	res, is := state.Res.(*resource.ComposeModule)
	if !is {
		return nil
	}

	// Get relate namespace
	ns, err := findNamespace(ctx, p.s, state.ParentResources, res.NsRef.Identifiers)
	if err != nil {
		return err
	}
	if ns == nil {
		return errors.New("@todo couldn't resolve namespace")
	}
	p.es.Set(res, res.NsRef.ResourceType, ns.ID, res.NsRef.Identifiers.StringSlice()...)

	// Check if the current module exits
	f := moduleFilterFromGeneric(GenericFilter(res.Identifiers()))
	f.NamespaceID = ns.ID
	mod, err := loadComposeModule(ctx, p.s, f)
	if err != nil {
		return err
	}
	if mod != nil {
		p.es.SetExists(res)
		p.es.Set(res, res.ResourceType(), mod.ID, res.Identifiers().StringSlice()...)
	}

	// Go over missing deps and handle those
	for _, m := range state.MissingDeps {
		switch m.ResourceType {
		case resource.COMPOSE_MODULE_RESOURCE_TYPE:
			f := moduleFilterFromGeneric(GenericFilter(m.Identifiers))
			f.NamespaceID = ns.ID
			mod, err := loadComposeModule(ctx, p.s, f)
			if err != nil {
				return err
			} else if mod == nil {
				continue
			}

			p.es.Set(res, resource.COMPOSE_MODULE_RESOURCE_TYPE, mod.ID, m.Identifiers.StringSlice()...)
		}
	}

	return nil
}

func encodeComposeModule(ctx context.Context, ectx *encodingContext, s store.Storer, state resRefs, res *resource.ComposeModule) (resRefs, error) {
	var (
		mod    = res.Res
		rState = make(resRefs)
	)

	mod.ID = state.Get(res)
	if mod.ID <= 0 {
		mod.ID = nextID()
	}
	if mod.CreatedAt.IsZero() {
		mod.CreatedAt = time.Now()
	}

	rState.Set(resource.COMPOSE_MODULE_RESOURCE_TYPE, mod.ID, res.Identifiers().StringSlice()...)

	// Namespace...
	// A module can exist under a single namespace, so this is good enough for now
	for _, v := range state[resource.COMPOSE_NAMESPACE_RESOURCE_TYPE] {
		mod.NamespaceID = v
		break
	}

	if !ectx.partial && !ectx.exists {
		for i, f := range mod.Fields {
			f.ID = nextID()
			f.ModuleID = mod.ID
			f.Place = i
			f.DeletedAt = nil

			if f.Kind == "Record" {
				refM := f.Options.String("module")
				mID := state[resource.COMPOSE_MODULE_RESOURCE_TYPE][refM]
				if mID > 0 {
					f.Options["module"] = mID
				}
			}
		}

		err := store.CreateComposeModule(ctx, s, mod)
		if err != nil {
			return nil, err
		}
		err = store.CreateComposeModuleField(ctx, s, mod.Fields...)
		if err != nil {
			return nil, err
		}
	}

	return rState, nil
}

func moduleFilterFromGeneric(gf genericFilter) types.ModuleFilter {
	f := types.ModuleFilter{
		Name:   gf.Name,
		Handle: gf.Ref,
	}
	if gf.ID > 0 {
		f.Query = fmt.Sprintf("moduleID=%d", gf.ID)
	}

	return f
}

func findModule(ctx context.Context, s store.Storer, rr []resource.Interface, ii resource.Identifiers) (ns *types.Module, err error) {
	// Try to find it in the parent resources
	var modRes *resource.ComposeModule
	walkResources(rr, func(r resource.Interface) error {
		mR, ok := r.(*resource.ComposeModule)
		if !ok {
			return nil
		}

		if mR.Identifiers().HasAny(r.Identifiers().StringSlice()...) {
			modRes = mR
		}
		return nil
	})

	// Found it
	if modRes != nil {
		return modRes.Res, nil
	}

	// Go in the store
	f := moduleFilterFromGeneric(GenericFilter(ii))
	return loadComposeModule(ctx, s, f)
}
func loadComposeModule(ctx context.Context, s store.Storer, f types.ModuleFilter) (*types.Module, error) {
	mdd, f, err := store.SearchComposeModules(ctx, s, f)
	if err != nil {
		return nil, err
	}

	if len(mdd) > 0 {
		return mdd[0], nil
	}

	return nil, nil
}

func loadComposeModuleFields(ctx context.Context, s store.Storer, mod *types.Module) (types.ModuleFieldSet, error) {
	if mod.ID <= 0 {
		return mod.Fields, nil
	}

	f := types.ModuleFieldFilter{
		ModuleID: []uint64{mod.ID},
	}
	ff, f, err := store.SearchComposeModuleFields(ctx, s, f)
	if err != nil {
		return nil, err
	}

	return ff, nil
}
