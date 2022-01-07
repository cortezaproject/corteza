package gig

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/app"
	"github.com/cortezaproject/corteza-server/compose/rest"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/pkg/api/server"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/csv"
	"github.com/cortezaproject/corteza-server/pkg/envoy/directory"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	es "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/pkg/envoy/yaml"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/gig"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/objstore/plain"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"github.com/go-chi/chi"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/afero"
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

	// wrapper around time.Now() that will aid service testing
	now = func() *time.Time {
		c := time.Now().Round(time.Second)
		return &c
	}
)

func init() {
	helpers.RecursiveDotEnvLoad()
}

// random string, 10 chars long by default
func rs(a ...int) string {
	var l = 10
	if len(a) > 0 {
		l = a[0]
	}

	return string(rand.Bytes(l))
}

func InitTestApp() {
	if testApp == nil {
		ctx := cli.Context()

		testApp = helpers.NewIntegrationTestApp(ctx, func(app *app.CortezaApp) (err error) {
			service.DefaultStore = app.Store
			service.DefaultObjectStore, err = plain.NewWithAfero(afero.NewMemMapFs(), "test")
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
		rest.MountRoutes(r)
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
func (h helper) mockPermissions(rules ...*rbac.Rule) {
	h.noError(rbac.Global().Grant(
		// TestService we use does not have any backend storage,
		context.Background(),
		rules...,
	))
}

// Unwraps error before it passes it to the tester
func (h helper) noError(err error) {
	for errors.Unwrap(err) != nil {
		err = errors.Unwrap(err)
	}

	h.a.NoError(err)
}

func setup(t *testing.T) (context.Context, gig.Service, helper, store.Storer) {
	h := newHelper(t)
	s := service.DefaultStore

	u := &sysTypes.User{
		ID: id.Next(),
	}
	u.SetRoles(auth.BypassRoles().IDs()...)

	ctx := auth.SetIdentityToContext(context.Background(), u)
	svc := gig.NewService(nil)

	return ctx, svc, h, s
}

func setupWithImportGig(t *testing.T) (context.Context, gig.Service, helper, store.Storer, gig.Gig) {
	ctx, svc, h, s := setup(t)
	g, err := svc.Create(ctx, gig.UpdatePayload{
		Worker: gig.WorkerImport(s),
	})
	h.a.NoError(err)

	return ctx, svc, h, s, g
}

func setupWithExportGig(t *testing.T) (context.Context, gig.Service, helper, store.Storer, gig.Gig) {
	ctx, svc, h, s := setup(t)
	g, err := svc.Create(ctx, gig.UpdatePayload{
		Worker: gig.WorkerExport(s),
	})
	h.a.NoError(err)

	return ctx, svc, h, s, g
}

func setupWithNoopGig(t *testing.T) (context.Context, gig.Service, helper, store.Storer, gig.Gig) {
	ctx, svc, h, s := setup(t)
	g, err := svc.Create(ctx, gig.UpdatePayload{
		Worker: gig.WorkerNoop(),
	})
	h.a.NoError(err)

	return ctx, svc, h, s, g
}

func noopGig(ctx context.Context, svc gig.Service, h helper) gig.Gig {
	g, err := svc.Create(ctx, gig.UpdatePayload{
		Worker: gig.WorkerNoop(),
	})
	h.a.NoError(err)

	return g
}

func scenario(t *testing.T) string {
	return t.Name()[5:]
}

func testSource(t *testing.T, pp ...string) string {
	return path.Join(append([]string{"testdata", scenario(t), "sources"}, pp...)...)
}

func loadScenario(ctx context.Context, s store.Storer, t *testing.T, h helper) {
	loadScenarioWithName(ctx, s, t, h, scenario(t))
}

func loadScenarioWithName(ctx context.Context, s store.Storer, t *testing.T, h helper, scenario string) {
	cleanup(ctx, h, s)
	parseEnvoy(ctx, s, h, path.Join("testdata", scenario, "base"))
}

func cleanup(ctx context.Context, h helper, s store.Storer) {
	h.noError(s.TruncateComposeNamespaces(ctx))
	h.noError(s.TruncateComposeModules(ctx))
	h.noError(s.TruncateComposeModuleFields(ctx))
	h.noError(s.TruncateComposeRecords(ctx, nil))
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

func makeRemoteServer(t *testing.T, pp ...string) (*httptest.Server, string) {
	dir := []string{}
	if len(pp) > 1 {
		dir = pp[:len(pp)-1]
	}
	srv := httptest.NewServer(http.FileServer(http.Dir(testSource(t, dir...))))
	url := srv.URL + "/" + path.Join(pp...)
	return srv, url
}
