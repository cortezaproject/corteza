package messaging

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func TestPermissionsRead(t *testing.T) {
	h := newHelper(t)
	h.allow(types.MessagingPermissionResource, "access")
	h.allow(types.MessagingPermissionResource, "grant")
	h.deny(types.MessagingPermissionResource, "channel.group.create")

	h.apiInit().
		Get(fmt.Sprintf("/permissions/%d/rules", h.roleID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}
