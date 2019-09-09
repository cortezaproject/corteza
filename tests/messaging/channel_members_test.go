package messaging

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/messaging/types"
	sysType "github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
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
	h.allow(ch.PermissionResource(), "members.manage")

	invitee := &sysType.User{ID: factory.Sonyflake.NextID()}

	h.apiInit().
		Put(fmt.Sprintf("/channels/%d/members/%d", ch.ID, invitee.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	h.repoChAssertMember(ch, invitee, types.ChannelMembershipTypeInvitee)
}
