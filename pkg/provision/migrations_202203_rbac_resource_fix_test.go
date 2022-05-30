package provision

import (
	"github.com/stretchr/testify/require"
	"testing"

	composeTypes "github.com/cortezaproject/corteza-server/compose/types"
	federationTypes "github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
)

func Test_migratePost202203RbacRules(t *testing.T) {
	rx := &resIndex{
		fields:         make(map[uint64]*composeTypes.ModuleField),
		modules:        make(map[uint64]*composeTypes.Module),
		charts:         make(map[uint64]*composeTypes.Chart),
		pages:          make(map[uint64]*composeTypes.Page),
		exposedModules: make(map[uint64]*federationTypes.ExposedModule),
		sharedModules:  make(map[uint64]*federationTypes.SharedModule),
	}

	tcc := []struct {
		wantOp   actionType
		rule     *rbac.Rule
		wantRule *rbac.Rule
	}{
		{actionRemove, rbac.AllowRule(0, "corteza::compose:module/*/123", "op"), rbac.AllowRule(0, "corteza::compose:module/*/123", "op")},
		{actionRemove, rbac.AllowRule(0, "corteza::compose:module-field/*/*/123", "op"), rbac.AllowRule(0, "corteza::compose:module-field/*/*/123", "op")},
		{actionRemove, rbac.AllowRule(0, "corteza::compose:record/*/*/234", "op"), rbac.AllowRule(0, "corteza::compose:record/*/*/234", "op")},
	}
	for _, tc := range tcc {
		t.Run(tc.rule.String(), func(t *testing.T) {
			require.Equal(t, tc.wantOp, migratePost202203RbacRule(tc.rule, rx))
			if tc.wantRule != nil {
				require.Equal(t, tc.wantRule.String(), tc.rule.String())
			}
		})
	}
}
