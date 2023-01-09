package rest

import (
	"context"
	"github.com/cortezaproject/corteza/server/pkg/filter"

	"github.com/cortezaproject/corteza/server/pkg/payload"
	"github.com/cortezaproject/corteza/server/system/rest/request"
	"github.com/cortezaproject/corteza/server/system/service"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	DataPrivacy struct {
		privacy service.DataPrivacyService
		ac      privacyAccessController
	}

	privacyAccessController interface {
		CanGrant(context.Context) bool
	}

	privacyConnectionSetPayload struct {
		Filter types.DalConnectionFilter     `json:"filter"`
		Set    types.PrivacyDalConnectionSet `json:"set"`
	}
)

func (DataPrivacy) New() *DataPrivacy {
	return &DataPrivacy{
		privacy: service.DefaultDataPrivacy,
		ac:      service.DefaultAccessControl,
	}
}

func (ctrl DataPrivacy) ConnectionList(ctx context.Context, r *request.DataPrivacyConnectionList) (interface{}, error) {
	var (
		err error
		set types.PrivacyDalConnectionSet

		f = types.DalConnectionFilter{
			ConnectionID: payload.ParseUint64s(r.ConnectionID),
			Handle:       r.Handle,
			Type:         r.Type,

			Deleted: r.Deleted,
		}
	)

	set, f, err = ctrl.privacy.FindConnections(ctx, f)
	return ctrl.makeFilterConnectionPayload(ctx, set, f, err)
}

func (ctrl DataPrivacy) makeFilterConnectionPayload(_ context.Context, rr types.PrivacyDalConnectionSet, f types.DalConnectionFilter, err error) (*privacyConnectionSetPayload, error) {
	if err != nil {
		return nil, err
	}

	if len(rr) == 0 {
		rr = make([]*types.PrivacyDalConnection, 0)
	}

	return &privacyConnectionSetPayload{Filter: f, Set: rr}, nil
}

func (ctrl DataPrivacy) RequestList(ctx context.Context, r *request.DataPrivacyRequestList) (interface{}, error) {
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

	set, f, err := ctrl.privacy.FindRequests(ctx, f)
	return ctrl.makeFilterRequestPayload(ctx, set, f, err)
}

func (ctrl DataPrivacy) makeFilterRequestPayload(_ context.Context, rr types.DataPrivacyRequestSet, f types.DataPrivacyRequestFilter, err error) (*dataPrivacyRequestSetPayload, error) {
	if err != nil {
		return nil, err
	}

	if len(rr) == 0 {
		rr = make([]*types.DataPrivacyRequest, 0)
	}

	return &dataPrivacyRequestSetPayload{Filter: f, Set: rr}, nil
}

func (ctrl DataPrivacy) RequestCreate(ctx context.Context, r *request.DataPrivacyRequestCreate) (interface{}, error) {
	req := &types.DataPrivacyRequest{
		Kind:    types.CastToRequestKind(r.Kind),
		Payload: r.Payload,
	}

	return ctrl.privacy.CreateRequest(ctx, req)
}

func (ctrl DataPrivacy) RequestRead(ctx context.Context, r *request.DataPrivacyRequestRead) (interface{}, error) {
	return ctrl.privacy.FindRequestByID(ctx, r.RequestID)
}

func (ctrl DataPrivacy) RequestUpdateStatus(ctx context.Context, r *request.DataPrivacyRequestUpdateStatus) (interface{}, error) {
	req := &types.DataPrivacyRequest{
		ID:     r.RequestID,
		Status: types.CastToRequestStatus(r.Status),
	}

	return ctrl.privacy.UpdateRequestStatus(ctx, req)
}

func (ctrl DataPrivacy) RequestCommentList(ctx context.Context, r *request.DataPrivacyRequestCommentList) (interface{}, error) {
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

	set, f, err := ctrl.privacy.FindRequestComments(ctx, f)
	return ctrl.makeFilterRequestCommentPayload(ctx, set, f, err)
}

func (ctrl DataPrivacy) makeFilterRequestCommentPayload(_ context.Context, rr types.DataPrivacyRequestCommentSet, f types.DataPrivacyRequestCommentFilter, err error) (*dataPrivacyRequestCommentSetPayload, error) {
	if err != nil {
		return nil, err
	}

	if len(rr) == 0 {
		rr = make([]*types.DataPrivacyRequestComment, 0)
	}

	return &dataPrivacyRequestCommentSetPayload{Filter: f, Set: rr}, nil
}

func (ctrl DataPrivacy) RequestCommentCreate(ctx context.Context, r *request.DataPrivacyRequestCommentCreate) (interface{}, error) {
	req := &types.DataPrivacyRequestComment{
		RequestID: r.RequestID,
		Comment:   r.Comment,
	}

	return ctrl.privacy.CreateRequestComment(ctx, req)
}
