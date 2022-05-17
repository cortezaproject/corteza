package service

import (
	"context"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	dataPrivacy struct {
		ac    dataPrivacyAccessController
		store store.Storer
	}

	dataPrivacyAccessController interface {
	}

	DataPrivacyService interface {
		FindByID(ctx context.Context, requestID uint64) (*types.DataPrivacyRequest, error)
		Find(context.Context, types.DataPrivacyRequestFilter) (types.DataPrivacyRequestSet, types.DataPrivacyRequestFilter, error)
		Create(ctx context.Context, request *types.DataPrivacyRequest) (*types.DataPrivacyRequest, error)
	}
)

func DataPrivacy(s store.Storer, ac dataPrivacyAccessController) *dataPrivacy {
	return &dataPrivacy{
		ac:    ac,
		store: s,
	}
}

func (svc dataPrivacy) FindByID(ctx context.Context, requestID uint64) (r *types.DataPrivacyRequest, err error) {
	// todo: actionlog

	err = func() error {
		if requestID == 0 {
			// fixme: err
			return RoleErrInvalidID()
		}

		//if r, err = svc.Lo(ctx, requestID); err != nil {
		//	return err
		//}

		// Todo: access control

		return nil
	}()

	// Todo: actionlog
	return r, nil

}

func (svc dataPrivacy) Find(ctx context.Context, filter types.DataPrivacyRequestFilter) (rr types.DataPrivacyRequestSet, f types.DataPrivacyRequestFilter, err error) {
	// todo: actionlog

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.DataPrivacyRequest) (bool, error) {
		// Todo: access control

		return true, nil
	}

	err = func() error {
		// Todo: access control

		if filter.Deleted > 0 {
			// If list with deleted or suspended users is requested
			// user must have access permissions to system (ie: is admin)
			//
			// not the best solution but ATM it allows us to have at least
			// some kind of control over who can see deleted or archived roles
			//if !svc.ac.CanAccess(ctx) {
			//	return RoleErrNotAllowedToListRoles()
			//}
		}

		//if rr, f, err = store.Search(ctx, svc.store, filter); err != nil {
		//	return err
		//}

		return nil
	}()

	return rr, f, nil
}

func (svc dataPrivacy) Create(ctx context.Context, new *types.DataPrivacyRequest) (r *types.DataPrivacyRequest, err error) {
	// todo: actionlog

	err = func() (err error) {
		// Todo: access control

		// todo: event before create if needed

		new.ID = nextID()
		new.CreatedAt = *now()

		//if err = store.Create(ctx, svc.store, new); err != nil {
		//	return
		//}

		r = new

		// todo: event after create if needed
		//svc.eventbus.Dispatch(ctx, )
		return
	}()

	return r, nil

}
