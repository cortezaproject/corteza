package rest

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	Queue struct {
		svc queueService
		ac  templateAccessController
	}

	queuePayload struct {
		*types.Queue
	}

	queueSetPayload struct {
		Filter types.QueueFilter `json:"filter"`
		Set    []*queuePayload   `json:"set"`
	}

	queueService interface {
		FindByID(ctx context.Context, ID uint64) (q *types.Queue, err error)
		Create(ctx context.Context, new *types.Queue) (q *types.Queue, err error)
		Update(ctx context.Context, upd *types.Queue) (q *types.Queue, err error)
		DeleteByID(ctx context.Context, ID uint64) (err error)
		UndeleteByID(ctx context.Context, ID uint64) (err error)
		Search(ctx context.Context, filter types.QueueFilter) (q types.QueueSet, f types.QueueFilter, err error)
	}
)

func (Queue) New() *Queue {
	return &Queue{
		svc: service.DefaultQueue,
		ac:  service.DefaultAccessControl,
	}
}

func (ctrl *Queue) List(ctx context.Context, r *request.QueuesList) (interface{}, error) {
	var (
		err error
		f   = types.QueueFilter{
			Query:   r.Query,
			Deleted: filter.State(r.Deleted),
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	set, filter, err := ctrl.svc.Search(ctx, f)

	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *Queue) Create(ctx context.Context, r *request.QueuesCreate) (interface{}, error) {
	var (
		err error
		q   = &types.Queue{
			Consumer: r.Consumer,
			Queue:    r.Queue,
			Meta:     r.Meta,
		}
	)

	q, err = ctrl.svc.Create(ctx, q)

	return ctrl.makePayload(ctx, q, err)
}

func (ctrl *Queue) Read(ctx context.Context, r *request.QueuesRead) (interface{}, error) {
	return ctrl.svc.FindByID(ctx, r.QueueID)
}

func (ctrl *Queue) Update(ctx context.Context, r *request.QueuesUpdate) (interface{}, error) {
	var (
		err error
		q   = &types.Queue{
			ID:       r.QueueID,
			Consumer: r.Consumer,
			Queue:    r.Queue,
			Meta:     r.Meta,
		}
	)

	q, err = ctrl.svc.Update(ctx, q)

	return ctrl.makePayload(ctx, q, err)
}

func (ctrl *Queue) Delete(ctx context.Context, r *request.QueuesDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.DeleteByID(ctx, r.QueueID)
}

func (ctrl *Queue) Undelete(ctx context.Context, r *request.QueuesUndelete) (interface{}, error) {
	return api.OK(), ctrl.svc.UndeleteByID(ctx, r.QueueID)
}

func (ctrl *Queue) makePayload(ctx context.Context, q *types.Queue, err error) (*queuePayload, error) {
	if err != nil || q == nil {
		return nil, err
	}

	qq := &queuePayload{
		Queue: q,
	}

	return qq, nil
}

func (ctrl *Queue) makeFilterPayload(ctx context.Context, nn types.QueueSet, f types.QueueFilter, err error) (*queueSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &queueSetPayload{Filter: f, Set: make([]*queuePayload, len(nn))}

	for i := range nn {
		msp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return msp, nil
}
