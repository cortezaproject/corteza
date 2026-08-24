package rest

import (
	"context"
	"testing"

	"github.com/crusttech/human/server/pkg/rbac"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/stretchr/testify/require"
)

type deleteACStub struct {
	permissionsAccessController

	held    rbac.RuleSet
	granted rbac.RuleSet
}

func (s *deleteACStub) FindRulesByRoleID(context.Context, uint64) (rbac.RuleSet, error) {
	return s.held, nil
}

func (s *deleteACStub) Grant(_ context.Context, rr ...*rbac.Rule) error {
	s.granted = rr
	return nil
}

// A role's rules span every component, and this endpoint is registered once per
// component behind that component's own grant check. Clearing the lot from here
// reached past what the caller was granted, and the service refused the foreign
// resources — after the handler had already flipped them to Inherit.
func TestPermissionsDeleteClearsOwnComponentOnly(t *testing.T) {
	foreign := rbac.AllowRule(1, "corteza::compose:module/2/3", "read")

	ac := &deleteACStub{held: rbac.RuleSet{
		rbac.AllowRule(1, "corteza::system/", "user.create"),
		rbac.AllowRule(1, "corteza::system:user/2", "read"),
		foreign,
	}}

	_, err := Permissions{ac: ac}.Delete(context.Background(), &request.PermissionsDelete{RoleID: 1})
	require.NoError(t, err)

	require.Len(t, ac.granted, 2, "only this component's rules belong in the grant")
	for _, r := range ac.granted {
		require.Equal(t, rbac.Inherit, r.Access)
		require.Equal(t, "corteza::system", rbac.ResourceComponent(r.Resource))
	}

	require.Equal(t, rbac.Allow, foreign.Access, "another component's rule must be left alone")
}
