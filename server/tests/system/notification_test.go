package system

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/service"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/cortezaproject/corteza/server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func (h helper) clearNotifications() {
	h.noError(store.TruncateNotifications(context.Background(), service.DefaultStore))
}

func (h helper) makeNotification() *types.Notification {
	return h.notification(h.cUser.ID, types.NotificationKindSimple, false)
}

func (h helper) makeNotificationByUserID(userID uint64) *types.Notification {
	return h.notification(userID, types.NotificationKindSimple, false)
}

func (h helper) makeReadNotification() *types.Notification {
	ntf := h.makeNotification()
	now := time.Now()
	ntf.ReadAt = &now
	h.noError(store.UpdateNotification(context.Background(), service.DefaultStore, ntf))
	return ntf
}

func (h helper) makeDeletedNotification() *types.Notification {
	return h.notification(h.cUser.ID, types.NotificationKindSimple, true)
}

func (h helper) notification(userID uint64, notificationKind types.NotificationKind, deleted bool) *types.Notification {
	ntf := &types.Notification{
		Kind: notificationKind,
		Config: types.NotificationConfig{
			Simple: types.SimpleNotificationConfig{
				Title:       "Test notification",
				Description: "This is a test notification",
			},
		},
		Recipient: userID,
		CreatedBy: h.cUser.ID,
	}
	ntf.ID = id.Next()
	ntf.CreatedAt = time.Now()
	if deleted {
		now := time.Now()
		ntf.DeletedAt = &now
	}

	h.noError(store.CreateNotification(context.Background(), service.DefaultStore, ntf))
	return ntf
}

func (h helper) lookupNotificationByID(ID uint64) *types.Notification {
	ntf, err := store.LookupNotificationByID(context.Background(), service.DefaultStore, ID)
	h.noError(err)
	return ntf
}

func TestNotificationCreate(t *testing.T) {
	h := newHelper(t)
	h.clearNotifications()

	ntfTemplate := h.makeNotification()

	h.apiInit().
		Post("/notification/").
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		JSON(map[string]interface{}{
			"kind":      string(ntfTemplate.Kind),
			"config":    ntfTemplate.Config,
			"recipient": strconv.FormatUint(ntfTemplate.Recipient, 10),
		}).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.kind`, string(ntfTemplate.Kind))).
		Assert(jsonpath.Equal(`$.response.config.simple.title`, ntfTemplate.Config.Simple.Title)).
		Assert(jsonpath.Equal(`$.response.config.simple.description`, ntfTemplate.Config.Simple.Description)).
		End()
}

func TestNotificationRead(t *testing.T) {
	h := newHelper(t)
	h.clearNotifications()

	ntf := h.makeNotification()

	h.apiInit().
		Get(fmt.Sprintf("/notification/%d", ntf.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.kind`, string(ntf.Kind))).
		Assert(jsonpath.Equal(`$.response.config.simple.title`, ntf.Config.Simple.Title)).
		End()
}

func TestNotificationReadForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearNotifications()

	ntf := h.makeNotificationByUserID(id.Next())

	h.apiInit().
		Get(fmt.Sprintf("/notification/%d", ntf.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("notification.errors.notFound")).
		End()
}

func TestNotificationList(t *testing.T) {
	h := newHelper(t)
	h.clearNotifications()

	// Create 3 notifications: 2 for current user (1 read, 1 unread) and 1 for another user
	h.makeNotification()
	h.makeReadNotification()
	h.makeNotificationByUserID(id.Next())

	h.apiInit().
		Get("/notification/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len("$.response.set", 1)).
		End()
}

func TestNotificationListUnreadOnly(t *testing.T) {
	h := newHelper(t)
	h.clearNotifications()

	// Create one unread notification
	h.makeNotification()

	// Create one read notification
	h.makeReadNotification()

	// Test with read=0 which should return only unread notifications
	h.apiInit().
		Get("/notification/").
		Query("read", "0").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len("$.response.set", 1)).
		End()
}

