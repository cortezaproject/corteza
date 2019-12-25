package messaging

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/messaging/service"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func TestChannelAlterState(t *testing.T) {
	h := newHelper(t)
	p := service.DefaultPermissions.(*permissions.TestService)
	ch := h.repoMakePublicCh()

	stateUrl := fmt.Sprintf("/channels/%d/state", ch.ID)

	testState := func(state string) {
		p.ClearGrants()
		h.mockPermissionsWithAccess()
		h.allow(types.ChannelPermissionResource.AppendWildcard(), permissions.Operation(state))

		h.apiInit().
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
