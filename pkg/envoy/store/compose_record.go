package store

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/cortezaproject/corteza-server/compose/service/values"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	composeRecordState struct {
		cfg *EncoderConfig

		res *resource.ComposeRecord

		relNS  *types.Namespace
		relMod *types.Module
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
		return errors.New("@todo couldn't resolve namespace")
	}

	n.relMod = findComposeModuleR(state.ParentResources, n.res.ModRef.Identifiers)
	if n.relMod == nil && n.relNS.ID > 0 {
		n.relMod, err = findComposeModuleS(ctx, s, n.relNS.ID, makeGenericFilter(n.res.ModRef.Identifiers))
		if err != nil {
			return err
		}
		if n.relMod == nil {
			return errors.New("@todo couldn't resolve module")
		}

		// Preload existing fields
		n.relMod.Fields, err = findComposeModuleFieldsS(ctx, s, n.relMod)
		if err != nil {
			return err
		}
	}

	// Add missing refs
	for _, f := range n.relMod.Fields {
		switch f.Kind {
		case "Record":
			refM := f.Options.String("module")
			if refM != "" && refM != "0" {
				// Make a reference with that module's records
				n.res.AddRef(resource.COMPOSE_RECORD_RESOURCE_TYPE, refM)
			}
		}
	}

	// Can't do anything else, since the NS doesn't yet exist
	if n.relNS.ID <= 0 {
		return nil
	}

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

	// Some pointing
	rm := n.res.RecMap
	im := n.res.IDMap

	return n.res.Walker(func(r *resource.ComposeRecordRaw) error {
		rec := &types.Record{
			ID:          im[r.ID],
			NamespaceID: nsID,
			ModuleID:    mod.ID,
			CreatedAt:   time.Now(),
		}
		exists := rm[r.ID] != nil

		if rec.ID <= 0 && exists {
			rec.ID = rm[r.ID].ID
		} else {
			rec.ID = nextID()
		}

		im[r.ID] = rec.ID

		if state.Conflicting {
			return nil
		}

		// Sys values
		//
		// @todo
		for k, v := range r.SysValues {
			if v == "" {
				continue
			}

			switch k {
			case "createdAt":
				// @todo set time
				rec.CreatedAt = time.Now()

			case "updatedAt":
				// @todo set time
				rec.UpdatedAt = nil

			case "deletedAt":
				// @todo set time
				rec.DeletedAt = nil
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

		// @todo validator
		rec.Values = rvSanitizer.Run(mod, rvs)

		// Create a new record
		if !exists {
			err = store.CreateComposeRecord(ctx, s, mod, rec)
			if err != nil {
				return err
			}
			return nil
		}

		// Update existing
		err = store.UpdateComposeRecord(ctx, s, mod, rec)
		if err != nil {
			return err
		}

		return nil
	})
}
