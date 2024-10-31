package federation

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/cortezaproject/corteza/server/app"
	"github.com/cortezaproject/corteza/server/federation/rest"
	"github.com/cortezaproject/corteza/server/federation/service"
	"github.com/cortezaproject/corteza/server/pkg/api/server"
	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/cli"
	"github.com/cortezaproject/corteza/server/pkg/eventbus"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/pkg/logger"
	sysTypes "github.com/cortezaproject/corteza/server/system/types"
	"github.com/cortezaproject/corteza/server/tests/helpers"
	"github.com/go-chi/chi/v5"
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
		token  []byte
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
	if testApp == nil {
		ctx := cli.Context()

		testApp = helpers.NewIntegrationTestApp(ctx, func(app *app.CortezaApp) (err error) {
			app.Opt.Federation.Enabled = true
			service.DefaultLogger = app.Log
			service.DefaultStore = app.Store

			eventbus.Set(eventBus)
			return nil
		})

	}

	if r == nil {
		r = chi.NewRouter()
		r.Use(server.BaseMiddleware(false, logger.Default())...)
		helpers.BindAuthMiddleware(r)
		r.Group(rest.MountRoutes(testApp.Opt.Limit))
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
	ctx := context.Background()
	h.token, err = auth.TokenIssuer.Issue(ctx, auth.WithIdentity(h.cUser))
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

// Unwraps error before it passes it to the tester
func (h helper) noError(err error) {
	for errors.Unwrap(err) != nil {
		err = errors.Unwrap(err)
	}

	h.a.NoError(err)
}
