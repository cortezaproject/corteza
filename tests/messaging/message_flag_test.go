package messaging

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"

	"github.com/cortezaproject/corteza-server/messaging/types"
)

func (h helper) apiMessageSetFlag(msg *types.Message, method, flag string) *apitest.Response {
	return h.apiInit().
		Method(method).
		URL(fmt.Sprintf("/channels/%d/messages/%d/%s", msg.ChannelID, msg.ID, flag)).
		Expect(h.t).
		Status(http.StatusOK)
}

func TestMessageFlag(t *testing.T) {
	h := newHelper(t)
	msg := h.makeMessage("flag target", h.repoMakePublicCh(), h.cUser)

	initialState := h.lookupFlagByMessageID(msg.ID)
	h.a.Len(initialState, 0)
	h.a.False(initialState.IsPinned())
	h.a.False(initialState.IsBookmarked(h.cUser.ID))

	// Pin flag (for everyone)

	h.apiMessageSetFlag(msg, "POST", "pin").
		End()
	h.a.True(h.lookupFlagByMessageID(msg.ID).IsPinned())

	h.apiMessageSetFlag(msg, "DELETE", "pin").
		End()
	h.a.False(h.lookupFlagByMessageID(msg.ID).IsPinned())

	// Bookmark flag (per user)

	h.apiMessageSetFlag(msg, "POST", "bookmark").
		End()
	h.a.True(h.lookupFlagByMessageID(msg.ID).IsBookmarked(h.cUser.ID))

	h.apiMessageSetFlag(msg, "DELETE", "bookmark").
		End()
	h.a.False(h.lookupFlagByMessageID(msg.ID).IsBookmarked(h.cUser.ID))

	// Custom flags (aka reactions)
	hasReaction := func(flag string) bool {
		ff, _ := h.lookupFlagByMessageID(msg.ID).Filter(func(f *types.MessageFlag) (b bool, e error) {
			return f.Flag == flag && f.UserID == h.cUser.ID && f.DeletedAt == nil, nil
		})

		return len(ff) > 0
	}

	h.apiMessageSetFlag(msg, "POST", "reaction/foo").
		End()

	h.a.True(hasReaction("foo"), "expecting message to have reaction")

	h.apiMessageSetFlag(msg, "DELETE", "reaction/foo").
		End()

	h.a.False(hasReaction("foo"), "expecting message not to have reaction")
}
