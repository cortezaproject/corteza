package dal

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/steinfletcher/apitest"

	"github.com/cortezaproject/corteza-server/app"
	"github.com/cortezaproject/corteza-server/auth/handlers"
	"github.com/cortezaproject/corteza-server/auth/request"
	composeRest "github.com/cortezaproject/corteza-server/compose/rest"
	composeService "github.com/cortezaproject/corteza-server/compose/service"
	composeTypes "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/api/server"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/objstore/plain"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/drivers/sqlite"
	"github.com/cortezaproject/corteza-server/system/rest"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"github.com/go-chi/chi/v5"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

type (
	helper struct {
		t *testing.T
		b *testing.B
		a *require.Assertions

		cUser  *types.User
		roleID uint64
		token  []byte
	}
)

var (
	testApp  *app.CortezaApp
	r        chi.Router
	testUser *types.User
	hh       *handlers.AuthHandlers
)

func init() {
	helpers.RecursiveDotEnvLoad()
}

func TestMain(m *testing.M) {
	InitTestApp()
	os.Exit(m.Run())
}

func InitTestApp() {
	var sm *request.SessionManager

	if testApp == nil {
		ctx := cli.Context()

		testApp = helpers.NewIntegrationTestApp(ctx, func(app *app.CortezaApp) (err error) {
			service.CurrentSettings.Auth.External.Enabled = true
			service.DefaultObjectStore, err = plain.NewWithAfero(afero.NewMemMapFs(), "test")
			if err != nil {
				return err
			}

			service.DefaultStore, err = sqlite.ConnectInMemory(ctx)
			if err != nil {
				return err
			}

			// Tests should be executed w/o any locales
			locale.SetGlobal(locale.Static(&locale.Language{Tag: language.Und}))

			sm = request.NewSessionManager(service.DefaultStore, app.Opt.Auth, service.DefaultLogger)

			return nil
		})

		if err := store.TruncateUsers(context.Background(), service.DefaultStore); err != nil {
			panic(fmt.Errorf("could not cleanup users: %v", err))
		}

		if err := testApp.Activate(ctx); err != nil {
			panic(fmt.Errorf("could not activate corteza: %v", err))
		}

	}

	hh = &handlers.AuthHandlers{
		Log:            zap.NewNop(),
		AuthService:    service.DefaultAuth,
		SessionManager: sm,
	}

	if r == nil {
		r = chi.NewRouter()
		r.Use(server.BaseMiddleware(false, logger.Default())...)

		helpers.BindAuthMiddleware(r)
		r.Route("/system", func(r chi.Router) {
			r.Group(rest.MountRoutes())
		})
		r.Route("/compose", func(r chi.Router) {
			r.Group(composeRest.MountRoutes())
		})
		hh.MountHttpRoutes(r)
	}
}

func newHelperT(t *testing.T) helper {
	h := newHelper(t, require.New(t))
	h.t = t
	return h
}

func newHelperB(b *testing.B) helper {
	h := newHelper(b, require.New(b))
	h.b = b
	return h
}

func newHelper(_ require.TestingT, a *require.Assertions) helper {
	var (
		h = helper{
			a:      a,
			roleID: id.Next(),
		}

		ctx = context.Background()
	)

	if testUser == nil {
		testUser = &types.User{
			Handle:    "test_user",
			Name:      "test_user",
			ID:        id.Next(),
			CreatedAt: time.Now(),
		}

		err := store.CreateUser(ctx, service.DefaultStore, testUser)
		if err != nil {
			panic(err)
		}

	}
	h.cUser = testUser

	h.cUser.SetRoles(h.roleID)
	helpers.UpdateRBAC(h.roleID)
	h.identityToHelper(ctx, h.cUser)
	h.mockPermissionsWithAccess()

	return h
}

// apiInit basics, initialize, set handler, add auth
func (h helper) apiInit() *apitest.APITest {
	InitTestApp()

	return apitest.
		New().
		Handler(r).
		Intercept(helpers.ReqHeaderRawAuthBearer(h.token))
}

func (h helper) MyRole() uint64 {
	return h.roleID
}

// Returns context w/ security details
func (h helper) secCtx() context.Context {
	return auth.SetIdentityToContext(context.Background(), h.cUser)
}

