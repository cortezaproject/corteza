package system

import (
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/settings"
	"github.com/cortezaproject/corteza-server/system/service"
	tt "github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"github.com/jmoiron/sqlx/types"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func TestSettingsList(t *testing.T) {
	h := newHelper(t)
	h.allow(tt.SystemPermissionResource, "settings.read")
	h.allow(tt.SystemPermissionResource, "settings.manage")

	err := service.DefaultSettings.With(h.secCtx()).BulkSet(settings.ValueSet{
		&settings.Value{Name: "t_sys_k1.s1", Value: types.JSONText(`"t_sys_v1"`)},
		&settings.Value{Name: "t_sys_k1.s2", Value: types.JSONText(`"t_sys_v2"`)},
	})
	h.a.NoError(err)

	h.apiInit().
		Get("/settings/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(`$.response[? @.name=="t_sys_k1.s1"]`)).
		Assert(jsonpath.Present(`$.response[? @.value=="t_sys_v1"]`)).
		Assert(jsonpath.Present(`$.response[? @.name=="t_sys_k1.s2"]`)).
		Assert(jsonpath.Present(`$.response[? @.value=="t_sys_v2"]`)).
		End()
}

func TestSettingsList_noPermissions(t *testing.T) {
	h := newHelper(t)
	h.deny(tt.SystemPermissionResource, "settings.read")

	h.apiInit().
		Get("/settings/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to read settings")).
		End()
}

func TestSettingsUpdate(t *testing.T) {
	h := newHelper(t)
	h.allow(tt.SystemPermissionResource, "settings.manage")
	h.allow(tt.SystemPermissionResource, "settings.read")

	err := service.DefaultSettings.With(h.secCtx()).BulkSet(settings.ValueSet{
		&settings.Value{Name: "t_sys_k1.s1", Value: types.JSONText(`"t_sys_v1"`)},
	})
	h.a.NoError(err)

	h.apiInit().
		Patch("/settings/").
		JSON(`{"values":[{"name":"t_sys_k1.s1","value":"t_sys_v1_edited"},{"name":"t_sys_k2.s1","value":"t_sys_v2_new"}]}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	s, err := service.DefaultSettings.With(h.secCtx()).Get("t_sys_k1.s1", 0)
	h.a.NoError(err)
	h.a.Equal(`"t_sys_v1_edited"`, s.Value.String(), "existing key should be updated")

	s, err = service.DefaultSettings.With(h.secCtx()).Get("t_sys_k2.s1", 0)
	h.a.NoError(err)
	h.a.Equal(`"t_sys_v2_new"`, s.Value.String(), "new key should be added")
}

func TestSettingsUpdate_noPermissions(t *testing.T) {
	h := newHelper(t)
	h.deny(tt.SystemPermissionResource, "settings.manage")

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
	h.allow(tt.SystemPermissionResource, "settings.read")
	h.allow(tt.SystemPermissionResource, "settings.manage")

	err := service.DefaultSettings.With(h.secCtx()).BulkSet(settings.ValueSet{
		&settings.Value{Name: "t_sys_k1.s1", Value: types.JSONText(`"t_sys_v1"`)},
	})
	h.a.NoError(err)

	h.apiInit().
		Get("/settings/t_sys_k1.s1").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(`$.response.name=="t_sys_k1.s1"`)).
		Assert(jsonpath.Present(`$.response.value=="t_sys_v1"`)).
		End()

	h.apiInit().
		Get("/settings/t_sys_k1.missing").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response`, nil)).
		End()
}

func TestSettingsGet_noPermissions(t *testing.T) {
	h := newHelper(t)
	h.deny(tt.SystemPermissionResource, "settings.read")

	h.apiInit().
		Get("/settings/t_sys_k1.s1").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to read settings")).
		End()
}
