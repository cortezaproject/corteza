package messaging

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func TestChannelAlterState(t *testing.T) {
	h := newHelper(t)
	ch := h.repoMakePublicCh()

	stateUrl := fmt.Sprintf("/channels/%d/state", ch.ID)

	testState := func(state string) {
		h.mockPermissionsWithAccess()
		h.allow(types.ChannelRBACResource.AppendWildcard(), rbac.Operation(state))

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
