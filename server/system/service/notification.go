package service

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	intAuth "github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
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
		Find(context.Context, types.NotificationFilter) (types.NotificationSet, types.NotificationFilter, error)
		FindByID(context.Context, uint64) (*types.Notification, error)

		Create(context.Context, *types.Notification) (*types.Notification, error)
		Update(context.Context, *types.Notification) (*types.Notification, error)
		Delete(context.Context, uint64) error

		MarkAsRead(context.Context, uint64) error
		MarkAllAsRead(context.Context) error
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

func (svc notification) Find(ctx context.Context, filter types.NotificationFilter) (nn types.NotificationSet, f types.NotificationFilter, err error) {
	var (
		raProps = &notificationActionProps{filter: &filter}
	)

	err = func() (err error) {
		// Get current user ID
		currentUserID := intAuth.GetIdentityFromContext(ctx).Identity()

		// Override recipient filter to ensure users can only see their own notifications
		filter.Recipient = currentUserID

		nn, f, err = store.SearchNotifications(ctx, svc.store, filter)
		if err != nil {
			return err
		}

		// Ensure we never return nil, but an empty slice instead
		if nn == nil {
			nn = make(types.NotificationSet, 0)
		}

		return nil
	}()

	return nn, f, svc.recordAction(ctx, raProps, NotificationActionSearch, err)
}

func (svc notification) FindByID(ctx context.Context, ID uint64) (n *types.Notification, err error) {
	var (
		raProps = &notificationActionProps{notification: &types.Notification{ID: ID}}
	)

	err = func() (err error) {
		if ID == 0 {
			return NotificationErrInvalidID()
		}

		n, err = store.LookupNotificationByID(ctx, svc.store, ID)
		if err != nil {
			return err
		}

		// Check if the notification belongs to the current user
		currentUserID := intAuth.GetIdentityFromContext(ctx).Identity()
		if n.Recipient != currentUserID {
			return NotificationErrNotFound()
		}

		raProps.setNotification(n)

		return nil
	}()

	return n, svc.recordAction(ctx, raProps, NotificationActionLookup, err)
}

func (svc notification) Create(ctx context.Context, new *types.Notification) (n *types.Notification, err error) {
	var (
		raProps = &notificationActionProps{new: new}
	)

	err = func() (err error) {
		// Check if current user has permission to assign notifications to others
		if err := svc.checkAssignee(ctx, new); err != nil {
			return err
		}

		n = new
		n.ID = nextID()
		n.CreatedAt = *now()

		// Set creator
		currentUserID := intAuth.GetIdentityFromContext(ctx).Identity()
		n.CreatedBy = currentUserID

		if err = store.CreateNotification(ctx, svc.store, n); err != nil {
			return err
		}

		// Send the notification via websocket
		if svc.notificationSender != nil {
			// Send only to the recipient
			_ = svc.notificationSender.Send("notification", n, n.Recipient)
		}

		return nil
	}()

	return n, svc.recordAction(ctx, raProps, NotificationActionCreate, err)
}

func (svc notification) Update(ctx context.Context, upd *types.Notification) (n *types.Notification, err error) {
	var (
		raProps = &notificationActionProps{updated: upd}
	)

	err = func() (err error) {
		if upd.ID == 0 {
			return NotificationErrInvalidID()
		}

		if n, err = store.LookupNotificationByID(ctx, svc.store, upd.ID); err != nil {
			return
		}

		// Check if the notification belongs to the current user
		currentUserID := intAuth.GetIdentityFromContext(ctx).Identity()
		if n.Recipient != currentUserID {
			return NotificationErrNotFound()
		}

		// Check if the recipient is being changed and if so, check permissions
		if upd.Recipient != 0 && upd.Recipient != n.Recipient {
			// Create a temporary notification for the permission check
			tempNotification := &types.Notification{
				Recipient: upd.Recipient,
			}

			if err := svc.checkAssignee(ctx, tempNotification); err != nil {
				return err
			}

			n.Recipient = upd.Recipient
		}

		// Assign changed values
		n.Kind = upd.Kind
		n.Config = upd.Config
		n.UpdatedAt = now()

		if err = store.UpdateNotification(ctx, svc.store, n); err != nil {
			return err
		}

		return nil
	}()

	return n, svc.recordAction(ctx, raProps, NotificationActionUpdate, err)
}

func (svc notification) Delete(ctx context.Context, ID uint64) (err error) {
	var (
		n *types.Notification

		raProps = &notificationActionProps{notification: &types.Notification{ID: ID}}
	)

	err = func() (err error) {
		if ID == 0 {
			return NotificationErrInvalidID()
		}

		if n, err = svc.FindByID(ctx, ID); err != nil {
			return NotificationErrNotFound()
		}

		// Check if the notification belongs to the current user
		currentUserID := intAuth.GetIdentityFromContext(ctx).Identity()
		if n.Recipient != currentUserID {
			return NotificationErrNotFound()
		}

		raProps.setNotification(n)

		n.DeletedAt = now()

		return store.UpdateNotification(ctx, svc.store, n)
	}()

	return svc.recordAction(ctx, raProps, NotificationActionDelete, err)
}

func (svc notification) MarkAsRead(ctx context.Context, ID uint64) (err error) {
	var (
		n *types.Notification

		raProps = &notificationActionProps{notification: &types.Notification{ID: ID}}
	)

	err = func() (err error) {
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

		raProps.setNotification(n)

		// Mark as read
		now := time.Now()
		n.ReadAt = &now
		n.UpdatedAt = &now

		if err = store.UpdateNotification(ctx, svc.store, n); err != nil {
			return err
		}

		return nil
	}()

	return svc.recordAction(ctx, raProps, NotificationActionMarkAsRead, err)
}

func (svc notification) MarkAllAsRead(ctx context.Context) (err error) {
	var (
		currentUserID = intAuth.GetIdentityFromContext(ctx).Identity()
		raProps       = &notificationActionProps{notification: &types.Notification{Recipient: currentUserID}}
	)

	err = func() (err error) {
		var (
			cursor *filter.PagingCursor
			nn     types.NotificationSet
			f      types.NotificationFilter
			now    = time.Now()
		)

		return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
			// Process unread notifications in batches
			for {
				// Query unread notifications with pagination
				f = types.NotificationFilter{
					Recipient: currentUserID,
					Read:      filter.StateExcluded, // Only unread notifications
					Paging: filter.Paging{
						Limit:      100, // Process in batches of 100
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

				// Update cursor for next page or break if no more pages
				cursor = f.PageCursor
				if cursor == nil {
					break
				}
			}

			return nil
		})
	}()

	return svc.recordAction(ctx, raProps, NotificationActionMarkAllAsRead, err)
}

func (svc notification) checkAssignee(ctx context.Context, n *types.Notification) (err error) {
	// Check if user is assigning to someone else
	if svc.checkAssignTo(ctx, n) {
		if !svc.ac.CanAssignNotification(ctx) {
			return NotificationErrNotAllowedToAssign()
		}
	}

	return nil
}

// checkAssignTo compares current user with notification.Recipient and return bool
func (svc notification) checkAssignTo(ctx context.Context, n *types.Notification) (valid bool) {
	return n.Recipient != svc.currentUser(ctx)
}

func (svc notification) currentUser(ctx context.Context) uint64 {
	return intAuth.GetIdentityFromContext(ctx).Identity()
}
