package commands

import (
	"context"
	"github.com/cortezaproject/corteza/server/compose/service"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/cli"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/seeder"

	"github.com/spf13/cobra"
)

type (
	serviceInitializer interface {
		InitServices(ctx context.Context) error
	}

	seederService interface {
		CreateUser(seeder.Params) ([]uint64, error)
		CreateRecord(seeder.RecordParams) ([]uint64, error)
		DeleteAllUser() error
		DeleteAllRecord(*types.Module) error
		DeleteAll(*seeder.RecordParams) (err error)
	}
)

var (
	svc seederService
)

func SeedRecords(ctx context.Context, app serviceInitializer) (cmd *cobra.Command) {
	var (
		namespaceID int
		moduleID    int

		params = seeder.RecordParams{}
	)
	cmd = &cobra.Command{
		Use:   "records",
		Short: "Seed records",

		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			if err = app.InitServices(cli.Context()); err != nil {
				return err
			}

			cmd.Println("compose:PersistentPreRunE")
			if err = service.Activate(ctx); err != nil {
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
				params.NamespaceID = uint64(namespaceID)
				params.ModuleID = uint64(moduleID)

				userIDs, err := svc.CreateRecord(params)
				cli.HandleError(err)

				cmd.Printf("                     Created    %d    records", len(userIDs))
				cmd.Println()
			},
		},
		&cobra.Command{
			Use:   "delete",
			Short: "Delete records",
			Args:  cobra.MaximumNArgs(0),
			Run: func(cmd *cobra.Command, args []string) {
				cli.HandleError(svc.DeleteAllRecord(&types.Module{}))

				cmd.Println("                     Deleted    all    records")

			},
		},
	)

	cmdCreate := cmd.Commands()[0]
	cmdCreate.Flags().IntVarP(&params.Limit, "limit", "l", 1, "How many records to be created")
	// @todo: Can be improved by adding one flag which accept string(handle) or id(int) for namespace and module
	cmdCreate.Flags().IntVarP(&namespaceID, "namespace-id", "n", 0, "namespace id for recode creation")
	cmdCreate.Flags().StringVarP(&params.NamespaceHandle, "namespace-handle", "a", "", "namespace handle for recode creation")
	cmdCreate.Flags().IntVarP(&moduleID, "module-id", "m", 0, "module id for recode creation")
	cmdCreate.Flags().StringVarP(&params.ModuleHandle, "module-handle", "b", "", "module handle for recode creation")

	return cmd
}
