package seeder

import (
	"github.com/spf13/cobra"
)

// BaseCommand returns a base command for seeder
//
// Used by app init procedure
func BaseCommand(sub ...*cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "seeder",
		Short: "Seeds fake data",
		Long:  "Generates fake data for user and records",
	}

	cmd.AddCommand(sub...)
	return cmd
}
