package messaging

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"net/http"
	"testing"
)

func TestPermissionsRead(t *testing.T) {
	h := newHelper(t)
	h.allow(types.MessagingRBACResource, "grant")
	h.deny(types.MessagingRBACResource, "channel.group.create")

	h.apiInit().
		Get(fmt.Sprintf("/permissions/%d/rules", h.roleID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}