func (h *helper) identityToHelper(ctx context.Context, u *types.User) {
	var err error
	h.cUser = u

	h.token, err = auth.TokenIssuer.Issue(ctx, auth.WithIdentity(h.cUser))
	if err != nil {
		panic(err)
	}
}

func (h helper) mockPermissions(rules ...*rbac.Rule) {
	h.a.NoError(rbac.Global().Grant(
		// TestService we use does not have any backend storage,
		context.Background(),
		rules...,
	))
}

// Prepends allow access rule for system service for everyone
func (h helper) mockPermissionsWithAccess(rules ...*rbac.Rule) {
	h.mockPermissions(rules...)
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Resource utilities

func (h helper) createNamespace(name string) *composeTypes.Namespace {
	ns := &composeTypes.Namespace{Name: name, Slug: name}
	ns.ID = id.Next()
	ns.CreatedAt = time.Now()
	h.a.NoError(store.CreateComposeNamespace(context.Background(), service.DefaultStore, ns))
	return ns
}

func (h helper) createSensitivityLevel(res *types.DalSensitivityLevel) *types.DalSensitivityLevel {
	if res.ID == 0 {
		res.ID = id.Next()
	}

	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}

	h.a.NoError(service.DefaultStore.CreateDalSensitivityLevel(context.Background(), res))
	h.a.NoError(service.DefaultDalSensitivityLevel.ReloadSensitivityLevels(context.Background(), service.DefaultStore))
	return res
}

func (h helper) createDalConnection(res *types.DalConnection) *types.DalConnection {
	if res.ID == 0 {
		res.ID = id.Next()
	}

	if res.Meta.Name == "" {
		res.Meta.Name = "Test Connection"
	}
	if res.Handle == "" {
		res.Handle = "test_connection"
	}
	if res.Type == "" {
		res.Type = types.DalConnectionResourceType
	}
	if res.Meta.Ownership == "" {
		res.Meta.Ownership = "tester"
	}

	if res.Config.DAL == nil {
		res.Config.DAL = &types.ConnectionConfigDAL{
			Type: "corteza::dal:connection:dsn",
			Params: map[string]any{
				"dsn": "sqlite3://file::memory:?cache=shared&mode=memory",
			},
		}

	}

	if res.Config.DAL.ModelIdent == "" {
		res.Config.DAL.ModelIdent = "compose_records"
	}

	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}
	if res.CreatedBy == 0 {
		res.CreatedBy = h.cUser.ID
	}

	h.a.NoError(service.DefaultStore.CreateDalConnection(context.Background(), res))
	h.a.NoError(service.DefaultDalConnection.ReloadConnections(context.Background()))
	return res
}

func (h helper) getPrimaryConnection() *types.DalConnection {
	cc, _, err := store.SearchDalConnections(context.Background(), service.DefaultStore, types.DalConnectionFilter{Type: types.DalPrimaryConnectionResourceType})
	h.a.NoError(err)

	if len(cc) != 1 {
		h.a.FailNow("invalid state: no or too many primary connections")
	}

	return cc[0]
}

