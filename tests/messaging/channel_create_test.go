package messaging

import (
	"net/http"
	"strconv"
	"testing"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"

	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func TestChannelCreateDenied(t *testing.T) {
	h := newHelper(t)
	h.deny(types.MessagingPermissionResource, "channel.public.create")

	h.testAPI().
		Post("/channels/").
		JSON(`{"name":"test channel","type":"public"}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertError("messaging.service.NoPermissions")).
		End()
}

func TestChannelCreate(t *testing.T) {
	h := newHelper(t)
	h.allow(types.MessagingPermissionResource, "channel.public.create")

	h.testAPI().
		Post("/channels/").
		JSON(`{"name":"test channel","type":"public"}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(`$.response.name`)).
		Assert(jsonpath.Present(`$.response.channelID`)).
		// Creator should be a member
		Assert(jsonpath.Contains(`$.response.members`, strconv.FormatUint(h.cUser.ID, 10))).
		End()
}
