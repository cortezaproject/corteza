package rest

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	DataPrivacy struct {
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

func (DataPrivacy) New() *DataPrivacy {
	return &DataPrivacy{
		dataPrivacy: service.DefaultDataPrivacy,
		ac:          service.DefaultAccessControl,
	}
}

func (ctrl DataPrivacy) ListRequests(ctx context.Context, r *request.DataPrivacyListRequests) (interface{}, error) {
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

func (ctrl DataPrivacy) makeFilterPayload(_ context.Context, rr types.DataPrivacyRequestSet, f types.DataPrivacyRequestFilter, err error) (*dataPrivacyRequestSetPayload, error) {
	if err != nil {
		return nil, err
	}

	if len(rr) == 0 {
		rr = make([]*types.DataPrivacyRequest, 0)
	}

	return &dataPrivacyRequestSetPayload{Filter: f, Set: rr}, nil
}

func (ctrl DataPrivacy) CreateRequest(ctx context.Context, r *request.DataPrivacyCreateRequest) (interface{}, error) {
	req := &types.DataPrivacyRequest{
		Name: r.Name,
		Kind: types.RequestKind(r.Kind),
	}

	return ctrl.dataPrivacy.CreateRequest(ctx, req)
}

func (ctrl DataPrivacy) UpdateRequest(ctx context.Context, r *request.DataPrivacyUpdateRequest) (interface{}, error) {
	req := &types.DataPrivacyRequest{
		ID:     r.RequestID,
		Status: types.RequestStatus(r.Status),
	}

	return ctrl.dataPrivacy.UpdateRequest(ctx, req)
}

func (ctrl DataPrivacy) UpdateRequestStatus(ctx context.Context, r *request.DataPrivacyUpdateRequestStatus) (interface{}, error) {
	req := &types.DataPrivacyRequest{
		ID:     r.RequestID,
		Status: types.RequestStatus(r.Status),
	}

	return ctrl.dataPrivacy.UpdateRequestStatus(ctx, req)
}

func (ctrl DataPrivacy) ReadRequest(ctx context.Context, r *request.DataPrivacyReadRequest) (interface{}, error) {
	return ctrl.dataPrivacy.FindRequestByID(ctx, r.RequestID)
}

func (ctrl DataPrivacy) ListResponsesOfRequest(ctx context.Context, request *request.DataPrivacyListResponsesOfRequest) (interface{}, error) {
	panic("implement me")
}

func (ctrl DataPrivacy) CreateResponseForRequest(ctx context.Context, request *request.DataPrivacyCreateResponseForRequest) (interface{}, error) {
	panic("implement me")
}
