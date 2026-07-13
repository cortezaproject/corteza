package service

import (
	"context"
	"fmt"
	"sort"

	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/dal"

	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

type (
	dalSensitivityLevelServices struct {
		dal dalSensitivityLevelManager
	}

	dalSensitivityLevelAccessController interface {
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
		services:  &dalSensitivityLevelServices{dal: dal},
	}
}

func (svc *dalSensitivityLevel) onLookup(ctx context.Context, ID uint64, aProps *dalSensitivityLevelActionProps) (*types.DalSensitivityLevel, error) {
	if ID == 0 {
		return nil, DalSensitivityLevelErrInvalidID()
	}

	res, err := store.LookupDalSensitivityLevelByID(ctx, svc.store, ID)
	if err != nil {
		return nil, DalSensitivityLevelErrInvalidID().Wrap(err)
	}

	aProps.setSensitivityLevel(res)

	if !svc.ac.CanManageDalSensitivityLevel(ctx) {
		return nil, DalSensitivityLevelErrNotAllowedToManage(aProps)
	}

	return res, nil
}

func (svc *dalSensitivityLevel) onSearch(ctx context.Context, filter types.DalSensitivityLevelFilter, aProps *dalSensitivityLevelActionProps) (types.DalSensitivityLevelSet, types.DalSensitivityLevelFilter, error) {
	filter.Check = func(res *types.DalSensitivityLevel) (bool, error) {
		if !svc.ac.CanManageDalSensitivityLevel(ctx) {
			return false, nil
		}
		return true, nil
	}

	if !svc.ac.CanManageDalSensitivityLevel(ctx) {
		return nil, filter, DalSensitivityLevelErrNotAllowedToManage()
	}

	set, f, err := store.SearchDalSensitivityLevels(ctx, svc.store, filter)
	if err != nil {
		return nil, f, err
	}

	return set, f, nil
}

func (svc *dalSensitivityLevel) onCreate(ctx context.Context, new *types.DalSensitivityLevel) error {
	if new.Meta.Name == "" {
		return DalSensitivityLevelErrMissingName()
	}

	if !svc.ac.CanManageDalSensitivityLevel(ctx) {
		return DalSensitivityLevelErrNotAllowedToManage()
	}

	new.CreatedAt = *now()
	new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()

	return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
		ups, err := svc.prepare(ctx, s, new)
		if err != nil {
			return err
		}

		new.ID = nextID()

		if err = store.UpsertDalSensitivityLevel(ctx, s, ups...); err != nil {
			return err
		}

		return dalSensitivityLevelReplace(ctx, svc.services.dal, new)
	})
}

func (svc *dalSensitivityLevel) onUpdate(ctx context.Context, s store.Storer, upd *types.DalSensitivityLevel, res *types.DalSensitivityLevel, aProps *dalSensitivityLevelActionProps, before func() error, after func() error) error {
	if upd.Meta.Name == "" {
		return DalSensitivityLevelErrMissingName()
	}

	if !svc.ac.CanManageDalSensitivityLevel(ctx) {
		return DalSensitivityLevelErrNotAllowedToManage(aProps)
	}

	upd.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

	ups, err := svc.prepare(ctx, s, upd)
	if err != nil {
		return err
	}

	if err = store.UpsertDalSensitivityLevel(ctx, s, ups...); err != nil {
		return err
	}

	return dalSensitivityLevelReplace(ctx, svc.services.dal, upd)
}

func (svc *dalSensitivityLevel) onDelete(ctx context.Context, s store.Storer, res *types.DalSensitivityLevel, aProps *dalSensitivityLevelActionProps) error {
	if !svc.ac.CanManageDalSensitivityLevel(ctx) {
		return DalSensitivityLevelErrNotAllowedToManage(aProps)
	}

	if !svc.services.dal.InUseSensitivityLevel(res.ID).Empty() {
		return DalSensitivityLevelErrDeleteInUse()
	}

	res.DeletedAt = now()
	res.DeletedBy = a.GetIdentityFromContext(ctx).Identity()

	ups, err := svc.prepare(ctx, s, res)
	if err != nil {
		return err
	}

	if err = store.UpsertDalSensitivityLevel(ctx, s, ups...); err != nil {
		return err
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

	if err = dalSensitivityLevelReplace(ctx, svc.services.dal, uu...); err != nil {
		return err
	}

	return dalSensitivityLevelRemove(ctx, svc.services.dal, dd...)
}

func (svc *dalSensitivityLevel) onUndelete(ctx context.Context, s store.Storer, res *types.DalSensitivityLevel, aProps *dalSensitivityLevelActionProps) error {
	if !svc.ac.CanManageDalSensitivityLevel(ctx) {
		return DalSensitivityLevelErrNotAllowedToManage(aProps)
	}

	res.DeletedAt = nil
	res.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

	if err := store.UpdateDalSensitivityLevel(ctx, s, res); err != nil {
		return err
	}

	return dalSensitivityLevelReplace(ctx, svc.services.dal, res)
}

func (svc *dalSensitivityLevel) onReloadSensitivityLevels(ctx context.Context, aProps *dalSensitivityLevelActionProps, s store.Storer) error {
	return dalSensitivityLevelReload(ctx, svc.store, svc.services.dal)
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
		for i, s := range set {
			if s.ID == sl.ID {
				set[i] = sl
				break
			}
		}

		if !deleting && !updating {
			set = append(set, sl)
		}

		sort.Sort(set)
	}

	return set, err
}

func dalSensitivityLevelReload(ctx context.Context, s store.Storer, dsm dalSensitivityLevelManager) (err error) {
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
