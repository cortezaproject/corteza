package service

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	a "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service/event"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	dataPrivacy struct {
		actionlog actionlog.Recorder

		ac       dataPrivacyAccessController
		eventbus eventDispatcher

		store store.Storer
	}

	dataPrivacyAccessController interface {
		CanSearchDataPrivacyRequest(context.Context) bool
		CanCreateDataPrivacyRequest(context.Context) bool
		CanReadDataPrivacyRequest(context.Context, *types.DataPrivacyRequest) bool
		CanApproveDataPrivacyRequest(context.Context, *types.DataPrivacyRequest) bool
	}

	DataPrivacyService interface {
		FindRequestByID(ctx context.Context, requestID uint64) (*types.DataPrivacyRequest, error)
		FindRequest(context.Context, types.DataPrivacyRequestFilter) (types.DataPrivacyRequestSet, types.DataPrivacyRequestFilter, error)
		CreateRequest(ctx context.Context, request *types.DataPrivacyRequest) (*types.DataPrivacyRequest, error)
		UpdateRequest(ctx context.Context, request *types.DataPrivacyRequest) (*types.DataPrivacyRequest, error)
		UpdateRequestStatus(ctx context.Context, request *types.DataPrivacyRequest) (*types.DataPrivacyRequest, error)
	}
)

func DataPrivacy(s store.Storer, ac dataPrivacyAccessController, al actionlog.Recorder, eb eventDispatcher) *dataPrivacy {
	return &dataPrivacy{
		actionlog: al,
		ac:        ac,
		eventbus:  eb,
		store:     s,
	}
}

func (svc dataPrivacy) FindRequestByID(ctx context.Context, requestID uint64) (r *types.DataPrivacyRequest, err error) {
	var (
		raProps = &dataPrivacyActionProps{dataPrivacyRequest: &types.DataPrivacyRequest{ID: requestID}}
	)

	err = func() error {
		if requestID == 0 {
			return DataPrivacyErrInvalidID()
		}

		r, err = store.LookupDataPrivacyRequestByID(ctx, svc.store, requestID)
		if r, err = svc.procRequest(ctx, r, err); err != nil {
			return err
		}

		raProps.setDataPrivacyRequest(r)

		if !svc.ac.CanReadDataPrivacyRequest(ctx, r) {
			return DataPrivacyErrNotAllowedToRead()
		}

		return nil
	}()

	return r, svc.recordAction(ctx, raProps, DataPrivacyActionLookup, err)
}

func (svc dataPrivacy) procRequest(_ context.Context, r *types.DataPrivacyRequest, err error) (*types.DataPrivacyRequest, error) {
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, DataPrivacyErrNotFound()
		}

		return nil, err
	}

	return r, nil
}

func (svc dataPrivacy) FindRequest(ctx context.Context, filter types.DataPrivacyRequestFilter) (rr types.DataPrivacyRequestSet, f types.DataPrivacyRequestFilter, err error) {
	var (
		raProps = &dataPrivacyActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(req *types.DataPrivacyRequest) (bool, error) {
		if !svc.ac.CanReadDataPrivacyRequest(ctx, req) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if !svc.ac.CanSearchDataPrivacyRequest(ctx) {
			return DataPrivacyErrNotAllowedToSearch()
		}

		if filter.Deleted > 0 {
			// If list with deleted or suspended resource is requested
			// user must have access permissions to system (ie: is admin)
			//
			// not the best solution but ATM it allows us to have at least
			// some kind of control over who can see deleted or archived resource
			//if !svc.ac.CanAccess(ctx) {
			//	return {Resource}ErrNotAllowedToList{Resource}s()
			//}
		}

		if rr, f, err = store.SearchDataPrivacyRequests(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return rr, f, svc.recordAction(ctx, raProps, DataPrivacyActionSearch, err)
}

func (svc dataPrivacy) CreateRequest(ctx context.Context, new *types.DataPrivacyRequest) (r *types.DataPrivacyRequest, err error) {
	var (
		raProps = &dataPrivacyActionProps{new: new}
	)

	err = func() (err error) {
		if !svc.ac.CanCreateDataPrivacyRequest(ctx) {
			return DataPrivacyErrNotAllowedToCreate()
		}

		if err = svc.eventbus.WaitFor(ctx, event.DataPrivacyRequestBeforeCreate(new, r)); err != nil {
			return
		}

		new.ID = nextID()
		new.Status = types.RequestStatusPending
		new.RequestedAt = *now()
		new.RequestedBy = a.GetIdentityFromContext(ctx).Identity()
		new.CreatedAt = *now()

		if err = store.CreateDataPrivacyRequest(ctx, svc.store, new); err != nil {
			return
		}

		r = new

		svc.eventbus.Dispatch(ctx, event.DataPrivacyRequestAfterCreate(new, r))
		return
	}()

	return r, svc.recordAction(ctx, raProps, DataPrivacyActionCreate, err)
}

func (svc dataPrivacy) UpdateRequest(ctx context.Context, upd *types.DataPrivacyRequest) (r *types.DataPrivacyRequest, err error) {
	var (
		raProps = &dataPrivacyActionProps{update: upd}
	)

	err = func() (err error) {
		if upd.ID == 0 {
			return DataPrivacyErrInvalidID()
		}

		// Todo: check status update access control

		// Todo: access control

		if r, err = store.LookupDataPrivacyRequestByID(ctx, svc.store, upd.ID); err != nil {
			return
		}

		raProps.setDataPrivacyRequest(r)

		if err = svc.eventbus.WaitFor(ctx, event.DataPrivacyRequestBeforeUpdate(upd, r)); err != nil {
			return
		}

		r.UpdatedAt = now()

		// Assign changed values
		if err = store.UpdateDataPrivacyRequest(ctx, svc.store, r); err != nil {
			return err
		}

		svc.eventbus.Dispatch(ctx, event.DataPrivacyRequestAfterUpdate(upd, r))

		return nil
	}()

	return r, svc.recordAction(ctx, raProps, DataPrivacyActionUpdate, err)
}

func (svc dataPrivacy) UpdateRequestStatus(ctx context.Context, upd *types.DataPrivacyRequest) (r *types.DataPrivacyRequest, err error) {
	var (
		raProps = &dataPrivacyActionProps{update: upd}
	)

	err = func() (err error) {
		if upd.ID == 0 {
			return DataPrivacyErrInvalidID()
		}

		if upd.Status == types.RequestStatusPending {
			return DataPrivacyErrInvalidStatus()
		}

		if upd.Status == types.RequestStatusApprove || upd.Status == types.RequestStatusReject {
			if !svc.ac.CanApproveDataPrivacyRequest(ctx, upd) {
				return DataPrivacyErrNotAllowedToApprove()
			}
		}

		if r, err = store.LookupDataPrivacyRequestByID(ctx, svc.store, upd.ID); err != nil {
			return
		}

		raProps.setDataPrivacyRequest(r)

		if err = svc.eventbus.WaitFor(ctx, event.DataPrivacyRequestBeforeUpdate(upd, r)); err != nil {
			return
		}

		r.CompletedAt = now()
		r.CompletedBy = a.GetIdentityFromContext(ctx).Identity()
		r.UpdatedAt = now()

		// Assign changed values
		if err = store.UpdateDataPrivacyRequest(ctx, svc.store, r); err != nil {
			return err
		}

		svc.eventbus.Dispatch(ctx, event.DataPrivacyRequestAfterUpdate(upd, r))

		return nil
	}()

	return r, svc.recordAction(ctx, raProps, DataPrivacyActionApprove, err)
}
