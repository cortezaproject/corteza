package service

import (
	"testing"

	labelTypes "github.com/cortezaproject/corteza/server/pkg/label/types"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/require"
)

func Test_rbacUserScope(t *testing.T) {
	u := &types.User{
		Email:    "vip@corp.example",
		Username: "vip",
		Handle:   "vip",
		Name:     "Director of Finance",
		Labels: map[string]labelTypes.LabelValue{
			"team":  {Val: "finance"},
			"sites": {Values: []string{"a", "b"}},
		},
	}

	t.Run("unconfirmed email is not exposed", func(t *testing.T) {
		scope := rbacUserScope(u)
		require.Equal(t, "", scope["email"])
	})

	t.Run("confirmed email and labels are exposed", func(t *testing.T) {
		u.EmailConfirmed = true
		scope := rbacUserScope(u)

		require.Equal(t, "vip@corp.example", scope["email"])
		require.Equal(t, map[string]interface{}{"team": "finance", "sites": []string{"a", "b"}}, scope["labels"])
	})

	t.Run("properties users can change by themselves are not exposed", func(t *testing.T) {
		scope := rbacUserScope(u)

		require.NotContains(t, scope, "name")
		require.NotContains(t, scope, "handle")
		require.NotContains(t, scope, "username")
	})
}

func Test_ignoredOnSelfUpdate(t *testing.T) {
	var (
		stored = &types.User{
			Email:       "jane@corp.example",
			Username:    "jane",
			UserGroupID: 42,
			Labels:      map[string]labelTypes.LabelValue{"team": {Val: "support"}},
		}
	)

	// partial user, as sent by the webapp when switching the theme
	require.Empty(t, ignoredOnSelfUpdate(stored, &types.User{Email: "jane@corp.example", Name: "Jane"}))

	// unchanged values
	require.Empty(t, ignoredOnSelfUpdate(stored, stored))

	require.Equal(t,
		[]string{"email", "username", "userGroupID", "kind", "labels"},
		ignoredOnSelfUpdate(stored, &types.User{
			Email:       "vip@corp.example",
			Username:    "vip",
			UserGroupID: 1,
			Kind:        types.UserKind("bot"),
			Labels:      map[string]labelTypes.LabelValue{"team": {Val: "finance"}},
		}),
	)
}
