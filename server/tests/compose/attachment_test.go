package compose

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/tests/helpers"
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
