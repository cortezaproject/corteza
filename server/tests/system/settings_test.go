package system

import (
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza/server/system/service"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/cortezaproject/corteza/server/tests/helpers"
	sqlTypes "github.com/jmoiron/sqlx/types"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func TestSettingsList(t *testing.T) {
	h := newHelper(t)
	helpers.AllowMe(h, types.ComponentRbacResource(), "settings.read")
	helpers.AllowMe(h, types.ComponentRbacResource(), "settings.manage")

	err := service.DefaultSettings.BulkSet(h.secCtx(), types.SettingValueSet{
		&types.SettingValue{Name: "t_sys_k1.s1", Value: sqlTypes.JSONText(`"t_sys_v1"`)},
		&types.SettingValue{Name: "t_sys_k1.s2", Value: sqlTypes.JSONText(`"t_sys_v2"`)},
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
	helpers.DenyMe(h, types.ComponentRbacResource(), "settings.read")

	h.apiInit().
		Get("/settings/").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("settings.errors.notAllowedToRead")).
		End()
}

func TestSettingsUpdate(t *testing.T) {
	h := newHelper(t)
	helpers.AllowMe(h, types.ComponentRbacResource(), "settings.manage")
	helpers.AllowMe(h, types.ComponentRbacResource(), "settings.read")

	err := service.DefaultSettings.BulkSet(h.secCtx(), types.SettingValueSet{
		&types.SettingValue{Name: "t_sys_k1.s1", Value: sqlTypes.JSONText(`"t_sys_v1"`)},
	})
	h.a.NoError(err)

	h.apiInit().
		Patch("/settings/").
		JSON(`{"values":[{"name":"t_sys_k1.s1","value":"t_sys_v1_edited"},{"name":"t_sys_k2.s1","value":"t_sys_v2_new"}]}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	s, err := service.DefaultSettings.Get(h.secCtx(), "t_sys_k1.s1", 0)
	h.a.NoError(err)
	h.a.Equal(`"t_sys_v1_edited"`, s.Value.String(), "existing key should be updated")

	s, err = service.DefaultSettings.Get(h.secCtx(), "t_sys_k2.s1", 0)
	h.a.NoError(err)
	h.a.Equal(`"t_sys_v2_new"`, s.Value.String(), "new key should be added")
}

func TestSettingsUpdate_noPermissions(t *testing.T) {
	h := newHelper(t)
	helpers.DenyMe(h, types.ComponentRbacResource(), "settings.manage")

	h.apiInit().
		Patch("/settings/").
		Header("Accept", "application/json").
		JSON(`{"values":[]}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("settings.errors.notAllowedToManage")).
		End()
}

func TestSettingsUpdate_validation(t *testing.T) {
	h := newHelper(t)
	helpers.AllowMe(h, types.ComponentRbacResource(), "settings.manage")
	helpers.AllowMe(h, types.ComponentRbacResource(), "settings.read")

	// Password constraints: The min password length should be 8 or more
	h.apiInit().
		Patch("/settings/").
		Header("Accept", "application/json").
		JSON(`{"values":[{"name":"auth.internal.password-constraints.min-length","value": 8}]}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	// Password constraints: The min password length should be not be less than 8
	h.apiInit().
		Patch("/settings/").
		Header("Accept", "application/json").
		JSON(`{"values":[{"name":"auth.internal.password-constraints.min-length","value":"7"}]}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("settings.errors.invalidPasswordMinLength")).
		End()

	// Password constraints: The min upper case count
	h.apiInit().
		Patch("/settings/").
		Header("Accept", "application/json").
		JSON(`{"values":[{"name":"auth.internal.password-constraints.min-upper-case","value": 2}]}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	// Password constraints: The min upper case count should not be a negative number
	h.apiInit().
		Patch("/settings/").
		Header("Accept", "application/json").
		JSON(`{"values":[{"name":"auth.internal.password-constraints.min-upper-case","value":-1}]}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("settings.errors.invalidPasswordMinUpperCase")).
		End()

	// Password constraints: The min upper case count should not be a negative number string
	h.apiInit().
		Patch("/settings/").
		Header("Accept", "application/json").
		JSON(`{"values":[{"name":"auth.internal.password-constraints.min-upper-case","value":"-1"}]}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("settings.errors.invalidPasswordMinUpperCase")).
		End()

	// Password constraints: The min lower case count should not be a negative number
	h.apiInit().
		Patch("/settings/").
		Header("Accept", "application/json").
		JSON(`{"values":[{"name":"auth.internal.password-constraints.min-lower-case","value":"-1"}]}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("settings.errors.invalidPasswordMinLowerCase")).
		End()

	// Password constraints: The min number of numeric characters should not be a negative number
	h.apiInit().
		Patch("/settings/").
		Header("Accept", "application/json").
		JSON(`{"values":[{"name":"auth.internal.password-constraints.min-num-count","value":"-1"}]}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("settings.errors.invalidPasswordMinNumCount")).
		End()

	// Password constraints: The min number of special characters should not be a negative number
	h.apiInit().
		Patch("/settings/").
		Header("Accept", "application/json").
		JSON(`{"values":[{"name":"auth.internal.password-constraints.min-special-count","value":"-1"}]}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("settings.errors.invalidPasswordMinSpecialCharCount")).
		End()

}

func TestSettingsGet(t *testing.T) {
	h := newHelper(t)
	helpers.AllowMe(h, types.ComponentRbacResource(), "settings.read")
	helpers.AllowMe(h, types.ComponentRbacResource(), "settings.manage")

	err := service.DefaultSettings.BulkSet(h.secCtx(), types.SettingValueSet{
		&types.SettingValue{Name: "t_sys_k1.s1", Value: sqlTypes.JSONText(`"t_sys_v1"`)},
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
	helpers.DenyMe(h, types.ComponentRbacResource(), "settings.read")

	h.apiInit().
		Get("/settings/t_sys_k1.s1").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("settings.errors.notAllowedToRead")).
		End()
}

func TestSettingsSet_noPermissions(t *testing.T) {
	h := newHelper(t)
	helpers.DenyMe(h, types.ComponentRbacResource(), "settings.read")

	h.apiInit().
		Get("/settings/t_sys_k1.s1").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("settings.errors.notAllowedToRead")).
		End()
}

func TestSettingsCurrent(t *testing.T) {
	h := newHelper(t)
	helpers.AllowMe(h, types.ComponentRbacResource(), "settings.read")

	h.apiInit().
		Get("/settings/current").
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Present(`$.response`)).
		End()
}
