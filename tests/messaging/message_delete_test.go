package messaging

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/messaging/types"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func TestMessagesDelete(t *testing.T) {
	h := newHelper(t)
	msg := h.repoMakeMessage("old", h.repoMakePublicCh(), h.cUser)

	h.apiInit().
		Delete(fmt.Sprintf("/channels/%d/messages/%d", msg.ChannelID, msg.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	_, err := h.repoMessage().FindByID(msg.ID)
	h.a.EqualError(err, "messaging.repository.MessageNotFound")
}

func TestMessagesDelete_forbiden(t *testing.T) {
	h := newHelper(t)

	msg := h.repoMakeMessage("old", h.repoMakePublicCh(), h.cUser)
	h.deny(types.ChannelPermissionResource.AppendWildcard(), "message.update.own")

	h.apiInit().
		Delete(fmt.Sprintf("/channels/%d/messages/%d", msg.ChannelID, msg.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("messaging.service.NoPermissions")).
		End()

	_, err := h.repoMessage().FindByID(msg.ID)
	h.a.Nil(err)
}

func TestMessagesDeleteOwnThreadMessage(t *testing.T) {
	// Covers deleting messages that reply to my own thread
	h := newHelper(t)

	msg := h.repoMakeMessage("old", h.repoMakePublicCh(), h.cUser)
	thrMsg := h.apiMessageCreateReply("thr1", msg)

	h.apiInit().
		Delete(fmt.Sprintf("/channels/%d/messages/%d", msg.ChannelID, thrMsg.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	_, err := h.repoMessage().FindByID(thrMsg.ID)
	h.a.EqualError(err, "messaging.repository.MessageNotFound")
}

func TestMessagesDeleteOwnThreadMessage_forbidenNotOwner(t *testing.T) {
	// Covers deleting someone elses messages that reply to my own thread
	h := newHelper(t)

	msg := h.repoMakeMessage("old", h.repoMakePublicCh(), h.cUser)

	nh := newHelper(t)
	thrMsg := nh.apiMessageCreateReply("thr1", msg)

	h.apiInit().
		Delete(fmt.Sprintf("/channels/%d/messages/%d", msg.ChannelID, thrMsg.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("messaging.service.NoPermissions")).
		End()

	_, err := h.repoMessage().FindByID(thrMsg.ID)
	h.a.Nil(err)
}

func TestMessagesDeleteThreadMessage(t *testing.T) {
	// Covers deleting messages that reply to someone elses thread
	h := newHelper(t)

	u := &sysTypes.User{ID: 10}
	msg := h.repoMakeMessage("old", h.repoMakePublicCh(), u)
	thrMsg := h.apiMessageCreateReply("thr1", msg)

	h.apiInit().
		Delete(fmt.Sprintf("/channels/%d/messages/%d", msg.ChannelID, thrMsg.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	_, err := h.repoMessage().FindByID(thrMsg.ID)
	h.a.EqualError(err, "messaging.repository.MessageNotFound")
}

func TestMessagesDeleteThreadMessage_forbidenNotOwner(t *testing.T) {
	// Covers deleting someone elses messages that reply to someone elses thread
	h := newHelper(t)

	u := &sysTypes.User{ID: 10}
	msg := h.repoMakeMessage("old", h.repoMakePublicCh(), u)

	nh := newHelper(t)
	thrMsg := nh.apiMessageCreateReply("thr1", msg)

	h.apiInit().
		Delete(fmt.Sprintf("/channels/%d/messages/%d", msg.ChannelID, thrMsg.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("messaging.service.NoPermissions")).
		End()

	_, err := h.repoMessage().FindByID(thrMsg.ID)
	h.a.Nil(err)
}
