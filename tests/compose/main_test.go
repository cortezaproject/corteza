package compose

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/filter"

	"github.com/cortezaproject/corteza-server/app"
	"github.com/cortezaproject/corteza-server/compose/rest"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/api/server"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/csv"
	"github.com/cortezaproject/corteza-server/pkg/envoy/directory"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	envoyStore "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/pkg/envoy/yaml"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/objstore/plain"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"github.com/go-chi/chi/v5"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/afero"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

type (
	helper struct {
		t *testing.T
		a *require.Assertions

		cUser  *sysTypes.User
		roleID uint64
		token  []byte
	}

	dalSvc interface {
		Purge(ctx context.Context)

		GetConnectionMeta(ctx context.Context, ID uint64) (cm dal.ConnectionConfig, err error)

		SearchModels(ctx context.Context) (out dal.ModelSet, err error)
		RemoveModel(ctx context.Context, connectionID, ID uint64) (err error)
		ReplaceModel(ctx context.Context, model *dal.Model) (err error)
		ReplaceModelAttribute(ctx context.Context, model *dal.Model, old, new *dal.Attribute, trans ...dal.TransformationFunction) (err error)
		SearchModelIssues(connectionID, resourceID uint64) (out []error)

		Create(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, vv ...dal.ValueGetter) error
		Update(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, rr ...dal.ValueGetter) (err error)
		Search(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, f filter.Filter) (dal.Iterator, error)
		Lookup(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, lookup dal.ValueGetter, dst dal.ValueSetter) (err error)
		Delete(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, pkv ...dal.ValueGetter) (err error)
		Truncate(ctx context.Context, m dal.ModelRef, operations dal.OperationSet) (err error)
	}
)

var (
	testApp  *app.CortezaApp
	r        chi.Router
	testUser *sysTypes.User

	eventBus = eventbus.New()
	defStore store.Storer
	defDal   dalSvc
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
			defStore = app.Store
			service.DefaultObjectStore, err = plain.NewWithAfero(afero.NewMemMapFs(), "test")
			if err != nil {
				return err
			}

			// Tests should be executed w/o any locales
			locale.SetGlobal(locale.Static(&locale.Language{Tag: language.Und}))

			eventbus.Set(eventBus)
			return nil
		})

		defDal = dal.Service()
	}

	if r == nil {
		r = chi.NewRouter()
		r.Use(server.BaseMiddleware(false, logger.Default())...)
		helpers.BindAuthMiddleware(r)
		r.Group(rest.MountRoutes())
	}
}

func TestMain(m *testing.M) {
	InitTestApp()
	os.Exit(m.Run())
}

func newHelper(t *testing.T) helper {
	ctx := context.Background()

	h := helper{
		t:      t,
		a:      require.New(t),
		roleID: id.Next(),
	}

	if testUser == nil {
		testUser = &sysTypes.User{
			Handle: "test_user",
			Name:   "test_user",
			ID:     id.Next(),
		}

		err := store.CreateUser(ctx, service.DefaultStore, testUser)
		if err != nil {
			panic(err)
		}

	}
	h.cUser = testUser

	h.cUser.SetRoles(h.roleID)
	helpers.UpdateRBAC(h.roleID)
	h.identityToHelper(h.cUser)

	return h
}