func TestNotificationListWithRead(t *testing.T) {
	h := newHelper(t)
	h.clearNotifications()

	// Create one unread notification
	h.makeNotification()

	// Create one read notification
	h.makeReadNotification()

	// Test with read=1 which should return both read and unread notifications
	h.apiInit().
		Get("/notification/").
		Query("read", "1").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len("$.response.set", 2)).
		End()
}

func TestNotificationListReadOnly(t *testing.T) {
	h := newHelper(t)
	h.clearNotifications()

	// Create one unread notification
	h.makeNotification()

	// Create one read notification
	h.makeReadNotification()

	// Test with read=2 which should return only read notifications
	h.apiInit().
		Get("/notification/").
		Query("read", "2").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len("$.response.set", 1)).
		End()
}

func TestNotificationListNoDeleted(t *testing.T) {
	h := newHelper(t)
	h.clearNotifications()

	h.makeNotification()
	h.makeDeletedNotification()

	h.apiInit().
		Get("/notification/").
		Query("deleted", "0").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len("$.response.set", 1)).
		End()
}

func TestNotificationListWithDeleted(t *testing.T) {
	h := newHelper(t)
	h.clearNotifications()

	h.makeNotification()
	h.makeDeletedNotification()

	h.apiInit().
		Get("/notification/").
		Query("deleted", "1").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len("$.response.set", 2)).
		End()
}

func TestNotificationListDeletedOnly(t *testing.T) {
	h := newHelper(t)
	h.clearNotifications()

	h.makeNotification()
	h.makeDeletedNotification()

	h.apiInit().
		Get("/notification/").
		Query("deleted", "2").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len("$.response.set", 1)).
		End()
}

