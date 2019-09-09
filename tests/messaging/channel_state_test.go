package messaging

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/internal/permissions"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func TestChannelAlterState(t *testing.T) {
	h := newHelper(t)
	ch := h.makePublicCh()

	stateUrl := fmt.Sprintf("/channels/%d/state", ch.ID)

	testState := func(state string) {
		p.ClearGrants()
		h.mockPermissionsWithAccess()
		h.allow(types.ChannelPermissionResource.AppendWildcard(), permissions.Operation(state))

		h.testAPI().
			Put(stateUrl).
			FormData("state", state).
			Expect(t).
			Status(http.StatusOK).
			Assert(helpers.AssertNoErrors).
			End()
	}

	testState("archive")
	testState("unarchive")
	testState("delete")
	testState("undelete")
}
