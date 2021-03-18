package store

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/service/values"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
	stypes "github.com/cortezaproject/corteza-server/system/types"
)

var (
	rvSanitizer = values.Sanitizer()
	rvValidator = values.Validator()
	rvFormatter = values.Formatter()
)

func NewComposeRecordFromResource(res *resource.ComposeRecord, cfg *EncoderConfig) resourceState {
	return &composeRecord{
		cfg: cfg,

		res: res,
	}
}

func (n *composeRecord) Prepare(ctx context.Context, pl *payload) (err error) {
	// Get related namespace
	n.relNS, err = findComposeNamespace(ctx, pl.s, pl.state.ParentResources, n.res.RefNs.Identifiers)
	if err != nil {
		return err
	}
	if n.relNS == nil {
		return resource.ComposeNamespaceErrUnresolved(n.res.RefNs.Identifiers)
	}

	n.missing = true
	n.relMod = resource.FindComposeModule(pl.state.ParentResources, n.res.RefMod.Identifiers)
	if n.relMod == nil && n.relNS.ID > 0 {
		n.relMod, err = findComposeModuleStore(ctx, pl.s, n.relNS.ID, makeGenericFilter(n.res.RefMod.Identifiers))
		if err != nil {
			return err
		}
		if n.relMod != nil {
			// Preload existing fields
			n.relMod.Fields, err = findComposeModuleFieldsStore(ctx, pl.s, n.relMod)
			if err != nil {
				return err
			}
		}
	}
	if n.relMod == nil {
		return resource.ComposeModuleErrUnresolved(n.res.RefMod.Identifiers)
	}

	// Preload sys users
	if n.res.UserFlakes == nil {
		n.res.UserFlakes = make(resource.UserstampIndex)
	}
	if len(n.res.UserFlakes) == 0 {
		// No users provided, let's try to fetch them
		uu, _, err := store.SearchUsers(ctx, pl.s, stypes.UserFilter{
			Paging: filter.Paging{
				Limit: 0,
			},
		})
		if err != nil {
			return err
		}
		n.res.UserFlakes.Add(uu...)
	}

	// Add missing refs
	for _, f := range n.relMod.Fields {
		switch f.Kind {
		case "Record":
			refM := f.Options.String("module")
			if refM == "" {
				refM = f.Options.String("moduleID")
			}
			if refM != "" && refM != "0" {
				// Make a reference with that module's records
				n.res.AddRef(resource.COMPOSE_RECORD_RESOURCE_TYPE, refM).Constraint(n.res.RefNs)
			}
		}
	}

	// Can't do anything else, since the NS doesn't yet exist
	if n.relNS.ID <= 0 {
		return nil
	}

	// Check if empty
	rr, _, err := store.SearchComposeRecords(ctx, pl.s, n.relMod, types.RecordFilter{
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
			r, err = store.LookupComposeRecordByID(ctx, pl.s, n.relMod, id)
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

// @note composeRecord.Encode can raise an error in case of unresolved user dependencies.
func (n *composeRecord) Encode(ctx context.Context, pl *payload) (err error) {
	// Namespace
	nsID := n.relNS.ID
	if nsID <= 0 {
		ns := resource.FindComposeNamespace(pl.state.ParentResources, n.res.RefNs.Identifiers)
		nsID = ns.ID
	}

	// Module
	mod := n.relMod
	if mod.ID <= 0 {
		m := resource.FindComposeModule(pl.state.ParentResources, n.res.RefMod.Identifiers)
		mod.ID = m.ID
	}

	// Aggregate all of the available users
	// ux will map { identifier: userID }
	ux := make(map[string]uint64)
	// Firstly the users provided by the record pl.state itself
	for _, ur := range n.res.UserFlakes {
		u := ur.U

		ux[strconv.FormatUint(u.ID, 10)] = u.ID
		ux[u.Handle] = u.ID
		ux[u.Name] = u.ID
		ux[u.Email] = u.ID
	}
	// Next all of the encoded users.
	// If identifiers overwrite eachother, that's fine.
	for _, ref := range pl.state.ParentResources {
		if ref.ResourceType() == resource.USER_RESOURCE_TYPE {
			refUsr := ref.(*resource.User)
			for i := range refUsr.Identifiers() {
				ux[i] = refUsr.SysID()
			}
		}
	}

	// Some pointing
	rm := n.res.RecMap
	im := n.res.IDMap

	createAcChecked := false
	updateAcChecked := false

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
		var old *types.Record
		if r.ID != "" {
			old = rm[r.ID]
			exists = old != nil
		}

		if rec.ID <= 0 && exists {
			rec.ID = rm[r.ID].ID
		} else {
			rec.ID = NextID()
		}

		im[r.ID] = rec.ID

		if pl.state.Conflicting {
			return nil
		}

		// Timestamps
		if r.Ts != nil {
			if r.Ts.CreatedAt != nil {
				rec.CreatedAt = *r.Ts.CreatedAt.T
			} else {
				rec.CreatedAt = *now()
			}
			if r.Ts.UpdatedAt != nil {
				rec.UpdatedAt = r.Ts.UpdatedAt.T
			}
			if r.Ts.DeletedAt != nil {
				rec.DeletedAt = r.Ts.DeletedAt.T
			}
		}

		// Userstamps
		if r.Us != nil {
			if r.Us.CreatedBy != nil {
				rec.CreatedBy = ux[r.Us.CreatedBy.Ref]
			}
			if r.Us.UpdatedBy != nil {
				rec.UpdatedBy = ux[r.Us.UpdatedBy.Ref]
			}
			if r.Us.DeletedBy != nil {
				rec.DeletedBy = ux[r.Us.DeletedBy.Ref]
			}
			if r.Us.OwnedBy != nil {
				rec.OwnedBy = ux[r.Us.OwnedBy.Ref]
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

			f := mod.Fields.FindByName(k)
			if f != nil && f.Kind == "User" {
				uID := ux[v]
				if uID == 0 {
					return resource.UserErrUnresolved(resource.MakeIdentifiers(v))
				}
				rv.Value = strconv.FormatUint(uID, 10)
				rv.Ref = uID
			}

			rvs = append(rvs, rv)
		}

		if err = service.RecordValueSanitazion(mod, rvs); err != nil {
			return err
		}

		rec.Values = rvs
		rec.Values.SetUpdatedFlag(true)

		// @todo owner?
		var rve *types.RecordValueErrorSet
		if old != nil {
			rec.Values, rve = service.RecordValueMerger(ctx, pl.composeAccessControl, mod, rec.Values, old.Values)
		} else {
			rec.Values, rve = service.RecordValueMerger(ctx, pl.composeAccessControl, mod, rec.Values, nil)
		}
		if !rve.IsValid() {
			return rve
		}

		rve = service.RecordPreparer(ctx, pl.s, rvSanitizer, rvValidator, rvFormatter, mod, rec)
		if !rve.IsValid() {
			return dfr(rve)
		}

		// AC
		//
		// AC needs to happen down here, because we are either creating or updating
		// records and we don't know that for sure in the Prepare method.
		//
		// @todo expand this when we allow record based AC
		if !exists && !createAcChecked {
			createAcChecked = true
			if !pl.composeAccessControl.CanCreateRecord(ctx, mod) {
				return fmt.Errorf("not allowed to create records for module %d", mod.ID)
			}
		} else if exists && !updateAcChecked {
			updateAcChecked = true
			if !pl.composeAccessControl.CanUpdateRecord(ctx, mod) {
				return fmt.Errorf("not allowed to update records for module %d", mod.ID)
			}
		}

		// Create a new record
		if !exists {
			err = store.CreateComposeRecord(ctx, pl.s, mod, rec)
			return dfr(err)
		}

		// Update existing
		err = store.UpdateComposeRecord(ctx, pl.s, mod, rec)
		return dfr(err)
	})
}
