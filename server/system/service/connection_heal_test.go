package service

import (
	"testing"

	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

func tpl(value, header string) types.ConnectionTemplate {
	return types.ConnectionTemplate{Value: value, HeaderName: header}
}

func TestHealCatalogAuthParams(t *testing.T) {
	catalog := types.ConnectionAuth{
		Method: "api_token",
		Params: map[string]types.ConnectionTemplate{
			"apiToken": tpl("{{apiToken}}", "Authorization"),
		},
	}

	t.Run("restores an empty stored value", func(t *testing.T) {
		conn := &types.Connection{}
		conn.Service.Auth.Params = map[string]types.ConnectionTemplate{
			"apiToken": tpl("", ""),
		}

		require.True(t, healCatalogAuthParams(conn, catalog))
		require.Equal(t, "{{apiToken}}", conn.Service.Auth.Params["apiToken"].Value)
		require.Equal(t, "Authorization", conn.Service.Auth.Params["apiToken"].HeaderName)
	})

	t.Run("adds a missing param", func(t *testing.T) {
		conn := &types.Connection{}
		require.True(t, healCatalogAuthParams(conn, catalog))
		require.Equal(t, "{{apiToken}}", conn.Service.Auth.Params["apiToken"].Value)
	})

	t.Run("never overwrites a non-empty value", func(t *testing.T) {
		conn := &types.Connection{}
		conn.Service.Auth.Params = map[string]types.ConnectionTemplate{
			"apiToken": tpl("user-edited", "X-Custom"),
		}

		require.False(t, healCatalogAuthParams(conn, catalog))
		require.Equal(t, "user-edited", conn.Service.Auth.Params["apiToken"].Value)
		require.Equal(t, "X-Custom", conn.Service.Auth.Params["apiToken"].HeaderName)
	})

	t.Run("ignores empty catalog values", func(t *testing.T) {
		conn := &types.Connection{}
		conn.Service.Auth.Params = map[string]types.ConnectionTemplate{
			"apiToken": tpl("", ""),
		}
		empty := types.ConnectionAuth{Params: map[string]types.ConnectionTemplate{"apiToken": tpl("", "")}}

		require.False(t, healCatalogAuthParams(conn, empty))
	})
}
