package system

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/tests/helpers"
)

// Reminders of other users can not be taken over by assigning them to yourself
func TestReminderUpdateForeign(t *testing.T) {
	h := newHelper(t)
	h.clearReminders()

	victimID := id.Next()
	rm := h.makeReminderByUserID(victimID)

	h.apiInit().
		Put(fmt.Sprintf("/reminder/%d", rm.ID)).
		Header("Accept", "application/json").
		FormData("resource", "hijacked:resource").
		FormData("assignedTo", strconv.FormatUint(h.cUser.ID, 10)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("reminder.errors.notAllowedToRead")).
		End()

	stored := h.lookupReminderByID(rm.ID)
	h.a.Equal(victimID, stored.AssignedTo)
	h.a.Equal("test:resource", stored.Resource)
}

func TestReminderUpdateOwn(t *testing.T) {
	h := newHelper(t)
	h.clearReminders()
	rm := h.makeReminder()

	h.apiInit().
		Put(fmt.Sprintf("/reminder/%d", rm.ID)).
		Header("Accept", "application/json").
		FormData("resource", "changed:resource").
		FormData("assignedTo", strconv.FormatUint(h.cUser.ID, 10)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	h.a.Equal("changed:resource", h.lookupReminderByID(rm.ID).Resource)
}
