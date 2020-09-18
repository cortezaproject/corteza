package messaging

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"

	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func (h helper) apiChCreate(name string, t types.ChannelType) *apitest.Response {
	payload, err := json.Marshal(struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}{name, t.String()})

	h.a.NoError(err)

	return h.apiInit().
		Post("/channels/").
		JSON(string(payload)).
		Expect(h.t).
		Status(http.StatusOK)
}

func (h helper) apiChPubCreate(name string) *apitest.Response {
	return h.apiChCreate(name, types.ChannelTypePublic)
}

func TestChannelCreateDenied(t *testing.T) {
	h := newHelper(t)
	h.deny(types.MessagingRBACResource, "channel.public.create")

	h.apiChPubCreate("should not be created").
		Assert(helpers.AssertError("not allowed to create channels")).
		End()
}

func TestChannelCreate(t *testing.T) {
	h := newHelper(t)
	h.allow(types.MessagingRBACResource, "channel.public.create")

	h.apiChPubCreate("test channel").
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(`$.response.name`)).
		Assert(jsonpath.Present(`$.response.channelID`)).
		// Creator should be a member
		Assert(jsonpath.Contains(`$.response.members`, strconv.FormatUint(h.cUser.ID, 10))).
		End()
}

func TestChannelCreateWithShortName(t *testing.T) {
	h := newHelper(t)
	h.allow(types.MessagingRBACResource, "channel.public.create")

	h.apiChPubCreate("").
		Status(http.StatusOK).
		Assert(helpers.AssertError("name not set")).
		End()
}

func TestChannelCreateWithLongName(t *testing.T) {
	h := newHelper(t)
	h.allow(types.MessagingRBACResource, "channel.public.create")

	h.apiChPubCreate(strings.Repeat("X ", 1000)).
		Status(http.StatusOK).
		Assert(helpers.AssertError("name too long")).
		End()

}
