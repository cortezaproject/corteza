package store

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/service/values"
	composeTypes "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
)

var (
	rvSanitizer = values.Sanitizer()
	rvValidator = values.Validator()
	rvFormatter = values.Formatter()
)

func NewComposeRecordFromResource(res *resource.ComposeRecord, cfg *EncoderConfig) resourceState {
	return &composeRecord{
		cfg:         cfg,
		fieldModRef: make(map[string]resource.Identifiers),
		externalRef: make(map[string]map[string]uint64),
		recMap:      make(map[string]*composeTypes.Record),

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
		uu, _, err := store.SearchUsers(ctx, pl.s, systemTypes.UserFilter{
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
	preloadRefs := make(resource.RefSet, 0, int(len(n.relMod.Fields)/2)+1)
	for _, f := range n.relMod.Fields {
		switch f.Kind {
		case "Record":
			refM := f.Options.String("module")
			if refM == "" {
				refM = f.Options.String("moduleID")
			}
			if refM != "" && refM != "0" {
				// Make a reference with that module's records
				ref := n.res.AddRef(composeTypes.RecordResourceType, refM).Constraint(n.res.RefNs)

				n.fieldModRef[f.Name] = ref.Identifiers
				preloadRefs = append(preloadRefs, ref)
			}
		}
	}

	// Can't do anything else, since the NS doesn't yet exist
	if n.relNS.ID == 0 {
		return nil
	}

	// Preload potential references
	//
	// This is a fairly primitive approach, try to think of something a bit nicer
	for _, ref := range preloadRefs {
		mod, err := findComposeModuleStore(ctx, pl.s, n.relNS.ID, makeGenericFilter(ref.Identifiers))
		if err != nil && err != store.ErrNotFound {
			return err
		}
		auxMap := make(map[string]uint64)
		for _, i := range ref.Identifiers {
			n.externalRef[i] = auxMap
		}

		// Preload all records
		rr, _, err := store.SearchComposeRecords(ctx, pl.s, mod, composeTypes.RecordFilter{
			ModuleID:    mod.ID,
			NamespaceID: mod.NamespaceID,
			Paging: filter.Paging{
				Limit: 0,
			},
		})
		if err != nil {
			return err
		}

		for _, r := range rr {
			auxMap[strconv.FormatUint(r.ID, 10)] = r.ID
		}
	}

	// Can't work with own record because the module doesn't yet exist
	if n.relMod.ID == 0 {
		return nil
	}

	// Preload own records
	rr, _, err := store.SearchComposeRecords(ctx, pl.s, n.relMod, composeTypes.RecordFilter{
		ModuleID:    n.relMod.ID,
		NamespaceID: n.relNS.ID,
		Paging: filter.Paging{
			Limit: 0,
		},
	})
	if err != nil && err != store.ErrNotFound {
		return err
	}
	n.missing = len(rr) == 0

	// Map existing records so we can perform updates
	// Map to xref map for easier use later
	auxMap := make(map[string]uint64)
	for _, i := range n.res.RefMod.Identifiers {
		n.externalRef[i] = auxMap
	}
	for _, r := range rr {
		key := strconv.FormatUint(r.ID, 10)
		n.recMap[key] = r

		// Map IDs to xref map
		auxMap[key] = r.ID
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
		if u.Handle != "" {
			ux[u.Handle] = u.ID
		}
		if u.Name != "" {
			ux[u.Name] = u.ID
		}
		if u.Email != "" {
			ux[u.Email] = u.ID
		}
	}
	// Next all of the encoded users.
	// If identifiers overwrite eachother, that's fine.
	for _, ref := range pl.state.ParentResources {
		if ref.ResourceType() == systemTypes.UserResourceType {
			refUsr := ref.(*resource.User)
			for _, i := range refUsr.Identifiers() {
				ux[i] = refUsr.SysID()
			}
		}
	}

	// Some pointing
	rm := n.recMap
	im := n.res.IDMap

	getKey := func(i int, k string) string {
		if k == "" {
			return strconv.FormatInt(int64(i), 10)
		}

		return k
	}

	checkXRef := func(ii resource.Identifiers, ref string) (uint64, error) {
		var auxMap map[string]uint64
		for _, ri := range ii {
			if mp, ok := n.externalRef[ri]; ok {
				auxMap = mp
				break
			}
		}

		if auxMap == nil || len(auxMap) == 0 {
			return 0, fmt.Errorf("referenced record not resolved: %s", resource.ComposeRecordErrUnresolved(resource.MakeIdentifiers(ref)))
		}

		return auxMap[ref], nil
	}

	i := -1
	return n.res.Walker(func(r *resource.ComposeRecordRaw) error {
		i++

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

		rec := &composeTypes.Record{
			ID:          im[getKey(i, r.ID)],
			NamespaceID: nsID,
			ModuleID:    mod.ID,
			CreatedAt:   time.Now(),
		}

		exists := false
		var old *composeTypes.Record
		if r.ID != "" {
			old = rm[r.ID]
			exists = old != nil
		}

		if rec.ID == 0 && exists {
			rec.ID = rm[r.ID].ID
		}
		if rec.ID == 0 {
			rec.ID = NextID()
		}

		im[getKey(i, r.ID)] = rec.ID

		if pl.state.Conflicting {
			return nil
		}

		err := func() error {
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
			rec.CreatedBy = pl.invokerID
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

			rec.OwnedBy = service.CalcRecordOwner(old.OwnedBy, rec.OwnedBy, pl.invokerID)

			rvs := make(composeTypes.RecordValueSet, 0, len(r.Values))
			for k, v := range r.Values {
				rv := &composeTypes.RecordValue{
					RecordID: rec.ID,
					Name:     k,
					Value:    v,
					Updated:  true,
				}

				f := mod.Fields.FindByName(k)
				if f != nil && v != "" {
					switch f.Kind {
					case "User":
						uID := ux[v]
						if uID == 0 {
							return resource.UserErrUnresolved(resource.MakeIdentifiers(v))
						}
						rv.Value = strconv.FormatUint(uID, 10)
						rv.Ref = uID

					case "Record":
						refIdentifiers, ok := n.fieldModRef[f.Name]
						if !ok {
							return fmt.Errorf("module field record reference not resoled: %s", f.Name)
						}

						// if self...
						if n.res.RefMod.Identifiers.HasAny(refIdentifiers) {
							rID := im[v]

							// Check if its in the store
							if rID == 0 {
								// Check if we have an xref
								rID, err = checkXRef(refIdentifiers, v)
								if err != nil {
									return err
								}
							}

							if rID == 0 {
								return resource.ComposeRecordErrUnresolved(resource.MakeIdentifiers(v))
							}
							rv.Value = strconv.FormatUint(rID, 10)
							rv.Ref = rID
						} else {
							// not self...
							rID := uint64(0)
							refRes := resource.FindComposeRecordResource(pl.state.ParentResources, refIdentifiers)

							if refRes != nil {
								// check if parent has it
								rID = refRes.IDMap[v]
							}

							if rID == 0 {
								// Check if we have an xref
								rID, err = checkXRef(refIdentifiers, v)
								if err != nil {
									return err
								}
							}

							if rID == 0 {
								return fmt.Errorf("referenced record not resolved: %s", resource.ComposeRecordErrUnresolved(resource.MakeIdentifiers(v)))
							}

							rv.Value = strconv.FormatUint(rID, 10)
							rv.Ref = rID
						}
					}
				}

				rvs = append(rvs, rv)
			}

			if err = service.RecordValueSanitization(mod, rvs); err != nil {
				return err
			}

			rec.Values = rvs
			rec.Values.SetUpdatedFlag(true)

			// @todo owner?
			var (
				rve *composeTypes.RecordValueErrorSet

				canAccessField = func(f *composeTypes.ModuleField) bool {
					return true
				}
			)

			if old != nil {
				rec.Values = old.Values.Merge(mod.Fields, rec.Values, canAccessField)
			}

			rec.Values, _ = rec.Values.Filter(func(v *composeTypes.RecordValue) (bool, error) {
				return mod.Fields.HasName(v.Name), nil
			})

			rve = service.RecordValueUpdateOpCheck(ctx, nil, mod, rec.Values)
			if !rve.IsValid() {
				return rve
			}

			rve = service.RecordPreparer(ctx, pl.s, rvSanitizer, rvValidator, rvFormatter, mod, rec)
			if !rve.IsValid() {
				return rve
			}

			// Create a new record
			if !exists {
				err = store.CreateComposeRecord(ctx, pl.s, mod, rec)
				return err
			}

			// Update existing
			err = store.UpdateComposeRecord(ctx, pl.s, mod, rec)
			return err
		}()

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
	})
}
