package compose

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza/server/compose/service"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/tests/helpers"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"net/http"
	"testing"
	"time"
)

func (h helper) clearAttachment() {
	h.clearNamespaces()
	h.noError(store.TruncateAttachments(context.Background(), service.DefaultStore))
}

func (h helper) repoMakeAttachment(namespaceID uint64, ss ...string) *types.Attachment {
	var res = &types.Attachment{
		ID:          id.Next(),
		CreatedAt:   time.Now(),
		Kind:        types.RecordAttachment,
		NamespaceID: namespaceID,
	}

	if len(ss) > 0 {
		res.Name = ss[0]
	} else {
		res.Name = "n_" + rs()
	}

	h.a.NoError(store.CreateComposeAttachment(context.Background(), service.DefaultStore, res))

	return res
}

func (h helper) lookupAttachmentByID(ID uint64) *types.Attachment {
	res, err := store.LookupComposeAttachmentByID(context.Background(), service.DefaultStore, ID)
	h.noError(err)
	return res
}

func TestAttachmentRead(t *testing.T) {
	h := newHelper(t)
	h.clearAttachment()

	ns := h.makeNamespace("some-namespace")
	a := h.repoMakeAttachment(ns.ID)
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/attachment/%s/%d", ns.ID, a.Kind, a.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.name`, a.Name)).
		Assert(jsonpath.Equal(`$.response.attachmentID`, fmt.Sprintf("%d", a.ID))).
		End()
}

func TestAttachmentDelete(t *testing.T) {
	h := newHelper(t)
	h.clearAttachment()

	ns := h.makeNamespace("some-namespace")
	a := h.repoMakeAttachment(ns.ID)
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "update")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/attachment/%s/%d", ns.ID, a.Kind, a.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	a = h.lookupAttachmentByID(a.ID)
	h.a.NotNil(a)
	h.a.NotNil(a.DeletedAt)
}

func TestAttachmentReadForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearAttachment()

	ns := h.makeNamespace("some-namespace")
	a := h.repoMakeAttachment(ns.ID)
	helpers.DenyMe(h, types.NamespaceRbacResource(0), "read")

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/attachment/%s/%d", ns.ID, a.Kind, a.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("attachment.errors.notAllowedToReadNamespace")).
		End()
}

func TestAttachmentReadForeignNamespace(t *testing.T) {
	h := newHelper(t)
	h.clearAttachment()

	ns := h.makeNamespace("some-namespace")
	foreign := h.makeNamespace("foreign-namespace")
	a := h.repoMakeAttachment(foreign.ID)
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/attachment/%s/%d", ns.ID, a.Kind, a.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("attachment.errors.notFound")).
		End()
}

func TestAttachmentReadSpoofedKind(t *testing.T) {
	h := newHelper(t)
	h.clearAttachment()

	ns := h.makeNamespace("some-namespace")
	a := h.repoMakeAttachment(ns.ID)
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/attachment/%s/%d", ns.ID, types.PageAttachment, a.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("attachment.errors.notFound")).
		End()
}

func TestAttachmentDeleteForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearAttachment()

	ns := h.makeNamespace("some-namespace")
	a := h.repoMakeAttachment(ns.ID)
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.DenyMe(h, types.NamespaceRbacResource(0), "update")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/attachment/%s/%d", ns.ID, a.Kind, a.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("attachment.errors.notAllowedToUpdateNamespace")).
		End()

	a = h.lookupAttachmentByID(a.ID)
	h.a.NotNil(a)
	h.a.Nil(a.DeletedAt)
}

func TestAttachmentDeleteForeignNamespace(t *testing.T) {
	h := newHelper(t)
	h.clearAttachment()

	ns := h.makeNamespace("some-namespace")
	foreign := h.makeNamespace("foreign-namespace")
	a := h.repoMakeAttachment(foreign.ID)
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "update")

	h.apiInit().
		Delete(fmt.Sprintf("/namespace/%d/attachment/%s/%d", ns.ID, a.Kind, a.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("attachment.errors.notFound")).
		End()

	a = h.lookupAttachmentByID(a.ID)
	h.a.NotNil(a)
	h.a.Nil(a.DeletedAt)
}

// Private attachments must not be served when public kind is put in the URL
func TestAttachmentOriginalSpoofedKindAnonymous(t *testing.T) {
	h := newHelper(t)
	h.clearAttachment()

	ns := h.makeNamespace("some-namespace")
	a := h.repoMakeAttachment(ns.ID)

	for _, kind := range []string{types.PageAttachment, types.IconAttachment, types.NamespaceAttachment} {
		InitTestApp()

		apitest.
			New().
			Handler(r).
			Get(fmt.Sprintf("/namespace/%d/attachment/%s/%d/original/secret.pdf", ns.ID, kind, a.ID)).
			Expect(t).
			Status(http.StatusNotFound).
			End()
	}
}

// Private attachments must not be served without a signature
func TestAttachmentOriginalUnsignedAnonymous(t *testing.T) {
	h := newHelper(t)
	h.clearAttachment()

	ns := h.makeNamespace("some-namespace")
	a := h.repoMakeAttachment(ns.ID)

	InitTestApp()

	apitest.
		New().
		Handler(r).
		Get(fmt.Sprintf("/namespace/%d/attachment/%s/%d/original/secret.pdf", ns.ID, a.Kind, a.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("missing signature")).
		End()
}
