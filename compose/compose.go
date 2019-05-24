package compose

import (
	"context"

	"github.com/go-chi/chi"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/titpetric/factory"
	"go.uber.org/zap"

	migrate "github.com/cortezaproject/corteza-server/compose/db"
	"github.com/cortezaproject/corteza-server/compose/internal/service"
	"github.com/cortezaproject/corteza-server/compose/rest"
	"github.com/cortezaproject/corteza-server/internal/db"
	"github.com/cortezaproject/corteza-server/internal/logger"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/cli/flags"
)

const (
	compose = "compose"
)

type (
	Compose struct {
		log *zap.Logger

		// General
		logOpt        *flags.LogOpt
		smtpOpt       *flags.SMTPOpt
		jwtOpt        *flags.JWTOpt
		httpClientOpt *flags.HttpClientOpt

		// Compose specific
		dbOpt        *flags.DBOpt
		provisionOpt *flags.ProvisionOpt
	}
)

func init() {
	logger.Init(zap.DebugLevel)
}

func InitCompose() *Compose {
	return &Compose{
		log: logger.Default().Named(compose),
	}
}

// Command produces cobra.Command
func (m *Compose) Command(ctx context.Context) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:              "corteza-server-compose",
		TraverseChildren: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			cli.InitGeneralServices(m.logOpt, m.smtpOpt, m.jwtOpt, m.httpClientOpt)

			return m.StartServices(ctx)
		},
	}

	m.BindGlobalFlags(cmd)

	srv := api.NewServer(m.log)
	serveApiCmd := srv.Command(ctx, compose, m.ApiServerPreRun)

	// Bind all flags we need for serving compose
	m.BindApiServerFlags(serveApiCmd)

	srv.MountRoutes(m.ApiServerRoutes)

	cmd.AddCommand(
		serveApiCmd,
		cli.SetupProvisionSubcommands(ctx, m),
	)

	m.AddCommands(cmd, ctx)

	return
}

// AddCommands - other commands that this subservice needs
func (m *Compose) AddCommands(cmd *cobra.Command, ctx context.Context) {}

// Binds all global flags
func (m *Compose) BindGlobalFlags(cmd *cobra.Command) {
	m.logOpt = flags.Log(cmd)
	m.smtpOpt = flags.SMTP(cmd)
	m.jwtOpt = flags.JWT(cmd)
	m.httpClientOpt = flags.HttpClient(cmd)
}

// BindApiServerFlags sets & binds all API server flags
func (m *Compose) BindApiServerFlags(cmd *cobra.Command) {
	m.dbOpt = flags.DB(cmd, compose)
	m.provisionOpt = flags.Provision(cmd, compose)
}

// StartServices
func (m *Compose) StartServices(ctx context.Context) (err error) {
	_, err = db.TryToConnect(ctx, m.log, compose, m.dbOpt.DSN, m.dbOpt.Profiler)
	if err != nil {
		return errors.Wrap(err, "could not connect to database")
	}

	if m.provisionOpt.Database {
		err = m.ProvisionMigrateDatabase(ctx)
		if err != nil {
			return
		}
	}

	err = service.Init(ctx)
	if err != nil {
		return
	}

	return
}

// ApiServerPreRun is executed before serve-api command runs REST API server
//
// Should initialize all that needs to run in the background
func (m Compose) ApiServerPreRun(ctx context.Context) error {
	service.DefaultPermissions.Watch(ctx)
	return nil
}

// ApiServerRoutes mounts api server routes
func (m *Compose) ApiServerRoutes(r chi.Router) {
	rest.MountRoutes(r)
}

// ProvisionMigrateDatabase migrates database to new version
//
// This is ran by default on serve-api (when not explicitly disabled with --compose-provision-database=false)
// or on demand with "provision migrate-database"
func (m Compose) ProvisionMigrateDatabase(ctx context.Context) error {
	var db, err = factory.Database.Get(compose)
	if err != nil {
		return err
	}

	db = db.With(ctx)
	// Disable profiler for migrations
	db.Profiler = nil

	return migrate.Migrate(db)
}

// ProvisionAccessControl resets access-control rules for roles admin (2) and everyone (1)
//
// Run with emand with "provision access-control-rules"
func (m Compose) ProvisionAccessControl(ctx context.Context) error {
	var ac = service.DefaultAccessControl
	return ac.Grant(ctx, ac.DefaultRules()...)
}
