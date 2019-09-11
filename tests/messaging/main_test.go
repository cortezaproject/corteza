package messaging

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

	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/internal/permissions"
	"github.com/cortezaproject/corteza-server/internal/rand"
	"github.com/cortezaproject/corteza-server/internal/store"
	"github.com/cortezaproject/corteza-server/messaging"
	migrate "github.com/cortezaproject/corteza-server/messaging/db"
	"github.com/cortezaproject/corteza-server/messaging/rest"
	"github.com/cortezaproject/corteza-server/messaging/service"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/logger"
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
	cfg *cli.Config
	r   chi.Router
	p   = permissions.NewTestService()
)

func db() *factory.DB {
	return factory.Database.MustGet("messaging").With(context.Background())
}

func InitConfig() {
	var err error

	if cfg != nil {
		return
	}

	helpers.RecursiveDotEnvLoad()

	ctx := context.Background()
	log, _ := zap.NewDevelopment()

	cfg = messaging.Configure()
	cfg.Log = log

	cfg.Init()

	auth.SetupDefault(string(rand.Bytes(32)), 10)

	if err = cfg.RootCommandDBSetup.Run(ctx, nil, cfg); err != nil {
		panic(err)
	} else if err := migrate.Migrate(factory.Database.MustGet("messaging"), log); err != nil {
		panic(err)
	}

	logger.SetDefault(log)
	service.DefaultPermissions = p
	if service.DefaultStore, err = store.NewWithAfero(afero.NewMemMapFs(), "test"); err != nil {
		panic(err)
	}

	cfg.InitServices(ctx, cfg)
}

func InitApp() {
	InitConfig()
	helpers.InitAuth()

	if r != nil {
		return
	}

	r = chi.NewRouter()
	r.Use(api.Base(logger.Default())...)
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

// Prepends allow access rule for messaging service for everyone
func (h helper) mockPermissionsWithAccess(rules ...*permissions.Rule) {
	rules = append(
		rules,
		permissions.AllowRule(permissions.EveryoneRoleID, types.MessagingPermissionResource, "access"),
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
