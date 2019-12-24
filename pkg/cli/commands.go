package cli

import (
	"github.com/spf13/cobra"
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
)

// RootCommand creates root command with a simple persistent-pre-run callback
//
// Callback is called when not executed with help subcommand
func RootCommand(ppRunEfn func() error) *cobra.Command {
	// Make a copy
	var cmd = rootCommand

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) (err error) {
		if cmd.Name() == "help" {
			// Do not run hooks on parts on help command
			return nil
		}

		return ppRunEfn()
	}

	return &cmd
}

func ServeCommand(runEfn func() error) *cobra.Command {
	// Make a copy
	var cmd = serveApiCommand

	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
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
