package store

import (
	"context"
	"strconv"
	"time"

	"github.com/cortezaproject/corteza-server/compose/service/values"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	composeRecordState struct {
		cfg *EncoderConfig

		res *resource.ComposeRecord

		relNS  *types.Namespace
		relMod *types.Module
		relUsr map[string]uint64

		// Little helper flag for conditional encoding
		missing bool
	}
)

var (
	rvSanitizer = values.Sanitizer()
	rvValidator = values.Validator()
)

func NewComposeRecordState(res *resource.ComposeRecord, cfg *EncoderConfig) resourceState {
	return &composeRecordState{
		cfg: cfg,

		res: res,
	}
}

func (n *composeRecordState) Prepare(ctx context.Context, s store.Storer, state *envoy.ResourceState) (err error) {
	// Get related namespace
	n.relNS, err = findComposeNamespaceRS(ctx, s, state.ParentResources, n.res.NsRef.Identifiers)
	if err != nil {
		return err
	}
	if n.relNS == nil {
		return composeNamespaceErrUnresolved(n.res.NsRef.Identifiers)
	}

	n.missing = true
	n.relMod = findComposeModuleR(state.ParentResources, n.res.ModRef.Identifiers)
	if n.relMod == nil && n.relNS.ID > 0 {
		n.relMod, err = findComposeModuleS(ctx, s, n.relNS.ID, makeGenericFilter(n.res.ModRef.Identifiers))
		if err != nil {
			return err
		}
		if n.relMod != nil {
			// Preload existing fields
			n.relMod.Fields, err = findComposeModuleFieldsS(ctx, s, n.relMod)
			if err != nil {
				return err
			}
		}
	}

	if n.relMod == nil {
		return composeModuleErrUnresolved(n.res.ModRef.Identifiers)
	}

	// Sys users
	n.relUsr = make(map[string]uint64)
	if err = resolveUserRefs(ctx, s, state.ParentResources, n.res.UserRefs(), n.relUsr); err != nil {
		return err
	}

	// Add missing refs
	for _, f := range n.relMod.Fields {
		switch f.Kind {
		case "Record":
			refM := f.Options.String("module")
			if refM != "" && refM != "0" {
				// Make a reference with that module's records
				n.res.AddRef(resource.COMPOSE_RECORD_RESOURCE_TYPE, refM).Constraint(n.res.NsRef)
			}
		}
	}

	// Can't do anything else, since the NS doesn't yet exist
	if n.relNS.ID <= 0 {
		return nil
	}

	// Check if empty
	rr, _, err := store.SearchComposeRecords(ctx, s, n.relMod, types.RecordFilter{
		ModuleID:    n.relMod.ID,
		NamespaceID: n.relNS.ID,
		Paging:      filter.Paging{Limit: 1},
	})
	if err != nil && err != store.ErrNotFound {
		return err
	}
	n.missing = len(rr) == 0

	// Try to get existing records
	//
	// @todo handle large amounts of
	for rID := range n.res.IDMap {
		var r *types.Record
		// @todo support for labels
		if refy.MatchString(rID) {
			id, _ := strconv.ParseUint(rID, 10, 64)
			r, err = store.LookupComposeRecordByID(ctx, s, n.relMod, id)
			if err == store.ErrNotFound {
				continue
			} else if err != nil {
				return err
			}
			if r != nil {
				n.res.RecMap[rID] = r
			}
		} else {
			continue
		}

		n.res.RecMap[rID] = r
	}

	return nil
}

func (n *composeRecordState) Encode(ctx context.Context, s store.Storer, state *envoy.ResourceState) (err error) {
	// Namespace
	nsID := n.relNS.ID
	if nsID <= 0 {
		ns := findComposeNamespaceR(state.ParentResources, n.res.NsRef.Identifiers)
		nsID = ns.ID
	}

	// Module
	mod := n.relMod
	if mod.ID <= 0 {
		m := findComposeModuleR(state.ParentResources, n.res.ModRef.Identifiers)
		mod.ID = m.ID
	}

	// Sys users
	for idf, ID := range n.relUsr {
		if ID != 0 {
			continue
		}
		u := findUserR(ctx, state.ParentResources, resource.MakeIdentifiers(idf))
		n.relUsr[idf] = u.ID
	}

	// Some pointing
	rm := n.res.RecMap
	im := n.res.IDMap

	return n.res.Walker(func(r *resource.ComposeRecordRaw) error {
		// So we don't have to worry about nil
		cfg := r.Config
		if cfg == nil {
			cfg = &resource.EnvoyConfig{}
		}

		if cfg.SkipIf != "" {
			evl, err := exprP.NewEvaluable(cfg.SkipIf)
			if err != nil {
				return err
			}
			// @todo expand this
			skip, err := evl.EvalBool(ctx, map[string]interface{}{
				"missing": n.missing,
			})
			if err != nil {
				return err
			}

			if skip {
				return nil
			}
		}

		// Simple wrapper to do some post-processing steps
		dfr := func(err error) error {
			if n.cfg.Defer != nil {
				n.cfg.Defer()
			}

			if err != nil {
				if n.cfg.DeferNok != nil {
					return n.cfg.DeferNok(err)
				}
				return err
			} else if n.cfg.DeferOk != nil {
				n.cfg.DeferOk()
			}

			return nil
		}

		rec := &types.Record{
			ID:          im[r.ID],
			NamespaceID: nsID,
			ModuleID:    mod.ID,
			CreatedAt:   time.Now(),
		}

		exists := false
		if r.ID != "" {
			exists = rm[r.ID] != nil
		}

		if rec.ID <= 0 && exists {
			rec.ID = rm[r.ID].ID
		} else {
			rec.ID = NextID()
		}

		im[r.ID] = rec.ID

		if state.Conflicting {
			return nil
		}

		// Timestamps
		if r.Ts != nil {
			if r.Ts.CreatedAt != "" {
				t := toTime(r.Ts.CreatedAt)
				if t != nil {
					rec.CreatedAt = *t
				}
			}
			if r.Ts.UpdatedAt != "" {
				rec.UpdatedAt = toTime(r.Ts.UpdatedAt)
			}
			if r.Ts.DeletedAt != "" {
				rec.DeletedAt = toTime(r.Ts.DeletedAt)
			}
		}
		// Userstamps
		if r.Us != nil {
			if r.Us.CreatedBy != "" {
				rec.CreatedBy = n.relUsr[r.Us.CreatedBy]
			}
			if r.Us.UpdatedBy != "" {
				rec.UpdatedBy = n.relUsr[r.Us.UpdatedBy]
			}
			if r.Us.DeletedBy != "" {
				rec.DeletedBy = n.relUsr[r.Us.DeletedBy]
			}
			if r.Us.OwnedBy != "" {
				rec.OwnedBy = n.relUsr[r.Us.OwnedBy]
			}
		}

		rvs := make(types.RecordValueSet, 0, len(r.Values))
		for k, v := range r.Values {
			rv := &types.RecordValue{
				RecordID: rec.ID,
				Name:     k,
				Value:    v,
				Updated:  true,
			}

			rvs = append(rvs, rv)
		}

		rec.Values = rvSanitizer.Run(mod, rvs)
		rve := rvValidator.Run(ctx, s, mod, rec)
		if !rve.IsValid() {
			return dfr(rve)
		}

		// Create a new record
		if !exists {
			err = store.CreateComposeRecord(ctx, s, mod, rec)
			return dfr(err)
		}

		// Update existing
		err = store.UpdateComposeRecord(ctx, s, mod, rec)
		return dfr(err)
	})
}
