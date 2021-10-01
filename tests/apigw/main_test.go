package apigw

import (
	"context"
	"errors"
	"os"
	"path"
	"testing"

	"github.com/cortezaproject/corteza-server/app"
	"github.com/cortezaproject/corteza-server/pkg/api/server"
	"github.com/cortezaproject/corteza-server/pkg/apigw"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/csv"
	"github.com/cortezaproject/corteza-server/pkg/envoy/directory"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	es "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/pkg/envoy/yaml"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/label"
	ltype "github.com/cortezaproject/corteza-server/pkg/label/types"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/sqlite3"
	"github.com/cortezaproject/corteza-server/system/rest"
	"github.com/cortezaproject/corteza-server/system/service"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"github.com/go-chi/chi"
	_ "github.com/joho/godotenv/autoload"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/require"
)

type (
	helper struct {
		t *testing.T
		a *require.Assertions

		cUser  *sysTypes.User
		roleID uint64
		token  string
	}
)

var (
	testApp *app.CortezaApp
	r       chi.Router

	eventBus = eventbus.New()
)

func init() {
	helpers.RecursiveDotEnvLoad()
}

func InitTestApp() {
	ctx := cli.Context()
	if testApp == nil {

		testApp = helpers.NewIntegrationTestApp(ctx, func(app *app.CortezaApp) (err error) {
			service.DefaultStore, err = sqlite3.ConnectInMemory(ctx)
			if err != nil {
				return err
			}

			eventbus.Set(eventBus)
			return nil
		})
	}

	if r == nil {
		r = chi.NewRouter()
		r.Use(server.BaseMiddleware(false, logger.Default())...)
		helpers.BindAuthMiddleware(r)

		// Sys routes for route management tests
		rest.MountRoutes(r)

		// API gw routes
		apigw.Setup(options.Apigw(), service.DefaultLogger, service.DefaultStore)
		err := apigw.Service().Reload(ctx)
		if err != nil {
			panic(err)
		}

		r.Handle("/*", apigw.Service())
	}
}

func TestMain(m *testing.M) {
	InitTestApp()
	os.Exit(m.Run())
}

func newHelper(t *testing.T) helper {
	h := helper{
		t:      t,
		a:      require.New(t),
		roleID: id.Next(),
		cUser: &sysTypes.User{
			ID: id.Next(),
		},
	}

	h.cUser.SetRoles(h.roleID)
	helpers.UpdateRBAC(h.roleID)

	var err error
	h.token, err = auth.DefaultJwtHandler.Generate(context.Background(), h.cUser)
	if err != nil {
		panic(err)
	}

	return h
}

func (h helper) MyRole() uint64 {
	return h.roleID
}

// Returns context w/ security details
func (h helper) secCtx() context.Context {
	return auth.SetIdentityToContext(context.Background(), h.cUser)
}

// apitest basics, initialize, set handler, add auth
func (h helper) apiInit() *apitest.APITest {
	InitTestApp()

	return apitest.
		New().
		Handler(r).
		Intercept(helpers.ReqHeaderRawAuthBearer(h.token))

}

func setupScenario(t *testing.T) (context.Context, helper, store.Storer) {
	ctx, h, s := setup(t)
	loadScenario(ctx, s, t, h)
	_ = apigw.Service().Reload(ctx)

	return ctx, h, s
}

func setup(t *testing.T) (context.Context, helper, store.Storer) {
	h := newHelper(t)
	s := service.DefaultStore

	u := &sysTypes.User{
		ID: id.Next(),
	}
	u.SetRoles(auth.BypassRoles().IDs()...)

	ctx := auth.SetIdentityToContext(context.Background(), u)

	return ctx, h, s
}

// Unwraps error before it passes it to the tester
func (h helper) noError(err error) {
	for errors.Unwrap(err) != nil {
		err = errors.Unwrap(err)
	}

	h.a.NoError(err)
}

func (h helper) setLabel(res label.LabeledResource, name, value string) {
	h.a.NoError(store.UpsertLabel(h.secCtx(), service.DefaultStore, &ltype.Label{
		Kind:       res.LabelResourceKind(),
		ResourceID: res.LabelResourceID(),
		Name:       name,
		Value:      value,
	}))
}

func loadScenario(ctx context.Context, s store.Storer, t *testing.T, h helper) {
	loadScenarioWithName(ctx, s, t, h, t.Name()[5:])
}

func loadScenarioWithName(ctx context.Context, s store.Storer, t *testing.T, h helper, scenario string) {
	cleanup(ctx, h, s)
	parseEnvoy(ctx, s, h, path.Join("testdata", scenario))
}

func cleanup(ctx context.Context, h helper, s store.Storer) {
	h.noError(s.TruncateApigwFilters(ctx))
	h.noError(s.TruncateApigwRoutes(ctx))
}

func parseEnvoy(ctx context.Context, s store.Storer, h helper, path string) {
	nn, err := directory.Decode(
		ctx,
		path,
		yaml.Decoder(),
		csv.Decoder(),
	)
	if err != nil {
		h.t.Fatalf("failed to decode scenario data: %v", err)
	}

	crs := resource.ComposeRecordShaper()
	nn, err = resource.Shape(nn, crs)
	h.a.NoError(err)

	// import into the store
	se := es.NewStoreEncoder(s, nil)
	bld := envoy.NewBuilder(se)
	g, err := bld.Build(ctx, nn...)
	h.a.NoError(err)
	err = envoy.Encode(ctx, g, se)
	h.a.NoError(err)
}
