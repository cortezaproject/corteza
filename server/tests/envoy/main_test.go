package envoy

import (
	"context"
	"os"
	"testing"
	"time"

	automationEnvoy "github.com/cortezaproject/corteza/server/automation/envoy"
	composeEnvoy "github.com/cortezaproject/corteza/server/compose/envoy"
	systemEnvoy "github.com/cortezaproject/corteza/server/system/envoy"

	"github.com/cortezaproject/corteza/server/pkg/cli"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/service"
	sysTypes "github.com/cortezaproject/corteza/server/system/types"
	"github.com/cortezaproject/corteza/server/tests/helpers"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
)

var (
	defaultEnvoy *envoyx.Service
	defaultStore store.Storer
	defaultDal   dal.FullService
)

func init() {
	helpers.RecursiveDotEnvLoad()
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

	if defaultDal == nil {
		initDalSvc(ctx)
	}

	if defaultEnvoy == nil {
		initSvc(ctx)
	}
}

func initStore(ctx context.Context) {
	var err error
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

func initDalSvc(ctx context.Context) {
	var err error
	defaultDal, err = dal.New(zap.NewNop(), true)
	if err != nil {
		panic(err)
	}

	conn := &sysTypes.DalConnection{
		ID:     id.Next(),
		Handle: sysTypes.DalPrimaryConnectionHandle,
		Type:   sysTypes.DalPrimaryConnectionResourceType,

		Meta: sysTypes.ConnectionMeta{
			Name: "Primary Database",
		},

		Config: sysTypes.ConnectionConfig{
			DAL: &sysTypes.ConnectionConfigDAL{
				ModelIdent: "compose_record",
			},
		},

		CreatedAt: time.Now(),
	}

	cw, err := service.MakeDalConnection(conn, defaultStore.ToDalConn())
	if err != nil {
		panic(err)
	}
	err = defaultDal.ReplaceConnection(ctx, cw, true)
	if err != nil {
		panic(err)
	}
}

func cleanup(t *testing.T) {
	var (
		ctx = context.Background()
	)

	err := collect(
		store.TruncateComposeCharts(ctx, defaultStore),
		store.TruncateComposeModules(ctx, defaultStore),
		store.TruncateComposeModuleFields(ctx, defaultStore),
		store.TruncateComposeNamespaces(ctx, defaultStore),
		store.TruncateComposePages(ctx, defaultStore),

		store.TruncateApplications(ctx, defaultStore),
		store.TruncateApigwRoutes(ctx, defaultStore),
		store.TruncateApigwFilters(ctx, defaultStore),
		store.TruncateAuthClients(ctx, defaultStore),
		store.TruncateQueues(ctx, defaultStore),
		store.TruncateReports(ctx, defaultStore),
		store.TruncateRoles(ctx, defaultStore),
		store.TruncateTemplates(ctx, defaultStore),
		store.TruncateUsers(ctx, defaultStore),
		store.TruncateDalConnections(ctx, defaultStore),
		store.TruncateDalSensitivityLevels(ctx, defaultStore),
		store.TruncateRbacRules(ctx, defaultStore),
		store.TruncateResourceTranslations(ctx, defaultStore),

		store.TruncateAutomationWorkflows(ctx, defaultStore),
		store.TruncateAutomationTriggers(ctx, defaultStore),

		truncateRecords(ctx),
	)
	if err != nil {
		t.Fatalf("failed to decode scenario data: %v", err)
	}
}

func truncateRecords(ctx context.Context) error {
	models, err := defaultDal.SearchModels(ctx)
	if err != nil {
		return err
	}
	for _, model := range models {
		err = defaultDal.Truncate(ctx, model.ToFilter(), nil)
		if err != nil {
			return err
		}
	}
	return nil
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
	defaultEnvoy.AddDecoder(envoyx.DecodeTypeStore,
		composeEnvoy.StoreDecoder{},
		systemEnvoy.StoreDecoder{},
		automationEnvoy.StoreDecoder{},
	)

	defaultEnvoy.AddEncoder(envoyx.EncodeTypeIo,
		composeEnvoy.YamlEncoder{},
		systemEnvoy.YamlEncoder{},
		automationEnvoy.YamlEncoder{},
	)
	defaultEnvoy.AddEncoder(envoyx.EncodeTypeStore,
		composeEnvoy.StoreEncoder{},
		systemEnvoy.StoreEncoder{},
		automationEnvoy.StoreEncoder{},
	)
	defaultEnvoy.AddEncoder(envoyx.EncodeTypeIo,
		composeEnvoy.CsvEncoder{},
	)
}
