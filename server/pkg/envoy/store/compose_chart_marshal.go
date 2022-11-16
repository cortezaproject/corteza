package store

import (
	"context"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza/server/store"
)

func newComposeChartFromResource(res *resource.ComposeChart, cfg *EncoderConfig) resourceState {
	return &composeChart{
		cfg: mergeConfig(cfg, res.Config()),

		res: res,
	}
}

func (n *composeChart) Prepare(ctx context.Context, pl *payload) (err error) {
	// Reset old identifiers
	n.res.Res.ID = 0
	n.res.Res.NamespaceID = 0

	// Get related namespace
	n.relNS, err = findComposeNamespace(ctx, pl.s, pl.state.ParentResources, n.res.RefNs.Identifiers)
	if err != nil {
		return err
	}
	if n.relNS == nil {
		return resource.ComposeNamespaceErrUnresolved(n.res.RefNs.Identifiers)
	}

	// Get related modules
	n.relMods = make(types.ModuleSet, len(n.res.RefMods))
	for i, rMod := range n.res.RefMods {
		n.relMods[i], err = findComposeModule(ctx, pl.s, n.relNS.ID, pl.state.ParentResources, rMod.Identifiers)
		if err != nil {
			return err
		}
		if n.relMods[i] == nil {
			return resource.ComposeModuleErrUnresolved(rMod.Identifiers)
		}
	}

	// Try to get the original chart
	n.chr, err = findComposeChartStore(ctx, pl.s, n.relNS.ID, makeGenericFilter(n.res.Identifiers()))
	if err != nil {
		return err
	}

	if n.chr != nil {
		n.res.Res.ID = n.chr.ID
		n.res.Res.NamespaceID = n.chr.NamespaceID
	}
	return nil
}

func (n *composeChart) Encode(ctx context.Context, pl *payload) (err error) {
	res := n.res.Res
	exists := n.chr != nil && n.chr.ID > 0

	// Determine the ID
	if res.ID <= 0 && exists {
		res.ID = n.chr.ID
	}
	if res.ID <= 0 {
		res.ID = NextID()
	}

	ts := n.res.Timestamps()
	if ts != nil {
		if ts.CreatedAt != nil {
			res.CreatedAt = *ts.CreatedAt.T
		} else {
			res.CreatedAt = *now()
		}
		if ts.UpdatedAt != nil {
			res.UpdatedAt = ts.UpdatedAt.T
		}
		if ts.DeletedAt != nil {
			res.DeletedAt = ts.DeletedAt.T
		}
	}

	// Namespace
	res.NamespaceID = n.relNS.ID
	if res.NamespaceID <= 0 {
		ns := resource.FindComposeNamespace(pl.state.ParentResources, n.res.RefNs.Identifiers)
		res.NamespaceID = ns.ID
	}
	if res.NamespaceID <= 0 {
		return resource.ComposeNamespaceErrUnresolved(n.res.RefNs.Identifiers)
	}

	// Report modules
	for i, r := range res.Config.Reports {
		relMod := n.relMods[i]
		if relMod == nil {
			relMod = resource.FindComposeModule(pl.state.ParentResources, n.res.RefMods[i].Identifiers)
		}
		if relMod == nil || relMod.ID <= 0 {
			return resource.ComposeModuleErrUnresolved(n.res.RefMods[i].Identifiers)
		}

		r.ModuleID = relMod.ID
	}

	// Evaluate the resource skip expression
	// @todo expand available parameters; similar implementation to compose/types/record@Dict
	if skip, err := basicSkipEval(ctx, n.cfg, !exists); err != nil {
		return err
	} else if skip {
		return nil
	}

	// Create fresh chart
	if !exists {
		return store.CreateComposeChart(ctx, pl.s, res)
	}

	// Update existing chart
	switch n.cfg.OnExisting {
	case resource.Skip:
		return nil

	case resource.MergeLeft:
		res = mergeComposeChart(n.chr, res)

	case resource.MergeRight:
		res = mergeComposeChart(res, n.chr)
	}

	err = store.UpdateComposeChart(ctx, pl.s, res)
	if err != nil {
		return err
	}

	n.res.Res = res
	return nil
}
