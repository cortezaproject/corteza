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
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/require"
)

func (h helper) clearNamespaces() {
	h.noError(store.TruncateComposeAttachments(context.Background(), service.DefaultStore))
	h.noError(store.TruncateComposeNamespaces(context.Background(), service.DefaultStore))
}

func (h helper) makeNamespace(name string) *types.Namespace {
	ns := &types.Namespace{Name: name, Slug: name}
	ns.ID = id.Next()
	ns.CreatedAt = time.Now()
	h.noError(store.CreateComposeNamespace(context.Background(), service.DefaultStore, ns))
	return ns
}

func (h helper) lookupNamespaceByID(ID uint64) *types.Namespace {
	ns, err := store.LookupComposeNamespaceByID(context.Background(), service.DefaultStore, ID)
	h.noError(err)
	return ns
}

func TestNamespaceRead(t *testing.T) {
	h := newHelper(t)
	h.clearNamespaces()

	ns := h.makeNamespace("some-namespace")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d", ns.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.name`, ns.Name)).
		Assert(jsonpath.Equal(`$.response.namespaceID`, fmt.Sprintf("%d", ns.ID))).
		End()
}

func TestNamespaceReadByHandle(t *testing.T) {
	h := newHelper(t)
	h.clearNamespaces()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	ns := h.makeNamespace("some-namespace-" + string(rand.Bytes(20)))

	nsbh, err := service.DefaultNamespace.FindByHandle(h.secCtx(), ns.Slug)

	h.noError(err)
	h.a.NotNil(nsbh)
	h.a.Equal(nsbh.ID, ns.ID)
	h.a.Equal(nsbh.Slug, ns.Slug)
}

func TestNamespaceList(t *testing.T) {
	h := newHelper(t)
	h.clearNamespaces()

	h.makeNamespace("ns1")
	h.makeNamespace("ns2")

	h.apiInit().
		Get("/namespace/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestNamespaceList_filterForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearNamespaces()

	h.makeNamespace("namespace")
	f := h.makeNamespace("namespace_forbidden")

	helpers.DenyMe(h, types.NamespaceRbacResource(f.ID), "read")

	h.apiInit().
		Get("/namespace/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.NotPresent(`$.response.set[? @.name=="namespace_forbidden"]`)).
		End()
}

func TestNamespaceCreateForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearNamespaces()

	h.apiInit().
		Post("/namespace/").
		Header("Accept", "application/json").
		FormData("name", "some-namespace").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to create namespaces")).
		End()
}

func TestNamespaceCreate(t *testing.T) {
	h := newHelper(t)
	h.clearNamespaces()

	helpers.AllowMe(h, types.ComponentRbacResource(), "namespace.create")

	h.apiInit().
		Post("/namespace/").
		FormData("name", "some-namespace").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestNamespaceUpdateForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearNamespaces()

	ns := h.makeNamespace("some-namespace")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d", ns.ID)).
		Header("Accept", "application/json").
		FormData("name", "changed-name").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to update this namespace")).
		End()
}

func TestNamespaceUpdate(t *testing.T) {
	h := newHelper(t)
	h.clearNamespaces()

	ns := h.makeNamespace("some-namespace")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "update")

	h.apiInit().
		Post(fmt.Sprintf("/namespace/%d", ns.ID)).
		FormData("name", "changed-name").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	ns = h.lookupNamespaceByID(ns.ID)
	h.a.NotNil(ns)
	h.a.Equal("changed-name", ns.Name)
}

func TestNamespaceDeleteForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearNamespaces()

	ns := h.makeNamespace("some-namespace")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.DenyMe(h, types.NamespaceRbacResource(0), "delete")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d", ns.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("not allowed to delete this namespace")).
		End()
}

func TestNamespaceDelete(t *testing.T) {
	h := newHelper(t)
	h.clearNamespaces()

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "delete")

	ns := h.makeNamespace("some-namespace")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d", ns.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	ns = h.lookupNamespaceByID(ns.ID)
	h.a.NotNil(ns.DeletedAt)
}

func TestNamespaceLabels(t *testing.T) {
	h := newHelper(t)
	h.clearNamespaces()

	helpers.AllowMe(h, types.ComponentRbacResource(), "namespace.create")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read", "delete", "update")

	var (
		ID uint64
	)

	t.Run("create", func(t *testing.T) {
		var (
			req     = require.New(t)
			payload = &types.Namespace{}
		)

		helpers.SetLabelsViaAPI(h.apiInit(), t,
			"/namespace/",
			types.Namespace{Labels: map[string]string{"foo": "bar", "bar": "42"}},
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
			payload = &types.Namespace{}
		)

		helpers.SetLabelsViaAPI(h.apiInit(), t,
			fmt.Sprintf("/namespace/%d", ID),
			&types.Namespace{Labels: map[string]string{"foo": "baz", "baz": "123"}},
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
			set = types.NamespaceSet{}
		)

		helpers.SearchWithLabelsViaAPI(h.apiInit(), t, "/namespace/", &set, url.Values{"labels": []string{"baz=123"}})
		req.NotEmpty(set)
		req.NotNil(set.FindByID(ID))
		req.NotNil(set.FindByID(ID).Labels)
	})
}
