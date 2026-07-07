package service

import (
	"context"
	"time"

	"github.com/crusttech/human/server/pkg/actionlog"
	intAuth "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
	"go.uber.org/zap"
)

type (
	notification struct {
		ac                 notificationAccessController
		log                *zap.Logger
		actionlog          actionlog.Recorder
		store              store.Storer
		notificationSender notificationSender
	}

	notificationSender interface {
		Send(kind string, payload interface{}, userIDs ...uint64) error
	}

	notificationAccessController interface {
		CanAssignNotification(ctx context.Context) bool
	}

	NotificationService interface {
		Search(context.Context, types.NotificationFilter) (types.NotificationSet, types.NotificationFilter, error)
		FindByID(context.Context, uint64) (*types.Notification, error)

		Create(context.Context, *types.Notification) (*types.Notification, error)
		Update(context.Context, *types.Notification) (*types.Notification, error)
		DeleteByID(context.Context, uint64) error

		MarkAsRead(context.Context, uint64) error
		MarkAsUnread(context.Context, uint64) error
		MarkAllAsRead(context.Context) error
		MarkAllAsUnread(context.Context) error
	}
)

func Notification(ctx context.Context, log *zap.Logger, ns notificationSender) NotificationService {
	return &notification{
		ac:                 DefaultAccessControl,
		log:                log,
		store:              DefaultStore,
		notificationSender: ns,
	}
}

// onSearch is the custom body for the generated Search. The generated method
// owns the action-log scaffold + recordAction; the recipient-scoping (read
// access is recipient-ownership, non-RBAC) lives here.
func (svc *notification) onSearch(ctx context.Context, filter types.NotificationFilter, aProps *notificationActionProps) (nn types.NotificationSet, f types.NotificationFilter, err error) {
	// Get current user ID
	currentUserID := intAuth.GetIdentityFromContext(ctx).Identity()

	// Override recipient filter to ensure users can only see their own notifications
	filter.Recipient = currentUserID

	nn, f, err = store.SearchNotifications(ctx, svc.store, filter)
	if err != nil {
		return nn, f, err
	}

	// Ensure we never return nil, but an empty slice instead
	if nn == nil {
		nn = make(types.NotificationSet, 0)
	}

	return nn, f, nil
}

// onLookup is the custom body for the generated FindByID. The generated method
// owns the action-log scaffold + recordAction; the recipient-ownership check
// (non-RBAC read access) lives here.
func (svc *notification) onLookup(ctx context.Context, ID uint64, aProps *notificationActionProps) (n *types.Notification, err error) {
	if ID == 0 {
		return nil, NotificationErrInvalidID()
	}

	n, err = store.LookupNotificationByID(ctx, svc.store, ID)
	if err != nil {
		return nil, err
	}

	// Check if the notification belongs to the current user
	currentUserID := intAuth.GetIdentityFromContext(ctx).Identity()
	if n.Recipient != currentUserID {
		return nil, NotificationErrNotFound()
	}

	aProps.setNotification(n)

	return n, nil
}

// onCreate is the custom body for the generated Create. The generated method
// owns the action-log scaffold; the assignee permission check, creator
// stamping and websocket delivery live here.
func (svc *notification) onCreate(ctx context.Context, new *types.Notification) (err error) {
	// Check if current user has permission to assign notifications to others
	if err := svc.checkAssignee(ctx, new); err != nil {
		return err
	}

	new.ID = nextID()
	new.CreatedAt = *now()

	// Set creator
	currentUserID := intAuth.GetIdentityFromContext(ctx).Identity()
	new.CreatedBy = currentUserID

	if err = store.CreateNotification(ctx, svc.store, new); err != nil {
		return err
	}

	// Send the notification via websocket
	if svc.notificationSender != nil {
		// Send only to the recipient
		if err = svc.notificationSender.Send("notification", new, new.Recipient); err != nil {
			return err
		}
	}

	return nil
}

// onUpdate is the custom body for the generated Update. The generated method
// owns the action-log scaffold + recordAction; the recipient-ownership check,
// assignee re-check on recipient change and field copy live here.
func (svc *notification) onUpdate(ctx context.Context, s store.Storer, upd, res *types.Notification, aProps *notificationActionProps, _ func() error, _ func() error) error {
	if upd.ID == 0 {
		return NotificationErrInvalidID()
	}

	// Check if the notification belongs to the current user
	currentUserID := intAuth.GetIdentityFromContext(ctx).Identity()
	if res.Recipient != currentUserID {
		return NotificationErrNotFound()
	}

	// Check if the recipient is being changed and if so, check permissions
	if upd.Recipient != 0 && upd.Recipient != res.Recipient {
		tempNotification := &types.Notification{
			Recipient: upd.Recipient,
		}

		if err := svc.checkAssignee(ctx, tempNotification); err != nil {
			return err
		}

		res.Recipient = upd.Recipient
	}

	// Assign changed values
	res.Kind = upd.Kind
	res.Config = upd.Config
	res.UpdatedAt = now()

	return store.UpdateNotification(ctx, s, res)
}

// onDelete is the custom body for the generated DeleteByID. The generated
// method owns the action-log scaffold + recordAction; the recipient-ownership
// check, soft-delete and websocket delivery live here.
func (svc *notification) onDelete(ctx context.Context, s store.Storer, res *types.Notification, aProps *notificationActionProps) error {
	// Check if the notification belongs to the current user
	currentUserID := intAuth.GetIdentityFromContext(ctx).Identity()
	if res.Recipient != currentUserID {
		return NotificationErrNotFound()
	}

	aProps.setNotification(res)

	res.DeletedAt = now()

	if err := store.UpdateNotification(ctx, s, res); err != nil {
		return err
	}

	// Send the deleted notification via websocket so client can update UI
	if svc.notificationSender != nil {
		if err := svc.notificationSender.Send("notification.delete", res, res.Recipient); err != nil {
			return err
		}
	}

	return nil
}

