package system

import (
	"context"
	"embed"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/cortezaproject/corteza-server/app"
	"github.com/cortezaproject/corteza-server/auth/handlers"
	"github.com/cortezaproject/corteza-server/auth/request"
	"github.com/cortezaproject/corteza-server/auth/saml"
	"github.com/cortezaproject/corteza-server/pkg/api/server"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/label"
	ltype "github.com/cortezaproject/corteza-server/pkg/label/types"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/objstore/plain"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/sqlite3"
	"github.com/cortezaproject/corteza-server/system/rest"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"github.com/go-chi/chi"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/afero"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

//go:embed static
var mockData embed.FS

type (
	helper struct {
		t  *testing.T
		a  *require.Assertions
		sp *saml.SamlSPService

		cUser  *types.User
		roleID uint64
		token  string
		data   embed.FS
	}
)

var (
	testApp *app.CortezaApp
	r       chi.Router
	hh      *handlers.AuthHandlers
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
	var sm *request.SessionManager

	if testApp == nil {
		ctx := cli.Context()

		testApp = helpers.NewIntegrationTestApp(ctx, func(app *app.CortezaApp) (err error) {
			service.CurrentSettings.Auth.External.Enabled = true
			service.DefaultObjectStore, err = plain.NewWithAfero(afero.NewMemMapFs(), "test")
			if err != nil {
				return err
			}

			service.DefaultStore, err = sqlite3.ConnectInMemory(ctx)
			if err != nil {
				return err
			}

			sm = request.NewSessionManager(service.DefaultStore, app.Opt.Auth, service.DefaultLogger)

			return nil
		})
	}

	sp, _ := loadSAMLService(context.Background())

	hh = &handlers.AuthHandlers{
		SamlSPService:  *sp,
		Log:            zap.NewNop(),
		AuthService:    service.DefaultAuth,
		SessionManager: sm,
	}

	if r == nil {
		r = chi.NewRouter()
		r.Use(server.BaseMiddleware(false, logger.Default())...)

		helpers.BindAuthMiddleware(r)
		rest.MountRoutes(r)
		hh.MountHttpRoutes(r)
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
		cUser: &types.User{
			ID: id.Next(),
		},
		data: mockData,
	}

	h.cUser.SetRoles(h.roleID)
	helpers.UpdateRBAC(h.roleID)
	h.mockPermissionsWithAccess()

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

func (h helper) assertBody(e string, s io.ReadCloser) {
	b, err := io.ReadAll(s)
	h.a.NoError(err)
	h.a.Equal(e, string(b))
}

func (h helper) clearTemplates() {
	h.noError(store.TruncateTemplates(context.Background(), service.DefaultStore))
}

func readStaticFile(f string) []byte {
	c, _ := mockData.ReadFile(f)
	return c
}
