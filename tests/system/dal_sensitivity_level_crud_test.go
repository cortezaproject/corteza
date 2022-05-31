package system

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func (h helper) clearSensitivityLevels() {
	h.noError(store.TruncateDalSensitivityLevels(context.Background(), service.DefaultStore))
}

func (h helper) createSensitivityLevel(res *types.DalSensitivityLevel) *types.DalSensitivityLevel {
	if res.ID == 0 {
		res.ID = id.Next()
	}

	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}

	h.a.NoError(service.DefaultStore.CreateDalSensitivityLevel(context.Background(), res))
	return res
}

func Test_dal_sensitivity_level_list(t *testing.T) {
	h := newHelper(t)
	defer h.clearSensitivityLevels()

	helpers.AllowMe(h, types.ComponentRbacResource(), "dal-sensitivity-level.manage")

	h.apiInit().
		Get("/dal/sensitivity-levels/").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func Test_dal_sensitivity_level_list_forbidden(t *testing.T) {
	h := newHelper(t)
	defer h.clearSensitivityLevels()

	h.apiInit().
		Get("/dal/sensitivity-levels/").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("dalSensitivityLevel.errors.notAllowedToManage")).
		End()
}

func Test_dal_sensitivity_level_create(t *testing.T) {
	h := newHelper(t)
	defer h.clearSensitivityLevels()

	helpers.AllowMe(h, types.ComponentRbacResource(), "dal-sensitivity-level.manage")

	h.apiInit().
		Post("/dal/sensitivity-levels/").
		Body(loadScenarioRequest(t, "generic.json")).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.sensitivityLevelID")).
		End()
}

func Test_dal_sensitivity_level_create_forbidden(t *testing.T) {
	h := newHelper(t)
	defer h.clearSensitivityLevels()

	h.apiInit().
		Post("/dal/sensitivity-levels/").
		Body(loadScenarioRequest(t, "generic.json")).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("dalSensitivityLevel.errors.notAllowedToManage")).
		End()
}

func Test_dal_sensitivity_level_update(t *testing.T) {
	h := newHelper(t)
	defer h.clearSensitivityLevels()

	sl := h.createSensitivityLevel(&types.DalSensitivityLevel{
		Handle: "test_sl",
	})

	helpers.AllowMe(h, types.ComponentRbacResource(), "dal-sensitivity-level.manage")

	h.apiInit().
		Put(fmt.Sprintf("/dal/sensitivity-levels/%d", sl.ID)).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Body(loadScenarioRequest(t, "generic.json")).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.sensitivityLevelID")).
		Assert(jsonpath.Equal("$.response.handle", "private_edited")).
		End()
}

func Test_dal_sensitivity_level_update_forbidden(t *testing.T) {
	h := newHelper(t)
	defer h.clearSensitivityLevels()

	sl := h.createSensitivityLevel(&types.DalSensitivityLevel{
		Handle: "test_sl",
	})

	h.apiInit().
		Put(fmt.Sprintf("/dal/sensitivity-levels/%d", sl.ID)).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Body(loadScenarioRequest(t, "generic.json")).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("dalSensitivityLevel.errors.notAllowedToManage")).
		End()
}

func Test_dal_sensitivity_level_read(t *testing.T) {
	h := newHelper(t)
	defer h.clearSensitivityLevels()

	sl := h.createSensitivityLevel(&types.DalSensitivityLevel{
		Handle: "test_sl",
	})

	helpers.AllowMe(h, types.ComponentRbacResource(), "dal-sensitivity-level.manage")

	h.apiInit().
		Get(fmt.Sprintf("/dal/sensitivity-levels/%d", sl.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.sensitivityLevelID")).
		End()
}

func Test_dal_sensitivity_level_read_forbiden(t *testing.T) {
	h := newHelper(t)
	defer h.clearSensitivityLevels()

	sl := h.createSensitivityLevel(&types.DalSensitivityLevel{
		Handle: "test_sl",
	})

	h.apiInit().
		Get(fmt.Sprintf("/dal/sensitivity-levels/%d", sl.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("dalSensitivityLevel.errors.notAllowedToManage")).
		End()
}

func Test_dal_sensitivity_level_delete(t *testing.T) {
	h := newHelper(t)
	defer h.clearSensitivityLevels()

	sl := h.createSensitivityLevel(&types.DalSensitivityLevel{
		Handle: "test_sl",
	})

	helpers.AllowMe(h, types.ComponentRbacResource(), "dal-sensitivity-level.manage")

	h.apiInit().
		Delete(fmt.Sprintf("/dal/sensitivity-levels/%d", sl.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal("$.success.message", "OK")).
		End()
}

func Test_dal_sensitivity_level_delete_forbidden(t *testing.T) {
	h := newHelper(t)
	defer h.clearSensitivityLevels()

	sl := h.createSensitivityLevel(&types.DalSensitivityLevel{
		Handle: "test_sl",
	})

	h.apiInit().
		Delete(fmt.Sprintf("/dal/sensitivity-levels/%d", sl.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("dalSensitivityLevel.errors.notAllowedToManage")).
		End()
}

func Test_dal_sensitivity_level_undelete(t *testing.T) {
	h := newHelper(t)
	defer h.clearSensitivityLevels()

	sl := h.createSensitivityLevel(&types.DalSensitivityLevel{
		Handle:    "test_sl",
		DeletedAt: &h.cUser.CreatedAt,
		DeletedBy: h.cUser.ID,
	})

	helpers.AllowMe(h, types.ComponentRbacResource(), "dal-sensitivity-level.manage")

	h.apiInit().
		Post(fmt.Sprintf("/dal/sensitivity-levels/%d/undelete", sl.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal("$.success.message", "OK")).
		End()
}

func Test_dal_sensitivity_level_undelete_forbidden(t *testing.T) {
	h := newHelper(t)
	defer h.clearSensitivityLevels()

	sl := h.createSensitivityLevel(&types.DalSensitivityLevel{
		Handle:    "test_sl",
		DeletedAt: &h.cUser.CreatedAt,
		DeletedBy: h.cUser.ID,
	})

	h.apiInit().
		Post(fmt.Sprintf("/dal/sensitivity-levels/%d/undelete", sl.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("dalSensitivityLevel.errors.notAllowedToManage")).
		End()
}
