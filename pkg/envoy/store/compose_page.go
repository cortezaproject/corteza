package store

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	composePageState struct {
		cfg *EncoderConfig

		res *resource.ComposePage
		pg  *types.Page

		relNS     *types.Namespace
		relMod    *types.Module
		relParent *types.Page

		relMods   map[string]*types.Module
		relCharts map[string]*types.Chart
	}
)

func NewComposePageState(res *resource.ComposePage, cfg *EncoderConfig) resourceState {
	return &composePageState{
		cfg: mergeConfig(cfg, res.Config()),

		res: res,

		relMods:   make(map[string]*types.Module),
		relCharts: make(map[string]*types.Chart),
	}
}

func (n *composePageState) Prepare(ctx context.Context, s store.Storer, state *envoy.ResourceState) (err error) {
	// Initial values
	if n.res.Res.CreatedAt.IsZero() {
		n.res.Res.CreatedAt = time.Now()
	}

	// Get relate namespace
	n.relNS, err = findComposeNamespaceRS(ctx, s, state.ParentResources, n.res.NsRef.Identifiers)
	if err != nil {
		return err
	}
	if n.relNS == nil {
		return composeNamespaceErrUnresolved(n.res.NsRef.Identifiers)
	}

	// Get related module
	// If this isn't a record page, there is no related module
	if n.res.ModRef != nil {
		n.relMod, err = findComposeModuleRS(ctx, s, n.relNS.ID, state.ParentResources, n.res.ModRef.Identifiers)
		if err != nil {
			return err
		}
		if n.relMod == nil {
			return composeModuleErrUnresolved(n.res.ModRef.Identifiers)
		}
	}

	// Get parent page
	if n.res.ParentRef != nil {
		n.relParent, err = findComposePageRS(ctx, s, n.relNS.ID, state.ParentResources, n.res.ParentRef.Identifiers)
		if err != nil {
			return err
		}
		if n.relParent == nil {
			return composePageErrUnresolved(n.res.ParentRef.Identifiers)
		}
	}

	// Get other related modules
	for _, mr := range n.res.ModRefs {
		mod, err := findComposeModuleRS(ctx, s, n.relNS.ID, state.ParentResources, mr.Identifiers)
		if err != nil {
			return err
		}
		if mod == nil {
			return composeModuleErrUnresolved(mr.Identifiers)
		}
		for id := range mr.Identifiers {
			n.relMods[id] = mod
		}
	}

	// Get related charts
	for _, cr := range n.res.ChartRefs {
		chr, err := findComposeChartRS(ctx, s, n.relNS.ID, state.ParentResources, cr.Identifiers)
		if err != nil {
			return err
		}
		if chr == nil {
			return composeChartErrUnresolved(cr.Identifiers)
		}
		for id := range cr.Identifiers {
			n.relCharts[id] = chr
		}
	}

	// Try to get the original page
	n.pg, err = findComposePageS(ctx, s, n.relNS.ID, makeGenericFilter(n.res.Identifiers()))
	if err != nil {
		return err
	}

	if n.pg != nil {
		n.res.Res.ID = n.pg.ID
		n.res.Res.NamespaceID = n.pg.NamespaceID
	}
	return nil
}

func (n *composePageState) Encode(ctx context.Context, s store.Storer, state *envoy.ResourceState) (err error) {
	res := n.res.Res
	exists := n.pg != nil && n.pg.ID > 0

	// Determine the ID
	if res.ID <= 0 && exists {
		res.ID = n.pg.ID
	}
	if res.ID <= 0 {
		res.ID = NextID()
	}

	// This is not possible, but let's do it anyway
	if state.Conflicting {
		return nil
	}

	// Timestamps
	ts := n.res.Timestamps()
	if ts != nil {
		if ts.CreatedAt != "" {
			t := toTime(ts.CreatedAt)
			if t != nil {
				res.CreatedAt = *t
			}
		}
		if ts.UpdatedAt != "" {
			res.UpdatedAt = toTime(ts.UpdatedAt)
		}
		if ts.DeletedAt != "" {
			res.DeletedAt = toTime(ts.DeletedAt)
		}
	}

	// Namespace
	res.NamespaceID = n.relNS.ID
	if res.NamespaceID <= 0 {
		ns := findComposeNamespaceR(state.ParentResources, n.res.NsRef.Identifiers)
		res.NamespaceID = ns.ID
	}

	if res.NamespaceID <= 0 {
		return composeNamespaceErrUnresolved(n.res.NsRef.Identifiers)
	}

	// Module?
	if n.res.ModRef != nil {
		res.ModuleID = n.relMod.ID
		if res.ModuleID <= 0 {
			mod := findComposeModuleR(state.ParentResources, n.res.ModRef.Identifiers)
			res.ModuleID = mod.ID
		}
	}

	// Parent?
	if n.res.ParentRef != nil {
		res.SelfID = n.relParent.ID
		if res.SelfID <= 0 {
			mod := findComposePageR(state.ParentResources, n.res.ParentRef.Identifiers)
			res.SelfID = mod.ID
		}
	}

	// Blocks
	getModID := func(id string) uint64 {
		mod := n.relMods[id]
		if mod == nil || mod.ID <= 0 {
			mod = findComposeModuleR(state.ParentResources, resource.MakeIdentifiers(id))
			if mod == nil || mod.ID <= 0 {
				return 0
			}
		}
		return mod.ID
	}

	for _, b := range res.Blocks {
		switch b.Kind {
		case "RecordList":
			id, _ := b.Options["module"].(string)
			if id == "" {
				continue
			}
			mID := getModID(id)
			if mID <= 0 {
				return composeModuleErrUnresolved(resource.MakeIdentifiers(id))
			}
			b.Options["moduleID"] = strconv.FormatUint(mID, 10)
			delete(b.Options, "module")

		case "Calendar":
			ff, _ := b.Options["feeds"].([]interface{})
			for _, f := range ff {
				feed, _ := f.(map[string]interface{})
				fOpts, _ := (feed["options"]).(map[string]interface{})
				id, _ := fOpts["module"].(string)
				if id == "" {
					continue
				}
				mID := getModID(id)
				if mID <= 0 {
					return composeModuleErrUnresolved(resource.MakeIdentifiers(id))
				}
				fOpts["moduleID"] = strconv.FormatUint(mID, 10)
				delete(fOpts, "module")
			}

		case "Chart":
			id, _ := b.Options["chart"].(string)
			if id == "" {
				continue
			}
			chr := n.relCharts[id]
			if chr == nil || chr.ID <= 0 {
				ii := resource.MakeIdentifiers(id)
				chr = findComposeChartR(state.ParentResources, ii)
				if chr == nil || chr.ID <= 0 {
					return composeChartErrUnresolved(ii)
				}
			}
			b.Options["chartID"] = strconv.FormatUint(chr.ID, 10)
			delete(b.Options, "chart")

		case "Metric":
			mm, _ := b.Options["metrics"].([]interface{})
			for _, m := range mm {
				mops, _ := m.(map[string]interface{})
				id, _ := mops["module"].(string)
				if id == "" {
					continue
				}
				mID := getModID(id)
				if mID <= 0 {
					return composeModuleErrUnresolved(resource.MakeIdentifiers(id))
				}
				mops["moduleID"] = strconv.FormatUint(mID, 10)
				delete(mops, "module")

			}
		}
	}

	// Evaluate the resource skip expression
	// @todo expand available parameters; similar implementation to compose/types/record@Dict
	if skip, err := basicSkipEval(ctx, n.cfg, !exists); err != nil {
		return err
	} else if skip {
		return nil
	}

	// Create a fresh page
	if !exists {
		return store.CreateComposePage(ctx, s, res)
	}

	// Update existing page
	switch n.cfg.OnExisting {
	case resource.Skip:
		return nil

	case resource.MergeLeft:
		res = mergeComposePage(n.pg, res)

	case resource.MergeRight:
		res = mergeComposePage(res, n.pg)
	}

	err = store.UpdateComposePage(ctx, s, res)
	if err != nil {
		return err
	}

	n.res.Res = res
	return nil
}

