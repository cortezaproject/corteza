package cli

import (
	"context"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/internal/db"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/cli/flags"
	"github.com/cortezaproject/corteza-server/pkg/logger"
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

		// Logger name for internal services, defaults to ServiceName
		LoggerName string
		Log        *zap.Logger

		// General options/flags
		LogOpt        *flags.LogOpt
		SmtpOpt       *flags.SMTPOpt
		JwtOpt        *flags.JWTOpt
		HttpClientOpt *flags.HttpClientOpt

		// Per-service options/flags
		DbOpt        *flags.DBOpt
		ProvisionOpt *flags.ProvisionOpt

		// DB Connection name, defaults to ServiceName
		DatabaseName string

		// Root command name, , defaults to "corteza-server-<ServiceName>"
		RootCommandName string

		// Flags that are bond to root command, no (per-service) prefixed
		RootCommandBaseFlags FlagBinders

		// Prefix for flags for root command, defaults to ServiceName
		RootCommandFlagsPrefix string

		// Flags that are bond to root command, (per-service) prefixed
		RootCommandPrefixedFlags FlagBinders

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

		// Prefix for "serve-api command flags", defaults to ServiceName
		ApiServerFlagsPrefix string

		// Additional command flags for API server
		ApiServerAdtFlags FlagBinders

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
		ProvisionAccessControl Runners
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

	c.Log = c.Log.Named(c.LoggerName)

	if c.RootCommandName == "" {
		c.RootCommandName = "corteza-server-" + c.ServiceName
	}

	if c.ApiServerCommandName == "" {
		c.ApiServerCommandName = "serve-api"
	}

	if c.ApiServerFlagsPrefix == "" {
		c.ApiServerFlagsPrefix = c.ServiceName
	}

	if c.DatabaseName == "" {
		c.DatabaseName = c.ServiceName
	}

	if c.RootCommandDBSetup == nil {
		c.RootCommandDBSetup = Runners{func(ctx context.Context, cmd *cobra.Command, c *Config) (err error) {
			_, err = db.TryToConnect(ctx, c.Log, c.DatabaseName, c.DbOpt.DSN, c.DbOpt.Profiler)
			if err != nil {
				return errors.Wrap(err, "could not connect to database")
			}

			return
		}}
	}

	// Flags, not prefixed with service name
	if c.RootCommandBaseFlags == nil {
		c.RootCommandBaseFlags = FlagBinders{
			func(cmd *cobra.Command, c *Config) {
				c.LogOpt = flags.Log(cmd)
				c.SmtpOpt = flags.SMTP(cmd)
				c.JwtOpt = flags.JWT(cmd)
				c.HttpClientOpt = flags.HttpClient(cmd)
			},
		}
	}

	if c.RootCommandFlagsPrefix == "" {
		c.RootCommandFlagsPrefix = c.ServiceName
	}

	// Flags, prefixed with service name
	if c.RootCommandPrefixedFlags == nil {
		c.RootCommandPrefixedFlags = FlagBinders{
			func(cmd *cobra.Command, c *Config) {
				c.DbOpt = flags.DB(cmd, c.ServiceName)
				c.ProvisionOpt = flags.Provision(cmd, c.ServiceName)
			},
		}
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
			InitGeneralServices(c.LogOpt, c.SmtpOpt, c.JwtOpt, c.HttpClientOpt)

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

	c.RootCommandBaseFlags.Bind(cmd, c)
	c.RootCommandPrefixedFlags.Bind(cmd, c)

	serveApiCmd := c.ApiServer.Command(ctx, c.ApiServerCommandName, c.ApiServerFlagsPrefix, func(ctx context.Context) (err error) {
		return c.ApiServerPreRun.Run(ctx, cmd, c)
	})

	// Bind all flags we need for serving the API
	c.ApiServerAdtFlags.Bind(serveApiCmd, c)

	cmd.AddCommand(serveApiCmd)

	if len(c.ProvisionMigrateDatabase) > 0 || len(c.ProvisionAccessControl) > 0 {
		var (
			provisionCmd = &cobra.Command{
				Use:   "provision",
				Short: "Provision tasks",
			}
		)

		// Add only commands with defined callbacks
		if len(c.ProvisionMigrateDatabase) > 0 {
			provisionCmd.AddCommand(&cobra.Command{
				Use:   "access-control-rules",
				Short: "Reset access control rules & roles",

				RunE: func(cmd *cobra.Command, args []string) error {
					return c.ProvisionAccessControl.Run(ctx, nil, c)
				},
			})
		}

		// Add only commands with defined callbacks
		if len(c.ProvisionAccessControl) > 0 {
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