func (h *helper) identityToHelper(u *sysTypes.User) {
	var err error
	h.cUser = u

	ctx := context.Background()
	h.token, err = auth.TokenIssuer.Issue(ctx, auth.WithIdentity(h.cUser))
	if err != nil {
		panic(err)
	}
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

func collect(ee ...error) error {
	for _, e := range ee {
		if e != nil {
			return e
		}
	}
	return nil
}

func cleanup(t *testing.T) {
	var (
		ctx = context.Background()
	)

	err := collect(
		store.TruncateComposeNamespaces(ctx, defStore),
		store.TruncateComposePages(ctx, defStore),
		store.TruncateComposeModuleFields(ctx, defStore),
		store.TruncateComposeModules(ctx, defStore),
	)
	if err != nil {
		t.Fatalf("failed to decode scenario data: %v", err)
	}

	err = truncateRecords(ctx)
	if err != nil {
		t.Fatalf("failed to truncate records: %v", err)
	}

	defDal.Purge(ctx)
}

func truncateRecords(ctx context.Context) error {
	models, err := defDal.SearchModels(ctx)
	if err != nil {
		return err
	}
	for _, model := range models {
		err = defDal.Truncate(ctx, model.ToFilter(), nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func scenario(t *testing.T) string {
	return t.Name()[5:]
}

func loadScenario(ctx context.Context, s store.Storer, t *testing.T, h helper) {
	loadScenarioWithName(ctx, s, t, h, scenario(t))
}

func loadScenarioWithName(ctx context.Context, s store.Storer, t *testing.T, h helper, scenario string) {
	cleanup(t)
	parseEnvoy(ctx, s, h, path.Join("testdata", scenario, "data_model"))
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
	se := envoyStore.NewStoreEncoder(s, dal.Service(), &envoyStore.EncoderConfig{})
	bld := envoy.NewBuilder(se)
	g, err := bld.Build(ctx, nn...)
	h.a.NoError(err)
	err = envoy.Encode(ctx, g, se)
	h.a.NoError(err)
}

func bypassRBAC(ctx context.Context) context.Context {
	u := &sysTypes.User{
		ID: id.Next(),
	}

	u.SetRoles(auth.BypassRoles().IDs()...)

	return auth.SetIdentityToContext(ctx, u)
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

func grantImportExport(h helper) {
	helpers.AllowMe(h, types.ComponentRbacResource(), "namespace.create")
	helpers.AllowMe(h, types.ComponentRbacResource(), "module.create")
	helpers.AllowMe(h, types.ComponentRbacResource(), "page.create")
	helpers.AllowMe(h, types.ComponentRbacResource(), "chart.create")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "modules.search")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "charts.search")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "pages.search")
	helpers.AllowMe(h, types.NamespaceRbacResource(0), "read")
	helpers.AllowMe(h, types.ModuleRbacResource(0, 0), "read")
	helpers.AllowMe(h, types.PageRbacResource(0, 0), "read")
	helpers.AllowMe(h, types.ChartRbacResource(0, 0), "read")
	helpers.AllowMe(h, sysTypes.ComponentRbacResource(), "roles.search")
	helpers.AllowMe(h, sysTypes.RoleRbacResource(0), "read")
}

func namespaceExportSafe(t *testing.T, h helper, namespaceID uint64) []byte {
	bb, err := namespaceExport(t, h, namespaceID)
	h.a.NoError(err)

	return bb
}

func namespaceExport(t *testing.T, h helper, namespaceID uint64) ([]byte, error) {
	out := h.apiInit().
		Get(fmt.Sprintf("/namespace/%d/export/out.zip?jwt=%s", namespaceID, h.token)).
		Expect(t).
		Status(http.StatusOK).
		End()

	defer out.Response.Body.Close()

	bb, err := ioutil.ReadAll(out.Response.Body)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(string(bb), "Error:") {
		return nil, errors.New(string(bb))
	}

	return bb, nil
}

func namespaceImportInitPathSafe(t *testing.T, h helper, pp ...string) uint64 {
	sessionID, err := namespaceImportInitPath(t, h, pp...)
	h.a.NoError(err)

	return sessionID
}

func namespaceImportInitPath(t *testing.T, h helper, pp ...string) (uint64, error) {
	f, err := os.Open(testSource(t, pp...))
	if err != nil {
		return 0, err
	}

	defer f.Close()
	bb, err := ioutil.ReadAll(f)
	if err != nil {
		return 0, err
	}

	return namespaceImportInit(t, h, bb)
}

func namespaceImportInitSafe(t *testing.T, h helper, arch []byte) uint64 {
	sessionID, err := namespaceImportInit(t, h, arch)
	h.a.NoError(err)

	return sessionID
}

func namespaceImportInit(t *testing.T, h helper, arch []byte) (uint64, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("upload", "archive.zip")
	h.noError(err)

	_, err = part.Write(arch)
	h.noError(err)
	h.noError(writer.Close())

	out := h.apiInit().
		Post("/namespace/import").
		Header("Accept", "application/json").
		Body(body.String()).
		ContentType(writer.FormDataContentType()).
		Expect(h.t).
		Status(http.StatusOK).
		End()

	defer out.Response.Body.Close()
	bb, err := ioutil.ReadAll(out.Response.Body)
	h.a.NoError(err)

	var aux struct {
		Error struct {
			Message string
		}
		Response struct {
			ID uint64 `json:"sessionID,string"`
		}
	}
	h.a.NoError(json.Unmarshal(bb, &aux))

	if aux.Error.Message != "" {
		return 0, errors.New(aux.Error.Message)
	}

	return aux.Response.ID, nil
}

func namespaceImportRun(ctx context.Context, s store.Storer, t *testing.T, h helper, sessionID uint64, name, slug string) (*types.Namespace, types.ModuleSet, types.PageSet, types.ChartSet) {
	h.apiInit().
		Post(fmt.Sprintf("/namespace/import/%d", sessionID)).
		Header("Accept", "application/json").
		FormData("name", name).
		FormData("slug", slug).
		Expect(h.t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End()

	ns, mm, pp, cc, _, err := fetchEntireNamespace(ctx, s, slug)
	h.a.NoError(err)

	return ns, mm, pp, cc
}

func testSource(t *testing.T, pp ...string) string {
	return path.Join(append([]string{"testdata", scenario(t)}, pp...)...)
}
