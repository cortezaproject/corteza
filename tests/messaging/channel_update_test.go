package messaging

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/messaging/rest/request"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func (h helper) chUpdate(ch *request.ChannelUpdate) *apitest.Response {
	payload, err := json.Marshal(ch)
	h.a.NoError(err)

	return h.testAPI().
		Put(fmt.Sprintf("/channels/%v", ch.ChannelID)).
		JSON(string(payload)).
		Expect(h.t).
		Status(http.StatusOK)
}

func channelToRequest(ch *types.Channel) *request.ChannelUpdate {
	req := &request.ChannelUpdate{
		ChannelID:        ch.ID,
		Name:             ch.Name,
		Topic:            ch.Topic,
		MembershipPolicy: ch.MembershipPolicy,
		Type:             ch.Type.String(),
	}

	return req
}

func TestChannelUpdateNonexistent(t *testing.T) {
	h := newHelper(t)

	req := &request.ChannelUpdate{
		ChannelID: factory.Sonyflake.NextID(),
		Name:      "some name",
	}

	h.chUpdate(req).
		Assert(helpers.AssertError("messaging.repository.ChannelNotFound")).
		End()

}

func TestChannelUpdateDenied(t *testing.T) {
	h := newHelper(t)
	ch := h.makePublicCh()

	h.deny(ch.PermissionResource(), "update")

	req := channelToRequest(ch)
	req.Name = "Updated name"

	h.chUpdate(req).
		Assert(helpers.AssertError("messaging.service.NoPermissions")).
		End()
}

func TestChannelUpdate(t *testing.T) {
	h := newHelper(t)
	ch := h.makePublicCh()

	h.allow(ch.PermissionResource(), "update")

	req := channelToRequest(ch)
	req.Name = "Updated name"

	h.chUpdate(req).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Equal(`$.response.name`, req.Name)).
		Assert(jsonpath.Equal(`$.response.channelID`, strconv.FormatUint(ch.ID, 10))).
		End()
}
