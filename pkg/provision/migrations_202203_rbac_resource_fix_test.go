package provision

import (
	composeTypes "github.com/cortezaproject/corteza-server/compose/types"
	federationTypes "github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_migratePost202203RbacRules(t *testing.T) {
	type (
		testingExpect struct {
			wantOp   actionType
			rule     *rbac.Rule
			wantRule *rbac.Rule
		}
	)
	var (
		ctx = cli.Context()
		req = require.New(t)

		rx = &resIndex{
			fields:         make(map[uint64]*composeTypes.ModuleField),
			modules:        make(map[uint64]*composeTypes.Module),
			charts:         make(map[uint64]*composeTypes.Chart),
			pages:          make(map[uint64]*composeTypes.Page),
			exposedModules: make(map[uint64]*federationTypes.ExposedModule),
			sharedModules:  make(map[uint64]*federationTypes.SharedModule),
		}
		tcc = []testingExpect{
			{
				actionRemove,
				rbac.AllowRule(0, "corteza::compose:module/*/123", "op"),
				nil,
			},
			{
				actionRemove,
				rbac.AllowRule(0, "corteza::compose:module-field/*/*/123", "op"),
				nil,
			},
			{
				actionRemove,
				rbac.AllowRule(0, "corteza::compose:record/*/*/234", "op"),
				nil,
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.rule.String(), func(t *testing.T) {
			req.Equal(tc.wantOp, migratePost202203RbacRule(ctx, nil, tc.rule, rx))
			if tc.wantRule != nil {
				req.Equal(tc.wantRule.String(), tc.rule.String())
			}
		})
	}
}
