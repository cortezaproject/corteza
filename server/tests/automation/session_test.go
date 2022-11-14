package automation

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/automation/service"
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func (h helper) clearSessions() {
	h.noError(store.TruncateAutomationSessions(context.Background(), service.DefaultStore))
}

func (h helper) repoMakeSession(wf *types.Workflow, ss ...string) *types.Session {
	var r = &types.Session{
		ID:         id.Next(),
		CreatedAt:  time.Now(),
		WorkflowID: wf.ID,
	}

	if len(ss) > 1 {
		r.ResourceType = ss[1]
	} else {
		r.ResourceType = "h_" + rs()

	}

	h.a.NoError(store.CreateAutomationSession(context.Background(), service.DefaultStore, r))

	return r
}

func TestSessionList(t *testing.T) {
	h := newHelper(t)
	h.clearSessions()

	helpers.AllowMe(h, types.ComponentRbacResource(), "sessions.search")

	wf := h.repoMakeWorkflow()
	h.repoMakeSession(wf)
	h.repoMakeSession(wf)

	h.apiInit().
		Get("/sessions/").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Len(`$.response.set`, 2)).
		End()
}

func TestSessionRead(t *testing.T) {
	h := newHelper(t)
	h.clearTriggers()

	wf := h.repoMakeWorkflow()
	s := h.repoMakeSession(wf)

	helpers.AllowMe(h, types.ComponentRbacResource(), "sessions.search")
	helpers.AllowMe(h, types.WorkflowRbacResource(0), "sessions.manage")

	h.apiInit().
		Get(fmt.Sprintf("/sessions/%d", s.ID)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()
}

func TestSessionTrace(t *testing.T) {
	t.Skip("pending implementation")
}

func TestSessionDelete(t *testing.T) {
	t.Skip("pending implementation")
}
