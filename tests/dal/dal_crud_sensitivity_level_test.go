package dal

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func Test_dal_crud_sensitivity_level_list(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	helpers.AllowMe(h, types.ComponentRbacResource(), "dal-sensitivity-level.manage")

	h.apiInit().
		Get("/system/dal/sensitivity-levels/").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func Test_dal_crud_sensitivity_level_list_forbidden(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	h.apiInit().
		Get("/system/dal/sensitivity-levels/").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("dalSensitivityLevel.errors.notAllowedToManage")).
		End()
}

func Test_dal_crud_sensitivity_level_create(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	helpers.AllowMe(h, types.ComponentRbacResource(), "dal-sensitivity-level.manage")

	h.apiInit().
		Post("/system/dal/sensitivity-levels/").
		Body(loadRequestFromGenerics(t, "ok_sensitivity_level.json")).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.sensitivityLevelID")).
		End()
}

func Test_dal_crud_sensitivity_level_create_forbidden(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	h.apiInit().
		Post("/system/dal/sensitivity-levels/").
		Body(loadRequestFromGenerics(t, "ok_sensitivity_level.json")).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("dalSensitivityLevel.errors.notAllowedToManage")).
		End()
}

func Test_dal_crud_sensitivity_level_update(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	sl := h.createSensitivityLevel(&types.DalSensitivityLevel{
		Handle: "test_sl",
	})

	helpers.AllowMe(h, types.ComponentRbacResource(), "dal-sensitivity-level.manage")

	h.apiInit().
		Put(fmt.Sprintf("/system/dal/sensitivity-levels/%d", sl.ID)).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Body(loadRequestFromGenerics(t, "ok_sensitivity_level_update.json")).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.sensitivityLevelID")).
		Assert(jsonpath.Equal("$.response.handle", "private_edited")).
		End()
}

func Test_dal_crud_sensitivity_level_update_forbidden(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	sl := h.createSensitivityLevel(&types.DalSensitivityLevel{
		Handle: "test_sl",
	})

	h.apiInit().
		Put(fmt.Sprintf("/system/dal/sensitivity-levels/%d", sl.ID)).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Body(loadRequestFromGenerics(t, "ok_sensitivity_level_update.json")).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("dalSensitivityLevel.errors.notAllowedToManage")).
		End()
}

func Test_dal_crud_sensitivity_level_read(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	sl := h.createSensitivityLevel(&types.DalSensitivityLevel{
		Handle: "test_sl",
	})

	helpers.AllowMe(h, types.ComponentRbacResource(), "dal-sensitivity-level.manage")

	h.apiInit().
		Get(fmt.Sprintf("/system/dal/sensitivity-levels/%d", sl.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.sensitivityLevelID")).
		End()
}

func Test_dal_crud_sensitivity_level_read_forbiden(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	sl := h.createSensitivityLevel(&types.DalSensitivityLevel{
		Handle: "test_sl",
	})

	h.apiInit().
		Get(fmt.Sprintf("/system/dal/sensitivity-levels/%d", sl.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("dalSensitivityLevel.errors.notAllowedToManage")).
		End()
}

func Test_dal_crud_sensitivity_level_delete(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	sl := h.createSensitivityLevel(&types.DalSensitivityLevel{
		Handle: "test_sl",
	})

	helpers.AllowMe(h, types.ComponentRbacResource(), "dal-sensitivity-level.manage")

	h.apiInit().
		Delete(fmt.Sprintf("/system/dal/sensitivity-levels/%d", sl.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal("$.success.message", "OK")).
		End()
}

func Test_dal_crud_sensitivity_level_delete_forbidden(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	sl := h.createSensitivityLevel(&types.DalSensitivityLevel{
		Handle: "test_sl",
	})

	h.apiInit().
		Delete(fmt.Sprintf("/system/dal/sensitivity-levels/%d", sl.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("dalSensitivityLevel.errors.notAllowedToManage")).
		End()
}

func Test_dal_crud_sensitivity_level_undelete(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	sl := h.createSensitivityLevel(&types.DalSensitivityLevel{
		Handle:    "test_sl",
		DeletedAt: &h.cUser.CreatedAt,
		DeletedBy: h.cUser.ID,
	})

	helpers.AllowMe(h, types.ComponentRbacResource(), "dal-sensitivity-level.manage")

	h.apiInit().
		Post(fmt.Sprintf("/system/dal/sensitivity-levels/%d/undelete", sl.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal("$.success.message", "OK")).
		End()
}

func Test_dal_crud_sensitivity_level_undelete_forbidden(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	sl := h.createSensitivityLevel(&types.DalSensitivityLevel{
		Handle:    "test_sl",
		DeletedAt: &h.cUser.CreatedAt,
		DeletedBy: h.cUser.ID,
	})

	h.apiInit().
		Post(fmt.Sprintf("/system/dal/sensitivity-levels/%d/undelete", sl.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("dalSensitivityLevel.errors.notAllowedToManage")).
		End()
}
