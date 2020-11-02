package tmp

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	composeModulePreproc struct {
		is *importState
		s  store.Storer
	}
)

func NewComposeModulePreproc(is *importState, s store.Storer) envoy.Processor {
	return &composeModulePreproc{
		is: is,
		s:  s,
	}
}

func (p *composeModulePreproc) Process(ctx context.Context, state *envoy.ExecState) error {
	// @todo can we/should we have the same pattern as with decoder's CanDecode?
	res, is := state.Res.(*resource.ComposeModule)
	if !is {
		return nil
	}

	nsID, err := p.namespace(ctx, res, state)
	if err != nil {
		return err
	} else if nsID <= 0 {
		// If the namespace doesn't exist, no underlying resource is able to exist.
		// @todo generate an error set to show as warnings?
		return nil
	}

	// Check if the current module exits
	f := res.SearchQuery()
	f.NamespaceID = nsID
	mod, err := loadComposeModule(ctx, p.s, f)
	if err != nil {
		return err
	}
	if mod != nil {
		p.is.Existint(res)
		p.is.AddRefMapping(res, res.ResourceType(), mod.ID, res.Identifiers().StringSlice()...)
	}

	// Go over missing deps and handle those
	for _, m := range state.MissingDeps {
		switch m.ResourceType {
		// @todo change this string to a better thing...
		case "compose:module":
			f := moduleFilterFromGeneric(GetGenericFilter(m.Identifiers))
			f.NamespaceID = nsID
			mod, err := loadComposeModule(ctx, p.s, f)
			if err != nil {
				return err
			} else if mod == nil {
				continue
			}

			p.is.AddRefMapping(res, "compose:module", mod.ID, m.Identifiers.StringSlice()...)
		}
	}

	return nil
}

func (p *composeModulePreproc) namespace(ctx context.Context, res *resource.ComposeModule, state *envoy.ExecState) (nsID uint64, err error) {
	nss := filterComposeNamespaceResources(state.ParentResources)
	if len(nss) > 0 {
		nsID = p.is.Existint(nss[0])
	} else {
		ns, idd, err := findMissingComposeNamespace(ctx, p.s, state.MissingDeps)
		if err != nil {
			return 0, err
		}
		if ns != nil {
			nsID = ns.ID
			p.is.AddRefMapping(res, "compose:namespace", nsID, idd.StringSlice()...)
		}
	}

	return nsID, nil
}
