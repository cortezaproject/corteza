package compose

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/crusttech/human/server/compose/service"
	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/store"
	systemTypes "github.com/crusttech/human/server/system/types"
	"github.com/crusttech/human/server/tests/helpers"
)

func (h helper) translationsOf(resource string) systemTypes.ResourceTranslationSet {
	set, _, err := store.SearchResourceTranslations(context.Background(), service.DefaultStore, systemTypes.ResourceTranslationFilter{Resource: resource})
	h.noError(err)
	return set
}

// Translations can be changed by users that can update the translated resource

func TestModuleTranslationsUpdateForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	ns := h.makeNamespace("some-namespace")
	m := h.makeModule(ns, "some-module")

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "read")
	helpers.DenyMe(h, types.ModuleRbacResource(0, 0), "update")

	h.apiInit().
		Patch(fmt.Sprintf("/namespace/%d/module/%d/translation", ns.ID, m.ID)).
		Header("Accept", "application/json").
		JSON(fmt.Sprintf(`{"translations":[{"resource":"%s","lang":"en","key":"name","message":"defaced"}]}`, m.ResourceTranslation())).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("module.errors.notAllowedToUpdate")).
		End()

	h.a.Empty(h.translationsOf(m.ResourceTranslation()))
}

// Translations sent to one resource must not change another one
func TestModuleTranslationsUpdateForeignResource(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	ns := h.makeNamespace("some-namespace")
	mine := h.makeModule(ns, "my-module")
	foreign := h.makeModule(h.makeNamespace("foreign-namespace"), "foreign-module")

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "read")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "update")

	h.apiInit().
		Patch(fmt.Sprintf("/namespace/%d/module/%d/translation", ns.ID, mine.ID)).
		Header("Accept", "application/json").
		JSON(fmt.Sprintf(`{"translations":[{"resource":"%s","lang":"en","key":"name","message":"defaced"}]}`, foreign.ResourceTranslation())).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertErrorP("translated resource")).
		End()

	h.a.Empty(h.translationsOf(foreign.ResourceTranslation()))
}

func TestPageTranslationsUpdateForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearPages()

	ns := h.makeNamespace("some-namespace")
	p := h.repoMakePage(ns, "some-page")

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.AllowMe(h, types.PageRbacResource(0, 0), "read")
	helpers.DenyMe(h, types.PageRbacResource(0, 0), "update")

	h.apiInit().
		Patch(fmt.Sprintf("/namespace/%d/page/%d/translation", ns.ID, p.ID)).
		Header("Accept", "application/json").
		JSON(fmt.Sprintf(`{"translations":[{"resource":"%s","lang":"en","key":"title","message":"defaced"}]}`, p.ResourceTranslation())).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("page.errors.notAllowedToUpdate")).
		End()

	h.a.Empty(h.translationsOf(p.ResourceTranslation()))
}

// Users that can update the module can change its translations
func TestModuleTranslationsUpdate(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	ns := h.makeNamespace("translations-ns-" + rs())
	m := h.makeModule(ns, "some-module")

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "read", "update")

	h.apiInit().
		Patch(fmt.Sprintf("/namespace/%d/module/%d/translation", ns.ID, m.ID)).
		Header("Accept", "application/json").
		JSON(fmt.Sprintf(`{"translations":[{"resource":"%s","lang":"und","key":"name","message":"translated"}]}`, m.ResourceTranslation())).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	h.a.NotEmpty(h.translationsOf(m.ResourceTranslation()))
}

// The translation manager permission works without update rights on the module
func TestModuleTranslationsUpdateWithManage(t *testing.T) {
	h := newHelper(t)
	h.clearModules()

	ns := h.makeNamespace("translations-ns-" + rs())
	m := h.makeModule(ns, "some-module")

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "read")
	helpers.DenyMe(h, types.ModuleRbacResource(0, 0), "update")
	helpers.AllowMe(h, types.ComponentRbacResource(), "resource-translations.manage")

	h.apiInit().
		Patch(fmt.Sprintf("/namespace/%d/module/%d/translation", ns.ID, m.ID)).
		Header("Accept", "application/json").
		JSON(fmt.Sprintf(`{"translations":[{"resource":"%s","lang":"en","key":"name","message":"translated"}]}`, m.ResourceTranslation())).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	h.a.NotEmpty(h.translationsOf(m.ResourceTranslation()))
}
