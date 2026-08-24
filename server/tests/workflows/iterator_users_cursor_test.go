package workflows

import (
	"context"
	"fmt"
	"testing"

	autTypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/wfexec"
	"github.com/cortezaproject/corteza/server/system/automation"
	sysTypes "github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/require"
)

// Test_iterator_users_cursor starts iteration from a page cursor and yields
// only the items that follow it
func Test_iterator_users_cursor(t *testing.T) {
	wfexec.MaxIteratorBufferSize = wfexec.DefaultMaxIteratorBufferSize
	defer func() {
		wfexec.MaxIteratorBufferSize = wfexec.DefaultMaxIteratorBufferSize
	}()

	var (
		ctx = bypassRBAC(context.Background())
		req = require.New(t)
	)

	req.NoError(defStore.TruncateUsers(ctx))

	loadNewScenario(ctx, t)

	// Take the cursor that points past the first two users
	_, uf, err := defStore.SearchUsers(ctx, sysTypes.UserFilter{Paging: filter.Paging{Limit: 2}})
	req.NoError(err)
	req.NotNil(uf.NextPage)

	input := &expr.Vars{}
	req.NoError(input.Set("cursor", expr.Must(expr.NewString(uf.NextPage.Encode()))))

	var (
		_, trace = mustExecWorkflow(ctx, t, "testing", autTypes.WorkflowExecParams{Input: input})
	)

	// 3 of the 5 users remain: 4x iterator, 3x continue, 1x terminator, 1x completed
	req.Len(trace, 9)

	for j := 0; j <= 2; j++ {
		usr, err := automation.NewUser(trace[j*2].Results.GetValue()["u"])
		req.NoError(err)
		req.Equal(fmt.Sprintf("u%d", j+3), usr.GetValue().Handle)
	}
}
