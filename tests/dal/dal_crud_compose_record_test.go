package dal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/dal/capabilities"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func Test_dal_crud_compose_record_create(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	ctx := h.secCtx()

	helpers.AllowMeModuleCRUD(h)
	helpers.AllowMeRecordCRUD(h)

	ns := h.createNamespace("test")
	module := createModuleFromGenerics(ctx, t, "ok_module.json", ns.ID, &types.ModelConfig{
		ConnectionID: 0,
		Capabilities: capabilities.FullCapabilities(),
	})

	// create
	h.apiInit().
		Post(fmt.Sprintf("/compose/namespace/%d/module/%d/record/", ns.ID, module.ID)).
		Body(loadRequestFromGenerics(t, "ok_record.json")).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	rsp := h.apiInit().
		Get(fmt.Sprintf("/compose/namespace/%d/module/%d/record/", ns.ID, module.ID)).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	aux := composeRecordSearchRestRsp{}
	dd := json.NewDecoder(rsp.Response.Body)
	h.a.NoError(dd.Decode(&aux))

	r := aux.Response.Set[0]
	h.a.NotEqual(uint64(0), r.ID)
	h.a.Equal(module.ID, r.ModuleID)
	h.a.Equal(module.NamespaceID, r.NamespaceID)
	h.a.Equal(h.cUser.ID, r.OwnedBy)

	h.a.NotEqual(time.Time{}, r.CreatedAt)
	h.a.Equal(h.cUser.ID, r.CreatedBy)

	h.a.Nil(r.UpdatedAt)
	h.a.Equal(uint64(0), r.UpdatedBy)

	h.a.Nil(r.DeletedAt)
	h.a.Equal(uint64(0), r.DeletedBy)

	h.a.Equal("me@test.tld", r.Values.Get("email", 0).Value)
	h.a.Equal("Me", r.Values.Get("name", 0).Value)
	h.a.Equal("42", r.Values.Get("a_number", 0).Value)
}

func Test_dal_crud_compose_record_update(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	ctx := h.secCtx()

	helpers.AllowMeModuleCRUD(h)
	helpers.AllowMeRecordCRUD(h)

	ns := h.createNamespace("test")
	module := createModuleFromGenerics(ctx, t, "ok_module.json", ns.ID, &types.ModelConfig{
		ConnectionID: 0,
		Capabilities: capabilities.FullCapabilities(),
	})

	record := createRecordFromGenerics(ctx, t, "ok_record.json", ns.ID, module.ID)
	h.apiInit().
		Post(fmt.Sprintf("/compose/namespace/%d/module/%d/record/%d", ns.ID, module.ID, record.ID)).
		Body(loadRequestFromGenerics(t, "ok_record_update.json")).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	rsp := h.apiInit().
		Get(fmt.Sprintf("/compose/namespace/%d/module/%d/record/", ns.ID, module.ID)).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	aux := composeRecordSearchRestRsp{}
	dd := json.NewDecoder(rsp.Response.Body)
	h.a.NoError(dd.Decode(&aux))

	r := aux.Response.Set[0]
	h.a.NotEqual(uint64(0), r.ID)
	h.a.Equal(module.ID, r.ModuleID)
	h.a.Equal(module.NamespaceID, r.NamespaceID)
	h.a.Equal(h.cUser.ID, r.OwnedBy)

	h.a.NotEqual(time.Time{}, r.CreatedAt)
	h.a.Equal(h.cUser.ID, r.CreatedBy)

	h.a.NotNil(r.UpdatedAt)
	h.a.Equal(h.cUser.ID, r.UpdatedBy)

	h.a.Nil(r.DeletedAt)
	h.a.Equal(uint64(0), r.DeletedBy)

	h.a.Equal("me+updated@test.tld", r.Values.Get("email", 0).Value)
	h.a.Equal("Me updated", r.Values.Get("name", 0).Value)
	h.a.Equal("43", r.Values.Get("a_number", 0).Value)
}

func Test_dal_crud_compose_record_delete(t *testing.T) {
	h := newHelperT(t)
	defer h.cleanupDal()

	ctx := h.secCtx()

	helpers.AllowMeModuleCRUD(h)
	helpers.AllowMeRecordCRUD(h)

	ns := h.createNamespace("test")
	module := createModuleFromGenerics(ctx, t, "ok_module.json", ns.ID, &types.ModelConfig{
		ConnectionID: 0,
		Capabilities: capabilities.FullCapabilities(),
	})

	record := createRecordFromGenerics(ctx, t, "ok_record.json", ns.ID, module.ID)

	h.apiInit().
		Delete(fmt.Sprintf("/compose/namespace/%d/module/%d/record/%d", ns.ID, module.ID, record.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	rsp := h.apiInit().
		Get(fmt.Sprintf("/compose/namespace/%d/module/%d/record/%d", ns.ID, module.ID, record.ID)).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	aux := composeRecordRestRsp{}
	dd := json.NewDecoder(rsp.Response.Body)
	h.a.NoError(dd.Decode(&aux))

	r := aux.Response
	h.a.NotEqual(uint64(0), r.ID)
	h.a.Equal(module.ID, r.ModuleID)
	h.a.Equal(module.NamespaceID, r.NamespaceID)
	h.a.Equal(h.cUser.ID, r.OwnedBy)

	h.a.NotEqual(time.Time{}, r.CreatedAt)
	h.a.Equal(h.cUser.ID, r.CreatedBy)

	h.a.Nil(r.UpdatedAt)
	h.a.Equal(uint64(0), r.UpdatedBy)

	h.a.NotNil(r.DeletedAt)
	h.a.Equal(h.cUser.ID, r.DeletedBy)
}
