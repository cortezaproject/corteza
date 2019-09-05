package messaging

import (
	"context"
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/steinfletcher/apitest"

	"github.com/cortezaproject/corteza-server/messaging"
	"github.com/cortezaproject/corteza-server/messaging/rest"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/cortezaproject/corteza-server/tests/helpers"
)

var (
	cfg *cli.Config
	r   chi.Router
)

func InitConfig() {
	if cfg != nil {
		return
	}

	helpers.RecursiveDotEnvLoad()

	ctx := context.Background()
	log, _ := zap.NewDevelopment()

	cfg = messaging.Configure()
	cfg.Log = log

	cfg.Init()

	if err := cfg.RootCommandDBSetup.Run(ctx, nil, cfg); err != nil {
		panic(err)
	}

	logger.SetDefault(log)
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

func NewApiTest(name string, user *types.User) *apitest.APITest {
	InitApp()

	return apitest.
		New(name).
		Handler(r).
		Intercept(helpers.ReqHeaderAuthBearer(user))
}

func TestMain(m *testing.M) {
	InitApp()
	os.Exit(m.Run())
}
