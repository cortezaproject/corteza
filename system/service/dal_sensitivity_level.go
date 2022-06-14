package service

import (
	"context"
	"fmt"
	"sort"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	a "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/dal"

	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/dalutils"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	dalSensitivityLevel struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        sensitivityLevelAccessController
		dal       dalSensitivityLevels
	}

	sensitivityLevelAccessController interface {
		CanManageDalSensitivityLevel(context.Context) bool
	}

	dalSensitivityLevels interface {
		ReloadSensitivityLevels(levels ...dal.SensitivityLevel) (err error)
		CreateSensitivityLevel(levels ...dal.SensitivityLevel) (err error)
		UpdateSensitivityLevel(levels ...dal.SensitivityLevel) (err error)
		DeleteSensitivityLevel(levels ...dal.SensitivityLevel) (err error)
	}
)

func SensitivityLevel(ctx context.Context, dal dalSensitivityLevels) *dalSensitivityLevel {
	return &dalSensitivityLevel{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		dal:       dal,
	}

}

func (svc *dalSensitivityLevel) FindByID(ctx context.Context, ID uint64) (q *types.DalSensitivityLevel, err error) {
	var (
		rProps = &dalSensitivityLevelActionProps{}
	)

	err = func() error {
		if ID == 0 {
			return DalSensitivityLevelErrInvalidID()
		}

		if q, err = store.LookupDalSensitivityLevelByID(ctx, svc.store, ID); err != nil {
			return DalSensitivityLevelErrInvalidID().Wrap(err)
		}

		rProps.setSensitivityLevel(q)

		if !svc.ac.CanManageDalSensitivityLevel(ctx) {
			return DalSensitivityLevelErrNotAllowedToManage(rProps)
		}

		return nil
	}()

	return q, svc.recordAction(ctx, rProps, DalSensitivityLevelActionLookup, err)
}

func (svc *dalSensitivityLevel) Create(ctx context.Context, new *types.DalSensitivityLevel) (q *types.DalSensitivityLevel, err error) {
	var (
		qProps = &dalSensitivityLevelActionProps{new: new}
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if !svc.ac.CanManageDalSensitivityLevel(ctx) {
			return DalSensitivityLevelErrNotAllowedToManage(qProps)
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

		q = new

		return dalutils.DalSensitivityLevelCreate(svc.dal, new)
	})

	return q, svc.recordAction(ctx, qProps, DalSensitivityLevelActionCreate, err)
}

func (svc *dalSensitivityLevel) Update(ctx context.Context, upd *types.DalSensitivityLevel) (q *types.DalSensitivityLevel, err error) {
	var (
		qProps = &dalSensitivityLevelActionProps{update: upd}
		qq     *types.DalSensitivityLevel
		e      error
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if qq, e = store.LookupDalSensitivityLevelByID(ctx, s, upd.ID); e != nil {
			return DalSensitivityLevelErrNotFound(qProps)
		}

		if !svc.ac.CanManageDalSensitivityLevel(ctx) {
			return DalSensitivityLevelErrNotAllowedToManage(qProps)
		}

		upd.UpdatedAt = now()
		upd.CreatedAt = qq.CreatedAt
		upd.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

		ups, err := svc.prepare(ctx, s, upd)
		if err != nil {
			return
		}
		err = store.UpsertDalSensitivityLevel(ctx, s, ups...)
		if err != nil {
			return
		}

		q = upd

		return dalutils.DalSensitivityLevelUpdate(svc.dal, upd)
	})

	return q, svc.recordAction(ctx, qProps, DalSensitivityLevelActionUpdate, err)
}

func (svc *dalSensitivityLevel) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		qProps = &dalSensitivityLevelActionProps{}
		q      *types.DalSensitivityLevel
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if ID == 0 {
			return DalSensitivityLevelErrInvalidID()
		}

		if q, err = store.LookupDalSensitivityLevelByID(ctx, s, ID); err != nil {
			return
		}

		if !svc.ac.CanManageDalSensitivityLevel(ctx) {
			return DalSensitivityLevelErrNotAllowedToManage(qProps)
		}

		qProps.setSensitivityLevel(q)

		q.DeletedAt = now()
		q.DeletedBy = a.GetIdentityFromContext(ctx).Identity()

		ups, err := svc.prepare(ctx, s, q)
		if err != nil {
			return
		}
		err = store.UpsertDalSensitivityLevel(ctx, s, ups...)
		if err != nil {
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

		if err = dalutils.DalSensitivityLevelUpdate(svc.dal, uu...); err != nil {
			return err
		}
		if err = dalutils.DalSensitivityLevelDelete(svc.dal, dd...); err != nil {
			return err
		}
		return nil
	})

	return svc.recordAction(ctx, qProps, DalSensitivityLevelActionDelete, err)
}

func (svc *dalSensitivityLevel) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		qProps = &dalSensitivityLevelActionProps{}
		q      *types.DalSensitivityLevel
	)

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if ID == 0 {
			return DalSensitivityLevelErrInvalidID()
		}

		if q, err = store.LookupDalSensitivityLevelByID(ctx, s, ID); err != nil {
			return
		}

		if !svc.ac.CanManageDalSensitivityLevel(ctx) {
			return DalSensitivityLevelErrNotAllowedToManage(qProps)
		}

		qProps.setSensitivityLevel(q)

		q.DeletedAt = nil
		q.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateDalSensitivityLevel(ctx, s, q); err != nil {
			return
		}

		return dalutils.DalSensitivityLevelCreate(svc.dal, q)
	})

	return svc.recordAction(ctx, qProps, DalSensitivityLevelActionDelete, err)
}

func (svc *dalSensitivityLevel) Search(ctx context.Context, filter types.DalSensitivityLevelFilter) (r types.DalSensitivityLevelSet, f types.DalSensitivityLevelFilter, err error) {
	var (
		aProps = &dalSensitivityLevelActionProps{search: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.DalSensitivityLevel) (bool, error) {
		if !svc.ac.CanManageDalSensitivityLevel(ctx) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if !svc.ac.CanManageDalSensitivityLevel(ctx) {
			return DalSensitivityLevelErrNotAllowedToManage()
		}

		if r, f, err = store.SearchDalSensitivityLevels(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return r, f, svc.recordAction(ctx, aProps, DalSensitivityLevelActionSearch, err)
}

func (svc *dalSensitivityLevel) ReloadSensitivityLevels(ctx context.Context, s store.Storer) (err error) {
	return dalutils.DalSensitivityLevelReload(ctx, s, svc.dal)
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

		// Normalize sensitivity level
		offset := 0
		for i := range set {
			if set[i].DeletedAt != nil {
				offset++
			}

			nxtLvl := i + 1 - offset
			if nxtLvl != set[i].Level {
				set[i].UpdatedAt = now()
				// Same user so we can cheat a bit
				set[i].UpdatedBy = sl.CreatedBy
			}

			set[i].Level = nxtLvl
		}
	}

	return set, err
}
