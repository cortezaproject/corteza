package system

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/service"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/cortezaproject/corteza/server/tests/helpers"
)

func (h helper) repoMakeAttachmentOf(kind string, ownerID uint64) *types.Attachment {
	res := &types.Attachment{
		ID:        id.Next(),
		OwnerID:   ownerID,
		Kind:      kind,
		Name:      "n_" + rs(),
		CreatedAt: time.Now(),
	}

	h.a.NoError(store.CreateAttachment(context.Background(), service.DefaultStore, res))
	return res
}

func TestAttachmentDeleteSettingsForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearAttachments()

	a := h.repoMakeAttachmentOf(types.AttachmentKindSettings, id.Next())

	h.apiInit().
		Delete(fmt.Sprintf("/attachment/%s/%d", a.Kind, a.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("attachment.errors.notAllowedToDelete")).
		End()

	h.a.Nil(h.lookupAttachmentByID(a.ID).DeletedAt)
}

func TestAttachmentDeleteSettings(t *testing.T) {
	h := newHelper(t)
	h.clearAttachments()

	a := h.repoMakeAttachmentOf(types.AttachmentKindSettings, id.Next())
	helpers.AllowMe(h, types.ComponentRbacResource(), "settings.manage")

	h.apiInit().
		Delete(fmt.Sprintf("/attachment/%s/%d", a.Kind, a.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	h.a.NotNil(h.lookupAttachmentByID(a.ID).DeletedAt)
}

func TestAttachmentDeleteForeignAvatarForbidden(t *testing.T) {
	h := newHelper(t)
	h.clearAttachments()

	a := h.repoMakeAttachmentOf(types.AttachmentKindAvatar, id.Next())

	h.apiInit().
		Delete(fmt.Sprintf("/attachment/%s/%d", a.Kind, a.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("attachment.errors.notAllowedToDelete")).
		End()

	h.a.Nil(h.lookupAttachmentByID(a.ID).DeletedAt)
}

func TestAttachmentDeleteOwnAvatar(t *testing.T) {
	h := newHelper(t)
	h.clearAttachments()

	a := h.repoMakeAttachmentOf(types.AttachmentKindAvatar, h.cUser.ID)

	h.apiInit().
		Delete(fmt.Sprintf("/attachment/%s/%d", a.Kind, a.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	h.a.NotNil(h.lookupAttachmentByID(a.ID).DeletedAt)
}

func TestAttachmentDeleteForeignAvatarWithUserUpdate(t *testing.T) {
	h := newHelper(t)
	h.clearAttachments()

	a := h.repoMakeAttachmentOf(types.AttachmentKindAvatar, id.Next())
	helpers.AllowMe(h, types.UserRbacResource(0), "update")

	h.apiInit().
		Delete(fmt.Sprintf("/attachment/%s/%d", a.Kind, a.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	h.a.NotNil(h.lookupAttachmentByID(a.ID).DeletedAt)
}

// Kind from the request must match the attachment
func TestAttachmentDeleteSpoofedKind(t *testing.T) {
	h := newHelper(t)
	h.clearAttachments()

	a := h.repoMakeAttachmentOf(types.AttachmentKindSettings, id.Next())

	h.apiInit().
		Delete(fmt.Sprintf("/attachment/%s/%d", types.AttachmentKindAvatar, a.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("attachment.errors.notFound")).
		End()

	h.a.Nil(h.lookupAttachmentByID(a.ID).DeletedAt)
}
