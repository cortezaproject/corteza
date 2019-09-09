package messaging

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func TestMessagesDelete(t *testing.T) {
	h := newHelper(t)
	msg := h.makeMessage("old", h.makePublicCh(), h.cUser)

	h.testAPI().
		Delete(fmt.Sprintf("/channels/%d/messages/%d", msg.ChannelID, msg.ID)).
		JSON(`{"message":"new"}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	_, err := h.repoMessage().FindByID(msg.ID)
	h.a.EqualError(err, "messaging.repository.MessageNotFound")
}
