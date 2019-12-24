package compose

import (
	"context"
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/afero"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/require"
	"github.com/titpetric/factory"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/compose"
	"github.com/cortezaproject/corteza-server/compose/rest"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/pkg/store/plain"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
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
	inited bool
	app    = &compose.App{}
	r      chi.Router
	p      = &permissions.TestService{}
)

func db() *factory.DB {
	return factory.Database.MustGet("compose", "default").With(context.Background())
}

// random string, 10 chars long by default
func rs(a ...int) string {
	var l = 10
	if len(a) > 0 {
		l = a[0]
	}

	return string(rand.Bytes(l))
}

func InitConfig() {
	var err error

	if inited {
		return
	}

	helpers.RecursiveDotEnvLoad()

	ctx := context.Background()
	log, _ := zap.NewDevelopment()
	logger.SetDefault(log)

	cli.HandleError(app.Connect(ctx))
	auth.SetupDefault(rs(32), 10)
	cli.HandleError(app.ProvisionMigrateDatabase(ctx))

	p = permissions.NewTestService(ctx, log, db(), "compose_permission_rules")

	logger.SetDefault(log)
	service.DefaultPermissions = p
	if service.DefaultStore, err = plain.NewWithAfero(afero.NewMemMapFs(), "test"); err != nil {
		panic(err)
	}

	cli.HandleError(app.Initialize(ctx))
	inited = true
}

func InitApp() {
	InitConfig()
	helpers.InitAuth()

	if r != nil {
		return
	}

	r = chi.NewRouter()
	r.Use(api.BaseMiddleware(logger.Default())...)
	helpers.BindAuthMiddleware(r)
	rest.MountRoutes(r)
}

func TestMain(m *testing.M) {
	InitApp()
	os.Exit(m.Run())
}

func newHelper(t *testing.T) helper {
	h := helper{
		t:      t,
		a:      require.New(t),
		roleID: factory.Sonyflake.NextID(),
		cUser: &sysTypes.User{
			ID: factory.Sonyflake.NextID(),
		},
	}

	h.cUser.SetRoles([]uint64{h.roleID})

	p.ClearGrants()
	h.mockPermissionsWithAccess()

	return h
}

// Returns context w/ security details
func (h helper) secCtx() context.Context {
	return auth.SetIdentityToContext(context.Background(), h.cUser)
}

// apitest basics, initialize, set handler, add auth
func (h helper) apiInit() *apitest.APITest {
	InitApp()

	return apitest.
		New().
		Handler(r).
		Intercept(helpers.ReqHeaderAuthBearer(h.cUser))

}

func (h helper) mockPermissions(rules ...*permissions.Rule) {
	h.a.NoError(p.Grant(
		// TestService we use does not have any backend storage,
		context.Background(),
		// We want to make sure we did not make a mistake with any of the mocked resources or actions
		service.DefaultAccessControl.Whitelist(),
		rules...,
	))
}

// Prepends allow access rule for compose service for everyone
func (h helper) mockPermissionsWithAccess(rules ...*permissions.Rule) {
	rules = append(
		rules,
		permissions.AllowRule(permissions.EveryoneRoleID, types.ComposePermissionResource, "access"),
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
