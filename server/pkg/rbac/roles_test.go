package rbac

import (
	"testing"

	"github.com/stretchr/testify/require"
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

		tres = NewResource("testResource")

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
				tres,
				[]*Role{{id: 1, kind: BypassRole}},
				[]*Role{{id: 1, kind: BypassRole}},
			},
			{
				"missing role",
				[]uint64{2},
				tres,
				[]*Role{{id: 1, kind: BypassRole}},
				[]*Role{},
			},
			{
				"dynamic role",
				[]uint64{1, 2},
				tres,
				[]*Role{
					{id: 1, kind: BypassRole},
					{id: 2, kind: ContextRole, check: dyCheck(true), crtypes: map[string]bool{tres.RbacResource(): true}},
					{id: 3, kind: ContextRole, check: dyCheck(false), crtypes: map[string]bool{tres.RbacResource(): true}},
					{id: 4, kind: ContextRole, check: dyCheck(true)},
				},
				[]*Role{{id: 1, kind: BypassRole}, {id: 2, kind: ContextRole}},
			},
			{
				"anonymous role mixed with matching contextual",
				[]uint64{1, 2},
				tres,
				[]*Role{
					{id: 1, kind: AnonymousRole},
					{id: 2, kind: ContextRole, check: dyCheck(true), crtypes: map[string]bool{tres.RbacResource(): true}},
					{id: 3, kind: ContextRole, check: dyCheck(false), crtypes: map[string]bool{tres.RbacResource(): true}},
					{id: 4, kind: ContextRole, check: dyCheck(true)},
				},
				[]*Role{{id: 1, kind: AnonymousRole}},
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			var (
				req = require.New(t)
			)

			req.Equal(partitionRoles(tc.output...), getSessionRoles(&session{rr: tc.sessionRoles}, tc.res, tc.preloadRoles, nil))
		})
	}
}

func Test_getContextRolesWithUserResolver(t *testing.T) {
	var (
		tres = NewResource("corteza::compose:record/123/456")

		// A check function that verifies user labels are accessible in scope
		labelCheck = func(scope map[string]interface{}) bool {
			user, ok := scope["user"].(map[string]interface{})
			if !ok {
				return false
			}
			labels, ok := user["labels"].(map[string]interface{})
			if !ok {
				return false
			}
			client, ok := labels["client"].(string)
			if !ok {
				return false
			}
			return client == "acme"
		}

		// A check function that verifies user email is accessible
		emailCheck = func(scope map[string]interface{}) bool {
			user, ok := scope["user"].(map[string]interface{})
			if !ok {
				return false
			}
			email, ok := user["email"].(string)
			if !ok {
				return false
			}
			return email == "test@example.com"
		}

		// Mock user resolver that returns user data
		mockResolver = func(userID uint64) map[string]interface{} {
			if userID == 42 {
				return map[string]interface{}{
					"email":    "test@example.com",
					"username": "testuser",
					"handle":   "test-handle",
					"name":     "Test User",
					"labels": map[string]interface{}{
						"client": "acme",
						"region": "eu",
					},
				}
			}
			return nil
		}
	)

	t.Run("context role with matching user labels", func(t *testing.T) {
		req := require.New(t)
		roles := []*Role{
			{id: 1, kind: ContextRole, check: labelCheck, crtypes: map[string]bool{"corteza::compose:record": true}},
		}
		result := getSessionRoles(&session{id: 42, rr: []uint64{1}}, tres, roles, mockResolver)
		req.NotNil(result[ContextRole])
		req.True(result[ContextRole][1], "context role should be active when user labels match")
	})

	t.Run("context role with non-matching user (no resolver data)", func(t *testing.T) {
		req := require.New(t)
		roles := []*Role{
			{id: 1, kind: ContextRole, check: labelCheck, crtypes: map[string]bool{"corteza::compose:record": true}},
		}
		// User 999 has no data from the resolver
		result := getSessionRoles(&session{id: 999, rr: []uint64{1}}, tres, roles, mockResolver)
		req.Nil(result[ContextRole], "context role should not be active when user has no resolver data")
	})

	t.Run("context role checking user email", func(t *testing.T) {
		req := require.New(t)
		roles := []*Role{
			{id: 1, kind: ContextRole, check: emailCheck, crtypes: map[string]bool{"corteza::compose:record": true}},
		}
		result := getSessionRoles(&session{id: 42, rr: []uint64{1}}, tres, roles, mockResolver)
		req.NotNil(result[ContextRole])
		req.True(result[ContextRole][1], "context role should be active when user email matches")
	})

	t.Run("nil resolver falls back to no user data", func(t *testing.T) {
		req := require.New(t)
		roles := []*Role{
			{id: 1, kind: ContextRole, check: labelCheck, crtypes: map[string]bool{"corteza::compose:record": true}},
		}
		result := getSessionRoles(&session{id: 42, rr: []uint64{1}}, tres, roles, nil)
		req.Nil(result[ContextRole], "context role should not be active when no resolver is provided")
	})
}
