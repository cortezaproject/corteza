package rest

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/automation/rest/request"
	"github.com/cortezaproject/corteza-server/automation/service"
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
)

type (
	Session struct {
		svc sessionService
	}

	sessionService interface {
		Search(ctx context.Context, filter types.SessionFilter) (types.SessionSet, types.SessionFilter, error)
		LookupByID(ctx context.Context, sessionID uint64) (*types.Session, error)
		Resume(sessionID, stateID uint64, i auth.Identifiable, input *expr.Vars) error
		PendingPrompts(context.Context) []*wfexec.PendingPrompt
	}

	sessionSetPayload struct {
		Filter types.SessionFilter `json:"filter"`
		Set    types.SessionSet    `json:"set"`
	}
)

func (Session) New() *Session {
	ctrl := &Session{}
	ctrl.svc = service.DefaultSession
	return ctrl
}

func (ctrl Session) List(ctx context.Context, r *request.SessionList) (interface{}, error) {
	var (
		err error
		f   = types.SessionFilter{
			WorkflowID:   payload.ParseUint64s(r.WorkflowID),
			SessionID:    payload.ParseUint64s(r.SessionID),
			EventType:    r.EventType,
			ResourceType: r.ResourceType,
			Completed:    filter.State(r.Completed),
			Status:       r.Status,
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

func (ctrl Session) Read(ctx context.Context, r *request.SessionRead) (interface{}, error) {
	return ctrl.svc.LookupByID(ctx, r.SessionID)
}

func (ctrl Session) Delete(ctx context.Context, r *request.SessionDelete) (interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

func (ctrl Session) Trace(ctx context.Context, trace *request.SessionTrace) (interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

func (ctrl Session) ListPrompts(ctx context.Context, r *request.SessionListPrompts) (interface{}, error) {
	return struct {
		Set []*wfexec.PendingPrompt `json:"set"`
	}{
		Set: ctrl.svc.PendingPrompts(ctx),
	}, nil
}

func (ctrl Session) ResumeState(ctx context.Context, r *request.SessionResumeState) (interface{}, error) {
	return api.OK(), ctrl.svc.Resume(r.SessionID, r.StateID, auth.GetIdentityFromContext(ctx), r.Input)
}

func (ctrl Session) DeleteState(ctx context.Context, r *request.SessionDeleteState) (interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

func (ctrl Session) makeFilterPayload(ctx context.Context, uu types.SessionSet, f types.SessionFilter, err error) (*sessionSetPayload, error) {
	if err != nil {
		return nil, err
	}

	if len(uu) == 0 {
		uu = make([]*types.Session, 0)
	}

	return &sessionSetPayload{Filter: f, Set: uu}, nil
}
