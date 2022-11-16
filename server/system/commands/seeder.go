package commands

import (
	"context"
	"github.com/cortezaproject/corteza/server/pkg/cli"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/seeder"

	"github.com/spf13/cobra"
)

type (
	seederService interface {
		CreateUser(seeder.Params) ([]uint64, error)
		CreateRecord(seeder.RecordParams) ([]uint64, error)
		DeleteAllUser() error
	}
)

var (
	svc seederService
)

func SeedUsers(ctx context.Context, app serviceInitializer) (cmd *cobra.Command) {
	var (
		limit int
	)
	cmd = &cobra.Command{
		Use:   "users",
		Short: "Seed users",
		Args:  cobra.MaximumNArgs(0),

		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			if err = app.InitServices(cli.Context()); err != nil {
				return err
			}
			svc = seeder.Seeder(ctx, seeder.DefaultStore, dal.Service(), seeder.Faker())
			return nil
		},
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "create",
			Short: "Create users",
			Args:  cobra.MaximumNArgs(0),
			Run: func(cmd *cobra.Command, args []string) {
				userIDs, err := svc.CreateUser(seeder.Params{Limit: limit})
				cli.HandleError(err)

				cmd.Printf("                     Created    %d    users", len(userIDs))
				cmd.Println()
			},
		},
		&cobra.Command{
			Use:   "delete",
			Short: "Delete users",
			Args:  cobra.MaximumNArgs(0),
			Run: func(cmd *cobra.Command, args []string) {
				cli.HandleError(svc.DeleteAllUser())

				cmd.Println("                     Deleted    all    users")

			},
		},
	)

	cmdCreate := cmd.Commands()[0]
	cmdCreate.Flags().IntVarP(&limit, "limit", "l", 1, "How many users to be created")

	return cmd
}
