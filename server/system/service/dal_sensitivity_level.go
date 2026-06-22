package service

import (
	"context"
	"fmt"
	"sort"

	"github.com/crusttech/human/server/pkg/actionlog"
	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/dal"

	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

// The CRUD skeleton (FindByID, Search, Create, Update, DeleteByID, UndeleteByID)
// is generated in dal_sensitivity_level.gen.go from system/dal_sensitivity_level.cue.
//
// Every op body is bespoke (store.Tx wrapper, svc.prepare normalization and the
// DAL ReplaceSensitivityLevel / RemoveSensitivityLevel side-effects), so the cue
// marks all ops as customBodyOps and the generated methods delegate to the
// on<Op> handlers below. Access is a single CanManageDalSensitivityLevel check
// (search/create are in customAccessOps so the scaffold does not emit a standard
// access check for them).
//
// This file owns the struct, access-controller interface, constructor, the
// on<Op> handlers and the resource-specific methods (ReloadSensitivityLevels,
// prepare and the package-level reload/replace/remove helpers).

type (
	dalSensitivityLevel struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        sensitivityLevelAccessController
		dal       dalSensitivityLevelManager
	}

	sensitivityLevelAccessController interface {
		CanManageDalSensitivityLevel(context.Context) bool
	}

	dalSensitivityLevelManager interface {
		ReplaceSensitivityLevel(levels ...dal.SensitivityLevel) (err error)
		RemoveSensitivityLevel(levels ...uint64) (err error)
		InUseSensitivityLevel(levelID uint64) (usage dal.SensitivityLevelUsage)
	}
)

func SensitivityLevel(ctx context.Context, dal dalSensitivityLevelManager) *dalSensitivityLevel {
	return &dalSensitivityLevel{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		dal:       dal,
	}

}

func (svc *dalSensitivityLevel) onLookup(ctx context.Context, ID uint64, rProps *dalSensitivityLevelActionProps) (q *types.DalSensitivityLevel, err error) {
	if ID == 0 {
		return nil, DalSensitivityLevelErrInvalidID()
	}

	if q, err = store.LookupDalSensitivityLevelByID(ctx, svc.store, ID); err != nil {
		return nil, DalSensitivityLevelErrInvalidID().Wrap(err)
	}

	rProps.setSensitivityLevel(q)

	if !svc.ac.CanManageDalSensitivityLevel(ctx) {
		return q, DalSensitivityLevelErrNotAllowedToManage(rProps)
	}

	return q, nil
}

func (svc *dalSensitivityLevel) onCreate(ctx context.Context, new *types.DalSensitivityLevel) (err error) {
	return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if new.Meta.Name == "" {
			return DalSensitivityLevelErrMissingName()
		}
		if !svc.ac.CanManageDalSensitivityLevel(ctx) {
			return DalSensitivityLevelErrNotAllowedToManage(&dalSensitivityLevelActionProps{new: new})
		}

		new.CreatedAt = *now()
		new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()
		ups, err := svc.prepare(ctx, s, new)
		if err != nil {
			return
		}
		new.ID = nextID()

		err = store.UpsertDalSensitivityLevel(ctx, s, ups...)
		if err != nil {
			return
		}

		return dalSensitivityLevelReplace(ctx, svc.dal, new)
	})
}

func (svc *dalSensitivityLevel) onUpdate(ctx context.Context, s store.Storer, upd, res *types.DalSensitivityLevel, qProps *dalSensitivityLevelActionProps, _ func() error, _ func() error) (err error) {
	if upd.Meta.Name == "" {
		return DalSensitivityLevelErrMissingName()
	}

	if !svc.ac.CanManageDalSensitivityLevel(ctx) {
		return DalSensitivityLevelErrNotAllowedToManage(qProps)
	}

	upd.UpdatedAt = now()
	upd.CreatedAt = res.CreatedAt
	upd.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

	ups, err := svc.prepare(ctx, s, upd)
	if err != nil {
		return
	}
	if err = store.UpsertDalSensitivityLevel(ctx, s, ups...); err != nil {
		return
	}

	// Reflect updated values back into res so the caller sees the final state.
	*res = *upd

	return dalSensitivityLevelReplace(ctx, svc.dal, upd)
}

func (svc *dalSensitivityLevel) onDelete(ctx context.Context, s store.Storer, res *types.DalSensitivityLevel, qProps *dalSensitivityLevelActionProps) (err error) {
	if !svc.ac.CanManageDalSensitivityLevel(ctx) {
		return DalSensitivityLevelErrNotAllowedToManage(qProps)
	}

	if !svc.dal.InUseSensitivityLevel(res.ID).Empty() {
		return DalSensitivityLevelErrDeleteInUse()
	}

	res.DeletedAt = now()
	res.DeletedBy = a.GetIdentityFromContext(ctx).Identity()

	ups, err := svc.prepare(ctx, s, res)
	if err != nil {
		return
	}
	if err = store.UpsertDalSensitivityLevel(ctx, s, ups...); err != nil {
		return
	}

	var (
		dd = make(types.DalSensitivityLevelSet, 0, len(ups)/2+1)
		uu = make(types.DalSensitivityLevelSet, 0, len(ups)/2+1)
	)

	for _, l := range ups {
		if l.DeletedAt != nil {
			dd = append(dd, l)
		} else {
			uu = append(uu, l)
		}
	}

	if err = dalSensitivityLevelReplace(ctx, svc.dal, uu...); err != nil {
		return err
	}
	if err = dalSensitivityLevelRemove(ctx, svc.dal, dd...); err != nil {
		return err
	}
	return nil
}

