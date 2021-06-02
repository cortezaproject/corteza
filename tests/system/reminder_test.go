package system

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func (h helper) clearReminders() {
	h.noError(store.TruncateReminders(context.Background(), service.DefaultStore))
}

func (h helper) makeReminder() *types.Reminder {
	rm := &types.Reminder{Resource: "test:resource", AssignedTo: h.cUser.ID}
	rm.ID = id.Next()
	rm.CreatedAt = time.Now()
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
		Assert(helpers.AssertError("not allowed to assign reminders to other users")).
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

func TestReminderList(t *testing.T) {
	h := newHelper(t)
	h.clearReminders()

	h.makeReminder()
	h.makeReminder()

	h.apiInit().
		Get("/reminder/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len("$.response.set", 2)).
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
		Assert(helpers.AssertError("not allowed to assign reminders to other users")).
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
