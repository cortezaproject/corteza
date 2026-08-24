package service

import (
	"testing"

	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

func tpl(value, header string) types.ConnectionTemplate {
	return types.ConnectionTemplate{Value: value, HeaderName: header}
}

// freshConn wraps a catalog auth block in a Connection so it can be passed to
// healCatalogAuthParams (which heals against a whole fresh connection definition).
func freshConn(auth types.ConnectionAuth) *types.Connection {
	return &types.Connection{Service: types.ConnectionService{Auth: auth}}
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

		require.True(t, healCatalogAuthParams(conn, freshConn(catalog)))
		require.Equal(t, "{{apiToken}}", conn.Service.Auth.Params["apiToken"].Value)
		require.Equal(t, "Authorization", conn.Service.Auth.Params["apiToken"].HeaderName)
	})

	t.Run("adds a missing param", func(t *testing.T) {
		conn := &types.Connection{}
		require.True(t, healCatalogAuthParams(conn, freshConn(catalog)))
		require.Equal(t, "{{apiToken}}", conn.Service.Auth.Params["apiToken"].Value)
	})

	t.Run("never overwrites a non-empty value", func(t *testing.T) {
		conn := &types.Connection{}
		conn.Service.Auth.Params = map[string]types.ConnectionTemplate{
			"apiToken": tpl("user-edited", "X-Custom"),
		}

		require.False(t, healCatalogAuthParams(conn, freshConn(catalog)))
		require.Equal(t, "user-edited", conn.Service.Auth.Params["apiToken"].Value)
		require.Equal(t, "X-Custom", conn.Service.Auth.Params["apiToken"].HeaderName)
	})

	t.Run("ignores empty catalog values", func(t *testing.T) {
		conn := &types.Connection{}
		conn.Service.Auth.Params = map[string]types.ConnectionTemplate{
			"apiToken": tpl("", ""),
		}
		empty := types.ConnectionAuth{Params: map[string]types.ConnectionTemplate{"apiToken": tpl("", "")}}

		require.False(t, healCatalogAuthParams(conn, freshConn(empty)))
	})

	t.Run("heals a matching auth option", func(t *testing.T) {
		conn := &types.Connection{}
		conn.Service.AuthOptions = []types.ConnectionAuth{
			{Method: "api_token", Params: map[string]types.ConnectionTemplate{"apiToken": tpl("", "")}},
		}
		fresh := &types.Connection{Service: types.ConnectionService{
			AuthOptions: []types.ConnectionAuth{catalog},
		}}

		require.True(t, healCatalogAuthParams(conn, fresh))
		require.Equal(t, "{{apiToken}}", conn.Service.AuthOptions[0].Params["apiToken"].Value)
	})
}
