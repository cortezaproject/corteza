package cli

import (
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/version"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
)

var (
	rootCommand = &cobra.Command{
		Use:              "corteza-server",
		Aliases:          []string{"corteza", "server"},
		TraverseChildren: true,
	}

	serveApiCommand = &cobra.Command{
		Use:     "serve-api",
		Aliases: []string{"serve"},

		Short: "Start HTTP server with REST API",
	}

	upgradeCommand = &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade tasks",
	}

	provisionCommand = &cobra.Command{
		Use:   "provision",
		Short: "Provision tasks",
	}

	versionCommand = &cobra.Command{
		Use:   "version",
		Short: "Version",
	}

	// List of commands that do not
	// run initialization (db connection, etc...)
	light = map[string]bool{
		"help":             true,
		versionCommand.Use: true,
	}
)

// RootCommand creates root command with a simple persistent-pre-run callback
//
// Callback is called when not executed with help subcommand
func RootCommand(ppRunEfn func() error) *cobra.Command {
	rootCommand.PersistentPreRunE = func(cmd *cobra.Command, args []string) (err error) {
		if light[cmd.Name()] {
			// Do not run hooks on parts on help command
			return nil
		}

		if ppRunEfn != nil {
			return ppRunEfn()
		}

		return nil
	}

	return rootCommand
}

func ServeCommand(runEfn func() error) *cobra.Command {
	serveApiCommand.RunE = func(cmd *cobra.Command, args []string) (err error) {
		if _, set := os.LookupEnv("LOG_LEVEL"); !set {
			// If LOG_LEVEL is not explicitly set, let's
			// set it to INFO so that it
			logger.DefaultLevel.SetLevel(zap.InfoLevel)
		}

		return runEfn()
	}

	return serveApiCommand
}

func UpgradeCommand(dbfn func() error) *cobra.Command {
	upgradeCommand.RunE = func(cmd *cobra.Command, args []string) (err error) {
		return dbfn()
	}

	return upgradeCommand
}

func ProvisionCommand(cffn func() error) *cobra.Command {
	provisionCommand.RunE = func(cmd *cobra.Command, args []string) (err error) {
		return cffn()
	}

	return provisionCommand
}

func VersionCommand() *cobra.Command {
	versionCommand.Run = func(cmd *cobra.Command, args []string) {
		cmd.Println(version.Version)
	}

	return versionCommand
}