func (svc *dalSensitivityLevel) onUndelete(ctx context.Context, s store.Storer, res *types.DalSensitivityLevel, qProps *dalSensitivityLevelActionProps) (err error) {
	if !svc.ac.CanManageDalSensitivityLevel(ctx) {
		return DalSensitivityLevelErrNotAllowedToManage(qProps)
	}

	res.DeletedAt = nil
	res.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

	if err = store.UpdateDalSensitivityLevel(ctx, s, res); err != nil {
		return
	}

	return dalSensitivityLevelReplace(ctx, svc.dal, res)
}

func (svc *dalSensitivityLevel) onSearch(ctx context.Context, filter types.DalSensitivityLevelFilter, aProps *dalSensitivityLevelActionProps) (r types.DalSensitivityLevelSet, f types.DalSensitivityLevelFilter, err error) {
	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.DalSensitivityLevel) (bool, error) {
		if !svc.ac.CanManageDalSensitivityLevel(ctx) {
			return false, nil
		}

		return true, nil
	}

	if !svc.ac.CanManageDalSensitivityLevel(ctx) {
		return nil, f, DalSensitivityLevelErrNotAllowedToManage()
	}

	if r, f, err = store.SearchDalSensitivityLevels(ctx, svc.store, filter); err != nil {
		return nil, f, err
	}

	return r, f, nil
}

func (svc *dalSensitivityLevel) ReloadSensitivityLevels(ctx context.Context, s store.Storer) (err error) {
	return dalSensitivityLevelReload(ctx, svc.store, svc.dal)
}

func (svc *dalSensitivityLevel) prepare(ctx context.Context, s store.Storer, sl *types.DalSensitivityLevel) (_ types.DalSensitivityLevelSet, err error) {
	set, _, err := store.SearchDalSensitivityLevels(ctx, s, types.DalSensitivityLevelFilter{})
	if err != nil {
		return
	}

	updating := sl.ID != 0
	deleting := sl.DeletedAt != nil

	// Validation
	{
		// Assure unique level
		for _, crt := range set {
			if crt.Level == sl.Level && crt.ID != sl.ID {
				return nil, fmt.Errorf("invalid sensitivity level: duplicated level value %d", sl.Level)
			}
		}

		var current *types.DalSensitivityLevel
		for _, crt := range set {
			if crt.ID == sl.ID {
				current = crt
				break
			}
		}

		if (updating || deleting) && current == nil {
			return nil, fmt.Errorf("cannot update sensitivity level %s: does not exist", sl.Handle)
		} else if !updating && current != nil {
			return nil, fmt.Errorf("cannot create sensitivity level %s: already exists", sl.Handle)
		}
	}

	// Preparations
	{
		// Make sure to properly update
		for i, s := range set {
			if s.ID == sl.ID {
				set[i] = sl
				break
			}
		}

		// Make sure it's in there
		if !deleting && !updating {
			set = append(set, sl)
		}

		// Sort by level for easier normalization
		sort.Sort(set)

		// @todo uncomment sensitivity level normalization after we redo the user interface
		// // Normalize sensitivity level
		// offset := 0
		// for i := range set {
		// 	if set[i].DeletedAt != nil {
		// 		offset++
		// 	}

		// 	nxtLvl := i + 1 - offset
		// 	if nxtLvl != set[i].Level {
		// 		set[i].UpdatedAt = now()
		// 		// Same user so we can cheat a bit
		// 		set[i].UpdatedBy = sl.CreatedBy
		// 	}

		// 	set[i].Level = nxtLvl
		// }
	}

	return set, err
}

func dalSensitivityLevelReload(ctx context.Context, s store.Storer, dsm dalSensitivityLevelManager) (err error) {
	// Get all available sensitivityLevels
	ll, _, err := store.SearchDalSensitivityLevels(ctx, s, types.DalSensitivityLevelFilter{})
	if err != nil {
		return
	}

	return dalSensitivityLevelReplace(ctx, dsm, ll...)
}

func dalSensitivityLevelReplace(ctx context.Context, dsm dalSensitivityLevelManager, ll ...*types.DalSensitivityLevel) (err error) {
	levels := make(dal.SensitivityLevelSet, len(ll))
	for i, l := range ll {
		levels[i] = dal.MakeSensitivityLevel(l.ID, l.Level, l.Handle)
	}

	return dsm.ReplaceSensitivityLevel(levels...)
}

func dalSensitivityLevelRemove(ctx context.Context, dsm dalSensitivityLevelManager, ll ...*types.DalSensitivityLevel) (err error) {
	dd := make([]uint64, len(ll))
	for i, l := range ll {
		dd[i] = l.ID
	}

	return dsm.RemoveSensitivityLevel(dd...)
}
