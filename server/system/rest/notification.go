package rest

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/payload"
	"github.com/cortezaproject/corteza/server/system/rest/request"
	"github.com/cortezaproject/corteza/server/system/service"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	Notification struct {
		notification service.NotificationService
	}

	notificationSetPayload struct {
		Filter types.NotificationFilter `json:"filter"`
		Set    types.NotificationSet    `json:"set,omitempty"`
	}
)

func (Notification) New() *Notification {
	ctrl := &Notification{}
	ctrl.notification = service.DefaultNotification
	return ctrl
}

func (ctrl *Notification) List(ctx context.Context, r *request.NotificationList) (interface{}, error) {
	var (
		err error
		f   = types.NotificationFilter{
			Kind:    []types.NotificationKind{types.NotificationKind(r.Kind)},
			Read:    filter.State(r.Read),
			Deleted: filter.State(r.Deleted),
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	if len(r.NotificationID) > 0 {
		f.NotificationID = payload.ParseUint64s(r.NotificationID)
	}

	set, f, err := ctrl.notification.Find(ctx, f)
	if err != nil {
		return nil, err
	}

	// Ensure we have a valid empty set rather than nil
	if set == nil {
		set = make(types.NotificationSet, 0)
	}

	return ctrl.makeFilterPayload(ctx, set, f, nil)
}

func (ctrl *Notification) Create(ctx context.Context, r *request.NotificationCreate) (interface{}, error) {
	ntf := &types.Notification{
		Kind:      types.NotificationKind(r.Kind),
		Config:    types.NotificationConfig{},
		Recipient: r.Recipient,
	}

	// Convert sqlxTypes.JSONText to NotificationConfig
	if len(r.Config) > 0 {
		if err := r.Config.Unmarshal(&ntf.Config); err != nil {
			return nil, err
		}
	}

	return ctrl.notification.Create(ctx, ntf)
}

func (ctrl *Notification) Update(ctx context.Context, r *request.NotificationUpdate) (interface{}, error) {
	ntf := &types.Notification{
		ID:        r.NotificationID,
		Kind:      types.NotificationKind(r.Kind),
		Config:    types.NotificationConfig{},
		Recipient: r.Recipient,
	}

	// Convert sqlxTypes.JSONText to NotificationConfig
	if len(r.Config) > 0 {
		if err := r.Config.Unmarshal(&ntf.Config); err != nil {
			return nil, err
		}
	}

	return ctrl.notification.Update(ctx, ntf)
}

func (ctrl *Notification) Read(ctx context.Context, r *request.NotificationRead) (interface{}, error) {
	return ctrl.notification.FindByID(ctx, r.NotificationID)
}

func (ctrl *Notification) Delete(ctx context.Context, r *request.NotificationDelete) (interface{}, error) {
	return api.OK(), ctrl.notification.Delete(ctx, r.NotificationID)
}

func (ctrl *Notification) MarkAsRead(ctx context.Context, r *request.NotificationMarkAsRead) (interface{}, error) {
	return api.OK(), ctrl.notification.MarkAsRead(ctx, r.NotificationID)
}

func (ctrl *Notification) MarkAllAsRead(ctx context.Context, r *request.NotificationMarkAllAsRead) (interface{}, error) {
	return api.OK(), ctrl.notification.MarkAllAsRead(ctx)
}

func (ctrl *Notification) makeFilterPayload(ctx context.Context, nn types.NotificationSet, f types.NotificationFilter, err error) (*notificationSetPayload, error) {
	if err != nil {
		return nil, err
	}

	// Ensure we never return nil, but an empty slice instead
	if nn == nil {
		nn = make(types.NotificationSet, 0)
	}

	return &notificationSetPayload{Filter: f, Set: nn}, nil
}
