package tests

import (
	"context"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/pkg/rand"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
)

func testNotifications(t *testing.T, s store.Notifications) {
	var (
		ctx = context.Background()

		makeNew = func(nn ...string) *types.Notification {
			// minimum data set for new notification
			thisID := id.Next()
			return &types.Notification{
				ID:   thisID,
				Kind: types.NotificationKindSimple,
				Config: types.NotificationConfig{
					Simple: types.SimpleNotificationConfig{
						Title:       "Test Notification",
						Description: "This is a test notification",
					},
				},
				Recipient: thisID,
				CreatedAt: time.Now(),
				CreatedBy: thisID,
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.Notification) {
			req := require.New(t)
			req.NoError(s.TruncateNotifications(ctx))
			notification := makeNew()
			req.NoError(s.CreateNotification(ctx, notification))
			return req, notification
		}

		truncAndFill = func(t *testing.T, l int) (*require.Assertions, types.NotificationSet) {
			req := require.New(t)
			req.NoError(s.TruncateNotifications(ctx))

			set := make([]*types.Notification, l)

			for i := 0; i < l; i++ {
				set[i] = makeNew(string(rand.Bytes(10)))
			}

			req.NoError(s.CreateNotification(ctx, set...))
			return req, set
		}
	)

	t.Run("create", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.CreateNotification(ctx, makeNew()))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, notification := truncAndCreate(t)
		fetched, err := s.LookupNotificationByID(ctx, notification.ID)
		req.NoError(err)
		req.Equal(notification.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
		req.Nil(fetched.UpdatedAt)
		req.Nil(fetched.DeletedAt)
	})

	t.Run("update", func(t *testing.T) {
		req, notification := truncAndCreate(t)
		req.NoError(s.UpdateNotification(ctx, notification))
	})

	t.Run("search", func(t *testing.T) {
		t.Run("by ID", func(t *testing.T) {
			req, prefill := truncAndFill(t, 5)

			set, f, err := s.SearchNotifications(ctx, types.NotificationFilter{NotificationID: []uint64{prefill[0].ID}})
			req.NoError(err)
			req.Equal([]uint64{prefill[0].ID}, f.NotificationID)
			req.Len(set, 1)
		})

		t.Run("by kind", func(t *testing.T) {
			req, prefill := truncAndFill(t, 5)

			set, _, err := s.SearchNotifications(ctx, types.NotificationFilter{Kind: []types.NotificationKind{prefill[0].Kind}})
			req.NoError(err)
			req.Len(set, 5) // All have the same kind
		})

		t.Run("by recipient", func(t *testing.T) {
			req, prefill := truncAndFill(t, 5)
			set, f, err := s.SearchNotifications(ctx, types.NotificationFilter{Recipient: prefill[0].Recipient})
			req.NoError(err)
			req.Equal(prefill[0].Recipient, f.Recipient)
			req.Len(set, 1)
		})

		t.Run("unread only", func(t *testing.T) {
			req, prefill := truncAndFill(t, 5)

			// Check that all notifications are initially unread
			filter := types.NotificationFilter{
				Unread: filter.StateExclusive, // Only unread notifications
			}
			set, _, err := s.SearchNotifications(ctx, filter)
			req.NoError(err)
			req.Len(set, 5, "Expected all 5 notifications to be unread")

			// Mark one notification as read
			now := time.Now()
			prefill[0].ReadAt = &now
			req.NoError(s.UpdateNotification(ctx, prefill[0]))

			// Search for unread notifications again
			set, _, err = s.SearchNotifications(ctx, filter)
			req.NoError(err)

			// We expect 4 unread notifications (5 total - 1 read)
			req.Len(set, 4, "Expected 4 unread notifications after marking one as read")
		})

		t.Run("deleted only", func(t *testing.T) {
			req, prefill := truncAndFill(t, 5)

			// Check that no notifications are initially deleted
			filter := types.NotificationFilter{
				Deleted: filter.StateExclusive, // Only deleted notifications
			}
			set, _, err := s.SearchNotifications(ctx, filter)
			req.NoError(err)
			req.Len(set, 0, "Expected no deleted notifications")

			// Mark one notification as deleted
			now := time.Now()
			prefill[0].DeletedAt = &now
			req.NoError(s.UpdateNotification(ctx, prefill[0]))

			// Search for deleted notifications
			set, _, err = s.SearchNotifications(ctx, filter)
			req.NoError(err)

			// We expect 1 deleted notification
			req.Len(set, 1, "Expected 1 deleted notification")
		})

		t.Run("with check", func(t *testing.T) {
			req, prefill := truncAndFill(t, 5)
			set, _, err := s.SearchNotifications(ctx, types.NotificationFilter{
				Check: func(notification *types.Notification) (bool, error) {
					// simple check that matches with the first notification from prefill
					return notification.ID == prefill[0].ID, nil
				},
			})
			req.NoError(err)
			req.Len(set, 1)
			req.Equal(prefill[0].ID, set[0].ID)
		})
	})

	t.Run("delete", func(t *testing.T) {
		req, notification := truncAndCreate(t)
		req.NoError(s.DeleteNotificationByID(ctx, notification.ID))

		// Verify it can't be found
		_, err := s.LookupNotificationByID(ctx, notification.ID)
		req.Error(err)
	})
}
