package service

import (
	"context"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"math"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/apigw"
	"github.com/cortezaproject/corteza-server/pkg/apigw/profiler"
	a "github.com/cortezaproject/corteza-server/pkg/auth"

	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	apigwRoute struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        routeAccessController
	}

	routeAccessController interface {
		CanGrant(context.Context) bool
		CanSearchApigwRoutes(ctx context.Context) bool

		CanCreateApigwRoute(context.Context) bool
		CanReadApigwRoute(context.Context, *types.ApigwRoute) bool
		CanUpdateApigwRoute(context.Context, *types.ApigwRoute) bool
		CanDeleteApigwRoute(context.Context, *types.ApigwRoute) bool
	}
)

func Route() *apigwRoute {
	return &apigwRoute{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	}
}

func (svc *apigwRoute) FindByID(ctx context.Context, ID uint64) (q *types.ApigwRoute, err error) {
	var (
		rProps = &apigwRouteActionProps{}
	)

	err = func() error {
		if ID == 0 {
			return ApigwRouteErrInvalidID()
		}

		if q, err = store.LookupApigwRouteByID(ctx, svc.store, ID); err != nil {
			return ApigwRouteErrInvalidID().Wrap(err)
		}

		rProps.setRoute(q)

		if !svc.ac.CanReadApigwRoute(ctx, q) {
			return ApigwRouteErrNotAllowedToRead(rProps)
		}

		return nil
	}()

	return q, svc.recordAction(ctx, rProps, ApigwRouteActionLookup, err)
}

func (svc *apigwRoute) Create(ctx context.Context, new *types.ApigwRoute) (q *types.ApigwRoute, err error) {
	var (
		qProps = &apigwRouteActionProps{new: new}
	)

	err = func() (err error) {
		if !svc.ac.CanCreateApigwRoute(ctx) {
			return ApigwRouteErrNotAllowedToCreate(qProps)
		}

		new.ID = nextID()
		new.CreatedAt = *now()
		new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()

		// todo
		new.Group = 0

		if err = store.CreateApigwRoute(ctx, svc.store, new); err != nil {
			return err
		}

		q = new

		// send the signal to reload all routes
		if new.Enabled {
			if err = apigw.Service().Reload(ctx); err != nil {
				return err
			}
		}

		return nil
	}()

	return q, svc.recordAction(ctx, qProps, ApigwRouteActionCreate, err)
}

func (svc *apigwRoute) Update(ctx context.Context, upd *types.ApigwRoute) (q *types.ApigwRoute, err error) {
	var (
		qProps = &apigwRouteActionProps{update: upd}
		qq     *types.ApigwRoute
		e      error
	)

	err = func() (err error) {
		if qq, e = store.LookupApigwRouteByID(ctx, svc.store, upd.ID); e != nil {
			return ApigwRouteErrNotFound(qProps)
		}

		if !svc.ac.CanUpdateApigwRoute(ctx, qq) {
			return ApigwRouteErrNotAllowedToUpdate(qProps)
		}

		// temp todo - update itself with the same endpoint
		// if qq, e = store.LookupApigwRouteByEndpoint(ctx, svc.store, upd.Endpoint); e == nil && qq == nil {
		// 	return ApigwRouteErrExistsEndpoint(qProps)
		// }

		upd.UpdatedAt = now()
		upd.CreatedAt = qq.CreatedAt
		upd.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateApigwRoute(ctx, svc.store, upd); err != nil {
			return
		}

		q = upd

		// send the signal to reload all route
		if qq.Enabled != upd.Enabled || qq.Enabled && upd.Enabled {
			if err = apigw.Service().Reload(ctx); err != nil {
				return err
			}
		}

		return nil
	}()

	return q, svc.recordAction(ctx, qProps, ApigwRouteActionUpdate, err)
}

func (svc *apigwRoute) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		qProps = &apigwRouteActionProps{}
		q      *types.ApigwRoute
	)

	err = func() (err error) {
		if ID == 0 {
			return ApigwRouteErrInvalidID()
		}

		if q, err = store.LookupApigwRouteByID(ctx, svc.store, ID); err != nil {
			return
		}

		if !svc.ac.CanDeleteApigwRoute(ctx, q) {
			return ApigwRouteErrNotAllowedToDelete(qProps)
		}

		qProps.setRoute(q)

		q.DeletedAt = now()
		q.DeletedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateApigwRoute(ctx, svc.store, q); err != nil {
			return
		}

		// send the signal to reload all queues
		if q.Enabled {
			if err = apigw.Service().Reload(ctx); err != nil {
				return err
			}
		}

		return nil
	}()

	return svc.recordAction(ctx, qProps, ApigwRouteActionDelete, err)
}

