package cli

import (
	"context"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/internal/http"
	"github.com/cortezaproject/corteza-server/internal/logger"
	"github.com/cortezaproject/corteza-server/internal/mail"
	"github.com/cortezaproject/corteza-server/pkg/cli/flags"
)

// SetupProvisionCommands sets-up standard provision commands
// Deprecated: use SetupProvisionSubCommands
func SetupProvisionCommands(ac func() error, md func() error) *cobra.Command {
	var (
		cmd = &cobra.Command{
			Use:   "provision",
			Short: "Provision tasks",
		}
	)

	// Add only commands with defined callbacks
	if ac != nil {
		cmd.AddCommand(&cobra.Command{
			Use:   "access-control-rules",
			Short: "Reset access control rules & roles",

			RunE: func(cmd *cobra.Command, args []string) error {
				return ac()
			},
		})
	}

	// Add only commands with defined callbacks
	if md != nil {
		cmd.AddCommand(&cobra.Command{
			Use:   "migrate-database",
			Short: "Run database migration scripts",

			RunE: func(cmd *cobra.Command, args []string) error {
				return md()
			},
		})
	}

	return cmd
}

type (
	provisioner interface {
		ProvisionMigrateDatabase(ctx context.Context) error
		ProvisionAccessControl(ctx context.Context) error
	}
)

func SetupProvisionSubcommands(ctx context.Context, p provisioner) *cobra.Command {
	var (
		cmd = &cobra.Command{
			Use:   "provision",
			Short: "Provision tasks",
		}
	)

	// Add only commands with defined callbacks
	cmd.AddCommand(&cobra.Command{
		Use:   "access-control-rules",
		Short: "Reset access control rules & roles",

		RunE: func(cmd *cobra.Command, args []string) error {
			return p.ProvisionAccessControl(ctx)
		},
	})

	// Add only commands with defined callbacks
	cmd.AddCommand(&cobra.Command{
		Use:   "migrate-database",
		Short: "Run database migration scripts",

		RunE: func(cmd *cobra.Command, args []string) error {
			return p.ProvisionMigrateDatabase(ctx)
		},
	})

	return cmd
}

func InitGeneralServices(logOpt *flags.LogOpt, smtpOpt *flags.SMTPOpt, jwtOpt *flags.JWTOpt, httpClientOpt *flags.HttpClientOpt) {
	var logLevel = zap.InfoLevel
	_ = logLevel.Set(logOpt.Level)

	if logger.Default() == nil {
		logger.Init(logLevel)
	} else {
		logger.DefaultLevel.SetLevel(logLevel)
	}

	auth.SetupDefault(jwtOpt.Secret, jwtOpt.Expiry)
	mail.SetupDialer(smtpOpt.Host, smtpOpt.Port, smtpOpt.User, smtpOpt.Pass, smtpOpt.From)
	http.SetupDefaults(
		httpClientOpt.HttpClientTimeout,
		httpClientOpt.ClientTSLInsecure,
	)
}
