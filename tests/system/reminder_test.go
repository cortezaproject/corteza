package system

//
//import (
//	"context"
//	"fmt"
//	"net/http"
//	"testing"
//	"time"
//
//	"github.com/cortezaproject/corteza-server/pkg/auth"
//	"github.com/cortezaproject/corteza-server/system/repository"
//	"github.com/cortezaproject/corteza-server/system/types"
//	"github.com/cortezaproject/corteza-server/tests/helpers"
//	st "github.com/jmoiron/sqlx/types"
//	jsonpath "github.com/steinfletcher/apitest-jsonpath"
//)
//
//func (h helper) repoReminder() repository.ReminderRepository {
//	return repository.Reminder(context.Background(), db())
//}
//
//func (h helper) tPtr(tt time.Time) *time.Time {
//	return &tt
//}
//
//func (h helper) repoMakeReminder(resource string, payload st.JSONText, assignedTo uint64, remindAt *time.Time) *types.Reminder {
//	a, err := h.
//		repoReminder().
//		Create(&types.Reminder{
//			Resource:   resource,
//			Payload:    payload,
//			AssignedAt: time.Now(),
//			AssignedBy: auth.GetIdentityFromContext(h.secCtx()).Identity(),
//			AssignedTo: assignedTo,
//			RemindAt:   remindAt,
//		})
//
//	h.a.NoError(err)
//
//	return a
//}
//
//func (h helper) repoDismissReminder(rr *types.Reminder) {
//	tt := time.Now()
//	rr.DismissedAt = &tt
//	_, err := h.
//		repoReminder().
//		Update(rr)
//
//	h.a.NoError(err)
//}
//
//func TestReminderList(t *testing.T) {
//	h := newHelper(t)
//
//	h.repoMakeReminder("test_lst_namespace:*", nil, 0, nil)
//
//	h.apiInit().
//		Get("/reminder/").
//		Expect(t).
//		Status(http.StatusOK).
//		Assert(helpers.AssertNoErrors).
//		Assert(jsonpath.Present(`$.response.set[? @.resource=="test_lst_namespace:*"]`)).
//		End()
//}
//
//func TestReminderListSpecificResource(t *testing.T) {
//	h := newHelper(t)
//
//	h.repoMakeReminder("test_ls_e_namespace:*", nil, 0, nil)
//	h.repoMakeReminder("test_ls_ne_namespace:*", nil, 0, nil)
//
//	h.apiInit().
//		Get("/reminder/").
//		QueryParams(map[string]string{
//			"resource": "test_ls_e_namespace:*",
//		}).
//		Expect(t).
//		Status(http.StatusOK).
//		Assert(helpers.AssertNoErrors).
//		Assert(jsonpath.Present(`$.response.set[? @.resource=="test_ls_e_namespace:*"]`)).
//		Assert(jsonpath.NotPresent(`$.response.set[? @.resource=="test_ls_ne_namespace:*"]`)).
//		End()
//}
//
//func TestReminderListAssignee(t *testing.T) {
//	h := newHelper(t)
//
//	h.repoMakeReminder("test_la_yes_namespace:*", nil, 111, nil)
//	h.repoMakeReminder("test_la_no_namespace:*", nil, 222, nil)
//
//	h.apiInit().
//		Get("/reminder/").
//		QueryParams(map[string]string{
//			"assignedTo": "111",
//		}).
//		Expect(t).
//		Status(http.StatusOK).
//		Assert(helpers.AssertNoErrors).
//		Assert(jsonpath.Present(`$.response.set[? @.resource=="test_la_yes_namespace:*"]`)).
//		Assert(jsonpath.NotPresent(`$.response.set[? @.resource=="test_la_no_namespace:*"]`)).
//		End()
//}
//
//func TestReminderListDateRange(t *testing.T) {
//	h := newHelper(t)
//
//	nn := time.Now().Round(time.Millisecond).Round(time.Second)
//
//	// Out of range
//	h.repoMakeReminder("test_lrn_ol_namespace:*", nil, 0, h.tPtr(nn.Add(time.Hour*-2)))
//	h.repoMakeReminder("test_lrn_or_namespace:*", nil, 0, h.tPtr(nn.Add(time.Hour*2)))
//
//	// In range
//	h.repoMakeReminder("test_lrn_el_namespace:*", nil, 0, h.tPtr(nn.Add(time.Hour*-1)))
//	h.repoMakeReminder("test_lrn_cc_namespace:*", nil, 0, h.tPtr(nn))
//	h.repoMakeReminder("test_lrn_er_namespace:*", nil, 0, h.tPtr(nn.Add(time.Hour*1)))
//
//	h.apiInit().
//		Get("/reminder/").
//		QueryParams(map[string]string{
//			"scheduledFrom":  nn.Add(time.Hour * -1).Format(time.RFC3339),
//			"scheduledUntil": nn.Add(time.Hour * 1).Format(time.RFC3339),
//		}).
//		Expect(t).
//		Status(http.StatusOK).
//		Assert(helpers.AssertNoErrors).
//		Assert(jsonpath.Present(`$.response.set[? @.resource=="test_lrn_el_namespace:*"]`)).
//		Assert(jsonpath.Present(`$.response.set[? @.resource=="test_lrn_cc_namespace:*"]`)).
//		Assert(jsonpath.Present(`$.response.set[? @.resource=="test_lrn_er_namespace:*"]`)).
//		Assert(jsonpath.NotPresent(`$.response.set[? @.resource=="test_lrn_ol_namespace:*"]`)).
//		Assert(jsonpath.NotPresent(`$.response.set[? @.resource=="test_lrn_or_namespace:*"]`)).
//		End()
//}
//
//func TestReminderListOnlyScheduled(t *testing.T) {
//	h := newHelper(t)
//
//	tt := time.Now()
//	h.repoMakeReminder("test_lsch_yes_namespace:*", nil, 111, &tt)
//	h.repoMakeReminder("test_lsch_no_namespace:*", nil, 222, nil)
//
//	h.apiInit().
//		Get("/reminder/").
//		QueryParams(map[string]string{
//			"scheduledOnly": "true",
//		}).
//		Expect(t).
//		Status(http.StatusOK).
//		Assert(helpers.AssertNoErrors).
//		Assert(jsonpath.Present(`$.response.set[? @.resource=="test_lsch_yes_namespace:*"]`)).
//		Assert(jsonpath.NotPresent(`$.response.set[? @.resource=="test_lsch_no_namespace:*"]`)).
//		End()
//}
//
//func TestReminderListExcludeDismissed(t *testing.T) {
//	h := newHelper(t)
//
//	h.repoMakeReminder("test_ldsm_yes_namespace:*", nil, 0, nil)
//
//	rr := h.repoMakeReminder("test_ldsm_no_namespace:*", nil, 0, nil)
//	h.repoDismissReminder(rr)
//
//	h.apiInit().
//		Get("/reminder/").
//		QueryParams(map[string]string{
//			"excludeDismissed": "true",
//		}).
//		Expect(t).
//		Status(http.StatusOK).
//		Assert(helpers.AssertNoErrors).
//		Assert(jsonpath.Present(`$.response.set[? @.resource=="test_ldsm_yes_namespace:*"]`)).
//		Assert(jsonpath.NotPresent(`$.response.set[? @.resource=="test_ldsm_no_namespace:*"]`)).
//		End()
//}
//
//func TestReminderCreate(t *testing.T) {
//	h := newHelper(t)
//
//	h.apiInit().
//		Post("/reminder/").
//		JSON(fmt.Sprintf(`{"resource":"test_c_namespace:*","assignedTo":"%d","payload":{}}`, h.cUser.Identity())).
//		Expect(t).
//		Status(http.StatusOK).
//		Assert(helpers.AssertNoErrors).
//		Assert(jsonpath.Present(`$.response.resource=="test_c_namespace:*"`)).
//		End()
//}
//
//func TestReminderCreateSetAssignee(t *testing.T) {
//	h := newHelper(t)
//
//	h.allow(types.SystemRBACResource, "reminder.assign")
//
//	h.apiInit().
//		Post("/reminder/").
//		JSON(`{"resource":"test_c_as_namespace:*","assignedTo":"222","payload":{}}`).
//		Expect(t).
//		Status(http.StatusOK).
//		Assert(helpers.AssertNoErrors).
//		Assert(jsonpath.Present(`$.response.resource=="test_c_as_namespace:*"`)).
//		End()
//}
//
//func TestReminderCreateSetAssigneeForbidden(t *testing.T) {
//	h := newHelper(t)
//
//	h.deny(types.SystemRBACResource, "reminder.assign")
//
//	h.apiInit().
//		Post("/reminder/").
//		JSON(`{"resource":"test_rf_as_namespace:*","assignedTo":"222","payload":{}}`).
//		Expect(t).
//		Status(http.StatusOK).
//		Assert(helpers.AssertError("not allowed to assign reminders to other users")).
//		End()
//}
//
//func TestReminderUpdate(t *testing.T) {
//	h := newHelper(t)
//
//	r := h.repoMakeReminder("test_upd_namespace:*", nil, 0, nil)
//
//	h.apiInit().
//		Put(fmt.Sprintf("/reminder/%d", r.ID)).
//		JSON(fmt.Sprintf(`{"resource":"test_upd_u_namespace:*","assignedTo":"%d","payload":{}}`, h.cUser.Identity())).
//		Expect(t).
//		Status(http.StatusOK).
//		Assert(helpers.AssertNoErrors).
//		Assert(jsonpath.Present(`$.response.resource=="test_upd_u_namespace:*"`)).
//		Assert(jsonpath.NotPresent(`$.response.resource=="test_upd_namespace:*"`)).
//		End()
//
//	h.apiInit().
//		Get(fmt.Sprintf("/reminder/%d", r.ID)).
//		Expect(t).
//		Status(http.StatusOK).
//		Assert(helpers.AssertNoErrors).
//		Assert(jsonpath.Present(`$.response.resource=="test_upd_u_namespace:*"`)).
//		Assert(jsonpath.NotPresent(`$.response.resource=="test_upd_namespace:*"`)).
//		End()
//}
//
//func TestReminderUpdateSetAssignee(t *testing.T) {
//	h := newHelper(t)
//
//	h.allow(types.SystemRBACResource, "reminder.assign")
//	r := h.repoMakeReminder("test_upd_as_namespace:*", nil, 0, nil)
//
//	h.apiInit().
//		Put(fmt.Sprintf("/reminder/%d", r.ID)).
//		JSON(`{"resource":"test_upd_as_namespace:*","assignedTo":"111","payload":{}}`).
//		Expect(t).
//		Status(http.StatusOK).
//		Assert(helpers.AssertNoErrors).
//		Assert(jsonpath.Present(`$.response.resource=="test_upd_as_namespace:*"`)).
//		Assert(jsonpath.Present(`$.response.assignedTo=="111"`)).
//		End()
//
//	h.apiInit().
//		Get(fmt.Sprintf("/reminder/%d", r.ID)).
//		Expect(t).
//		Status(http.StatusOK).
//		Assert(helpers.AssertNoErrors).
//		Assert(jsonpath.Present(`$.response.resource=="test_upd_as_namespace:*"`)).
//		Assert(jsonpath.Present(`$.response.assignedTo=="111"`)).
//		End()
//}
//
//func TestReminderUpdateSetAssigneeForbidden(t *testing.T) {
//	h := newHelper(t)
//
//	h.deny(types.SystemRBACResource, "reminder.assign")
//	r := h.repoMakeReminder("test_updf_as_namespace:*", nil, 0, nil)
//
//	h.apiInit().
//		Put(fmt.Sprintf("/reminder/%d", r.ID)).
//		JSON(`{"resource":"test_updf_as_namespace:*","assignedTo":"111","payload":{}}`).
//		Expect(t).
//		Status(http.StatusOK).
//		Assert(helpers.AssertError("not allowed to assign reminders to other users")).
//		End()
//}
//
//func TestReminderRead(t *testing.T) {
//	h := newHelper(t)
//
//	r := h.repoMakeReminder("test_r_namespace:*", nil, 0, nil)
//
//	h.apiInit().
//		Get(fmt.Sprintf("/reminder/%d", r.ID)).
//		Expect(t).
//		Status(http.StatusOK).
//		Assert(helpers.AssertNoErrors).
//		Assert(jsonpath.Present(`$.response.resource=="test_r_namespace:*"`)).
//		End()
//}
//
//func TestReminderDelete(t *testing.T) {
//	h := newHelper(t)
//
//	r := h.repoMakeReminder("test_d_namespace:*", nil, 0, nil)
//
//	h.apiInit().
//		Delete(fmt.Sprintf("/reminder/%d", r.ID)).
//		Expect(t).
//		Status(http.StatusOK).
//		Assert(helpers.AssertNoErrors).
//		Assert(jsonpath.Present(`$.success.message=="OK"`)).
//		End()
//}
//
//func TestReminderDismiss(t *testing.T) {
//	h := newHelper(t)
//
//	r := h.repoMakeReminder("test_dsm_namespace:*", nil, 0, nil)
//
//	h.apiInit().
//		Patch(fmt.Sprintf("/reminder/%d/dismiss", r.ID)).
//		Expect(t).
//		Status(http.StatusOK).
//		Assert(jsonpath.Present(`$.success.message=="OK"`)).
//		End()
//
//	h.apiInit().
//		Get(fmt.Sprintf("/reminder/%d", r.ID)).
//		Expect(t).
//		Status(http.StatusOK).
//		Assert(jsonpath.Present(`$.response.dismissedAt!="null"`)).
//		End()
//}
//
//func TestReminderSnooze(t *testing.T) {
//	h := newHelper(t)
//
//	r := h.repoMakeReminder("test_snz_namespace:*", nil, 0, nil)
//
//	h.apiInit().
//		Patch(fmt.Sprintf("/reminder/%d/snooze", r.ID)).
//		JSON(`{"remindAt":"2001-01-01T01:00:00Z"}`).
//		Expect(t).
//		Status(http.StatusOK).
//		Assert(jsonpath.Present(`$.success.message=="OK"`)).
//		End()
//
//	h.apiInit().
//		Get(fmt.Sprintf("/reminder/%d", r.ID)).
//		Expect(t).
//		Status(http.StatusOK).
//		Assert(jsonpath.Present(`$.response.snoozeCount==1`)).
//		Assert(jsonpath.Present(`$.response.remindAt!="null"`)).
//		End()
//}
