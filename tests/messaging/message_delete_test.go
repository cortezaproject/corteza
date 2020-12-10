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
	msg := h.makeMessage("old", h.repoMakePublicCh(), h.cUser)

	h.apiInit().
		Delete(fmt.Sprintf("/channels/%d/messages/%d", msg.ChannelID, msg.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	m, err := h.lookupMessageByID(msg.ID)
	h.a.NoError(err)
	h.a.NotNil(m.DeletedAt)
}

func TestMessagesDelete_forbidden(t *testing.T) {
	t.Skipf("inconsistency between store backends")
	h := newHelper(t)

	msg := h.makeMessage("old", h.repoMakePublicCh(), h.cUser)
	h.deny(types.ChannelRBACResource.AppendWildcard(), "message.update.own")

	h.apiInit().
		Delete(fmt.Sprintf("/channels/%d/messages/%d", msg.ChannelID, msg.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("failed to complete transaction: messaging.service.NoPermissions")).
		End()

	_, err := h.lookupMessageByID(msg.ID)
	h.a.Nil(err)
}

func TestMessagesDeleteOwnThreadMessage(t *testing.T) {
	// Covers deleting messages that reply to my own thread
	h := newHelper(t)

	msg := h.makeMessage("old", h.repoMakePublicCh(), h.cUser)
	thrMsg := h.apiMessageCreateReply("thr1", msg)

	h.apiInit().
		Delete(fmt.Sprintf("/channels/%d/messages/%d", msg.ChannelID, thrMsg.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	m, err := h.lookupMessageByID(thrMsg.ID)
	h.a.NoError(err)
	h.a.NotNil(m.DeletedAt)
}

func TestMessagesDeleteOwnThreadMessage_forbiddenNotOwner(t *testing.T) {
	t.Skipf("inconsistency between store backends")

	// Covers deleting someone elses messages that reply to my own thread
	h := newHelper(t)

	msg := h.makeMessage("old", h.repoMakePublicCh(), h.cUser)

	nh := newHelper(t)
	thrMsg := nh.apiMessageCreateReply("thr1", msg)

	h.apiInit().
		Delete(fmt.Sprintf("/channels/%d/messages/%d", msg.ChannelID, thrMsg.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("failed to complete transaction: messaging.service.NoPermissions")).
		End()

	_, err := h.lookupMessageByID(thrMsg.ID)
	h.a.Nil(err)
}

func TestMessagesDeleteThreadMessage(t *testing.T) {
	// Covers deleting messages that reply to someone elses thread
	h := newHelper(t)

	u := &sysTypes.User{ID: 10}
	msg := h.makeMessage("old", h.repoMakePublicCh(), u)
	thrMsg := h.apiMessageCreateReply("thr1", msg)

	h.apiInit().
		Delete(fmt.Sprintf("/channels/%d/messages/%d", msg.ChannelID, thrMsg.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	m, err := h.lookupMessageByID(thrMsg.ID)
	h.a.NoError(err)
	h.a.NotNil(m.DeletedAt)
}

func TestMessagesDeleteThreadMessage_forbiddenNotOwner(t *testing.T) {
	t.Skipf("inconsistency between store backends")

	// Covers deleting someone else messages that reply to someone else thread
	h := newHelper(t)

	u := &sysTypes.User{ID: 10}
	msg := h.makeMessage("old", h.repoMakePublicCh(), u)

	nh := newHelper(t)
	thrMsg := nh.apiMessageCreateReply("thr1", msg)

	h.apiInit().
		Delete(fmt.Sprintf("/channels/%d/messages/%d", msg.ChannelID, thrMsg.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("failed to complete transaction: messaging.service.NoPermissions")).
		End()

	_, err := h.lookupMessageByID(thrMsg.ID)
	h.a.Nil(err)
}
