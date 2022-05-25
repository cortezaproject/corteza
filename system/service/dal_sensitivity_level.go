package service

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	a "github.com/cortezaproject/corteza-server/pkg/auth"

	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	dalSensitivityLevel struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        sensitivityLevelAccessController
		// dal       dalSensitivityLevels
	}

	sensitivityLevelAccessController interface {
		CanManageDalSensitivityLevel(context.Context) bool
	}

	// dalSensitivityLevels interface {
	// 	AddSensitivityLevel(ctx context.Context, sensitivityLevelID uint64, cp dal.SensitivityLevelParams, dft dal.SensitivityLevelDefaults, capabilities ...capabilities.Capability) (err error)
	// 	UpdateSensitivityLevel(ctx context.Context, sensitivityLevelID uint64, cp dal.SensitivityLevelParams, dft dal.SensitivityLevelDefaults, capabilities ...capabilities.Capability) (err error)
	// 	RemoveSensitivityLevel(ctx context.Context, sensitivityLevelID uint64) (err error)
	// }
)

func SensitivityLevel(ctx context.Context) (*dalSensitivityLevel, error) {
	out := &dalSensitivityLevel{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		// dal:                     dal,
	}

	return out, nil // out.reloadSensitivityLevels(ctx)
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

	err = func() (err error) {
		if !svc.ac.CanManageDalSensitivityLevel(ctx) {
			return DalSensitivityLevelErrNotAllowedToManage(qProps)
		}

		new.ID = nextID()
		new.CreatedAt = *now()
		new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.CreateDalSensitivityLevel(ctx, svc.store, new); err != nil {
			return err
		}

		q = new

		// return svc.dal.AddSensitivityLevel(ctx, new.ID, new.Config.SensitivityLevel, new.SensitivityLevelDefaults(), new.ActiveCapabilities()...)
		return
	}()

	return q, svc.recordAction(ctx, qProps, DalSensitivityLevelActionCreate, err)
}

func (svc *dalSensitivityLevel) Update(ctx context.Context, upd *types.DalSensitivityLevel) (q *types.DalSensitivityLevel, err error) {
	var (
		qProps = &dalSensitivityLevelActionProps{update: upd}
		qq     *types.DalSensitivityLevel
		e      error
	)

	err = func() (err error) {
		if qq, e = store.LookupDalSensitivityLevelByID(ctx, svc.store, upd.ID); e != nil {
			return DalSensitivityLevelErrNotFound(qProps)
		}

		if !svc.ac.CanManageDalSensitivityLevel(ctx) {
			return DalSensitivityLevelErrNotAllowedToManage(qProps)
		}

		upd.UpdatedAt = now()
		upd.CreatedAt = qq.CreatedAt
		upd.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateDalSensitivityLevel(ctx, svc.store, upd); err != nil {
			return
		}

		q = upd

		// return svc.dal.UpdateSensitivityLevel(ctx, upd.ID, upd.Config.SensitivityLevel, upd.SensitivityLevelDefaults(), upd.ActiveCapabilities()...)
		return
	}()

	return q, svc.recordAction(ctx, qProps, DalSensitivityLevelActionUpdate, err)
}

func (svc *dalSensitivityLevel) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		qProps = &dalSensitivityLevelActionProps{}
		q      *types.DalSensitivityLevel
	)

	err = func() (err error) {
		if ID == 0 {
			return DalSensitivityLevelErrInvalidID()
		}

		if q, err = store.LookupDalSensitivityLevelByID(ctx, svc.store, ID); err != nil {
			return
		}

		if !svc.ac.CanManageDalSensitivityLevel(ctx) {
			return DalSensitivityLevelErrNotAllowedToManage(qProps)
		}

		qProps.setSensitivityLevel(q)

		q.DeletedAt = now()
		q.DeletedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateDalSensitivityLevel(ctx, svc.store, q); err != nil {
			return
		}

		// return svc.dal.RemoveSensitivityLevel(ctx, q.ID)
		return
	}()

	return svc.recordAction(ctx, qProps, DalSensitivityLevelActionDelete, err)
}

func (svc *dalSensitivityLevel) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		qProps = &dalSensitivityLevelActionProps{}
		q      *types.DalSensitivityLevel
	)

	err = func() (err error) {
		if ID == 0 {
			return DalSensitivityLevelErrInvalidID()
		}

		if q, err = store.LookupDalSensitivityLevelByID(ctx, svc.store, ID); err != nil {
			return
		}

		if !svc.ac.CanManageDalSensitivityLevel(ctx) {
			return DalSensitivityLevelErrNotAllowedToManage(qProps)
		}

		qProps.setSensitivityLevel(q)

		q.DeletedAt = nil
		q.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateDalSensitivityLevel(ctx, svc.store, q); err != nil {
			return
		}

		// return svc.dal.AddSensitivityLevel(ctx, q.ID, q.Config.SensitivityLevel, q.SensitivityLevelDefaults(), q.ActiveCapabilities()...)
		return
	}()

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

// func (svc *dalSensitivityLevel) reloadSensitivityLevels(ctx context.Context) (err error) {
// 	// Get all available sensitivityLevels
// 	cc, _, err := store.SearchDalSensitivityLevels(ctx, svc.store, types.DalSensitivityLevelFilter{})
// 	if err != nil {
// 		return
// 	}

// 	for _, c := range cc {
// // 		if err = svc.dal.AddSensitivityLevel(ctx, c.ID, c.Config.SensitivityLevel, c.SensitivityLevelDefaults(), c.ActiveCapabilities()...); err != nil {
// return
// 			return
// 		}
// 	}

// 	return
// }
