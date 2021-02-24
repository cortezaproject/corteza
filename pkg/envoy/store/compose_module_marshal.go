package store

import (
	"context"
	"strconv"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
)

func NewComposeModuleFromResource(res *resource.ComposeModule, cfg *EncoderConfig) resourceState {
	return &composeModule{
		cfg: mergeConfig(cfg, res.Config()),

		res: res,

		recFields: make(map[string]uint64),
	}
}

// Prepare prepares the composeModule to be encoded
//
// Any validation, additional constraining should be performed here.
func (n *composeModule) Prepare(ctx context.Context, pl *payload) (err error) {
	// Get related namespace
	n.relNS, err = findComposeNamespaceRS(ctx, pl.s, pl.state.ParentResources, n.res.RefNs.Identifiers)
	if err != nil {
		return err
	}
	if n.relNS == nil {
		return resource.ComposeNamespaceErrUnresolved(n.res.RefNs.Identifiers)
	}

	// Get related record field modules
	for _, refMod := range n.res.RefMods {
		var mod *types.Module
		if n.relNS.ID > 0 {
			mod, err = findComposeModuleS(ctx, pl.s, n.relNS.ID, makeGenericFilter(refMod.Identifiers))
			if err != nil {
				return err
			}
		}
		if mod == nil {
			mod = resource.FindComposeModule(pl.state.ParentResources, refMod.Identifiers)
		}
		if mod == nil {
			return composeModuleErrUnresolvedRecordField(refMod.Identifiers)
		}

		for i := range refMod.Identifiers {
			n.recFields[i] = mod.ID
		}
	}

	// Can't do anything else, since the NS doesn't yet exist
	if n.relNS.ID <= 0 {
		return nil
	}

	// Try to get the original module
	n.mod, err = findComposeModuleS(ctx, pl.s, n.relNS.ID, makeGenericFilter(n.res.Identifiers()))
	if err != nil {
		return err
	}

	// Nothing else to do
	if n.mod == nil {
		return nil
	}

	// Get the original module fields
	// These are used later for some merging logic
	n.mod.Fields, err = findComposeModuleFieldsS(ctx, pl.s, n.mod)
	if err != nil {
		return err
	}

	if n.mod != nil {
		n.res.Res.ID = n.mod.ID
		n.res.Res.NamespaceID = n.mod.NamespaceID
	}
	return nil
}

// Encode encodes the composeModule to the store
//
// Encode is allowed to do some data manipulation, but no resource constraints
// should be changed.
func (n *composeModule) Encode(ctx context.Context, pl *payload) (err error) {
	res := n.res.Res
	exists := n.mod != nil && n.mod.ID > 0

	// Determine the ID
	if res.ID <= 0 && exists {
		res.ID = n.mod.ID
	}
	if res.ID <= 0 {
		res.ID = NextID()
	}

	if pl.state.Conflicting {
		return nil
	}
	// Timestamps
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

	// Fields
	var originalFields types.ModuleFieldSet
	if n.mod != nil && n.mod.Fields != nil {
		originalFields = n.mod.Fields
	} else {
		originalFields = make(types.ModuleFieldSet, 0)
	}
	for i, f := range res.Fields {
		of := originalFields.FindByName(f.Name)
		if of != nil {
			f.ID = of.ID
		} else {
			f.ID = NextID()
		}
		f.ModuleID = res.ID
		f.Place = i
		f.DeletedAt = nil
		f.CreatedAt = *now()

		if f.Kind == "Record" {
			refMod := f.Options.String("module")
			if refMod == "" {
				refMod = f.Options.String("moduleID")
			}
			modID := n.recFields[refMod]
			if modID <= 0 {
				ii := resource.MakeIdentifiers(refMod)
				mod := resource.FindComposeModule(pl.state.ParentResources, ii)
				if mod == nil || mod.ID <= 0 {
					return composeModuleErrUnresolvedRecordField(ii)
				}
				modID = mod.ID
			}

			f.Options["moduleID"] = strconv.FormatUint(modID, 10)
			delete(f.Options, "module")
		}
	}

	// Evaluate the resource skip expression
	// @todo expand available parameters; similar implementation to compose/types/record@Dict
	if skip, err := basicSkipEval(ctx, n.cfg, !exists); err != nil {
		return err
	} else if skip {
		return nil
	}

	// Create a fresh module
	if !exists {
		err = store.CreateComposeModule(ctx, pl.s, res)
		if err != nil {
			return err
		}

		err = store.CreateComposeModuleField(ctx, pl.s, res.Fields...)
		if err != nil {
			return err
		}

		return nil
	}

	// Update existing module
	switch n.cfg.OnExisting {
	case resource.Skip:
		return nil

	case resource.MergeLeft:
		res = mergeComposeModule(n.mod, res)
		res.Fields = mergeComposeModuleFields(n.mod.Fields, res.Fields)

	case resource.MergeRight:
		res = mergeComposeModule(res, n.mod)
		res.Fields = mergeComposeModuleFields(res.Fields, n.mod.Fields)
	}

	err = store.UpdateComposeModule(ctx, pl.s, res)
	if err != nil {
		return err
	}

	err = store.DeleteComposeModuleField(ctx, pl.s, n.mod.Fields...)
	if err != nil {
		return err
	}
	err = store.CreateComposeModuleField(ctx, pl.s, res.Fields...)
	if err != nil {
		return err
	}

	n.res.Res = res
	return nil
}