func makeConnectionDefinition(dsn string) *types.DalConnection {
	return &types.DalConnection{
		ID:   id.Next(),
		Type: types.DalConnectionResourceType,
		Config: types.ConnectionConfig{
			DAL: &types.ConnectionConfigDAL{
				ModelIdent: "compose_record",
			},
		},
	}
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Utilities

// ---

func loadRequestFromScenario(t *testing.T, name string) string {
	return loadRequestFrom(t, suiteFromT(t), name)
}
func loadRequestFromGenerics(t *testing.T, name string) string {
	return loadRequestFrom(t, "generic", name)
}
func loadRequestFrom(t *testing.T, suite, name string) string {
	f, err := os.Open(path.Join("testdata", suite, name))
	require.NoError(t, err)
	defer f.Close()

	bb, err := ioutil.ReadAll(f)
	require.NoError(t, err)

	return string(bb)
}
func loadRequestFromScenarioWithConnection(t *testing.T, req string, connectionID uint64) string {
	out := loadRequestFromScenario(t, req)

	aux := &composeTypes.Module{}
	require.NoError(t, json.Unmarshal([]byte(out), &aux))

	aux.Config.DAL.ConnectionID = connectionID

	a, err := json.Marshal(aux)
	require.NoError(t, err)

	return string(a)
}

func createRecordFromCase(ctx context.Context, t *testing.T, name string, namespaceID, moduleID uint64) (record *composeTypes.Record) {
	return createRecordFrom(ctx, t, suiteFromT(t), name, namespaceID, moduleID)
}
func createRecordFromGenerics(ctx context.Context, t *testing.T, name string, namespaceID, moduleID uint64) (record *composeTypes.Record) {
	return createRecordFrom(ctx, t, "generic", name, namespaceID, moduleID)
}
func createRecordFrom(ctx context.Context, t *testing.T, suite, name string, namespaceID, moduleID uint64) (record *composeTypes.Record) {
	raw := loadRequestFrom(t, suite, name)

	record = &composeTypes.Record{}
	require.NoError(t, json.Unmarshal([]byte(raw), &record))

	record.NamespaceID = namespaceID
	record.ModuleID = moduleID

	record, _, err := composeService.DefaultRecord.Create(ctx, record)
	require.NoError(t, err)

	return record
}

func createModuleFromCase(ctx context.Context, t *testing.T, name string, namespaceID uint64, config *composeTypes.ModuleConfig) (module *composeTypes.Module) {
	return createModuleFrom(ctx, t, suiteFromT(t), name, namespaceID, config)
}
func createModuleFromGenerics(ctx context.Context, t *testing.T, name string, namespaceID uint64, config *composeTypes.ModuleConfig) (module *composeTypes.Module) {
	return createModuleFrom(ctx, t, "generic", name, namespaceID, config)
}
func createModuleFrom(ctx context.Context, t *testing.T, suite, name string, namespaceID uint64, config *composeTypes.ModuleConfig) (module *composeTypes.Module) {
	raw := loadRequestFrom(t, suite, name)

	module = &composeTypes.Module{}
	require.NoError(t, json.Unmarshal([]byte(raw), &module))

	module.NamespaceID = namespaceID

	if config != nil {
		// let's be careful not to override whole config
		module.Config.DAL = config.DAL
	}

	module, err := composeService.DefaultModule.Create(ctx, module)
	require.NoError(t, err)

	return module
}

func createConnectionFromCase(ctx context.Context, t *testing.T, name string) (connection *types.DalConnection) {
	return createConnectionFrom(ctx, t, suiteFromT(t), name)
}
func createConnectionFromGenerics(ctx context.Context, t *testing.T, name string) (connection *types.DalConnection) {
	return createConnectionFrom(ctx, t, "generic", name)
}
func createConnectionFrom(ctx context.Context, t *testing.T, suite, name string) (connection *types.DalConnection) {
	raw := loadRequestFrom(t, suite, name)

	aux := &types.DalConnection{}
	require.NoError(t, json.Unmarshal([]byte(raw), &aux))

	connection, err := service.DefaultDalConnection.Create(ctx, aux)
	require.NoError(t, err)

	return connection
}

// ---

func suiteFromT(t *testing.T) string {
	return t.Name()[5:]
}

// random string, 10 chars long by default
func rs(a ...int) string {
	var l = 10
	if len(a) > 0 {
		l = a[0]
	}

	return string(rand.Bytes(l))
}

func loadScenarioSources(t wrapTest, driver, ext string) (src string) {
	scenarioName := t.Name()[5:]

	src = loadScenarioSource(t, "generic", fmt.Sprintf("_.%s", ext))
	src += "\n"
	src += loadScenarioSource(t, "generic", fmt.Sprintf("%s.%s", driver, ext))
	src += "\n"
	src += loadScenarioSource(t, scenarioName, fmt.Sprintf("_.%s", ext))
	src += "\n"
	src += loadScenarioSource(t, scenarioName, fmt.Sprintf("%s.%s", driver, ext))

	return
}

func loadScenarioSource(t wrapTest, scenarioName, srcName string) (src string) {
	f, err := os.Open(path.Join("testdata", scenarioName, srcName))
	if err != nil && os.IsNotExist(err) {
		return ""
	}
	require.NoError(t, err)
	defer f.Close()

	out, err := ioutil.ReadAll(bufio.NewReader(f))
	require.NoError(t, err)

	return string(out)
}
