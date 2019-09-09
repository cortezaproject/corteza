package messaging

import (
	"fmt"
	"net/http"
	"testing"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"

	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func TestMessagesUpdate(t *testing.T) {
	h := newHelper(t)
	msg := h.repoMakeMessage("old", h.repoMakePublicCh(), h.cUser)

	h.apiInit().
		Put(fmt.Sprintf("/channels/%d/messages/%d", msg.ChannelID, msg.ID)).
		JSON(`{"message":"new"}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.message`, `new`)).
		End()

	m := h.repoMsgExistingLoad(msg.ID)
	h.a.Equal(`new`, m.Message)
}
