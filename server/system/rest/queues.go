package rest

import (
	"context"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

type (
	Queue struct {
		svc queueService
		ac  queueAccessController
	}

	queuePayload struct {
		*types.Queue

		CanGrant       bool `json:"canGrant"`
		CanUpdateQueue bool `json:"canUpdateQueue"`
		CanDeleteQueue bool `json:"canDeleteQueue"`
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

	queueAccessController interface {
		CanGrant(context.Context) bool

		CanCreateQueue(context.Context) bool
		CanUpdateQueue(context.Context, *types.Queue) bool
		CanDeleteQueue(context.Context, *types.Queue) bool
	}
)

func (Queue) New() *Queue {
	return &Queue{
		svc: service.DefaultQueue,
		ac:  service.DefaultAccessControl,
	}
}

// makeFilter builds the search filter for the generated List controller.
func (ctrl *Queue) makeFilter(ctx context.Context, r *request.QueuesList) (types.QueueFilter, error) {
	var (
		err error
		f   = types.QueueFilter{
			Query:   r.Query,
			Deleted: filter.State(r.Deleted),
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return f, err
	}

	f.IncTotal = r.IncTotal

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return f, err
	}

	return f, nil
}

// beforeCreate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Create controller
// already mapped the plain-value params.
func (ctrl *Queue) beforeCreate(ctx context.Context, res *types.Queue, r *request.QueuesCreate) error {
	res.Meta = r.Meta
	return nil
}

// beforeUpdate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Update controller
// already mapped the plain-value params (and ID/UpdatedAt).
func (ctrl *Queue) beforeUpdate(ctx context.Context, res *types.Queue, r *request.QueuesUpdate) error {
	res.Meta = r.Meta
	return nil
}

func (ctrl *Queue) makePayload(ctx context.Context, q *types.Queue, err error) (*queuePayload, error) {
	if err != nil || q == nil {
		return nil, err
	}

	qq := &queuePayload{
		Queue: q,

		CanGrant: ctrl.ac.CanGrant(ctx),

		CanUpdateQueue: ctrl.ac.CanUpdateQueue(ctx, q),
		CanDeleteQueue: ctrl.ac.CanDeleteQueue(ctx, q),
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
