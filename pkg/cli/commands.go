package cli

import (
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/version"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
)

var (
	rootCommand = cobra.Command{
		Use:              "corteza-server",
		Aliases:          []string{"corteza", "server"},
		TraverseChildren: true,
	}

	serveApiCommand = cobra.Command{
		Use:     "serve-api",
		Aliases: []string{"serve"},

		Short: "Start HTTP server with REST API",
	}

	upgradeCommand = cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade tasks",
	}

	provisionCommand = cobra.Command{
		Use:   "provision",
		Short: "Provision tasks",
	}

	versionCommand = cobra.Command{
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
	// Make a copy
	var (
		cmd = rootCommand

		silent, debug bool
	)

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) (err error) {
		if light[cmd.Name()] {
			// Do not run hooks on parts on help command
			return nil
		}

		if debug {
			logger.DefaultLevel.SetLevel(zap.DebugLevel)
		} else if silent {
			logger.DefaultLevel.SetLevel(zap.FatalLevel)
		}

		return ppRunEfn()
	}

	cmd.Flags().BoolVarP(&silent, "silent", "s", false, "No output")
	cmd.Flags().BoolVarP(&debug, "debug", "d", false, "Debug")

	return &cmd
}

func ServeCommand(runEfn func() error) *cobra.Command {
	// Make a copy
	var cmd = serveApiCommand

	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		if _, set := os.LookupEnv("LOG_LEVEL"); !set {
			// If LOG_LEVEL is not explicitly set, let's
			// set it to INFO so that it
			logger.DefaultLevel.SetLevel(zap.InfoLevel)
		}

		return runEfn()
	}

	return &cmd
}

func UpgradeCommand(dbfn func() error) *cobra.Command {
	// Make a copies
	var dbCmd = upgradeCommand

	dbCmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		return dbfn()
	}

	return &dbCmd
}

func ProvisionCommand(cffn func() error) *cobra.Command {
	// Make a copies
	var cfCmd = provisionCommand

	cfCmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		return cffn()
	}

	return &cfCmd
}

func VersionCommand() *cobra.Command {
	// Make a copies
	var cfCmd = versionCommand

	cfCmd.Run = func(cmd *cobra.Command, args []string) {
		cmd.Println(version.Version)
	}

	return &cfCmd
}
