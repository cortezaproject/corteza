package rest

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	DataPrivacyRequest struct {
		dataPrivacy service.DataPrivacyService
		ac          dataPrivacyAccessController
	}

	dataPrivacyRequestSetPayload struct {
		Filter types.DataPrivacyRequestFilter `json:"filter"`
		Set    types.DataPrivacyRequestSet    `json:"set"`
	}

	dataPrivacyAccessController interface {
		CanGrant(context.Context) bool
	}

	DataPrivacyRequestComment struct {
		dataPrivacy service.DataPrivacyService
		ac          dataPrivacyAccessController
	}

	dataPrivacyRequestCommentSetPayload struct {
		Filter types.DataPrivacyRequestCommentFilter `json:"filter"`
		Set    types.DataPrivacyRequestCommentSet    `json:"set"`
	}
)

func (DataPrivacyRequest) New() *DataPrivacyRequest {
	return &DataPrivacyRequest{
		dataPrivacy: service.DefaultDataPrivacy,
		ac:          service.DefaultAccessControl,
	}
}

func (ctrl DataPrivacyRequest) List(ctx context.Context, r *request.DataPrivacyRequestList) (interface{}, error) {
	var (
		err error
		f   = types.DataPrivacyRequestFilter{
			RequestedBy: payload.ParseUint64s(r.RequestedBy),
			Query:       r.Query,
			Kind:        r.Kind,
			Status:      r.Status,
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	set, f, err := ctrl.dataPrivacy.FindRequests(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, f, err)
}

func (ctrl DataPrivacyRequest) makeFilterPayload(_ context.Context, rr types.DataPrivacyRequestSet, f types.DataPrivacyRequestFilter, err error) (*dataPrivacyRequestSetPayload, error) {
	if err != nil {
		return nil, err
	}

	if len(rr) == 0 {
		rr = make([]*types.DataPrivacyRequest, 0)
	}

	return &dataPrivacyRequestSetPayload{Filter: f, Set: rr}, nil
}

func (ctrl DataPrivacyRequest) Create(ctx context.Context, r *request.DataPrivacyRequestCreate) (interface{}, error) {
	req := &types.DataPrivacyRequest{
		Kind:    types.CastToRequestKind(r.Kind),
		Payload: r.Payload,
	}

	return ctrl.dataPrivacy.CreateRequest(ctx, req)
}

func (ctrl DataPrivacyRequest) UpdateStatus(ctx context.Context, r *request.DataPrivacyRequestUpdateStatus) (interface{}, error) {
	req := &types.DataPrivacyRequest{
		ID:     r.RequestID,
		Status: types.CastToRequestStatus(r.Status),
	}

	return ctrl.dataPrivacy.UpdateRequestStatus(ctx, req)
}

func (ctrl DataPrivacyRequest) Read(ctx context.Context, r *request.DataPrivacyRequestRead) (interface{}, error) {
	return ctrl.dataPrivacy.FindRequestByID(ctx, r.RequestID)
}

func (DataPrivacyRequestComment) New() *DataPrivacyRequestComment {
	return &DataPrivacyRequestComment{
		dataPrivacy: service.DefaultDataPrivacy,
		ac:          service.DefaultAccessControl,
	}
}

func (ctrl DataPrivacyRequestComment) List(ctx context.Context, r *request.DataPrivacyRequestCommentList) (interface{}, error) {
	var (
		err error
		f   = types.DataPrivacyRequestCommentFilter{
			RequestID: []uint64{r.RequestID},
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	set, f, err := ctrl.dataPrivacy.FindRequestComments(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, f, err)
}

func (ctrl DataPrivacyRequestComment) makeFilterPayload(_ context.Context, rr types.DataPrivacyRequestCommentSet, f types.DataPrivacyRequestCommentFilter, err error) (*dataPrivacyRequestCommentSetPayload, error) {
	if err != nil {
		return nil, err
	}

	if len(rr) == 0 {
		rr = make([]*types.DataPrivacyRequestComment, 0)
	}

	return &dataPrivacyRequestCommentSetPayload{Filter: f, Set: rr}, nil
}

func (ctrl DataPrivacyRequestComment) Create(ctx context.Context, r *request.DataPrivacyRequestCommentCreate) (interface{}, error) {
	req := &types.DataPrivacyRequestComment{
		RequestID: r.RequestID,
		Comment:   r.Comment,
	}

	return ctrl.dataPrivacy.CreateRequestComment(ctx, req)
}
