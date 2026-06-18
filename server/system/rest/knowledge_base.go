package rest

import (
	"context"

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

// makeFilter builds the search filter for the generated List controller.
func (ctrl *KnowledgeBase) makeFilter(ctx context.Context, r *request.KnowledgeBaseList) (types.KnowledgeBaseFilter, error) {
	var (
		err error
		f   = types.KnowledgeBaseFilter{
			Query:   r.Query,
			Handle:  r.Handle,
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
func (ctrl *KnowledgeBase) beforeCreate(ctx context.Context, res *types.KnowledgeBase, r *request.KnowledgeBaseCreate) error {
	if len(r.Context.Namespaces) > 0 {
		res.Context = &r.Context
	}

	return nil
}

// beforeUpdate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Update controller
// already mapped the plain-value params (and ID/UpdatedAt).
func (ctrl *KnowledgeBase) beforeUpdate(ctx context.Context, res *types.KnowledgeBase, r *request.KnowledgeBaseUpdate) error {
	if len(r.Context.Namespaces) > 0 {
		res.Context = &r.Context
	}

	return nil
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
