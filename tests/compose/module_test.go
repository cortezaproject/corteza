package compose

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/require"
)

func (h helper) clearModules() {
	h.clearNamespaces()
	h.noError(store.TruncateComposeModules(context.Background(), service.DefaultStore))
	h.noError(store.TruncateComposeModuleFields(context.Background(), service.DefaultStore))
}

func (h helper) makeModule(ns *types.Namespace, name string, ff ...*types.ModuleField) *types.Module {
	return h.createModule(&types.Module{
		Name:        name,
		NamespaceID: ns.ID,
		Fields:      ff,
		CreatedAt:   time.Now(),
	})
}

func (h helper) createModule(res *types.Module) *types.Module {
	res.ID = id.Next()
	res.CreatedAt = time.Now()
	h.noError(store.CreateComposeModule(context.Background(), service.DefaultStore, res))

	_ = res.Fields.Walk(func(f *types.ModuleField) error {
		f.ID = id.Next()
		f.ModuleID = res.ID
		f.CreatedAt = time.Now()
		return nil
	})

	h.noError(store.CreateComposeModuleField(context.Background(), service.DefaultStore, res.Fields...))

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

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	h.allow(types.ModuleRBACResource.AppendWildcard(), "read")
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

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	h.allow(types.ModuleRBACResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")
	c := h.makeModule(ns, "some-module")

	cbh, err := service.DefaultModule.With(h.secCtx()).FindByHandle(ns.ID, c.Handle)

	h.noError(err)
	h.a.NotNil(cbh)
	h.a.Equal(cbh.ID, c.ID)
	h.a.Equal(cbh.Handle, c.Handle)
}

func TestModuleList(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
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

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")

	h.createModule(&types.Module{
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

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")

	h.makeModule(ns, "module")
	f := h.makeModule(ns, "module_forbiden")

	h.deny(types.ModuleRBACResource.AppendID(f.ID), "read")

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/module/", ns.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.NotPresent(`$.response.set[? @.name=="module_forbiden"]`)).
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
		Assert(helpers.AssertError("not allowed to create modules")).
		End()
}

func TestModuleCreate(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	h.allow(types.NamespaceRBACResource.AppendWildcard(), "module.create")

	ns := h.makeNamespace("some-namespace")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/", ns.ID)).
		FormData("name", "some-module").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestModuleUpdateForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")
	m := h.makeModule(ns, "some-module")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, m.ID)).
		Header("Accept", "application/json").
		FormData("name", "changed-name").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to update this module")).
		End()
}

func TestModuleUpdate(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")
	m := h.makeModule(ns, "some-module")
	h.allow(types.ModuleRBACResource.AppendWildcard(), "update")

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

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")
	m := h.makeModule(ns, "some-module", &types.ModuleField{ID: id.Next(), Kind: "String", Name: "existing"})
	h.allow(types.ModuleRBACResource.AppendWildcard(), "update")

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

func TestModuleFieldsUpdate_removed(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")
	m := h.makeModule(ns, "some-module", &types.ModuleField{ID: id.Next(), Kind: "String", Name: "a"}, &types.ModuleField{ID: id.Next(), Kind: "String", Name: "b"})
	h.allow(types.ModuleRBACResource.AppendWildcard(), "update")

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

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")
	m := h.makeModule(ns, "some-module", &types.ModuleField{ID: id.Next(), Kind: "String", Name: "a"}, &types.ModuleField{ID: id.Next(), Kind: "String", Name: "b"})
	h.makeRecord(m, &types.RecordValue{Name: "a", Value: "va"}, &types.RecordValue{Name: "b", Value: "vb"})
	h.allow(types.ModuleRBACResource.AppendWildcard(), "update")

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

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	h.allow(types.NamespaceRBACResource.AppendWildcard(), "module.create")
	ns := h.makeNamespace("some-namespace")
	h.allow(types.ModuleRBACResource.AppendWildcard(), "read")
	h.allow(types.ModuleRBACResource.AppendWildcard(), "update")

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

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")
	m := h.makeModule(ns, "some-module", &types.ModuleField{ID: id.Next(), Kind: "String", Name: "existing"})
	h.makeRecord(m, &types.RecordValue{Name: "existing", Value: "value"})
	h.allow(types.ModuleRBACResource.AppendWildcard(), "update")

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

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	h.allow(types.ModuleRBACResource.AppendWildcard(), "read")
	ns := h.makeNamespace("some-namespace")
	m := h.makeModule(ns, "some-module")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/module/%d", ns.ID, m.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to delete this module")).
		End()
}

func TestModuleDelete(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	h.allow(types.ModuleRBACResource.AppendWildcard(), "read")
	h.allow(types.ModuleRBACResource.AppendWildcard(), "delete")

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

	h.allow(types.NamespaceRBACResource.AppendWildcard(), "read")
	h.allow(types.NamespaceRBACResource.AppendWildcard(), "module.create")
	h.allow(types.ModuleRBACResource.AppendWildcard(), "read")
	h.allow(types.ModuleRBACResource.AppendWildcard(), "update")
	h.allow(types.ModuleRBACResource.AppendWildcard(), "delete")

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
