package federation

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/cortezaproject/corteza-server/app"
	"github.com/cortezaproject/corteza-server/federation/rest"
	"github.com/cortezaproject/corteza-server/federation/service"
	"github.com/cortezaproject/corteza-server/pkg/api/server"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"github.com/go-chi/chi"
	_ "github.com/joho/godotenv/autoload"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type (
	helper struct {
		t *testing.T
		a *require.Assertions

		cUser  *sysTypes.User
		roleID uint64
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
			service.DefaultStore = app.Store
			rbac.SetGlobal(rbac.NewTestService(zap.NewNop(), app.Store))
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

	h.cUser.SetRoles([]uint64{h.roleID})

	rbac.Global().(*rbac.TestService).ClearGrants()
	h.mockPermissionsWithAccess()

	return h
}

// Returns context w/ security details
func (h helper) secCtx() context.Context {
	return auth.SetIdentityToContext(context.Background(), h.cUser)
}

// apitest basics, initialize, set handler, add auth
func (h helper) apiInit() *apitest.APITest {
	//InitTestApp()

	return apitest.
		New().
		Handler(r).
		Intercept(helpers.ReqHeaderAuthBearer(h.cUser))

}

func (h helper) mockPermissions(rules ...*rbac.Rule) {
	//h.noError(rbac.Global().Grant(
	//	// TestService we use does not have any backend storage,
	//	context.Background(),
	//	// We want to make sure we did not make a mistake with any of the mocked resources or actions
	//	service.DefaultAccessControl.Whitelist(),
	//	rules...,
	//))
}

// Prepends allow access rule for federation service for everyone
func (h helper) mockPermissionsWithAccess(rules ...*rbac.Rule) {
	rules = append(
		rules,
		//rbac.AllowRule(rbac.EveryoneRoleID, types.federationRBACResource, "access"),
	)

	h.mockPermissions(rules...)
}

// Set allow permision for test role
func (h helper) allow(r rbac.Resource, o rbac.Operation) {
	h.mockPermissions(rbac.AllowRule(h.roleID, r, o))
}

// set deny permission for test role
func (h helper) deny(r rbac.Resource, o rbac.Operation) {
	h.mockPermissions(rbac.DenyRule(h.roleID, r, o))
}

// Unwraps error before it passes it to the tester
func (h helper) noError(err error) {
	for errors.Unwrap(err) != nil {
		err = errors.Unwrap(err)
	}

	h.a.NoError(err)
}
