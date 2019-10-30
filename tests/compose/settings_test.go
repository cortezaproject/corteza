package compose

import (
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/compose/service"
	tt "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/settings"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"github.com/jmoiron/sqlx/types"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func TestSettingsList(t *testing.T) {
	h := newHelper(t)
	h.allow(tt.ComposePermissionResource, "settings.read")
	h.allow(tt.ComposePermissionResource, "settings.manage")

	err := service.DefaultSettings.BulkSet(h.secCtx(), settings.ValueSet{
		&settings.Value{Name: "t_cmp_k1.s1", Value: types.JSONText(`"t_cmp_v1"`)},
		&settings.Value{Name: "t_cmp_k1.s2", Value: types.JSONText(`"t_cmp_v2"`)},
	})
	h.a.NoError(err)

	h.apiInit().
		Get("/settings/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(`$.response[? @.name=="t_cmp_k1.s1"]`)).
		Assert(jsonpath.Present(`$.response[? @.value=="t_cmp_v1"]`)).
		Assert(jsonpath.Present(`$.response[? @.name=="t_cmp_k1.s2"]`)).
		Assert(jsonpath.Present(`$.response[? @.value=="t_cmp_v2"]`)).
		End()
}

func TestSettingsList_noPermissions(t *testing.T) {
	h := newHelper(t)
	h.deny(tt.ComposePermissionResource, "settings.read")

	h.apiInit().
		Get("/settings/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to read settings")).
		End()
}

func TestSettingsUpdate(t *testing.T) {
	h := newHelper(t)
	h.allow(tt.ComposePermissionResource, "settings.manage")
	h.allow(tt.ComposePermissionResource, "settings.read")

	err := service.DefaultSettings.BulkSet(h.secCtx(), settings.ValueSet{
		&settings.Value{Name: "t_cmp_k1.s1", Value: types.JSONText(`"t_cmp_v1"`)},
	})
	h.a.NoError(err)

	h.apiInit().
		Patch("/settings/").
		JSON(`{"values":[{"name":"t_cmp_k1.s1","value":"t_cmp_v1_edited"},{"name":"t_cmp_k2.s1","value":"t_cmp_v2_new"}]}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	s, err := service.DefaultSettings.Get(h.secCtx(), "t_cmp_k1.s1", 0)
	h.a.NoError(err)
	h.a.Equal(`"t_cmp_v1_edited"`, s.Value.String(), "existing key should be updated")

	s, err = service.DefaultSettings.Get(h.secCtx(), "t_cmp_k2.s1", 0)
	h.a.NoError(err)
	h.a.Equal(`"t_cmp_v2_new"`, s.Value.String(), "new key should be added")
}

func TestSettingsUpdate_noPermissions(t *testing.T) {
	h := newHelper(t)
	h.deny(tt.ComposePermissionResource, "settings.manage")

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
	h.allow(tt.ComposePermissionResource, "settings.read")
	h.allow(tt.ComposePermissionResource, "settings.manage")

	err := service.DefaultSettings.BulkSet(h.secCtx(), settings.ValueSet{
		&settings.Value{Name: "t_cmp_k1.s1", Value: types.JSONText(`"t_cmp_v1"`)},
	})
	h.a.NoError(err)

	h.apiInit().
		Get("/settings/t_cmp_k1.s1").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(`$.response.name=="t_cmp_k1.s1"`)).
		Assert(jsonpath.Present(`$.response.value=="t_cmp_v1"`)).
		End()

	h.apiInit().
		Get("/settings/t_cmp_k1.missing").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response`, nil)).
		End()
}

func TestSettingsGet_noPermissions(t *testing.T) {
	h := newHelper(t)
	h.deny(tt.ComposePermissionResource, "settings.read")

	h.apiInit().
		Get("/settings/t_cmp_k1.s1").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to read settings")).
		End()
}
