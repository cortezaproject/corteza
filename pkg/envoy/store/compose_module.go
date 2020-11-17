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
	composeModuleState struct {
		cfg *EncoderConfig

		res *resource.ComposeModule
		mod *types.Module

		relNS     *types.Namespace
		recFields map[string]uint64
	}
)

func NewComposeModuleState(res *resource.ComposeModule, cfg *EncoderConfig) resourceState {
	return &composeModuleState{
		cfg: cfg,

		res: res,
	}
}

func (n *composeModuleState) Prepare(ctx context.Context, s store.Storer, state *envoy.ResourceState) (err error) {
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

	// Get related record field modules
	for _, r := range n.res.ModRef {
		mod, err := findComposeModuleRS(ctx, s, n.relNS.ID, state.ParentResources, r.Identifiers)
		if err != nil {
			return err
		} else if mod == nil {
			return errors.New("[mod prepare] couldn't find related module; @todo error")
		}

		for i := range r.Identifiers {
			n.recFields[i] = mod.ID
		}
	}

	// Try to get the original module
	n.mod, err = findComposeModuleS(ctx, s, n.relNS.ID, makeGenericFilter(n.res.Identifiers()))
	if err != nil {
		return err
	}

	// Nothing else to do
	if n.mod == nil {
		return nil
	}

	// Get the original module fields
	// These are used later for some merging logic
	n.mod.Fields, err = findComposeModuleFieldsS(ctx, s, n.mod)
	if err != nil {
		return err
	}

	if n.mod != nil {
		n.res.Res.ID = n.mod.ID
	}
	return nil
}

func (n *composeModuleState) Encode(ctx context.Context, s store.Storer, state *envoy.ResourceState) (err error) {
	res := n.res.Res
	exists := n.mod != nil && n.mod.ID > 0

	// Determine the ID
	if res.ID <= 0 && exists {
		res.ID = n.mod.ID
	}
	if res.ID <= 0 {
		res.ID = nextID()
	}

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
		return errors.New("[module] couldn't find related namespace; @todo error")
	}

	// Fields
	for i, f := range res.Fields {
		f.ID = res.ID
		f.ModuleID = res.ID
		f.Place = i
		f.DeletedAt = nil
		f.CreatedAt = time.Now()

		if f.Kind == "Record" {
			refM := f.Options.String("module")
			mID := n.recFields[refM]
			if mID <= 0 {
				mod := findComposeModuleR(state.ParentResources, resource.MakeIdentifiers(refM))
				if mod == nil || mod.ID <= 0 {
					return errors.New("[module field] couldn't find related module; @todo error")
				}
				mID = mod.ID
			}

			f.Options["module"] = mID
		}
	}

	// Create a fresh module
	if !exists {
		err = store.CreateComposeModule(ctx, s, res)
		if err != nil {
			return err
		}

		err = store.CreateComposeModuleField(ctx, s, res.Fields...)
		if err != nil {
			return err
		}

		return nil
	}

	// Update existing module
	switch n.cfg.OnExisting {
	case Skip:
		return nil

	case MergeLeft:
		res = mergeComposeModule(n.mod, res)
		res.Fields = mergeComposeModuleFields(n.mod.Fields, res.Fields)

	case MergeRight:
		res = mergeComposeModule(res, n.mod)
		res.Fields = mergeComposeModuleFields(res.Fields, n.mod.Fields)
	}

	err = store.UpdateComposeModule(ctx, s, res)
	if err != nil {
		return err
	}

	err = store.DeleteComposeModuleField(ctx, s, n.mod.Fields...)
	if err != nil {
		return err
	}
	err = store.CreateComposeModuleField(ctx, s, res.Fields...)
	if err != nil {
		return err
	}

	n.res.Res = res
	return nil
}

// mergeComposeModuleFields merges b into a, prioritising a
func mergeComposeModuleFields(a, b types.ModuleFieldSet) types.ModuleFieldSet {
	ff := a.Clone()
	missing := make(types.ModuleFieldSet, 0, len(b))

	for _, fb := range b {
		for _, fa := range ff {
			if fb.ID != fa.ID && fb.Name != fa.Name {
				continue
			}

			if fa.Kind == "" {
				fa.Kind = fb.Kind
			}
			if fa.Name == "" {
				fa.Name = fb.Name
			}
			if fa.Label == "" {
				fa.Label = fb.Label
			}
			if fa.Options == nil {
				fa.Options = fb.Options
			}
			if fa.DefaultValue == nil || len(fa.DefaultValue) <= 0 {
				fa.DefaultValue = fb.DefaultValue
			}

			goto out
		}
		missing = append(missing, fb)
	out:
	}

	ff = append(ff, missing...)

	return ff
}

// mergeComposeModule merges b into a, prioritising a
func mergeComposeModule(a, b *types.Module) *types.Module {
	c := a.Clone()

	if c.Handle == "" {
		c.Handle = b.Handle
	}
	if c.Name == "" {
		c.Name = b.Name
	}

	// I'll just compare the entire struct for now
	if c.Meta == nil {
		c.Meta = b.Meta
	}

	return c
}

// findComposeModuleRS looks for the chart in the resources & the store
//
// Provided resources are prioritized.
func findComposeModuleRS(ctx context.Context, s store.Storer, nsID uint64, rr resource.InterfaceSet, ii resource.Identifiers) (mod *types.Module, err error) {
	mod = findComposeModuleR(rr, ii)
	if mod != nil {
		return mod, nil
	}

	// Go in the store
	return findComposeModuleS(ctx, s, nsID, makeGenericFilter(ii))
}

// findComposeModuleS looks for the module in the store
func findComposeModuleS(ctx context.Context, s store.Storer, nsID uint64, gf genericFilter) (mod *types.Module, err error) {
	if gf.id > 0 {
		mod, err = store.LookupComposeModuleByID(ctx, s, gf.id)
		if err != nil && err != store.ErrNotFound {
			return nil, err
		}

		if mod != nil {
			return
		}
	}

	if gf.handle != "" {
		mod, err = store.LookupComposeModuleByNamespaceIDHandle(ctx, s, nsID, gf.handle)
		if err != nil && err != store.ErrNotFound {
			return nil, err
		}

		if mod != nil {
			return
		}
	}

	return nil, nil
}

// findComposeModuleR looks for the module in the store
func findComposeModuleR(rr resource.InterfaceSet, ii resource.Identifiers) (ns *types.Module) {
	var modRes *resource.ComposeModule
	var ok bool

	rr.Walk(func(r resource.Interface) error {
		modRes, ok = r.(*resource.ComposeModule)
		if !ok {
			return nil
		}

		if !modRes.Identifiers().HasAny(r.Identifiers()) {
			modRes = nil
		}
		return nil
	})

	// Found it
	if modRes != nil {
		return modRes.Res
	}
	return nil
}

// findComposeModuleFieldsS looks for the module fields in the store
func findComposeModuleFieldsS(ctx context.Context, s store.Storer, mod *types.Module) (types.ModuleFieldSet, error) {
	if mod.ID <= 0 {
		return mod.Fields, nil
	}

	ff, _, err := store.SearchComposeModuleFields(ctx, s, types.ModuleFieldFilter{
		ModuleID: []uint64{mod.ID},
	})
	if err != nil {
		return nil, err
	}

	return ff, nil
}
