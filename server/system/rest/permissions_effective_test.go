package rest

import (
	"context"
	"fmt"
	"testing"

	"github.com/crusttech/human/server/pkg/rbac"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

type effectiveACStub struct {
	permissionsAccessController

	component  []rbac.Resource
	forced     string
	forcedErr  error
	forcedResp rbac.EffectiveSet
}

func (s *effectiveACStub) Effective(_ context.Context, rr ...rbac.Resource) rbac.EffectiveSet {
	s.component = rr
	return rbac.EffectiveSet{{Resource: "component"}}
}

func (s *effectiveACStub) EffectiveFor(_ context.Context, res string) (rbac.EffectiveSet, error) {
	s.forced = res
	return s.forcedResp, s.forcedErr
}

// The endpoint has always taken a resource and always ignored it, answering
// about the component whatever was asked. That made every per-resource question
// look like a denial, since the component carries no such operation.
func TestPermissionsEffectiveHonoursResource(t *testing.T) {
	t.Run("no resource asks about the component", func(t *testing.T) {
		ac := &effectiveACStub{}
		out, err := Permissions{ac: ac}.Effective(context.Background(), &request.PermissionsEffective{})

		require.NoError(t, err)
		require.Equal(t, rbac.EffectiveSet{{Resource: "component"}}, out)
		require.Equal(t, []rbac.Resource{types.Component{}}, ac.component)
		require.Empty(t, ac.forced, "a component question must not be routed per-resource")
	})

	t.Run("a resource is passed through", func(t *testing.T) {
		want := rbac.EffectiveSet{{Resource: "corteza::system:application/1", Operation: "access", Allow: true}}
		ac := &effectiveACStub{forcedResp: want}

		out, err := Permissions{ac: ac}.Effective(
			context.Background(),
			&request.PermissionsEffective{Resource: "corteza::system:application/1"},
		)

		require.NoError(t, err)
		require.Equal(t, want, out)
		require.Equal(t, "corteza::system:application/1", ac.forced)
		require.Nil(t, ac.component, "a per-resource question must not fall back to the component")
	})

	t.Run("a refused resource surfaces as an error, not an empty answer", func(t *testing.T) {
		ac := &effectiveACStub{forcedErr: fmt.Errorf("invalid resource path structure")}

		_, err := Permissions{ac: ac}.Effective(
			context.Background(),
			&request.PermissionsEffective{Resource: "corteza::system:module/*"},
		)

		require.Error(t, err)
	})
}
