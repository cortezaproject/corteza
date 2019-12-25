package system

import (
	"context"
	"os"
	"testing"

	"github.com/go-chi/chi"
	_ "github.com/joho/godotenv/autoload"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/require"
	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/corteza"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/app"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/system"
	"github.com/cortezaproject/corteza-server/system/rest"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

type (
	helper struct {
		t *testing.T
		a *require.Assertions

		cUser  *types.User
		roleID uint64
	}

	TestApp struct {
		helpers.TestApp
	}
)

var (
	testApp app.Runnable
	r       chi.Router
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

func db() *factory.DB {
	return factory.Database.MustGet(system.SERVICE, "default").With(context.Background())
}

func (app *TestApp) Initialize(ctx context.Context) (err error) {
	service.DefaultPermissions = permissions.NewTestService(ctx, app.Log, db(), "compose_permission_rules")
	return
}

func (app *TestApp) Activate(ctx context.Context) (err error) {
	service.DefaultPermissions.(*permissions.TestService).Reload(ctx)
	return
}

func InitTestApp() {
	if testApp == nil {
		testApp = helpers.NewIntegrationTestApp(
			system.SERVICE,
			&corteza.App{},
			&TestApp{},
			&system.App{},
		)
	}

	if r == nil {
		r = chi.NewRouter()
		r.Use(api.BaseMiddleware(logger.Default())...)
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
		roleID: factory.Sonyflake.NextID(),
		cUser: &types.User{
			ID: factory.Sonyflake.NextID(),
		},
	}

	h.cUser.SetRoles([]uint64{h.roleID})

	service.DefaultPermissions.(*permissions.TestService).ClearGrants()
	h.mockPermissionsWithAccess()

	return h
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
		Intercept(helpers.ReqHeaderAuthBearer(h.cUser))
}

func (h helper) mockPermissions(rules ...*permissions.Rule) {
	h.a.NoError(service.DefaultPermissions.(*permissions.TestService).Grant(
		// TestService we use does not have any backend storage,
		context.Background(),
		// We want to make sure we did not make a mistake with any of the mocked resources or actions
		service.DefaultAccessControl.Whitelist(),
		rules...,
	))
}

// Prepends allow access rule for messaging service for everyone
func (h helper) mockPermissionsWithAccess(rules ...*permissions.Rule) {
	rules = append(
		rules,
		permissions.AllowRule(permissions.EveryoneRoleID, types.SystemPermissionResource, "access"),
	)

	h.mockPermissions(rules...)
}

// Set allow permision for test role
func (h helper) allow(r permissions.Resource, o permissions.Operation) {
	h.mockPermissions(permissions.AllowRule(h.roleID, r, o))
}

// set deny permission for test role
func (h helper) deny(r permissions.Resource, o permissions.Operation) {
	h.mockPermissions(permissions.DenyRule(h.roleID, r, o))
}
