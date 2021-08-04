package service

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/apigw"
	a "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	apigwFunction struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        functionAccessController
	}

	functionAccessController interface {
		CanSearchApiGwFunctions(context.Context) bool

		CanCreateApiGwFunction(context.Context) bool
		CanReadApigwFunction(context.Context, *types.ApigwFunction) bool
		CanUpdateApigwFunction(context.Context, *types.ApigwFunction) bool
		CanDeleteApigwFunction(context.Context, *types.ApigwFunction) bool
	}
)

func Function() *apigwFunction {
	return (&apigwFunction{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	})
}

func (svc *apigwFunction) FindByID(ctx context.Context, ID uint64) (q *types.ApigwFunction, err error) {
	var (
		rProps = &apigwFunctionActionProps{}
	)

	err = func() error {
		if ID == 0 {
			return ApigwFunctionErrInvalidID()
		}

		if !svc.ac.CanSearchApiGwFunctions(ctx) {
			return ApigwFunctionErrNotAllowedToRead(rProps)
		}

		if q, err = store.LookupApigwFunctionByID(ctx, svc.store, ID); err != nil {
			return TemplateErrInvalidID().Wrap(err)
		}

		rProps.setFunction(q)

		return nil
	}()

	return q, svc.recordAction(ctx, rProps, ApigwFunctionActionLookup, err)
}

func (svc *apigwFunction) Create(ctx context.Context, new *types.ApigwFunction) (q *types.ApigwFunction, err error) {
	var (
		qProps = &apigwFunctionActionProps{function: new}
	)

	err = func() (err error) {
		if !svc.ac.CanCreateApiGwFunction(ctx) {
			return ApigwFunctionErrNotAllowedToCreate(qProps)
		}

		// Set new values after beforeCreate events are emitted
		new.ID = nextID()
		new.CreatedAt = *now()
		new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()

		if _, err = DefaultRoute.FindByID(ctx, new.Route); err != nil {
			return ApigwFunctionErrNotFound(qProps)
		}

		if err = store.CreateApigwFunction(ctx, svc.store, new); err != nil {
			return err
		}

		q = new
		// send the signal to reload all functions
		apigw.Service().Reload(ctx)

		return nil
	}()

	return q, svc.recordAction(ctx, qProps, ApigwFunctionActionCreate, err)
}

func (svc *apigwFunction) Update(ctx context.Context, upd *types.ApigwFunction) (q *types.ApigwFunction, err error) {
	var (
		qProps = &apigwFunctionActionProps{function: upd}
		qq     *types.ApigwFunction
		e      error
	)

	err = func() (err error) {
		if qq, e = store.LookupApigwFunctionByID(ctx, svc.store, upd.ID); e != nil {
			return ApigwFunctionErrNotFound(qProps)
		}

		if !svc.ac.CanUpdateApigwFunction(ctx, qq) {
			return ApigwFunctionErrNotAllowedToUpdate(qProps)
		}

		if _, err = DefaultRoute.FindByID(ctx, upd.Route); err != nil {
			return err
		}

		if qq, e = store.LookupApigwFunctionByID(ctx, svc.store, upd.ID); e == nil && qq == nil {
			return ApigwFunctionErrNotFound(qProps)
		}

		upd.UpdatedAt = now()
		upd.CreatedAt = qq.CreatedAt
		upd.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateApigwFunction(ctx, svc.store, upd); err != nil {
			return
		}

		q = upd

		// send the signal to reload all function
		apigw.Service().Reload(ctx)

		return nil
	}()

	return q, svc.recordAction(ctx, qProps, ApigwFunctionActionUpdate, err)
}

func (svc *apigwFunction) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		qProps = &apigwFunctionActionProps{}
		q      *types.ApigwFunction
	)

	err = func() (err error) {
		if q, err = store.LookupApigwFunctionByID(ctx, svc.store, ID); err != nil {
			return ApigwFunctionErrNotFound(qProps)
		}

		if !svc.ac.CanDeleteApigwFunction(ctx, q) {
			return ApigwFunctionErrNotAllowedToDelete(qProps)
		}

		qProps.setFunction(q)

		q.DeletedAt = now()
		q.DeletedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateApigwFunction(ctx, svc.store, q); err != nil {
			return
		}

		// send the signal to reload all queues
		apigw.Service().Reload(ctx)

		return nil
	}()

	return svc.recordAction(ctx, qProps, ApigwFunctionActionDelete, err)
}

func (svc *apigwFunction) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		qProps = &apigwFunctionActionProps{}
		q      *types.ApigwFunction
	)

	err = func() (err error) {
		if q, err = store.LookupApigwFunctionByID(ctx, svc.store, ID); err != nil {
			return ApigwFunctionErrNotFound(qProps)
		}

		if !svc.ac.CanDeleteApigwFunction(ctx, q) {
			return ApigwFunctionErrNotAllowedToDelete(qProps)
		}

		qProps.setFunction(q)

		q.DeletedAt = nil
		q.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateApigwFunction(ctx, svc.store, q); err != nil {
			return
		}

		// send the signal to reload all queues
		apigw.Service().Reload(ctx)

		return nil
	}()

	return svc.recordAction(ctx, qProps, ApigwFunctionActionDelete, err)
}

func (svc *apigwFunction) Search(ctx context.Context, filter types.ApigwFunctionFilter) (r types.ApigwFunctionSet, f types.ApigwFunctionFilter, err error) {
	var (
		aProps = &apigwFunctionActionProps{search: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.ApigwFunction) (bool, error) {
		if !svc.ac.CanReadApigwFunction(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if r, f, err = store.SearchApigwFunctions(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return r, f, svc.recordAction(ctx, aProps, ApigwFunctionActionSearch, err)
}

func (svc *apigwFunction) DefFunction(ctx context.Context, kind string) (l interface{}, err error) {
	var (
		qProps = &apigwFunctionActionProps{}
	)

	err = func() error {
		if !svc.ac.CanSearchApiGwFunctions(ctx) {
			return ApigwFunctionErrNotAllowedToRead(qProps)
		}

		// get the definitions from registry
		l = apigw.Service().Funcs(kind)

		return nil
	}()

	return l, svc.recordAction(ctx, qProps, ApigwFunctionActionSearch, err)

}

func (svc *apigwFunction) DefProxyAuth(ctx context.Context) (l interface{}, err error) {
	var (
		qProps = &apigwFunctionActionProps{}
	)

	err = func() error {
		if !svc.ac.CanSearchApiGwFunctions(ctx) {
			return ApigwFunctionErrNotAllowedToRead(qProps)
		}

		// get the definitions from registry
		l = apigw.Service().ProxyAuthDef()

		return nil
	}()

	return l, svc.recordAction(ctx, qProps, ApigwFunctionActionSearch, err)

}
