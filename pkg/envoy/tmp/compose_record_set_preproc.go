package tmp

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	composeRecordSetPreproc struct {
		is *importState
		s  store.Storer
	}
)

func NewComposeRecordSetPreproc(is *importState, s store.Storer) *composeRecordSetPreproc {
	return &composeRecordSetPreproc{
		is: is,
		s:  s,
	}
}

func (p *composeRecordSetPreproc) Process(ctx context.Context, state *envoy.ExecState) error {
	// @todo can we/should we have the same pattern as with decoder's CanDecode?
	res, is := state.Res.(*resource.ComposeRecordSet)
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

	modID, err := p.module(ctx, res, state)
	if err != nil {
		return err
	} else if modID <= 0 {
		// If the module doesn't exist, no underlying resource is able to exist.
		// @todo generate an error set to show as warnings?
		return nil
	}

	// @todo existing records, related records

	return nil
}

func (p *composeRecordSetPreproc) namespace(ctx context.Context, res *resource.ComposeRecordSet, state *envoy.ExecState) (nsID uint64, err error) {
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

func (p *composeRecordSetPreproc) module(ctx context.Context, res *resource.ComposeRecordSet, state *envoy.ExecState) (modID uint64, err error) {
	modd := filterComposeModuleResources(state.ParentResources)
	if len(modd) > 0 {
		modID = p.is.Existint(modd[0])
	} else {
		mod, idd, err := findMissingComposeModule(ctx, p.s, state.MissingDeps)
		if err != nil {
			return 0, err
		}
		if mod != nil {
			modID = mod.ID
			p.is.AddRefMapping(res, "compose:module", modID, idd.StringSlice()...)
		}
	}

	return modID, nil
}
