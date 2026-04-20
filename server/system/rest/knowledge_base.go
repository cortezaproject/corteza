package rest

import (
	"context"

	"github.com/crusttech/human/server/pkg/api"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

type (
	KnowledgeBase struct {
		svc knowledgeBaseService
		ac  knowledgeBaseAccessController
	}

	knowledgeBasePayload struct {
		*types.KnowledgeBase

		CanGrant               bool `json:"canGrant"`
		CanUpdateKnowledgeBase bool `json:"canUpdateKnowledgeBase"`
		CanDeleteKnowledgeBase bool `json:"canDeleteKnowledgeBase"`
	}

	knowledgeBaseSetPayload struct {
		Filter types.KnowledgeBaseFilter `json:"filter"`
		Set    []*knowledgeBasePayload   `json:"set"`
	}

	knowledgeBaseService interface {
		FindByID(ctx context.Context, ID uint64) (kb *types.KnowledgeBase, err error)
		Create(ctx context.Context, new *types.KnowledgeBase) (kb *types.KnowledgeBase, err error)
		Update(ctx context.Context, upd *types.KnowledgeBase) (kb *types.KnowledgeBase, err error)
		DeleteByID(ctx context.Context, ID uint64) (err error)
		UndeleteByID(ctx context.Context, ID uint64) (err error)
		Search(ctx context.Context, filter types.KnowledgeBaseFilter) (set types.KnowledgeBaseSet, f types.KnowledgeBaseFilter, err error)
	}

	knowledgeBaseAccessController interface {
		CanGrant(context.Context) bool

		CanCreateKnowledgeBase(context.Context) bool
		CanUpdateKnowledgeBase(context.Context, *types.KnowledgeBase) bool
		CanDeleteKnowledgeBase(context.Context, *types.KnowledgeBase) bool
	}
)

func (KnowledgeBase) New() *KnowledgeBase {
	return &KnowledgeBase{
		svc: service.DefaultKnowledgeBase,
		ac:  service.DefaultAccessControl,
	}
}

func (ctrl *KnowledgeBase) List(ctx context.Context, r *request.KnowledgeBaseList) (interface{}, error) {
	var (
		err error
		f   = types.KnowledgeBaseFilter{
			Query:   r.Query,
			Handle:  r.Handle,
			Deleted: filter.State(r.Deleted),
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	f.IncTotal = r.IncTotal

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	set, filter, err := ctrl.svc.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *KnowledgeBase) Create(ctx context.Context, r *request.KnowledgeBaseCreate) (interface{}, error) {
	var (
		err error
		kb  = &types.KnowledgeBase{
			Handle:      r.Handle,
			Title:       r.Title,
			Description: r.Description,
		}
	)

	if len(r.Context.Namespaces) > 0 {
		kb.Context = &r.Context
	}

	kb, err = ctrl.svc.Create(ctx, kb)
	return ctrl.makePayload(ctx, kb, err)
}

func (ctrl *KnowledgeBase) Read(ctx context.Context, r *request.KnowledgeBaseRead) (interface{}, error) {
	res, err := ctrl.svc.FindByID(ctx, r.KnowledgeBaseID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *KnowledgeBase) Update(ctx context.Context, r *request.KnowledgeBaseUpdate) (interface{}, error) {
	var (
		err error
		kb  = &types.KnowledgeBase{
			ID:          r.KnowledgeBaseID,
			Handle:      r.Handle,
			Title:       r.Title,
			Description: r.Description,
			UpdatedAt:   r.UpdatedAt,
		}
	)

	if len(r.Context.Namespaces) > 0 {
		kb.Context = &r.Context
	}

	kb, err = ctrl.svc.Update(ctx, kb)
	return ctrl.makePayload(ctx, kb, err)
}

func (ctrl *KnowledgeBase) Delete(ctx context.Context, r *request.KnowledgeBaseDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.DeleteByID(ctx, r.KnowledgeBaseID)
}

func (ctrl *KnowledgeBase) Undelete(ctx context.Context, r *request.KnowledgeBaseUndelete) (interface{}, error) {
	return api.OK(), ctrl.svc.UndeleteByID(ctx, r.KnowledgeBaseID)
}

func (ctrl *KnowledgeBase) makePayload(ctx context.Context, kb *types.KnowledgeBase, err error) (*knowledgeBasePayload, error) {
	if err != nil || kb == nil {
		return nil, err
	}

	return &knowledgeBasePayload{
		KnowledgeBase: kb,

		CanGrant: ctrl.ac.CanGrant(ctx),

		CanUpdateKnowledgeBase: ctrl.ac.CanUpdateKnowledgeBase(ctx, kb),
		CanDeleteKnowledgeBase: ctrl.ac.CanDeleteKnowledgeBase(ctx, kb),
	}, nil
}

func (ctrl *KnowledgeBase) makeFilterPayload(ctx context.Context, nn types.KnowledgeBaseSet, f types.KnowledgeBaseFilter, err error) (*knowledgeBaseSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &knowledgeBaseSetPayload{Filter: f, Set: make([]*knowledgeBasePayload, len(nn))}

	for i := range nn {
		msp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return msp, nil
}
