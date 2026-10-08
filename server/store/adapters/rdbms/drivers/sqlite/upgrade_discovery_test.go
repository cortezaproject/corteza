package sqlite_test

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/drivers/sqlite"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type discoveryConfig struct {
	Discovery struct {
		Public struct {
			Result []struct {
				Lang string `json:"lang"`
			} `json:"result"`
		} `json:"public"`
	} `json:"discovery"`
	RecordRevisions struct {
		Enabled bool `json:"enabled"`
	} `json:"recordRevisions"`
}

// bootWithModuleConfig stores a module with the given config and runs the
// store upgrade twice, as two boots do, then reads the config back.
func bootWithModuleConfig(t *testing.T, name, config string) (out discoveryConfig) {
	var (
		ctx = context.Background()
		req = require.New(t)
	)

	s, err := sqlite.Connect(ctx, "sqlite3://file:"+name+"?mode=memory&cache=shared")
	req.NoError(err)
	rs := s.(*rdbms.Store)
	req.NoError(store.Upgrade(ctx, zap.NewNop(), s))

	req.NoError(store.CreateComposeModule(ctx, s, &types.Module{ID: 1, NamespaceID: 9, Handle: "m1", Name: "m", CreatedAt: time.Now()}))
	_, err = rs.DB.ExecContext(ctx, `UPDATE compose_module SET config = ? WHERE id = 1`, config)
	req.NoError(err)

	for i := 0; i < 2; i++ {
		req.NoError(store.Upgrade(ctx, zap.NewNop(), s), "boot "+strconv.Itoa(i+1))
	}

	var raw string
	req.NoError(rs.DB.QueryRowxContext(ctx, `SELECT config FROM compose_module WHERE id = 1`).Scan(&raw))
	req.NoError(json.Unmarshal([]byte(raw), &out), raw)
	return
}

// A module's discovery settings keep every language they list, boot after
// boot, along with the rest of its config.
func TestUpgradeKeepsModuleDiscoveryLanguages(t *testing.T) {
	c := bootWithModuleConfig(t, "upgrade_discovery_languages",
		`{"discovery":{"public":{"result":[{"lang":"en","fields":["title"]},{"lang":"de","fields":["body"]}]}},"recordRevisions":{"enabled":true}}`)

	require.Len(t, c.Discovery.Public.Result, 2)
	require.Equal(t, "en", c.Discovery.Public.Result[0].Lang)
	require.Equal(t, "de", c.Discovery.Public.Result[1].Lang)
	require.True(t, c.RecordRevisions.Enabled)
}
