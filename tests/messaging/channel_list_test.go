package messaging

import (
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func TestChannelList(t *testing.T) {
	h := newHelper(t)

	h.testAPI().
		Get("/channels/").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}
