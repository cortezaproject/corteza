package messaging

import (
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func TestChannelList(t *testing.T) {

	NewApiTest("get list of channels", &types.User{ID: 5}).
		Get("/channels/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.NoErrors).
		End()
}

func TestChannelRead(t *testing.T) {
	NewApiTest("find single channel by ID", &types.User{ID: 5}).
		Get("/channels/324234").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.NoErrors).
		End()
}

func TestChannelCreate(t *testing.T) {
	t.Skip()
	NewApiTest("create channel", &types.User{ID: 5}).
		Post("/channels/").
		Body(`{"name":"test channel"}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.NoErrors).
		End()
}
