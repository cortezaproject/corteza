package cli

import (
	"context"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/internal/db"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/cli/options"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/sentry"
)

type (
	Runner  func(ctx context.Context, cmd *cobra.Command, c *Config) error
	Runners []Runner

	CommandMaker  func(ctx context.Context, c *Config) *cobra.Command
	CommandMakers []CommandMaker

	FlagBinder  func(cmd *cobra.Command, c *Config)
	FlagBinders []FlagBinder

	Mounter  func(r chi.Router)
	Mounters []Mounter

	Config struct {
		init bool

		// Service name (messaging, system...)
		// See comments on other fields for how it is used.
		ServiceName string

		// Prefix for ENV variables
		EnvPrefix string

		// Logger name for internal services, defaults to ServiceName
		LoggerName string
		Log        *zap.Logger

		// General options
		SmtpOpt       *options.SMTPOpt
		JwtOpt        *options.JWTOpt
		HttpClientOpt *options.HttpClientOpt
		DbOpt         *options.DBOpt
		ProvisionOpt  *options.ProvisionOpt
		SentryOpt     *options.SentryOpt
		StorageOpt    *options.StorageOpt
		ScriptRunner  *options.CorredorOpt

		// Services will be calling each other so we need
		// to keep the config opts spearated
		GRPCServerSystem *options.GRPCServerOpt
		// GRPCServerMessaging    *options.GRPCServerOpt
		// GRPCServerCompose    *options.GRPCServerOpt

		// DB Connection name, defaults to ServiceName
		DatabaseName string

		// Root command name, , defaults to "corteza-server-<ServiceName>"
		RootCommandName string

		// Database setup/connection procedure
		// Runner autobinds default runner that tries to connect using DbOpt.DSN
		RootCommandDBSetup Runners

		// All that needs to be initialized before any sub-comman is executed
		RootCommandPreRun Runners

		// ******************************************************************

		// API Server instance
		ApiServer *api.Server

		// API Server command name
		ApiServerCommandName string

		// Code that needs to be executed before HTTP server is started
		ApiServerPreRun Runners

		// Routers that we mount on HTTP server
		ApiServerRoutes Mounters

		// Sets-up all available subcommands.
		AdtSubCommands CommandMakers

		// Database migration code
		// This is used for "provision migrate-database" command and after db connection is
		// established (if --provision-migrate-database is enabled
		ProvisionMigrateDatabase Runners

		// Access control initial setup
		// Reapplies default access control rules for roles "everyone" [1] and "admin" [2]
		ProvisionConfig Runners

		// ******************************************************************

		// This callback behaves a bit differently and should be called manually
		// from wherever we need service initialized
		InitServices func(ctx context.Context, c *Config)
	}
)

func init() {
	// Have logger ready in case we need to log anything
	// before it gets properly initialized through InitGeneralServices
	logger.Init()
}

func (rr Runners) Run(ctx context.Context, cmd *cobra.Command, c *Config) (err error) {
	for i := range rr {
		err = rr[i](ctx, cmd, c)
		if err != nil {
			return
		}
	}

	return
}

func (rr Mounters) MountRoutes(r chi.Router) {
	for i := range rr {
		rr[i](r)
	}
}

func (bb FlagBinders) Bind(cmd *cobra.Command, c *Config) {
	for i := range bb {
		bb[i](cmd, c)
	}
}

func (mm CommandMakers) Make(ctx context.Context, c *Config) []*cobra.Command {
	var (
		valid = make([]*cobra.Command, 0)
		cmd   *cobra.Command
	)

	for i := range mm {
		if cmd = mm[i](ctx, c); cmd != nil {
			valid = append(valid, cmd)
		}
	}

	return valid
}

func CombineFlagBinders(rr ...FlagBinders) (out FlagBinders) {
	for i := range rr {
		out = append(out, rr[i]...)
	}

	return out
}

