package service

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	dalSchemaAlteration struct {
		actionlog actionlog.Recorder
		ac        dalSchemaAlterationAccessController

		store store.Storer
	}

	dalSchemaAlterationAccessController interface {
	}
)

func DalSchemaAlteration() *dalSchemaAlteration {
	return &dalSchemaAlteration{
		ac:        DefaultAccessControl,
		store:     DefaultStore,
		actionlog: DefaultActionlog,
	}
}

func (svc dalSchemaAlteration) FindByID(ctx context.Context, dalSchemaAlterationID uint64) (a *types.DalSchemaAlteration, err error) {
	var (
		uaProps = &dalSchemaAlterationActionProps{dalSchemaAlteration: &types.DalSchemaAlteration{ID: dalSchemaAlterationID}}
	)

	err = func() error {
		a, err = loadDalSchemaAlteration(ctx, svc.store, dalSchemaAlterationID)
		if err != nil {
			return err
		}

		uaProps.setDalSchemaAlteration(a)

		// if !svc.ac.CanReadDalSchemaAlteration(ctx, u) {
		// 	return DalSchemaAlterationErrNotAllowedToRead()
		// }

		return nil
	}()

	return a, svc.recordAction(ctx, uaProps, DalSchemaAlterationActionLookup, err)
}

// Search interacts with backend storage and
//
// @todo rename to Search() for consistency
func (svc dalSchemaAlteration) Search(ctx context.Context, filter types.DalSchemaAlterationFilter) (aa types.DalSchemaAlterationSet, f types.DalSchemaAlterationFilter, err error) {
	var (
		uaProps = &dalSchemaAlterationActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	// if !svc.ac.CanReadDalSchemaAlteration(ctx, res) {
	// 	return false, nil
	// }

	// 	return true, nil
	// }

	err = func() error {
		// if !svc.ac.CanSearchDalSchemaAlterations(ctx) {
		// 	return DalSchemaAlterationErrNotAllowedToSearch()
		// }

		aa, f, err = store.SearchDalSchemaAlterations(ctx, svc.store, filter)
		return err
	}()

	return aa, f, svc.recordAction(ctx, uaProps, DalSchemaAlterationActionSearch, err)
}

func (svc dalSchemaAlteration) DeleteByID(ctx context.Context, dalSchemaAlterationID uint64) (err error) {
	var (
		u       *types.DalSchemaAlteration
		uaProps = &dalSchemaAlterationActionProps{dalSchemaAlteration: &types.DalSchemaAlteration{ID: dalSchemaAlterationID}}
	)

	err = func() (err error) {
		if u, err = loadDalSchemaAlteration(ctx, svc.store, dalSchemaAlterationID); err != nil {
			return
		}

		// if !svc.ac.CanDeleteDalSchemaAlteration(ctx, u) {
		// 	return DalSchemaAlterationErrNotAllowedToDelete()
		// }

		// if err = svc.eventbus.WaitFor(ctx, event.DalSchemaAlterationBeforeDelete(nil, u)); err != nil {
		// 	return
		// }

		u.DeletedAt = now()
		if err = store.UpdateDalSchemaAlteration(ctx, svc.store, u); err != nil {
			return
		}

		// _ = svc.eventbus.WaitFor(ctx, event.DalSchemaAlterationAfterDelete(nil, u))
		return nil
	}()

	return svc.recordAction(ctx, uaProps, DalSchemaAlterationActionDelete, err)
}

func (svc dalSchemaAlteration) UndeleteByID(ctx context.Context, dalSchemaAlterationID uint64) (err error) {
	var (
		u       *types.DalSchemaAlteration
		uaProps = &dalSchemaAlterationActionProps{dalSchemaAlteration: &types.DalSchemaAlteration{ID: dalSchemaAlterationID}}
	)

	err = func() (err error) {
		if u, err = loadDalSchemaAlteration(ctx, svc.store, dalSchemaAlterationID); err != nil {
			return
		}

		uaProps.setDalSchemaAlteration(u)

		// if err = uniqueDalSchemaAlterationCheck(ctx, svc.store, u); err != nil {
		// 	return err
		// }

		// if !svc.ac.CanDeleteDalSchemaAlteration(ctx, u) {
		// 	return DalSchemaAlterationErrNotAllowedToDelete()
		// }

		u.DeletedAt = nil
		if err = store.UpdateDalSchemaAlteration(ctx, svc.store, u); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, uaProps, DalSchemaAlterationActionUndelete, err)

}

func loadDalSchemaAlteration(ctx context.Context, s store.DalSchemaAlterations, ID uint64) (res *types.DalSchemaAlteration, err error) {
	if ID == 0 {
		return nil, DalSchemaAlterationErrInvalidID()
	}

	if res, err = store.LookupDalSchemaAlterationByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, DalSchemaAlterationErrNotFound()
	}

	return
}

// // uniqueDalSchemaAlterationCheck verifies dalSchemaAlteration's email, dalSchemaAlterationname and handle
// func uniqueDalSchemaAlterationCheck(ctx context.Context, s store.Storer, u *types.DalSchemaAlteration) (err error) {
// 	isUnique := func(field string) bool {
// 		f := types.DalSchemaAlterationFilter{
// 			// If dalSchemaAlteration exists and is deleted -- not a dup
// 			Deleted: filter.StateExcluded,

// 			// If dalSchemaAlteration exists and is suspended -- duplicate
// 			Suspended: filter.StateInclusive,
// 		}

// 		f.Limit = 1

// 		switch field {
// 		case "email":
// 			if u.Email == "" {
// 				return true
// 			}

// 			f.Email = u.Email

// 		case "dalSchemaAlterationname":
// 			if u.DalSchemaAlterationname == "" {
// 				return true
// 			}

// 			f.DalSchemaAlterationname = u.DalSchemaAlterationname
// 		case "handle":
// 			if u.Handle == "" {
// 				return true
// 			}

// 			f.Handle = u.Handle
// 		}

// 		set, _, err := store.SearchDalSchemaAlterations(ctx, s, f)
// 		if err != nil || len(set) > 1 {
// 			// In case of error or multiple dalSchemaAlterations returned
// 			return false
// 		}

// 		return len(set) == 0 || set[0].ID == u.ID
// 	}

// 	if !isUnique("email") {
// 		return DalSchemaAlterationErrEmailNotUnique()
// 	}

// 	if !isUnique("dalSchemaAlterationname") {
// 		return DalSchemaAlterationErrDalSchemaAlterationnameNotUnique()
// 	}

// 	if !isUnique("handle") {
// 		return DalSchemaAlterationErrHandleNotUnique()
// 	}

// 	return nil
// }
