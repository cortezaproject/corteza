package rest

import (
	"context"

	"github.com/crusttech/human/server/pkg/api"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/payload"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
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

// makeFilter builds the search filter for the generated List controller.
func (ctrl *Notification) makeFilter(ctx context.Context, r *request.NotificationList) (types.NotificationFilter, error) {
	var (
		err error
		f   = types.NotificationFilter{
			Kind:    []types.NotificationKind{types.NotificationKind(r.Kind)},
			Read:    filter.State(r.Read),
			Deleted: filter.State(r.Deleted),
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return f, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return f, err
	}

	if len(r.NotificationID) > 0 {
		f.NotificationID = payload.ParseUint64s(r.NotificationID)
	}

	return f, nil
}

// beforeCreate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Create controller
// already mapped the plain-value params.
func (ctrl *Notification) beforeCreate(ctx context.Context, res *types.Notification, r *request.NotificationCreate) error {
	res.Kind = types.NotificationKind(r.Kind)
	res.Config = types.NotificationConfig{}

	// Convert sqlxTypes.JSONText to NotificationConfig
	if len(r.Config) > 0 {
		if err := r.Config.Unmarshal(&res.Config); err != nil {
			return err
		}
	}

	return nil
}

// beforeUpdate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Update controller
// already mapped the plain-value params (and ID).
func (ctrl *Notification) beforeUpdate(ctx context.Context, res *types.Notification, r *request.NotificationUpdate) error {
	res.Kind = types.NotificationKind(r.Kind)
	res.Config = types.NotificationConfig{}

	// Convert sqlxTypes.JSONText to NotificationConfig
	if len(r.Config) > 0 {
		if err := r.Config.Unmarshal(&res.Config); err != nil {
			return err
		}
	}

	return nil
}

func (ctrl *Notification) MarkAsRead(ctx context.Context, r *request.NotificationMarkAsRead) (interface{}, error) {
	return api.OK(), ctrl.notification.MarkAsRead(ctx, r.NotificationID)
}

func (ctrl *Notification) MarkAsUnread(ctx context.Context, r *request.NotificationMarkAsUnread) (interface{}, error) {
	return api.OK(), ctrl.notification.MarkAsUnread(ctx, r.NotificationID)
}

func (ctrl *Notification) MarkAllAsRead(ctx context.Context, r *request.NotificationMarkAllAsRead) (interface{}, error) {
	return api.OK(), ctrl.notification.MarkAllAsRead(ctx)
}

func (ctrl *Notification) MarkAllAsUnread(ctx context.Context, r *request.NotificationMarkAllAsUnread) (interface{}, error) {
	return api.OK(), ctrl.notification.MarkAllAsUnread(ctx)
}

func (ctrl *Notification) makePayload(ctx context.Context, ntf *types.Notification, err error) (*types.Notification, error) {
	if err != nil {
		return nil, err
	}

	return ntf, nil
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
