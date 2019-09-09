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

	ch := h.makePublicCh()

	h.testAPI().
		Get(fmt.Sprintf("/channels/%d/members", ch.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestChannelMemberJoinSelf(t *testing.T) {
	h := newHelper(t)

	ch := h.makePublicCh()

	h.testAPI().
		Put(fmt.Sprintf("/channels/%d/members/%d", ch.ID, h.cUser.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	h.chAssertMember(ch, h.cUser, types.ChannelMembershipTypeMember)
}

func TestChannelMemberLeaveSelf(t *testing.T) {
	h := newHelper(t)

	ch := h.makePublicCh()
	h.makeMember(ch, h.cUser)

	h.testAPI().
		Delete(fmt.Sprintf("/channels/%d/members/%d", ch.ID, h.cUser.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	h.chAssertNotMember(ch, h.cUser)
}

func TestChannelMemberInvite(t *testing.T) {
	t.SkipNow()

	h := newHelper(t)
	ch := h.makePublicCh()
	h.allow(ch.PermissionResource(), "members.manage")

	invitee := &sysType.User{ID: factory.Sonyflake.NextID()}

	h.testAPI().
		Put(fmt.Sprintf("/channels/%d/members/%d", ch.ID, invitee.ID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	h.chAssertMember(ch, invitee, types.ChannelMembershipTypeInvitee)
}
