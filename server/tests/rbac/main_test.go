package rbac

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	automationEnvoy "github.com/cortezaproject/corteza/server/automation/envoy"
	composeEnvoy "github.com/cortezaproject/corteza/server/compose/envoy"
	systemEnvoy "github.com/cortezaproject/corteza/server/system/envoy"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/require"

	"github.com/cortezaproject/corteza/server/pkg/cli"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/tests/helpers"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
)

type (
	sesWrap struct {
		identity uint64
		roles    []uint64
		context  context.Context
	}

	resWrap struct {
		resource string
	}

	testStorage struct {
		upserts []*rbac.Rule

		returnRuleSearch []*rbac.Rule
	}

	svcModFnc func(*rbac.Service)
)

var (
	defaultEnvoy *envoyx.Service
	defaultStore store.Storer
)

func init() {
	helpers.RecursiveDotEnvLoad()
	id.Init(cli.Context())
}

func TestMain(m *testing.M) {
	InitTestApp()
	os.Exit(m.Run())
}

func InitTestApp() {
	ctx := cli.Context()

	if defaultStore == nil {
		initStore(ctx)
	}

	if defaultEnvoy == nil {
		initSvc(ctx)
	}
}

func initStore(ctx context.Context) {
	var err error
	// dsn := "postgres://corteza:corteza@127.0.0.1:3402/testing?sslmode=disable"
	dsn := "sqlite3+debug://file::memory:?cache=shared&mode=memory"
	defaultStore, err = store.Connect(ctx, zap.NewNop(), dsn, true)
	if err != nil {
		panic(err)
	}

	err = store.Upgrade(ctx, zap.NewNop(), defaultStore)
	if err != nil {
		panic(err)
	}
}

func cleanup(b *testing.B) {
	var (
		ctx = context.Background()
	)

	err := collect(
		store.TruncateRbacRules(ctx, defaultStore),
		store.TruncateRoles(ctx, defaultStore),
		store.TruncateUsers(ctx, defaultStore),
	)
	if err != nil {
		b.Fatalf("failed to decode scenario data: %v", err)
	}
}

func collect(ee ...error) error {
	for _, e := range ee {
		if e != nil {
			return e
		}
	}
	return nil
}

func initSvc(ctx context.Context) {
	defaultEnvoy = envoyx.New()
	defaultEnvoy.AddDecoder(envoyx.DecodeTypeURI,
		composeEnvoy.YamlDecoder{},
		systemEnvoy.YamlDecoder{},
		automationEnvoy.YamlDecoder{},
	)

	defaultEnvoy.AddEncoder(envoyx.EncodeTypeStore,
		composeEnvoy.StoreEncoder{},
		systemEnvoy.StoreEncoder{},
		automationEnvoy.StoreEncoder{},
	)
}

func initState(t *testing.T, maxIndexSize int, things ...svcModFnc) (context.Context, *require.Assertions, *rbac.Service, *testStorage) {
	var (
		ctx = context.Background()
		req = require.New(t)
	)

	store := &testStorage{}
	svc, err := rbac.NewService(ctx, zap.NewNop(), defaultStore, rbac.Config{
		Synchronous: true,

		MaxIndexSize:    maxIndexSize,
		DecayFactor:     1,
		DecayInterval:   time.Hour * 4,
		CleanupInterval: time.Hour * 4,

		RuleStorage: store,
		RoleStorage: store,
	})
	req.NoError(err)

	for _, f := range things {
		f(svc)
	}

	return ctx, req, svc, store
}

func mustStats(req *require.Assertions, svc *rbac.Service) rbac.Stats {
	stats, err := svc.Stats()
	req.NoError(err)
	return stats
}

func must(req *require.Assertions, err error) {
	req.NoError(err)
}

func checkHitRatios(req *require.Assertions, stats rbac.Stats, hits, misses uint, lastHitsLastMisses ...[][]uint64) {
	req.Equal(hits, stats.CacheHits)
	req.Equal(misses, stats.CacheMisses)

	if len(lastHitsLastMisses) > 0 {
		for i := 0; i < len(lastHitsLastMisses[0]); i++ {
			req.Contains(stats.LastHits[i], fmt.Sprintf("%v", lastHitsLastMisses[0][i]))
		}
	}

	if len(lastHitsLastMisses) > 1 {
		for i := 0; i < len(lastHitsLastMisses[1]); i++ {
			req.Contains(stats.LastMisses[i], fmt.Sprintf("%v", lastHitsLastMisses[1][i]))
		}
	}
}

// Utils

func (ts *testStorage) SearchRbacRules(ctx context.Context, f rbac.RuleFilter) (rs rbac.RuleSet, rf rbac.RuleFilter, er error) {
	return ts.returnRuleSearch, f, nil
}

func (ts *testStorage) UpsertRbacRule(ctx context.Context, rr ...*rbac.Rule) (err error) {
	ts.upserts = append(ts.upserts, rr...)
	return
}

func (testStorage) DeleteRbacRule(ctx context.Context, rr ...*rbac.Rule) (err error) {
	return
}

func (testStorage) TruncateRbacRules(ctx context.Context) (err error) {
	return
}

func (testStorage) SearchRoles(ctx context.Context, f types.RoleFilter) (rs types.RoleSet, rf types.RoleFilter, err error) {
	return
}

func (sw sesWrap) Identity() uint64 {
	return sw.identity
}
func (sw sesWrap) Roles() []uint64 {
	return sw.roles
}
func (sw sesWrap) Context() context.Context {
	return sw.context
}

func (rw resWrap) RbacResource() string {
	return rw.resource
}

func svcWithRoles(roles ...*rbac.Role) svcModFnc {
	return func(s *rbac.Service) {
		s.UpdateRoles(roles...)
	}
}