func (c *Config) Init() {
	if c.init {
		return
	}

	if c.Log == nil {
		c.Log = logger.Default()
	}

	if c.LoggerName == "" {
		c.LoggerName = c.ServiceName
	}

	if c.EnvPrefix == "" {
		c.EnvPrefix = c.ServiceName
	}

	c.Log = c.Log.Named(c.LoggerName)

	if c.RootCommandName == "" {
		c.RootCommandName = "corteza-server-" + c.ServiceName
	}

	if c.ApiServerCommandName == "" {
		c.ApiServerCommandName = "serve-api"
	}

	if c.DatabaseName == "" {
		c.DatabaseName = c.ServiceName
	}

	c.SmtpOpt = options.SMTP(c.EnvPrefix)
	c.JwtOpt = options.JWT(c.EnvPrefix)
	c.HttpClientOpt = options.HttpClient(c.EnvPrefix)
	c.DbOpt = options.DB(c.ServiceName)
	c.ProvisionOpt = options.Provision(c.ServiceName)
	c.SentryOpt = options.Sentry(c.EnvPrefix)
	c.StorageOpt = options.Storage(c.EnvPrefix)
	c.ScriptRunner = options.Corredor(c.EnvPrefix)
	c.GRPCServerSystem = options.GRPCServer("system")
	// c.GRPCServerCompose = options.GRPCServer("compose")
	// c.GRPCServerMessagign = options.GRPCServer("messaging")

	if c.RootCommandDBSetup == nil {
		c.RootCommandDBSetup = Runners{func(ctx context.Context, cmd *cobra.Command, c *Config) (err error) {
			if c.DbOpt != nil {
				_, err = db.TryToConnect(ctx, c.Log, c.DatabaseName, *c.DbOpt)
				if err != nil {
					return errors.Wrap(err, "could not connect to database")
				}
			}

			return
		}}
	}

	if c.ApiServer == nil {
		c.ApiServer = api.NewServer(c.Log)
	}

	for i := range c.ApiServerRoutes {
		c.ApiServer.MountRoutes(c.ApiServerRoutes[i])
	}
}

// MakeCLI creates command line interface
//
// It tries to construct "serve-api" and "provision" sub-commands
// if configured properly (see Config struct)
func (c *Config) MakeCLI(ctx context.Context) (cmd *cobra.Command) {
	c.Init()

	cmd = &cobra.Command{
		Use:              c.RootCommandName,
		TraverseChildren: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			if err = sentry.Init(c.SentryOpt); err != nil {
				c.Log.Error("could not initialize Sentry", zap.Error(err))
			}

			defer sentry.Recover()

			InitGeneralServices(c.SmtpOpt, c.JwtOpt, c.HttpClientOpt)

			err = c.RootCommandDBSetup.Run(ctx, cmd, c)
			if err != nil {
				c.Log.Error("Failed to connect to the database", zap.Error(err))
				return nil
			}

			err = c.RootCommandPreRun.Run(ctx, cmd, c)
			if err != nil {
				c.Log.Error("Failed to run command pre-run scripts", zap.Error(err))
				return nil
			}

			return nil
		},
	}

	serveApiCmd := c.ApiServer.Command(ctx, c.ApiServerCommandName, c.EnvPrefix, func(ctx context.Context) (err error) {
		defer sentry.Recover()
		return c.ApiServerPreRun.Run(ctx, cmd, c)
	})

	cmd.AddCommand(serveApiCmd)

	if len(c.ProvisionMigrateDatabase) > 0 || len(c.ProvisionConfig) > 0 {
		var (
			provisionCmd = &cobra.Command{
				Use:   "provision",
				Short: "Provision tasks",
			}
		)

		// Add only commands with defined callbacks
		if len(c.ProvisionMigrateDatabase) > 0 {
			provisionCmd.AddCommand(&cobra.Command{
				Use:   "configuration",
				Short: "Create permissions & resources",

				RunE: func(cmd *cobra.Command, args []string) error {
					return c.ProvisionConfig.Run(ctx, nil, c)
				},
			})
		}

		// Add only commands with defined callbacks
		if len(c.ProvisionConfig) > 0 {
			provisionCmd.AddCommand(&cobra.Command{
				Use:   "migrate-database",
				Short: "Run database migration scripts",

				RunE: func(cmd *cobra.Command, args []string) error {
					return c.ProvisionMigrateDatabase.Run(ctx, nil, c)
				},
			})
		}

		cmd.AddCommand(provisionCmd)
	}

	cmd.AddCommand(c.AdtSubCommands.Make(ctx, c)...)

	return
}
