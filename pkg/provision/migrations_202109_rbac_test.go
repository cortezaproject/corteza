package provision

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/stretchr/testify/require"
)

func Test_migratePre202106RbacRule(t *testing.T) {
	tcc := []struct {
		wantOp   int
		rule     *rbac.Rule
		wantRule *rbac.Rule
	}{
		{-1, rbac.AllowRule(0, "messaging", "whatever"), nil},
		{-1, rbac.AllowRule(0, "foo:automation-script", "whatever"), nil},
		{1, rbac.AllowRule(0, "abc", "whatever"), rbac.AllowRule(0, "corteza::abc/", "whatever")},
		{1, rbac.AllowRule(0, "federation:module", "map"), rbac.AllowRule(0, "corteza::federation:shared-module/", "map")},
		{1, rbac.AllowRule(0, "federation:module", "manage"), rbac.AllowRule(0, "corteza::federation:exposed-module/", "manage")},
		{1, rbac.AllowRule(0, "compose:module:234", "record.read"), rbac.AllowRule(0, "corteza::compose:record/*/234/*", "read")},
		{1, rbac.AllowRule(0, "compose:module-field:234", "op"), rbac.AllowRule(0, "corteza::compose:module-field/*/*/234", "op")},
	}
	for _, tc := range tcc {
		t.Run(tc.rule.String(), func(t *testing.T) {
			require.Equal(t, tc.wantOp, migratePre202106RbacRule(tc.rule))
			if tc.wantRule != nil {
				require.Equal(t, tc.wantRule.String(), tc.rule.String())
			}
		})
	}
}
