package messaging

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func TestPermissionsRead(t *testing.T) {
	h := newHelper(t)
	h.allow("messaging", "access")
	h.allow("messaging", "grant")
	h.deny("messaging", "channel.group.create")

	h.apiInit().
		Get(fmt.Sprintf("/permissions/%d/rules", h.roleID)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}