func (svc *notification) onMarkAsRead(ctx context.Context, aProps *notificationActionProps, ID uint64) (err error) {
	var n *types.Notification

	if ID == 0 {
		return NotificationErrInvalidID()
	}

	if n, err = store.LookupNotificationByID(ctx, svc.store, ID); err != nil {
		return NotificationErrNotFound()
	}

	// Check if the notification belongs to the current user
	currentUserID := intAuth.GetIdentityFromContext(ctx).Identity()
	if n.Recipient != currentUserID {
		return NotificationErrNotAllowedToRead()
	}

	aProps.setNotification(n)

	// Mark as read
	now := time.Now()
	n.ReadAt = &now
	n.UpdatedAt = &now

	if err = store.UpdateNotification(ctx, svc.store, n); err != nil {
		return err
	}

	// Send the updated notification via websocket so client can update UI
	if svc.notificationSender != nil {
		if err = svc.notificationSender.Send("notification.read", n, n.Recipient); err != nil {
			return err
		}
	}

	return nil
}

func (svc *notification) onMarkAsUnread(ctx context.Context, aProps *notificationActionProps, ID uint64) (err error) {
	var n *types.Notification

	if ID == 0 {
		return NotificationErrInvalidID()
	}

	if n, err = store.LookupNotificationByID(ctx, svc.store, ID); err != nil {
		return NotificationErrNotFound()
	}

	// Check if the notification belongs to the current user
	currentUserID := intAuth.GetIdentityFromContext(ctx).Identity()
	if n.Recipient != currentUserID {
		return NotificationErrNotAllowedToRead()
	}

	aProps.setNotification(n)

	// Mark as unread
	n.ReadAt = nil
	now := time.Now()
	n.UpdatedAt = &now

	if err = store.UpdateNotification(ctx, svc.store, n); err != nil {
		return err
	}

	// Send the updated notification via websocket so client can update UI
	if svc.notificationSender != nil {
		if err = svc.notificationSender.Send("notification.unread", n, n.Recipient); err != nil {
			return err
		}
	}

	return nil
}

func (svc *notification) onMarkAllAsRead(ctx context.Context, aProps *notificationActionProps) (err error) {
	var (
		currentUserID = intAuth.GetIdentityFromContext(ctx).Identity()

		cursor *filter.PagingCursor
		nn     types.NotificationSet
		f      types.NotificationFilter
		now    = time.Now()
	)

	aProps.setNotification(&types.Notification{Recipient: currentUserID})

	return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
		// Process unread notifications in batches
		for {
			// Query unread notifications with pagination
			f = types.NotificationFilter{
				Recipient: currentUserID,
				Read:      filter.StateExcluded,
				Paging: filter.Paging{
					Limit:      100,
					PageCursor: cursor,
				},
			}

			nn, f, err = store.SearchNotifications(ctx, svc.store, f)
			if err != nil {
				return err
			}

			// Mark all notifications in this batch as read
			for _, n := range nn {
				n.ReadAt = &now
				n.UpdatedAt = &now
			}

			err = store.UpdateNotification(ctx, svc.store, nn...)
			if err != nil {
				return err
			}

			// Send the updated notifications via websocket so client can update UI
			if svc.notificationSender != nil && len(nn) > 0 {
				if err = svc.notificationSender.Send("notification.read.all", nn, currentUserID); err != nil {
					return err
				}
			}

			// Update cursor for next page or break if no more pages
			cursor = f.PageCursor
			if cursor == nil {
				break
			}
		}

		return nil
	})
}

func (svc *notification) onMarkAllAsUnread(ctx context.Context, aProps *notificationActionProps) (err error) {
	var (
		currentUserID = intAuth.GetIdentityFromContext(ctx).Identity()

		cursor *filter.PagingCursor
		nn     types.NotificationSet
		f      types.NotificationFilter
		now    = time.Now()
	)

	aProps.setNotification(&types.Notification{Recipient: currentUserID})

	return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
		for {
			f = types.NotificationFilter{
				Recipient: currentUserID,
				Read:      filter.StateExclusive, // only read notifications
				Paging: filter.Paging{
					Limit:      100,
					PageCursor: cursor,
				},
			}

			nn, f, err = store.SearchNotifications(ctx, svc.store, f)
			if err != nil {
				return err
			}

			for _, n := range nn {
				n.ReadAt = nil
				n.UpdatedAt = &now
			}

			if err = store.UpdateNotification(ctx, svc.store, nn...); err != nil {
				return err
			}

			if svc.notificationSender != nil && len(nn) > 0 {
				if err = svc.notificationSender.Send("notification.unread.all", nn, currentUserID); err != nil {
					return err
				}
			}

			cursor = f.PageCursor
			if cursor == nil {
				break
			}
		}

		return nil
	})
}

func (svc *notification) checkAssignee(ctx context.Context, n *types.Notification) (err error) {
	// Check if user is assigning to someone else
	if n.Recipient != svc.currentUser(ctx) {
		if !svc.ac.CanAssignNotification(ctx) {
			return NotificationErrNotAllowedToAssign()
		}
	}

	return nil
}

func (svc *notification) currentUser(ctx context.Context) uint64 {
	return intAuth.GetIdentityFromContext(ctx).Identity()
}
