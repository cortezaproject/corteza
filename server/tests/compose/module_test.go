package compose

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/compose/service"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/require"
)

func (h helper) clearModules() {
	h.clearNamespaces()
	h.noError(store.TruncateComposeModules(context.Background(), service.DefaultStore))
	h.noError(store.TruncateComposeModuleFields(context.Background(), service.DefaultStore))

	models, err := defDal.SearchModels(context.Background())
	h.noError(err)
	for _, m := range models {
		h.noError(defDal.RemoveModel(context.Background(), m.ConnectionID, m.ResourceID))
	}
}

func (h helper) makeModule(ns *types.Namespace, name string, ff ...*types.ModuleField) *types.Module {
	return h.createModule(ns, &types.Module{
		Name:        name,
		NamespaceID: ns.ID,
		Fields:      ff,
		CreatedAt:   time.Now(),
	})
}

func (h helper) createModule(ns *types.Namespace, res *types.Module) *types.Module {
	res.ID = id.Next()
	res.CreatedAt = time.Now()
	h.noError(store.CreateComposeModule(context.Background(), service.DefaultStore, res))

	if res.Config.DAL.ConnectionID == 0 {
		res.Config.DAL.ConnectionID = defDal.GetConnectionByID(0).ID
	}

	_ = res.Fields.Walk(func(f *types.ModuleField) error {
		f.ID = id.Next()
		f.ModuleID = res.ID
		f.NamespaceID = res.NamespaceID
		f.CreatedAt = time.Now()
		return nil
	})

	h.noError(store.CreateComposeModuleField(context.Background(), service.DefaultStore, res.Fields...))

	h.noError(service.DalModelReplace(context.Background(), service.DefaultStore, nil, defDal, ns, res))

	return res
}

func (h helper) lookupModuleByID(ID uint64) *types.Module {
	res, err := store.LookupComposeModuleByID(context.Background(), service.DefaultStore, ID)
	h.noError(err)

	res.Fields, _, err = store.SearchComposeModuleFields(context.Background(), service.DefaultStore, types.ModuleFieldFilter{ModuleID: []uint64{ID}})
	h.noError(err)

	return res
}