func (svc *apigwRoute) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		qProps = &apigwRouteActionProps{}
		q      *types.ApigwRoute
	)

	err = func() (err error) {
		if ID == 0 {
			return ApigwRouteErrInvalidID()
		}

		if q, err = store.LookupApigwRouteByID(ctx, svc.store, ID); err != nil {
			return
		}

		if !svc.ac.CanDeleteApigwRoute(ctx, q) {
			return ApigwRouteErrNotAllowedToUndelete(qProps)
		}

		qProps.setRoute(q)

		q.DeletedAt = nil
		q.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateApigwRoute(ctx, svc.store, q); err != nil {
			return
		}

		// send the signal to reload all queues
		if q.Enabled {
			if err = apigw.Service().Reload(ctx); err != nil {
				return err
			}
		}

		return nil
	}()

	return svc.recordAction(ctx, qProps, ApigwRouteActionDelete, err)
}

func (svc *apigwRoute) Search(ctx context.Context, filter types.ApigwRouteFilter) (r types.ApigwRouteSet, f types.ApigwRouteFilter, err error) {
	var (
		aProps = &apigwRouteActionProps{search: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.ApigwRoute) (bool, error) {
		if !svc.ac.CanReadApigwRoute(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if !svc.ac.CanSearchApigwRoutes(ctx) {
			return ApigwRouteErrNotAllowedToSearch()
		}

		if r, f, err = store.SearchApigwRoutes(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return r, f, svc.recordAction(ctx, aProps, ApigwRouteActionSearch, err)
}

// HitsAggregated fetches a list of hits from integration gateway profiler
func (svc *apigwRoute) Hits(ctx context.Context, filter types.ApigwProfilerFilter) (r types.ApigwProfilerHitSet, f types.ApigwProfilerFilter, err error) {

	f = filter
	r = make(types.ApigwProfilerHitSet, 0)

	uDec, err := base64.URLEncoding.DecodeString(filter.Path)

	if err != nil {
		return
	}

	filter.Path = string(uDec)

	if filter.Path == "" && filter.Hit == "" {
		err = errors.New("fetching all hits (no route and hit specified) not supported")
		return
	}

	var sorting = profiler.Sort{
		Hit:    filter.Hit,
		Path:   filter.Path,
		Before: filter.Before,
	}

	var (
		list = apigw.Service().Profiler().Hits(sorting)
	)

	var pp = ""

	for k, _ := range list {
		if filter.Hit != "" {
			pp = k
			break
		}

		if filter.Path != "" && k == filter.Path {
			pp = k
			break
		}
	}

	if pp == "" {
		return
	}

	for _, h := range list[pp] {
		hh := &types.ApigwProfilerHit{
			ID: h.ID,

			Route:   h.Route,
			Status:  h.Status,
			Request: *h.R,

			Ts: h.Ts,
			Tf: h.Tf,
			D:  h.D,
			Dr: float64(h.D.Microseconds()) / 1000,
		}

		// fetch body only on hit details
		if filter.Hit != "" {
			hh.Body, _ = ioutil.ReadAll(hh.Request.Body)
		}

		r = append(r, hh)
	}

	// sort
	profiler.SortHits(&r, &f)

	// filter sorted
	profiler.FilterHits(&r, &f)

	return
}

// HitsAggregated fetches a list of hits from integration gateway profiler
// and aggregates them with assigned filters
func (svc *apigwRoute) HitsAggregated(ctx context.Context, filter types.ApigwProfilerFilter) (r types.ApigwProfilerAggregationSet, f types.ApigwProfilerFilter, err error) {
	f = filter
	r = make(types.ApigwProfilerAggregationSet, 0)

	uDec, err := base64.URLEncoding.DecodeString(filter.Path)

	if err != nil {
		return
	}

	filter.Path = string(uDec)

	var (
		list = apigw.Service().Profiler().Hits(profiler.Sort{
			Path:   filter.Path,
			Before: filter.Before,
		})

		tsum, tmin, tmax time.Duration
		ssum, smin, smax int64
		i                uint64 = 1
	)

	for p, v := range list {
		tmin, tmax, tsum = time.Hour, 0, 0
		smin, smax, ssum = math.MaxInt64, 0, 0

		i = 0

		for _, vv := range v {
			var (
				d = *vv.D
				s = vv.R.ContentLength
			)

			if d < tmin {
				tmin = d
			}

			if d > tmax {
				tmax = d
			}

			if s < smin {
				smin = s
			}

			if s > smax {
				smax = s
			}

			tsum += d
			ssum += s
			i++
		}

		r = append(r, &types.ApigwProfilerAggregation{
			Path:  p,
			Count: i,
			Tmin:  float64(tmin.Microseconds()) / 1000,
			Tmax:  float64(tmax.Microseconds()) / 1000,
			Tavg:  float64(tsum.Microseconds()) / float64(i) / 1000,
			Smin:  smin,
			Smax:  smax,
			Savg:  float64(ssum) / float64(i),
		})
	}

	// sort
	profiler.SortAggregation(&r, &f)

	// filter
	profiler.FilterAggregation(&r, &f)

	return
}
