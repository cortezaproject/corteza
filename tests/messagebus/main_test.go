package messagebus

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/app"
	"github.com/cortezaproject/corteza-server/automation/service"
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/messagebus"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store/sqlite3"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"github.com/go-chi/chi"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
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
			rbac.SetGlobal(rbac.NewTestService(logger.Default(), app.Store))

			service.DefaultStore, err = sqlite3.ConnectInMemory(ctx)

			if err != nil {
				return err
			}

			eventbus.Set(eventBus)

			// messageBus := messagebus.New(&options.MessagebusOpt{LogEnabled: false}, zap.NewNop(), eventbus.Service())
			messageBus := messagebus.New(&options.MessagebusOpt{LogEnabled: false}, logger.Default(), eventbus.Service())
			messagebus.Set(messageBus)

			return nil
		})
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

func (h helper) mockPermissions(rules ...*rbac.Rule) {
	h.a.NoError(rbac.Global().Grant(
		// TestService we use does not have any backend storage,
		context.Background(),
		// We want to make sure we did not make a mistake with any of the mocked resources or actions
		service.DefaultAccessControl.Whitelist(),
		rules...,
	))
}

// Prepends allow access rule for messaging service for everyone
func (h helper) mockPermissionsWithAccess(rules ...*rbac.Rule) {
	rules = append(
		rules,
		rbac.AllowRule(rbac.EveryoneRoleID, types.AutomationRBACResource, "access"),
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

func (h helper) prepareRBAC() {}

func (h helper) prepareQueues(ctx context.Context, qs ...*messagebus.QueueSettings) {
	h.noError(testApp.Store.TruncateMessagebusQueuesettings(ctx))
	h.noError(testApp.Store.CreateMessagebusQueuesetting(ctx, qs...))
}

func (h helper) prepareMessages(ctx context.Context, qs ...*messagebus.QueueSettings) {
	h.noError(testApp.Store.TruncateMessagebusQueuemessages(ctx))
}

func (h helper) checkPersistedMessages(ctx context.Context, f messagebus.QueueMessageFilter) messagebus.QueueMessageSet {
	s, f, err := service.DefaultStore.SearchMessagebusQueuemessages(ctx, f)
	h.noError(err)

	return s
}

func (h helper) initMessagebus(ctx context.Context) {
	// re-init
	messagebus.Service().Init(ctx, service.DefaultStore)

	// set messagebus watchers again
	messagebus.Service().Listen(ctx)
}

func makeDelay(d time.Duration) *time.Duration {
	return &d
}

func now() *time.Time {
	c := time.Now().Round(time.Second)
	return &c
}