func TestModuleRead(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "modules.search")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "read")
	ns := h.makeNamespace("some-namespace")
	m := h.makeModule(ns, "some-module")

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, m.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.name`, m.Name)).
		Assert(jsonpath.Equal(`$.response.moduleID`, fmt.Sprintf("%d", m.ID))).
		End()
}

func TestModuleReadByHandle(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "modules.search")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "read")
	ns := h.makeNamespace("some-namespace")
	c := h.makeModule(ns, "some-module")

	cbh, err := service.DefaultModule.FindByHandle(h.secCtx(), ns.ID, c.Handle)

	h.noError(err)
	h.a.NotNil(cbh)
	h.a.Equal(cbh.ID, c.ID)
	h.a.Equal(cbh.Handle, c.Handle)
}

func TestModuleList(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "modules.search")
	ns := h.makeNamespace("some-namespace")

	h.makeModule(ns, "app")
	h.makeModule(ns, "app")

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/", ns.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestModuleListQuery(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "modules.search")
	ns := h.makeNamespace("some-namespace")

	h.createModule(ns, &types.Module{
		Name:        "name",
		Handle:      "handle",
		NamespaceID: ns.ID,
	})

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/", ns.ID)).
		Query("query", "handle").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.set != null")).
		End()
}

func TestModuleList_filterForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "modules.search")
	ns := h.makeNamespace("some-namespace")

	h.makeModule(ns, "module")
	f := h.makeModule(ns, "module_forbidden")

	helpers.DenyMe(h, types.ModuleRbacResource(f.NamespaceID, f.ID), "read")

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/", ns.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.NotPresent(`$.response.set[? @.name=="module_forbidden"]`)).
		End()
}

func TestModuleCreateForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	ns := h.makeNamespace("some-namespace")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/", ns.ID)).
		Header("Accept", "application/json").
		FormData("name", "some-module").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("module.errors.notAllowedToCreate")).
		End()
}

func TestModuleCreate(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "modules.search")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "module.create")

	ns := h.makeNamespace("some-namespace")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/", ns.ID)).
		FormData("name", "some-module").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestModuleCreateInvalidField(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "modules.search")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "module.create")

	ns := h.makeNamespace("some-namespace")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/", ns.ID)).
		JSON(`{ "name": "mod", "fields": [{ "name": "a", "kind": "Number" }] }`).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("module.errors.invalidHandle")).
		End()
}

func TestModuleUpdateForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "modules.search")
	ns := h.makeNamespace("some-namespace")
	m := h.makeModule(ns, "some-module")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, m.ID)).
		Header("Accept", "application/json").
		FormData("name", "changed-name").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("module.errors.notAllowedToUpdate")).
		End()
}

func TestModuleUpdate(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "modules.search")
	ns := h.makeNamespace("some-namespace")
	m := h.makeModule(ns, "some-module")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "update")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, m.ID)).
		FormData("name", "changed-name").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present("$.response.updatedAt")).
		End()

	m, err := store.LookupComposeModuleByID(context.Background(), service.DefaultStore, m.ID)
	h.noError(err)
	h.a.NotNil(m)
	h.a.Equal("changed-name", m.Name)
}

func TestModuleFieldsUpdate(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "modules.search")
	ns := h.makeNamespace("some-namespace")
	m := h.makeModule(ns, "some-module", &types.ModuleField{ID: id.Next(), Kind: "String", Name: "existing"})
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "update")

	f := m.Fields[0]
	fjs := fmt.Sprintf(`{ "name": "%s", "fields": [{ "fieldID": "%d", "name": "existing_edited", "kind": "Number" }, { "name": "new", "kind": "DateTime" }] }`, m.Name, f.ID)
	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, m.ID)).
		JSON(fjs).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	m = h.lookupModuleByID(m.ID)
	h.a.NotNil(m)
	h.a.NotNil(m.Fields)
	h.a.Len(m.Fields, 2)

	h.a.NotNil(m.Fields[0].UpdatedAt)
	h.a.Equal(m.Fields[0].Name, "existing_edited")
	h.a.Equal(m.Fields[0].Kind, "Number")
	h.a.Nil(m.Fields[1].UpdatedAt)
	h.a.Equal(m.Fields[1].Name, "new")
	h.a.Equal(m.Fields[1].Kind, "DateTime")
}

func TestModuleFieldsUpdate_invalidHandle(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "modules.search")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "module.create")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "update")

	ns := h.makeNamespace("some-namespace")
	mod := h.makeModule(ns, "mod", &types.ModuleField{
		Name: "a",
		Kind: "String",
	})

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, mod.ID)).
		JSON(`{ "name": "mod", "fields": [{ "name": "a", "kind": "String" }] }`).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestModuleUpdateWithReservedFieldName(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "modules.search")
	ns := h.makeNamespace("some-namespace")
	m := h.makeModule(ns, "some-module", &types.ModuleField{ID: id.Next(), Kind: "String", Name: "ownedBy"})
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "update")

	f := m.Fields[0]
	fjs := fmt.Sprintf(`{ "name": "%s", "fields": [{ "fieldID": "%d", "name": "ownedBy", "kind": "Number" }]}`, m.Name, f.ID)
	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, m.ID)).
		Header("Accept", "application/json").
		JSON(fjs).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	f = m.Fields[0]
	fjs = fmt.Sprintf(`{ "name": "%s", "fields": [{ "fieldID": "%d", "name": "updatedBy", "kind": "Number" }]}`, m.Name, f.ID)
	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, m.ID)).
		Header("Accept", "application/json").
		JSON(fjs).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("module.errors.fieldNameReserved")).
		End()
}

func TestModuleCreateWithReservedFieldName(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "modules.search")
	ns := h.makeNamespace("some-namespace")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "update")

	fjs := fmt.Sprintf(`{ "name": "foo", "fields": [{ "name": "ownedBy", "kind": "Number" }]}`)
	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/", ns.ID)).
		Header("Accept", "application/json").
		JSON(fjs).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("module.errors.fieldNameReserved")).
		End()
}

func TestModuleFieldsUpdate_defaults(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "modules.search")
	ns := h.makeNamespace("some-namespace")
	m := h.makeModule(ns, "some-module", &types.ModuleField{ID: id.Next(), Kind: "String", Name: "existing", Required: true, DefaultValue: types.RecordValueSet{&types.RecordValue{Value: "test"}}})

	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "update")

	f := m.Fields[0]
	fjs := fmt.Sprintf(`{ "name": "%s", "fields": [{ "fieldID": "%d", "name": "existing_edited", "kind": "String", "isRequired": true, "defaultValue": [{ "value": "test" }] }, { "name": "new", "kind": "Bool", "defaultValue": [{ "value": "1" }] }] }`, m.Name, f.ID)
	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, m.ID)).
		JSON(fjs).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	m = h.lookupModuleByID(m.ID)
	h.a.NotNil(m)
	h.a.NotNil(m.Fields)
	h.a.Len(m.Fields, 2)

	h.a.Equal(m.Fields[0].Kind, "String")
	h.a.Len(m.Fields[0].DefaultValue, 1)
	h.a.Equal("test", m.Fields[0].DefaultValue[0].Value)

	h.a.Equal(m.Fields[1].Kind, "Bool")
	h.a.Len(m.Fields[1].DefaultValue, 1)
	h.a.Equal("1", m.Fields[1].DefaultValue[0].Value)
}

func TestModuleFieldsDefaultValue(t *testing.T) {
	var ns *types.Namespace

	h := newHelper(t)
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "modules.search")

	prep := func() {
		h.clearModules()
		ns = h.makeNamespace("some-namespace")
	}

	t.Run("boolean; true", func(t *testing.T) {
		prep()

		m := h.makeModule(ns, "some-module", &types.ModuleField{ID: id.Next(), Kind: "Boolean", Name: "boolean"})
		helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "update")

		f := m.Fields[0]
		fjs := fmt.Sprintf(`{ "name": "%s", "fields": [{ "fieldID": "%d", "name": "boolean", "kind": "Boolean", "defaultValue": [{"name": "boolean", "value": "1"}] }] }`, m.Name, f.ID)
		h.apiInit().
			Post(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, m.ID)).
			JSON(fjs).
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End()

		m = h.lookupModuleByID(m.ID)
		h.a.NotNil(m)
		h.a.NotNil(m.Fields)
		h.a.Len(m.Fields, 1)

		h.a.NotNil(m.Fields[0].DefaultValue)
		h.a.Len(m.Fields[0].DefaultValue, 1)
		h.a.Equal("1", m.Fields[0].DefaultValue[0].Value)
	})

	t.Run("boolean; false", func(t *testing.T) {
		prep()

		m := h.makeModule(ns, "some-module", &types.ModuleField{ID: id.Next(), Kind: "Boolean", Name: "boolean"})
		helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "update")

		f := m.Fields[0]
		fjs := fmt.Sprintf(`{ "name": "%s", "fields": [{ "fieldID": "%d", "name": "boolean", "kind": "Boolean", "defaultValue": [{"name": "boolean", "value": ""}] }] }`, m.Name, f.ID)
		h.apiInit().
			Post(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, m.ID)).
			JSON(fjs).
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End()

		m = h.lookupModuleByID(m.ID)
		h.a.NotNil(m)
		h.a.NotNil(m.Fields)
		h.a.Len(m.Fields, 1)

		h.a.NotNil(m.Fields[0].DefaultValue)
		h.a.Len(m.Fields[0].DefaultValue, 1)
		h.a.Equal("", m.Fields[0].DefaultValue[0].Value)
	})

	t.Run("boolean; undefined", func(t *testing.T) {
		prep()

		m := h.makeModule(ns, "some-module", &types.ModuleField{ID: id.Next(), Kind: "Boolean", Name: "boolean"})
		helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "update")

		f := m.Fields[0]
		fjs := fmt.Sprintf(`{ "name": "%s", "fields": [{ "fieldID": "%d", "name": "boolean", "kind": "Boolean" }] }`, m.Name, f.ID)
		h.apiInit().
			Post(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, m.ID)).
			JSON(fjs).
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End()

		m = h.lookupModuleByID(m.ID)
		h.a.NotNil(m)
		h.a.NotNil(m.Fields)
		h.a.Len(m.Fields, 1)

		h.a.Nil(m.Fields[0].DefaultValue)
		h.a.Len(m.Fields[0].DefaultValue, 0)
	})

	t.Run("boolean; true; compact form", func(t *testing.T) {
		prep()

		m := h.makeModule(ns, "some-module", &types.ModuleField{ID: id.Next(), Kind: "Boolean", Name: "boolean"})
		helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "update")

		f := m.Fields[0]
		fjs := fmt.Sprintf(`{ "name": "%s", "fields": [{ "fieldID": "%d", "name": "boolean", "kind": "Boolean", "defaultValue": [{"value": "1"}] }] }`, m.Name, f.ID)
		h.apiInit().
			Post(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, m.ID)).
			JSON(fjs).
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End()

		m = h.lookupModuleByID(m.ID)
		h.a.NotNil(m)
		h.a.NotNil(m.Fields)
		h.a.Len(m.Fields, 1)

		h.a.NotNil(m.Fields[0].DefaultValue)
		h.a.Len(m.Fields[0].DefaultValue, 1)
		h.a.Equal("1", m.Fields[0].DefaultValue[0].Value)
	})

	t.Run("record; valid", func(t *testing.T) {
		prep()

		// Prep the related module
		refM := h.makeModule(ns, "ref-mod", &types.ModuleField{
			ID:   id.Next(),
			Kind: "String",
			Name: "string",
		})
		rec := h.makeRecord(refM, &types.RecordValue{
			Name:  "string",
			Value: "val",
		})

		// Prep the module and a field
		m := h.makeModule(ns, "some-module", &types.ModuleField{
			ID:   id.Next(),
			Kind: "String",
			Name: "string",
		})

		// RBAC (make sure to allow record read)
		helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "update")
		helpers.AllowMe(h, types.RecordRbacResource(0, 0, 0), "read")

		f := m.Fields[0]
		fjs := fmt.Sprintf(`{ "name": "%s", "fields": [{ "fieldID": "%d", "name": "record", "kind": "Record", "options": {"moduleID": "%d"}, "defaultValue": [{"value": "%d"}] }] }`, m.Name, f.ID, refM.ID, rec.ID)
		h.apiInit().
			Post(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, m.ID)).
			JSON(fjs).
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End()

		m = h.lookupModuleByID(m.ID)
		h.a.NotNil(m)
		h.a.NotNil(m.Fields)
		h.a.Len(m.Fields, 1)

		h.a.NotNil(m.Fields[0].DefaultValue)
		h.a.Len(m.Fields[0].DefaultValue, 1)
		h.a.Equal(strconv.FormatUint(rec.ID, 10), m.Fields[0].DefaultValue[0].Value)
	})
}

func TestModuleFieldsUpdate_removed(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "modules.search")
	ns := h.makeNamespace("some-namespace")
	m := h.makeModule(ns, "some-module", &types.ModuleField{ID: id.Next(), Kind: "String", Name: "a"}, &types.ModuleField{ID: id.Next(), Kind: "String", Name: "b"})
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "update")

	f := m.Fields[0]
	fjs := fmt.Sprintf(`{ "name": "%s", "fields": [{ "fieldID": "%d", "name": "a", "kind": "String" }] }`, m.Name, f.ID)
	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, m.ID)).
		JSON(fjs).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	m = h.lookupModuleByID(m.ID)
	h.a.NotNil(m)
	h.a.NotNil(m.Fields)
	h.a.Len(m.Fields, 1)

	h.a.NotNil(m.Fields[0].UpdatedAt)
	h.a.Equal(m.Fields[0].Name, "a")
}

func TestModuleFieldsUpdate_removedHasRecords(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "modules.search")
	ns := h.makeNamespace("some-namespace")
	m := h.makeModule(ns, "some-module", &types.ModuleField{ID: id.Next(), Kind: "String", Name: "a"}, &types.ModuleField{ID: id.Next(), Kind: "String", Name: "b"})
	h.makeRecord(m, &types.RecordValue{Name: "a", Value: "va"}, &types.RecordValue{Name: "b", Value: "vb"})
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "update")

	f := m.Fields[0]
	fjs := fmt.Sprintf(`{ "name": "%s", "fields": [{ "fieldID": "%d", "name": "a", "kind": "String" }] }`, m.Name, f.ID)
	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, m.ID)).
		JSON(fjs).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	m = h.lookupModuleByID(m.ID)
	h.a.NotNil(m)
	h.a.NotNil(m.Fields)
	h.a.Len(m.Fields, 1)

	h.a.NotNil(m.Fields[0].UpdatedAt)
	h.a.Equal(m.Fields[0].Name, "a")
}

func TestModuleFieldsUpdateExpressions(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "modules.search")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "module.create")
	ns := h.makeNamespace("some-namespace")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "read")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "update")

	var (
		m = &types.Module{
			NamespaceID: ns.ID,
		}

		f = &types.ModuleField{
			ID:   id.Next(),
			Kind: "String",
			Name: "existing",
			Expressions: types.ModuleFieldExpr{
				ValueExpr:  `"foo"`,
				Sanitizers: []string{"value"},
				Validators: []types.ModuleFieldValidator{
					{
						Test:  `value != ""`,
						Error: "error",
					},
				},
				DisableDefaultValidators: true,
				Formatters:               []string{`"foo"`},
				DisableDefaultFormatters: false,
			},
		}

		aux = struct{ Response types.Module }{}
	)

	m.Fields = append(m.Fields, f)

	// create module
	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/", ns.ID)).
		JSON(helpers.JSON(m)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End().
		JSON(&aux)

	m = &aux.Response

	h.a.NotEmpty(m.Fields)
	f = m.Fields.FindByName("existing")
	h.a.NotNil(f)
	h.a.Equal(`"foo"`, f.Expressions.ValueExpr)
	h.a.NotEmpty(f.Expressions.Sanitizers)
	h.a.NotEmpty(f.Expressions.Validators)
	h.a.True(f.Expressions.DisableDefaultValidators)
	h.a.NotEmpty(f.Expressions.Formatters)
	h.a.False(f.Expressions.DisableDefaultFormatters)

	f.Expressions = types.ModuleFieldExpr{
		ValueExpr:  `"bar"`,
		Sanitizers: []string{`""`},
		Validators: []types.ModuleFieldValidator{
			{
				Test:  `value == ""`,
				Error: "foo",
			},
		},
		DisableDefaultValidators: false,
		Formatters:               []string{`"bar"`},
		DisableDefaultFormatters: true,
	}

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, m.ID)).
		JSON(helpers.JSON(m)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, m.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End().
		JSON(&aux)

	h.a.NotEmpty(aux.Response.Fields)
	f = aux.Response.Fields.FindByName("existing")
	h.a.NotNil(f)
	h.a.Equal(`"bar"`, f.Expressions.ValueExpr)
	h.a.NotEmpty(f.Expressions.Sanitizers)
	h.a.NotEmpty(f.Expressions.Validators)
	h.a.False(f.Expressions.DisableDefaultValidators)
	h.a.NotEmpty(f.Expressions.Formatters)
	h.a.True(f.Expressions.DisableDefaultFormatters)
}

func TestModuleFieldsPreventUpdate_ifRecordExists(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "modules.search")
	ns := h.makeNamespace("some-namespace")
	m := h.makeModule(ns, "some-module", &types.ModuleField{ID: id.Next(), Kind: "String", Name: "existing"})
	h.makeRecord(m, &types.RecordValue{Name: "existing", Value: "value"})
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "update")

	f := m.Fields[0]
	fjs := fmt.Sprintf(`{ "name": "%s", "fields": [{ "fieldID": "%d", "name": "existing_edited", "kind": "Number" }, { "name": "new", "kind": "DateTime" }] }`, m.Name, f.ID)
	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, m.ID)).
		JSON(fjs).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	m = h.lookupModuleByID(m.ID)
	h.a.NotNil(m)
	h.a.NotNil(m.Fields)
	h.a.Len(m.Fields, 2)

	h.a.NotNil(m.Fields[0].UpdatedAt)
	h.a.Equal(m.Fields[0].Name, "existing")
	h.a.Equal(m.Fields[0].Kind, "String")
	h.a.Nil(m.Fields[1].UpdatedAt)
	h.a.Equal(m.Fields[1].Name, "new")
	h.a.Equal(m.Fields[1].Kind, "DateTime")
}

func TestModuleDeleteForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "modules.search")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "read")
	ns := h.makeNamespace("some-namespace")
	m := h.makeModule(ns, "some-module")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, m.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("module.errors.notAllowedToDelete")).
		End()
}

func TestModuleDelete(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "modules.search")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "read")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "delete")

	ns := h.makeNamespace("some-namespace")
	res := h.makeModule(ns, "some-module")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, res.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	res = h.lookupModuleByID(res.ID)
	h.a.NotNil(res.DeletedAt)
}

func TestModuleLabels(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "modules.search")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "module.create")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "read")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "update")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "delete")

	var (
		ns          = h.makeNamespace("some-namespace")
		fieldID, ID uint64
	)

	t.Run("create", func(t *testing.T) {
		var (
			req     = require.New(t)
			payload = &types.Module{}
		)

		helpers.SetLabelsViaAPI(h.apiInit(), t,
			fmt.Sprintf("/namespace/%d/module/", ns.ID),
			types.Module{Labels: map[string]string{"foo": "bar", "bar": "42"}},
			payload,
		)
		req.NotZero(payload.ID)

		h.a.Equal(payload.Labels["foo"], "bar",
			"labels must contain foo with value bar")
		h.a.Equal(payload.Labels["bar"], "42",
			"labels must contain bar with value 42")
		req.Equal(payload.Labels, helpers.LoadLabelsFromStore(t, service.DefaultStore, payload.LabelResourceKind(), payload.ID),
			"response must match stored labels")

		ID = payload.ID
	})

	t.Run("update", func(t *testing.T) {
		if ID == 0 {
			t.Skip("label/create test not ran")
		}

		var (
			req     = require.New(t)
			payload = &types.Module{}
		)

		helpers.SetLabelsViaAPI(h.apiInit(), t,
			fmt.Sprintf("/namespace/%d/module/%d", ns.ID, ID),
			types.Module{Labels: map[string]string{"foo": "baz", "baz": "123"}},
			payload,
		)
		req.NotZero(payload.ID)
		req.Nil(payload.UpdatedAt, "updatedAt must not change after changing labels")

		req.Equal(payload.Labels["foo"], "baz",
			"labels must contain foo with value baz")
		req.NotContains(payload.Labels, "bar",
			"labels must not contain bar")
		req.Equal(payload.Labels["baz"], "123",
			"labels must contain baz with value 123")
		req.Equal(payload.Labels, helpers.LoadLabelsFromStore(t, service.DefaultStore, payload.LabelResourceKind(), payload.ID),
			"response must match stored labels")
	})

	t.Run("search", func(t *testing.T) {
		if ID == 0 {
			t.Skip("label/create test not ran")
		}

		var (
			req = require.New(t)
			set = types.ModuleSet{}
		)

		helpers.SearchWithLabelsViaAPI(h.apiInit(), t,
			fmt.Sprintf("/namespace/%d/module/", ns.ID),
			&set,
			url.Values{"labels": []string{"baz=123"}},
		)
		req.NotEmpty(set)
		req.NotNil(set.FindByID(ID).Labels)
	})

	t.Run("field create", func(t *testing.T) {
		if ID == 0 {
			t.Skip("label/create test not ran")
		}

		var (
			ctx      = context.Background()
			req      = require.New(t)
			mod, err = store.LookupComposeModuleByID(ctx, service.DefaultStore, ID)

			payload = struct{ Response *types.Module }{}
		)

		req.NoError(err)
		req.NotNil(mod)
		mod.Fields = append(mod.Fields, &types.ModuleField{
			Kind:   "String",
			Name:   "labeled",
			Label:  "",
			Labels: map[string]string{"fldfoo": "fldbar"},
		})

		h.apiInit().
			Post(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, ID)).
			JSON(helpers.JSON(mod)).
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			Assert(jsonpath.Equal(`$.response.fields[0].labels.fldfoo`, "fldbar")).
			End().
			JSON(&payload)

		fieldID = payload.Response.Fields[0].ID

	})

	t.Run("field update", func(t *testing.T) {
		if fieldID == 0 {
			t.Skip("label/field create test not ran")
		}

		var (
			ctx      = context.Background()
			req      = require.New(t)
			mod, err = store.LookupComposeModuleByID(ctx, service.DefaultStore, ID)
		)

		req.NoError(err)
		req.NotNil(mod)
		mod.Fields = append(mod.Fields, &types.ModuleField{
			ID:     fieldID,
			Kind:   "String",
			Name:   "labeled",
			Label:  "",
			Labels: map[string]string{"fldfoo": "fldbaz"},
		})

		h.apiInit().
			Post(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, ID)).
			JSON(helpers.JSON(mod)).
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			Assert(jsonpath.Equal(`$.response.fields[0].labels.fldfoo`, "fldbaz")).
			End()
	})
}
