package compose

import (
	"context"
	"fmt"
	"github.com/crusttech/human/server/compose/service"
	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"net/http"
	"testing"
	"time"
)

func (h helper) clearAttachment() {
	h.clearNamespaces()
	h.noError(store.TruncateAttachments(context.Background(), service.DefaultStore))
}

func (h helper) repoMakeAttachment(ss ...string) *types.Attachment {
	var res = &types.Attachment{
		ID:        id.Next(),
		CreatedAt: time.Now(),
		Kind:      types.RecordAttachment,
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
	a := h.repoMakeAttachment()

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
	a := h.repoMakeAttachment()

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

func TestAttachmentServedOnlyUnderItsOwnKind(t *testing.T) {
	h := newHelper(t)
	h.clearAttachment()

	ns := h.makeNamespace("attachment kind namespace")
	page := h.repoMakePage(ns, "some-page")

	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.AllowMe(h, types.PageRbacResource(0, 0), "read", "update")

	var uploaded struct {
		Response struct {
			ID uint64 `json:"attachmentID,string"`
		} `json:"response"`
	}

	helpers.InitFileUpload(t, h.apiInit(),
		fmt.Sprintf("/namespace/%d/page/%d/attachment", ns.ID, page.ID),
		nil,
		[]byte("private"),
		"private.txt",
		"text/plain",
	).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End().
		JSON(&uploaded)

	a := h.lookupAttachmentByID(uploaded.Response.ID)
	unsigned := fmt.Sprintf("/namespace/%d/attachment/%s/%d/original/private.txt", ns.ID, types.PageAttachment, a.ID)

	h.apiInit().Get(unsigned).Expect(t).Status(http.StatusOK).Body("private").End()

	a.Kind = types.RecordAttachment
	h.noError(store.UpdateComposeAttachment(context.Background(), service.DefaultStore, a))

	h.apiInit().Get(unsigned).Expect(t).Status(http.StatusNotFound).End()

	h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/attachment/%s/%d/original/private.txt", ns.ID, types.RecordAttachment, a.ID)).
		Query("sign", auth.DefaultSigner.Sign(h.cUser.ID, ns.ID, a.ID)).
		Query("userID", fmt.Sprintf("%d", h.cUser.ID)).
		Expect(t).
		Status(http.StatusOK).
		Body("private").
		End()
}
