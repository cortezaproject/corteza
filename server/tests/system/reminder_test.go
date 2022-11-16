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

func (h helper) clearReminders() {
	h.noError(store.TruncateReminders(context.Background(), service.DefaultStore))
}

func (h helper) makeReminder() *types.Reminder {
	return h.reminder(h.cUser.ID, false)
}

func (h helper) makeReminderByUserID(userID uint64) *types.Reminder {
	return h.reminder(userID, false)
}

func (h helper) makeDeletedReminder() *types.Reminder {
	return h.reminder(h.cUser.ID, true)
}

func (h helper) reminder(userID uint64, deleted bool) *types.Reminder {
	rm := &types.Reminder{Resource: "test:resource", AssignedTo: userID}
	rm.ID = id.Next()
	rm.CreatedAt = time.Now()
	if deleted {
		now := time.Now()
		rm.DeletedAt = &now
	}

	h.noError(store.CreateReminder(context.Background(), service.DefaultStore, rm))
	return rm
}

func (h helper) lookupReminderByID(ID uint64) *types.Reminder {
	rm, err := store.LookupReminderByID(context.Background(), service.DefaultStore, ID)
	h.noError(err)
	return rm
}

func TestReminderCreate(t *testing.T) {
	h := newHelper(t)
	h.clearReminders()

	h.apiInit().
		Post("/reminder/").
		Header("Accept", "application/json").
		FormData("resource", "some:resource").
		FormData("assignedTo", strconv.FormatUint(h.cUser.ID, 10)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.resource`, "some:resource")).
		End()
}

func TestReminderAssign_forbidden(t *testing.T) {
	h := newHelper(t)
	h.clearReminders()

	h.apiInit().
		Post("/reminder/").
		Header("Accept", "application/json").
		FormData("resource", "some:resource").
		FormData("assignedTo", "404").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("reminder.errors.notAllowedToAssign")).
		End()
}

func TestReminderAssign(t *testing.T) {
	h := newHelper(t)
	h.clearReminders()

	helpers.AllowMe(h, types.ComponentRbacResource(), "reminder.assign")

	h.apiInit().
		Post("/reminder/").
		Header("Accept", "application/json").
		FormData("resource", "some:resource").
		FormData("assignedTo", "404").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.resource`, "some:resource")).
		End()
}

func TestReminderRead(t *testing.T) {
	h := newHelper(t)
	h.clearReminders()

	rm := h.makeReminder()

	h.apiInit().
		Get(fmt.Sprintf("/reminder/%d", rm.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.resource`, rm.Resource)).
		End()
}

// TestReminderReadForbidden checks only user themself can read reminder assigned to them
func TestReminderReadForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearReminders()

	rm := h.makeReminderByUserID(id.Next())

	h.apiInit().
		Get(fmt.Sprintf("/reminder/%d", rm.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("reminder.errors.notAllowedToRead")).
		End()
}

func TestReminderList(t *testing.T) {
	h := newHelper(t)
	h.clearReminders()

	h.makeReminder()
	h.makeDeletedReminder()
	h.makeDeletedReminder()

	h.apiInit().
		Get("/reminder/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len("$.response.set", 1)).
		End()
}

func TestReminderListIncludeDeleted(t *testing.T) {
	h := newHelper(t)
	h.clearReminders()

	h.makeReminder()
	h.makeDeletedReminder()
	h.makeDeletedReminder()

	h.apiInit().
		Get("/reminder/").
		Query("includeDeleted", "true").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len("$.response.set", 3)).
		End()
}

func TestReminderUpdateForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearReminders()

	rm := h.makeReminder()

	h.apiInit().
		Put(fmt.Sprintf("/reminder/%d", rm.ID)).
		Header("Accept", "application/json").
		FormData("resource", "changed:resource").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("reminder.errors.notAllowedToAssign")).
		End()
}

func TestReminderUpdate(t *testing.T) {
	h := newHelper(t)
	h.clearReminders()

	helpers.AllowMe(h, types.ComponentRbacResource(), "reminder.assign")

	rm := h.makeReminder()

	h.apiInit().
		Put(fmt.Sprintf("/reminder/%d", rm.ID)).
		Header("Accept", "application/json").
		FormData("resource", "changed:resource").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.resource`, "changed:resource")).
		End()
}

func TestReminderDelete(t *testing.T) {
	h := newHelper(t)
	h.clearReminders()

	rm := h.makeReminder()

	h.apiInit().
		Delete(fmt.Sprintf("/reminder/%d", rm.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestReminderDismiss(t *testing.T) {
	h := newHelper(t)
	h.clearReminders()

	rm := h.makeReminder()

	h.apiInit().
		Patch(fmt.Sprintf("/reminder/%d/dismiss", rm.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

// TestReminderDismissForbidden checks only user themself can dismiss reminder assigned to them
func TestReminderDismissForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearReminders()

	rm := h.makeReminderByUserID(id.Next())

	h.apiInit().
		Patch(fmt.Sprintf("/reminder/%d/dismiss", rm.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("reminder.errors.notAllowedToDismiss")).
		End()
}

func TestReminderSnooze(t *testing.T) {
	h := newHelper(t)
	h.clearReminders()

	rm := h.makeReminder()

	h.apiInit().
		Patch(fmt.Sprintf("/reminder/%d/snooze", rm.ID)).
		Header("Accept", "application/json").
		FormData("remindAt", time.Now().Format(time.RFC3339)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}
