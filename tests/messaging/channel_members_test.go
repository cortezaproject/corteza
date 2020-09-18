package messaging

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/id"
	sysType "github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"net/http"
	"testing"
)

func TestChannelMemberList(t *testing.T) {
	h := newHelper(t)

	ch := h.repoMakePublicCh()

	h.apiInit().
		Get(fmt.Sprintf("/channels/%d/members", ch.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestChannelMemberJoinSelf(t *testing.T) {
	h := newHelper(t)

	ch := h.repoMakePublicCh()

	h.apiInit().
		Put(fmt.Sprintf("/channels/%d/members/%d", ch.ID, h.cUser.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	h.repoChAssertMember(ch, h.cUser, types.ChannelMembershipTypeMember)
}

func TestChannelMemberLeaveSelf(t *testing.T) {
	h := newHelper(t)

	ch := h.repoMakePublicCh()
	h.repoMakeMember(ch, h.cUser)

	h.apiInit().
		Delete(fmt.Sprintf("/channels/%d/members/%d", ch.ID, h.cUser.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	h.repoChAssertNotMember(ch, h.cUser)
}

func TestChannelMemberInvite(t *testing.T) {
	t.SkipNow()

	h := newHelper(t)
	ch := h.repoMakePublicCh()
	h.allow(ch.RBACResource(), "members.manage")

	invitee := &sysType.User{ID: id.Next()}

	h.apiInit().
		Put(fmt.Sprintf("/channels/%d/members/%d", ch.ID, invitee.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	h.repoChAssertMember(ch, invitee, types.ChannelMembershipTypeInvitee)
}
