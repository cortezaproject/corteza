package messaging

import (
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/messaging/service"
	tt "github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/settings"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"github.com/jmoiron/sqlx/types"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func TestSettingsList(t *testing.T) {
	h := newHelper(t)
	h.allow(tt.MessagingPermissionResource, "settings.read")
	h.allow(tt.MessagingPermissionResource, "settings.manage")

	err := service.DefaultSettings.With(h.secCtx()).BulkSet(settings.ValueSet{
		&settings.Value{Name: "t_msg_k1.s1", Value: types.JSONText(`"t_msg_v1"`)},
		&settings.Value{Name: "t_msg_k1.s2", Value: types.JSONText(`"t_msg_v2"`)},
	})
	h.a.NoError(err)

	h.apiInit().
		Get("/settings/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(`$.response[? @.name=="t_msg_k1.s1"]`)).
		Assert(jsonpath.Present(`$.response[? @.value=="t_msg_v1"]`)).
		Assert(jsonpath.Present(`$.response[? @.name=="t_msg_k1.s2"]`)).
		Assert(jsonpath.Present(`$.response[? @.value=="t_msg_v2"]`)).
		End()
}

func TestSettingsList_noPermissions(t *testing.T) {
	h := newHelper(t)
	h.deny(tt.MessagingPermissionResource, "settings.read")

	h.apiInit().
		Get("/settings/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to read settings")).
		End()
}

func TestSettingsUpdate(t *testing.T) {
	h := newHelper(t)
	h.allow(tt.MessagingPermissionResource, "settings.manage")
	h.allow(tt.MessagingPermissionResource, "settings.read")

	err := service.DefaultSettings.With(h.secCtx()).BulkSet(settings.ValueSet{
		&settings.Value{Name: "t_msg_k1.s1", Value: types.JSONText(`"t_msg_v1"`)},
	})
	h.a.NoError(err)

	h.apiInit().
		Patch("/settings/").
		JSON(`{"values":[{"name":"t_msg_k1.s1","value":"t_msg_v1_edited"},{"name":"t_msg_k2.s1","value":"t_msg_v2_new"}]}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	s, err := service.DefaultSettings.With(h.secCtx()).Get("t_msg_k1.s1", 0)
	h.a.NoError(err)
	h.a.Equal(`"t_msg_v1_edited"`, s.Value.String(), "existing key should be updated")

	s, err = service.DefaultSettings.With(h.secCtx()).Get("t_msg_k2.s1", 0)
	h.a.NoError(err)
	h.a.Equal(`"t_msg_v2_new"`, s.Value.String(), "new key should be added")
}

func TestSettingsUpdate_noPermissions(t *testing.T) {
	h := newHelper(t)
	h.deny(tt.MessagingPermissionResource, "settings.manage")

	h.apiInit().
		Patch("/settings/").
		JSON(`{"values":[]}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to manage settings")).
		End()
}

func TestSettingsGet(t *testing.T) {
	h := newHelper(t)
	h.allow(tt.MessagingPermissionResource, "settings.read")
	h.allow(tt.MessagingPermissionResource, "settings.manage")

	err := service.DefaultSettings.With(h.secCtx()).BulkSet(settings.ValueSet{
		&settings.Value{Name: "t_msg_k1.s1", Value: types.JSONText(`"t_msg_v1"`)},
	})
	h.a.NoError(err)

	h.apiInit().
		Get("/settings/t_msg_k1.s1").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(`$.response.name=="t_msg_k1.s1"`)).
		Assert(jsonpath.Present(`$.response.value=="t_msg_v1"`)).
		End()

	h.apiInit().
		Get("/settings/t_msg_k1.missing").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response`, nil)).
		End()
}

func TestSettingsGet_noPermissions(t *testing.T) {
	h := newHelper(t)
	h.deny(tt.MessagingPermissionResource, "settings.read")

	h.apiInit().
		Get("/settings/t_msg_k1.s1").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to read settings")).
		End()
}
