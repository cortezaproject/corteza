package rest

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/filter"
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
			Deleted: filter.State(r.Deleted),
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	set, f, err := ctrl.dataPrivacy.FindRequest(ctx, f)
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
		Name: r.Name,
		Kind: types.RequestKind(r.Kind),
	}

	return ctrl.dataPrivacy.CreateRequest(ctx, req)
}

func (ctrl DataPrivacyRequest) Update(ctx context.Context, r *request.DataPrivacyRequestUpdate) (interface{}, error) {
	req := &types.DataPrivacyRequest{
		ID:     r.RequestID,
		Status: types.RequestStatus(r.Status),
	}

	return ctrl.dataPrivacy.UpdateRequest(ctx, req)
}

func (ctrl DataPrivacyRequest) UpdateStatus(ctx context.Context, r *request.DataPrivacyRequestUpdateStatus) (interface{}, error) {
	req := &types.DataPrivacyRequest{
		ID:     r.RequestID,
		Status: types.RequestStatus(r.Status),
	}

	return ctrl.dataPrivacy.UpdateRequestStatus(ctx, req)
}

func (ctrl DataPrivacyRequest) Read(ctx context.Context, r *request.DataPrivacyRequestRead) (interface{}, error) {
	return ctrl.dataPrivacy.FindRequestByID(ctx, r.RequestID)
}

func (ctrl DataPrivacyRequest) ListResponses(ctx context.Context, request *request.DataPrivacyRequestListResponses) (interface{}, error) {
	panic("implement me")
}

func (ctrl DataPrivacyRequest) CreateResponse(ctx context.Context, request *request.DataPrivacyRequestCreateResponse) (interface{}, error) {
	panic("implement me")
}
