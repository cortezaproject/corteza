package rest

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/messagebus"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/davecgh/go-spew/spew"
)

type (
	Queue struct {
		svc service.QueueService
		ac  templateAccessController
	}

	queuePayload struct {
		*messagebus.QueueSettings
	}

	queueSetPayload struct {
		Filter messagebus.QueueSettingsFilter `json:"filter"`
		Set    []*queuePayload                `json:"set"`
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
		f   = messagebus.QueueSettingsFilter{
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
		q   = &messagebus.QueueSettings{
			Handler: r.Handler,
			Queue:   r.Queue,
			// Meta:    r.Meta,
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
		q   = &messagebus.QueueSettings{
			ID:      r.QueueID,
			Handler: r.Handler,
			Queue:   r.Queue,
			// Meta:    r.Meta,
		}
	)

	q, err = ctrl.svc.Update(ctx, q)

	return ctrl.makePayload(ctx, q, err)
}

func (ctrl *Queue) Delete(ctx context.Context, r *request.QueuesDelete) (interface{}, error) {
	spew.Dump("DELETE", r)
	return api.OK(), ctrl.svc.DeleteByID(ctx, r.QueueID)
}

func (ctrl *Queue) Undelete(ctx context.Context, r *request.TemplateUndelete) (interface{}, error) {
	return api.OK(), ctrl.svc.UndeleteByID(ctx, r.TemplateID)
}

func (ctrl *Queue) makePayload(ctx context.Context, q *messagebus.QueueSettings, err error) (*queuePayload, error) {
	if err != nil || q == nil {
		return nil, err
	}

	qq := &queuePayload{
		QueueSettings: q,
	}

	return qq, nil
}

func (ctrl *Queue) makeFilterPayload(ctx context.Context, nn messagebus.QueueSettingsSet, f messagebus.QueueSettingsFilter, err error) (*queueSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &queueSetPayload{Filter: f, Set: make([]*queuePayload, len(nn))}

	for i := range nn {
		msp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return msp, nil
}
