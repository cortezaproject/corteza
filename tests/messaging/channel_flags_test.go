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

	flagCheck := func(ID uint64, flag string) {
		mm, err := h.repoChMember().Find(types.ChannelMemberFilter{
			ChannelID: []uint64{ID},
			MemberID:  []uint64{h.cUser.ID},
		})
		h.a.NoError(err)
		h.a.Len(mm, 1, "expecting 1 member")
		h.a.Equal(flag, string(mm[0].Flag), "expecting flags to match")
	}

	flagChannel := func(ID uint64, flag string) *apitest.Response {
		return h.apiInit().
			Put(fmt.Sprintf("/channels/%d/flag", ID)).
			FormData("flag", flag).
			Expect(t).
			Status(http.StatusOK)
	}

	flagChannelOK := func(ID uint64, flag string) {
		flagChannel(ID, flag).
			Assert(helpers.AssertNoErrors).
			End()

		flagCheck(ID, flag)
	}

	unflagChannel := func(ID uint64) {
		h.apiInit().
			Delete(fmt.Sprintf("/channels/%d/flag", ID)).
			Expect(t).
			Status(http.StatusOK).
			End()

		flagCheck(ID, string(types.ChannelMembershipFlagNone))
	}

	flagCheck(ch.ID, string(types.ChannelMembershipFlagNone))

	flagChannelOK(ch.ID, string(types.ChannelMembershipFlagPinned))
	flagChannelOK(ch.ID, string(types.ChannelMembershipFlagHidden))
	flagChannelOK(ch.ID, string(types.ChannelMembershipFlagIgnored))

	flagChannel(ch.ID, "foo").
		Assert(helpers.AssertError("invalid flag"))
	flagCheck(ch.ID, string(types.ChannelMembershipFlagIgnored))

	unflagChannel(ch.ID)
}
