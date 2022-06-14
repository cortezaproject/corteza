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

var (
	testApp *app.CortezaApp
	r       chi.Router
	hh      *handlers.AuthHandlers
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

	aux.ModelConfig.ConnectionID = connectionID

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

	record, err := composeService.DefaultRecord.Create(ctx, record)
	require.NoError(t, err)

	return record
}

func createModuleFromCase(ctx context.Context, t *testing.T, name string, namespaceID uint64, config *composeTypes.ModelConfig) (module *composeTypes.Module) {
	return createModuleFrom(ctx, t, suiteFromT(t), name, namespaceID, config)
}
func createModuleFromGenerics(ctx context.Context, t *testing.T, name string, namespaceID uint64, config *composeTypes.ModelConfig) (module *composeTypes.Module) {
	return createModuleFrom(ctx, t, "generic", name, namespaceID, config)
}
func createModuleFrom(ctx context.Context, t *testing.T, suite, name string, namespaceID uint64, config *composeTypes.ModelConfig) (module *composeTypes.Module) {
	raw := loadRequestFrom(t, suite, name)

	module = &composeTypes.Module{}
	require.NoError(t, json.Unmarshal([]byte(raw), &module))

	module.NamespaceID = namespaceID

	if config != nil {
		module.ModelConfig = *config
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

func loadScenarioSources(t aaaa, driver, ext string) (src string) {
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

func loadScenarioSource(t aaaa, scenarioName, srcName string) (src string) {
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
