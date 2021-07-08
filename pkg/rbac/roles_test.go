package rbac

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_partitionRoles(t *testing.T) {
	var (
		req = require.New(t)
		pr  = partitionRoles(
			&Role{id: 1, kind: BypassRole},
			&Role{id: 2, kind: BypassRole},
			&Role{id: 3, kind: BypassRole},
			&Role{id: 4, kind: ContextRole},
			&Role{id: 5, kind: CommonRole},
		)
	)

	req.Nil(pr[AuthenticatedRole])
	req.Nil(pr[AnonymousRole])
	req.NotNil(pr[BypassRole])
	req.NotNil(pr[ContextRole])
	req.NotNil(pr[CommonRole])
	req.Len(pr[BypassRole], 3)
	req.True(pr[BypassRole][1])
	req.True(pr[BypassRole][2])
	req.True(pr[BypassRole][3])
	req.Len(pr[ContextRole], 1)
	req.True(pr[ContextRole][4])
	req.Len(pr[CommonRole], 1)
	req.True(pr[CommonRole][5])
}

func Test_getContextRoles(t *testing.T) {
	var (
		dyCheck = func(r bool) ctxRoleCheckFn {
			return func(map[string]interface{}) bool {
				return r
			}
		}

		tcc = []struct {
			name         string
			sessionRoles []uint64
			res          Resource
			preloadRoles []*Role
			output       []*Role
		}{
			{
				"existing role",
				[]uint64{1},
				nil,
				[]*Role{{id: 1, kind: BypassRole}},
				[]*Role{{id: 1, kind: BypassRole}},
			},
			{
				"missing role",
				[]uint64{2},
				nil,
				[]*Role{{id: 1, kind: BypassRole}},
				[]*Role{},
			},
			{
				"dynamic role",
				[]uint64{1, 2},
				nil,
				[]*Role{
					{id: 1, kind: BypassRole},
					{id: 2, kind: ContextRole, check: dyCheck(true)},
					{id: 3, kind: ContextRole, check: dyCheck(false)},
				},
				[]*Role{{id: 1, kind: BypassRole}, {id: 2, kind: ContextRole}},
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			var (
				req = require.New(t)
			)

			req.Equal(partitionRoles(tc.output...), getContextRoles(tc.sessionRoles, tc.res, tc.preloadRoles))
		})
	}
}
