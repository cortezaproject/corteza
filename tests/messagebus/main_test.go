package messagebus

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/app"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/messagebus"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/store/sqlite3"
	"github.com/cortezaproject/corteza-server/system/service"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
	"github.com/go-chi/chi/v5"
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
			service.DefaultStore, err = sqlite3.ConnectInMemory(ctx)

			if err != nil {
				return err
			}

			eventbus.Set(eventBus)

			messageBus := messagebus.New(&options.MessagebusOpt{Enabled: true, LogEnabled: false}, logger.Default())
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

	h.cUser.SetRoles(h.roleID)
	helpers.UpdateRBAC(h.roleID)

	return h
}

// Returns context w/ security details
func (h helper) secCtx() context.Context {
	return auth.SetIdentityToContext(context.Background(), h.cUser)
}

// Unwraps error before it passes it to the tester
func (h helper) noError(err error) {
	for errors.Unwrap(err) != nil {
		err = errors.Unwrap(err)
	}

	h.a.NoError(err)
}

func (h helper) prepareQueues(ctx context.Context, qs ...*sysTypes.Queue) {
	h.noError(testApp.Store.TruncateQueues(ctx))
	h.noError(testApp.Store.CreateQueue(ctx, qs...))
}

func (h helper) prepareMessages(ctx context.Context, qs ...*sysTypes.Queue) {
	h.noError(testApp.Store.TruncateQueueMessages(ctx))
}

func (h helper) checkPersistedMessages(ctx context.Context, f sysTypes.QueueMessageFilter) sysTypes.QueueMessageSet {
	s, f, err := service.DefaultStore.SearchQueueMessages(ctx, f)
	h.noError(err)

	return s
}

func (h helper) initMessagebus(ctx context.Context) {
	// re-init
	messagebus.Service().Init(ctx, service.DefaultQueue)

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
