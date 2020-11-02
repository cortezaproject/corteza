package messaging

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"

	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func TestChannelSetFlag(t *testing.T) {
	h := newHelper(t)

	ch := h.repoMakePublicCh()
	h.repoMakeMember(ch, h.cUser)

	flagCheck := func(ch *types.Channel, flag string) {
		mm := h.lookupChMembership(ch)
		h.a.Len(mm, 1, "expecting 1 member")
		h.a.Equal(flag, string(mm[0].Flag), "expecting flags to match")
	}

	flagChannel := func(ch *types.Channel, flag string) *apitest.Response {
		return h.apiInit().
			Put(fmt.Sprintf("/channels/%d/flag", ch.ID)).
			FormData("flag", flag).
			Expect(t).
			Status(http.StatusOK)
	}

	flagChannelOK := func(ch *types.Channel, flag string) {
		flagChannel(ch, flag).
			Assert(helpers.AssertNoErrors).
			End()

		flagCheck(ch, flag)
	}

	unflagChannel := func(ch *types.Channel) {
		h.apiInit().
			Delete(fmt.Sprintf("/channels/%d/flag", ch.ID)).
			Expect(t).
			Status(http.StatusOK).
			End()

		flagCheck(ch, string(types.ChannelMembershipFlagNone))
	}

	flagCheck(ch, string(types.ChannelMembershipFlagNone))

	flagChannelOK(ch, string(types.ChannelMembershipFlagPinned))
	flagChannelOK(ch, string(types.ChannelMembershipFlagHidden))
	flagChannelOK(ch, string(types.ChannelMembershipFlagIgnored))

	flagChannel(ch, "foo").
		Header("Accept", "application/json").
		Assert(helpers.AssertError("invalid flag"))
	flagCheck(ch, string(types.ChannelMembershipFlagIgnored))

	unflagChannel(ch)
}
