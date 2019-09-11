package messaging

import (
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func TestPermissionsEffective(t *testing.T) {
	h := newHelper(t)
	h.allow(types.MessagingPermissionResource, "access")
	h.deny(types.MessagingPermissionResource, "channel.group.create")

	h.apiInit().
		Get("/permissions/effective").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}