// mergeComposePage merges b into a, prioritising a
func mergeComposePage(a, b *types.Page) *types.Page {
	c := a.Clone()

	c.SelfID = b.SelfID
	c.NamespaceID = b.NamespaceID
	c.ModuleID = b.ModuleID
	c.Weight = b.Weight
	if c.Handle == "" {
		c.Handle = b.Handle
	}
	if c.Title == "" {
		c.Title = b.Title
	}
	if c.Description == "" {
		c.Description = b.Description
	}
	if len(c.Blocks) <= 0 {
		c.Blocks = b.Blocks
	}
	if len(c.Children) <= 0 {
		c.Children = b.Children
	}

	if c.CreatedAt.IsZero() {
		c.CreatedAt = b.CreatedAt
	}
	if c.UpdatedAt == nil {
		c.UpdatedAt = b.UpdatedAt
	}
	if c.DeletedAt == nil {
		c.DeletedAt = b.DeletedAt
	}

	return c
}

// findComposePageRS looks for the page in the resources & the store
//
// Provided resources are prioritized.
func findComposePageRS(ctx context.Context, s store.Storer, nsID uint64, rr resource.InterfaceSet, ii resource.Identifiers) (pg *types.Page, err error) {
	pg = findComposePageR(rr, ii)
	if pg != nil {
		return pg, nil
	}

	if nsID <= 0 {
		return nil, nil
	}

	// Go in the store
	return findComposePageS(ctx, s, nsID, makeGenericFilter(ii))
}

// findComposePageS looks for the page in the store
func findComposePageS(ctx context.Context, s store.Storer, nsID uint64, gf genericFilter) (pg *types.Page, err error) {
	if nsID == 0 {
		return nil, nil
	}

	if gf.id > 0 {
		pg, err = store.LookupComposePageByID(ctx, s, gf.id)
		if err != nil && err != store.ErrNotFound {
			return nil, err
		}

		if pg != nil {
			return
		}
	}

	for _, i := range gf.identifiers {
		pg, err = store.LookupComposePageByNamespaceIDHandle(ctx, s, nsID, i)
		if err == store.ErrNotFound {
			var pp types.PageSet
			pp, _, err = store.SearchComposePages(ctx, s, types.PageFilter{
				NamespaceID: nsID,
				Title:       i,
				Paging: filter.Paging{
					Limit: 2,
				},
			})
			if len(pp) > 1 {
				return nil, resourceErrIdentifierNotUnique(i)
			}
			if len(pp) == 1 {
				pg = pp[0]
			}
		}

		if err != nil && err != store.ErrNotFound {
			return nil, err
		}

		if pg != nil {
			return
		}
	}

	return nil, nil
}

// findComposePageR looks for the page in the resources
func findComposePageR(rr resource.InterfaceSet, ii resource.Identifiers) (pg *types.Page) {
	var pgRes *resource.ComposePage

	rr.Walk(func(r resource.Interface) error {
		pr, ok := r.(*resource.ComposePage)
		if !ok {
			return nil
		}

		if pr.Identifiers().HasAny(ii) {
			pgRes = pr
		}
		return nil
	})

	// Found it
	if pgRes != nil {
		return pgRes.Res
	}
	return nil
}

func composePageErrUnresolved(ii resource.Identifiers) error {
	return fmt.Errorf("compose page unresolved %v", ii.StringSlice())
}
