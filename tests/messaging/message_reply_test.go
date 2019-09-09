package messaging

import (
	"fmt"
	"net/http"
	"testing"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"

	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func TestMessagesReply(t *testing.T) {
	h := newHelper(t)
	msg := h.makeMessage("old", h.makePublicCh(), h.cUser)

	rval := struct {
		Response struct {
			ID uint64 `json:"messageID,string"`
		}
	}{}

	h.testAPI().
		Debug().
		Post(fmt.Sprintf("/channels/%d/messages/%d/replies", msg.ChannelID, msg.ID)).
		JSON(`{"message":"new reply"}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(`$.response.messageID`)).
		Assert(jsonpath.Present(`$.response.replyTo`)).
		Assert(jsonpath.Equal(`$.response.message`, `new reply`)).
		End().
		JSON(&rval)

	m := h.msgExistingLoad(rval.Response.ID)
	h.a.Equal(`new reply`, m.Message)
	h.a.Equal(msg.ID, m.ReplyTo)

}