func TestNotificationUpdate(t *testing.T) {
	h := newHelper(t)
	h.clearNotifications()

	ntf := h.makeNotification()

	// Create update data
	updatedConfig := types.NotificationConfig{
		Simple: types.SimpleNotificationConfig{
			Title:       "Updated Title",
			Description: "Updated Description",
		},
	}

	h.apiInit().
		Put(fmt.Sprintf("/notification/%d", ntf.ID)).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		JSON(map[string]interface{}{
			"kind":   string(ntf.Kind),
			"config": updatedConfig,
		}).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.config.simple.title`, updatedConfig.Simple.Title)).
		Assert(jsonpath.Equal(`$.response.config.simple.description`, updatedConfig.Simple.Description)).
		End()
}

func TestNotificationUpdateForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearNotifications()

	ntf := h.makeNotificationByUserID(id.Next())

	// Create update data
	updatedConfig := types.NotificationConfig{
		Simple: types.SimpleNotificationConfig{
			Title:       "Updated Title",
			Description: "Updated Description",
		},
	}

	h.apiInit().
		Put(fmt.Sprintf("/notification/%d", ntf.ID)).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		JSON(map[string]interface{}{
			"kind":   string(ntf.Kind),
			"config": updatedConfig,
		}).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("notification.errors.notFound")).
		End()
}

func TestNotificationDelete(t *testing.T) {
	h := newHelper(t)
	h.clearNotifications()

	ntf := h.makeNotification()

	h.apiInit().
		Delete(fmt.Sprintf("/notification/%d", ntf.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	// Verify notification is marked as deleted
	deletedNtf := h.lookupNotificationByID(ntf.ID)
	if deletedNtf.DeletedAt == nil {
		t.Error("notification should be marked as deleted")
	}
}

func TestNotificationDeleteForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearNotifications()

	ntf := h.makeNotificationByUserID(id.Next())

	h.apiInit().
		Delete(fmt.Sprintf("/notification/%d", ntf.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("notification.errors.notFound")).
		End()
}

func TestNotificationMarkAsRead(t *testing.T) {
	h := newHelper(t)
	h.clearNotifications()

	ntf := h.makeNotification()

	h.apiInit().
		Patch(fmt.Sprintf("/notification/%d/read", ntf.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	// Verify notification is marked as read
	readNtf := h.lookupNotificationByID(ntf.ID)
	if readNtf.ReadAt == nil {
		t.Error("notification should be marked as read")
	}
}

func TestNotificationMarkAsReadForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearNotifications()

	ntf := h.makeNotificationByUserID(id.Next())

	h.apiInit().
		Patch(fmt.Sprintf("/notification/%d/read", ntf.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("notification.errors.notAllowedToRead")).
		End()
}

func TestNotificationMarkAllAsRead(t *testing.T) {
	h := newHelper(t)
	h.clearNotifications()

	// Create multiple unread notifications
	h.makeNotification()
	h.makeNotification()
	h.makeNotification()

	h.apiInit().
		Patch("/notification/all/read").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	// List read notifications and verify all are marked as read
	h.apiInit().
		Get("/notification/").
		Query("read", "2").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len("$.response.set", 3)).
		End()
}

func TestNotificationAssign_forbidden(t *testing.T) {
	h := newHelper(t)
	h.clearNotifications()

	// Try to create a notification for another user without the proper permission
	h.apiInit().
		Post("/notification/").
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		JSON(map[string]interface{}{
			"kind": string(types.NotificationKindSimple),
			"config": types.NotificationConfig{
				Simple: types.SimpleNotificationConfig{
					Title:       "Test notification",
					Description: "This is a test notification",
				},
			},
			"recipient": strconv.FormatUint(id.Next(), 10),
		}).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("notification.errors.notAllowedToAssign")).
		End()
}

func TestNotificationAssign(t *testing.T) {
	h := newHelper(t)
	h.clearNotifications()

	// Allow current user to assign notifications
	helpers.AllowMe(h, types.ComponentRbacResource(), "notification.assign")

	otherUserID := id.Next()

	// Create a notification for another user with the permission
	h.apiInit().
		Post("/notification/").
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		JSON(map[string]interface{}{
			"kind": string(types.NotificationKindSimple),
			"config": types.NotificationConfig{
				Simple: types.SimpleNotificationConfig{
					Title:       "Test notification",
					Description: "This is a test notification",
				},
			},
			"recipient": strconv.FormatUint(otherUserID, 10),
		}).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.kind`, string(types.NotificationKindSimple))).
		Assert(jsonpath.Equal(`$.response.recipient`, strconv.FormatUint(otherUserID, 10))).
		End()
}

func TestNotificationUpdateRecipient_forbidden(t *testing.T) {
	h := newHelper(t)
	h.clearNotifications()

	ntf := h.makeNotification()
	otherUserID := id.Next()

	// Try to update a notification with a different recipient without permission
	h.apiInit().
		Put(fmt.Sprintf("/notification/%d", ntf.ID)).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		JSON(map[string]interface{}{
			"kind":      string(ntf.Kind),
			"config":    ntf.Config,
			"recipient": strconv.FormatUint(otherUserID, 10),
		}).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("notification.errors.notAllowedToAssign")).
		End()
}

func TestNotificationUpdateRecipient(t *testing.T) {
	h := newHelper(t)
	h.clearNotifications()

	// Allow current user to assign notifications
	helpers.AllowMe(h, types.ComponentRbacResource(), "notification.assign")

	ntf := h.makeNotification()
	otherUserID := id.Next()

	// Update a notification with a different recipient with permission
	h.apiInit().
		Put(fmt.Sprintf("/notification/%d", ntf.ID)).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		JSON(map[string]interface{}{
			"kind":      string(ntf.Kind),
			"config":    ntf.Config,
			"recipient": strconv.FormatUint(otherUserID, 10),
		}).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.recipient`, strconv.FormatUint(otherUserID, 10))).
		End()
}
