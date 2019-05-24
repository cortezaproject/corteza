package monolith

import (
	"context"

	"github.com/go-chi/chi"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/compose"
	"github.com/cortezaproject/corteza-server/messaging"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/cli/flags"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/system"
)

type (
	// Sets up compose messaging & system subservices and runs them as one.
	Monolith struct {
		log *zap.Logger

		// General
		logOpt        *flags.LogOpt
		smtpOpt       *flags.SMTPOpt
		jwtOpt        *flags.JWTOpt
		httpClientOpt *flags.HttpClientOpt

		compose   subservice
		messaging subservice
		system    subservice
	}

	subservice interface {
		AddCommands(cmd *cobra.Command, ctx context.Context)
		BindApiServerFlags(cmd *cobra.Command)
		StartServices(ctx context.Context) (err error)
		ApiServerPreRun(ctx context.Context) error
		ApiServerRoutes(r chi.Router)
		ProvisionMigrateDatabase(ctx context.Context) error
		ProvisionAccessControl(ctx context.Context) error
	}

	runner  func(context.Context) error
	runners []runner
)

func (rr runners) run(ctx context.Context) (err error) {
	for i := range rr {
		err = rr[i](ctx)
		if err != nil {
			return
		}
	}

	return
}

func init() {
	logger.Init(zap.DebugLevel)
}

func InitMonolith() *Monolith {
	return &Monolith{
		log:       logger.Default(),
		compose:   compose.InitCompose(),
		messaging: messaging.InitMessaging(),
		system:    system.InitSystem(),
	}
}

// Command produces cobra.Command
func (m *Monolith) Command(ctx context.Context) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:              "corteza-server-monolith",
		TraverseChildren: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			cli.InitGeneralServices(m.logOpt, m.smtpOpt, m.jwtOpt, m.httpClientOpt)

			return m.StartServices(ctx)
		},
	}

	m.BindGlobalFlags(cmd)

	srv := api.NewServer(m.log)
	serveApiCmd := srv.Command(ctx, "", m.ApiServerPreRun)

	// Bind all flags we need for serving monolith
	m.BindApiServerFlags(serveApiCmd)

	srv.MountRoutes(m.ApiServerRoutes)

	cmd.AddCommand(
		serveApiCmd,
		cli.SetupProvisionSubcommands(ctx, m),
	)

	m.AddCommands(cmd, ctx)

	return
}

// AddCommands - other commands that this subservices need
//
// We wrap each seubservice's set into a subcommand so that we do not get naming collisions
func (m *Monolith) AddCommands(cmd *cobra.Command, ctx context.Context) {
	var (
		composeCmd   = &cobra.Command{Use: "compose"}
		messagingCmd = &cobra.Command{Use: "messaging"}
		systemCmd    = &cobra.Command{Use: "system"}
	)

	m.compose.AddCommands(composeCmd, ctx)
	if len(composeCmd.Commands()) > 0 {
		cmd.AddCommand(composeCmd)
	}

	m.messaging.AddCommands(messagingCmd, ctx)
	if len(messagingCmd.Commands()) > 0 {
		cmd.AddCommand(messagingCmd)
	}

	m.system.AddCommands(systemCmd, ctx)
	if len(systemCmd.Commands()) > 0 {
		cmd.AddCommand(systemCmd)
	}

}

// Binds all global flags
func (m *Monolith) BindGlobalFlags(cmd *cobra.Command) {
	m.logOpt = flags.Log(cmd)
	m.smtpOpt = flags.SMTP(cmd)
	m.jwtOpt = flags.JWT(cmd)
	m.httpClientOpt = flags.HttpClient(cmd)
}

// BindApiServerFlags sets & binds all API server flags
func (m *Monolith) BindApiServerFlags(cmd *cobra.Command) {
	m.compose.BindApiServerFlags(cmd)
	m.messaging.BindApiServerFlags(cmd)
	m.system.BindApiServerFlags(cmd)
}

func (m *Monolith) StartServices(ctx context.Context) (err error) {
	return (runners{
		m.compose.StartServices,
		m.messaging.StartServices,
		m.system.StartServices,
	}).run(ctx)
}

// ApiServerPreRun is executed before serve-api command runs REST API server
//
// Should initialize all that needs to run in the background
func (m Monolith) ApiServerPreRun(ctx context.Context) error {
	return (runners{
		m.compose.ApiServerPreRun,
		m.messaging.ApiServerPreRun,
		m.system.ApiServerPreRun,
	}).run(ctx)
}

// ApiServerRoutes mounts api server routes
func (m *Monolith) ApiServerRoutes(r chi.Router) {
	r.Route("/compose", m.compose.ApiServerRoutes)
	r.Route("/messaging", m.messaging.ApiServerRoutes)
	r.Route("/system", m.system.ApiServerRoutes)
}

// ProvisionMigrateDatabase migrates database to new version
//
// This is ran by default on serve-api (when not explicitly disabled with --compose-provision-database=false)
// or on demand with "provision migrate-database"
func (m Monolith) ProvisionMigrateDatabase(ctx context.Context) error {
	return (runners{
		m.compose.ProvisionMigrateDatabase,
		m.messaging.ProvisionMigrateDatabase,
		m.system.ProvisionMigrateDatabase,
	}).run(ctx)
}

// ProvisionAccessControl resets access-control rules for roles admin (2) and everyone (1)
//
// Run with emand with "provision access-control-rules"
func (m Monolith) ProvisionAccessControl(ctx context.Context) error {
	return (runners{
		m.compose.ProvisionAccessControl,
		m.messaging.ProvisionAccessControl,
		m.system.ProvisionAccessControl,
	}).run(ctx)
}
