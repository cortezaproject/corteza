package store

import (
	"context"
	"errors"
	"time"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	composePageState struct {
		cfg *EncoderConfig

		res *resource.ComposePage
		pg  *types.Page

		relNS  *types.Namespace
		relMod *types.Module
	}
)

func NewComposePageState(res *resource.ComposePage, cfg *EncoderConfig) resourceState {
	return &composePageState{
		cfg: cfg,

		res: res,
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
		return errors.New("@todo couldn't resolve namespace")
	}

	// Can't do anything else, since the NS doesn't yet exist
	if n.relNS.ID <= 0 {
		return nil
	}

	// Get related module
	// If this isn't a record page, there is no related module
	if n.res.ModRef != nil {
		n.relMod, err = findComposeModuleRS(ctx, s, n.relNS.ID, state.ParentResources, n.res.ModRef.Identifiers)
		if err != nil {
			return err
		}
		if n.relMod == nil {
			return errors.New("@todo couldn't resolve module")
		}
	}

	// Try to get the original page
	n.pg, err = findComposePageS(ctx, s, n.relNS.ID, makeGenericFilter(n.res.Identifiers()))
	if err != nil {
		return err
	}

	if n.pg != nil {
		n.res.Res.ID = n.pg.ID
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
		res.ID = nextID()
	}

	// This is not possible, but let's do it anyway
	if state.Conflicting {
		return nil
	}

	// Namespace
	res.NamespaceID = n.relNS.ID
	if res.NamespaceID <= 0 {
		ns := findComposeNamespaceR(state.ParentResources, n.res.NsRef.Identifiers)
		res.NamespaceID = ns.ID
	}

	if res.NamespaceID <= 0 {
		return errors.New("[chart] couldn't find related namespace; @todo error")
	}

	// Module?
	if n.res.ModRef != nil {
		res.ModuleID = n.relMod.ID
		if res.ModuleID <= 0 {
			mod := findComposeModuleR(state.ParentResources, n.res.ModRef.Identifiers)
			res.ModuleID = mod.ID
		}
	}

	// Create a fresh page
	if !exists {
		return store.CreateComposePage(ctx, s, res)
	}

	// Update existing page
	switch n.cfg.OnExisting {
	case Skip:
		return nil

	case MergeLeft:
		res = mergeComposePage(n.pg, res)

	case MergeRight:
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

	if c.SelfID <= 0 {
		c.SelfID = b.SelfID
	}
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

	return c
}

// findComposePageRS looks for the page in the resources & the store
//
// Provided resources are prioritized.
func findComposePageRS(ctx context.Context, s store.Storer, nsID uint64, rr resource.InterfaceSet, ii resource.Identifiers) (pg *types.Page, err error) {
	pg = findComposePageR(ctx, rr, ii)
	if pg != nil {
		return pg, nil
	}

	// Go in the store
	return findComposePageS(ctx, s, nsID, makeGenericFilter(ii))
}

// findComposePageS looks for the page in the store
func findComposePageS(ctx context.Context, s store.Storer, nsID uint64, gf genericFilter) (pg *types.Page, err error) {
	if gf.id > 0 {
		pg, err = store.LookupComposePageByID(ctx, s, gf.id)
		if err != nil && err != store.ErrNotFound {
			return nil, err
		}

		if pg != nil {
			return
		}
	}

	if gf.handle != "" {
		pg, err = store.LookupComposePageByNamespaceIDHandle(ctx, s, nsID, gf.handle)
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
func findComposePageR(ctx context.Context, rr resource.InterfaceSet, ii resource.Identifiers) (pg *types.Page) {
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
