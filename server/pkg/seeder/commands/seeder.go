package commands

import (
	"context"
	cService "github.com/cortezaproject/corteza/server/compose/service"
	cTypes "github.com/cortezaproject/corteza/server/compose/types"
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
		DeleteAllRecord(*cTypes.Module) error
		DeleteAll(*seeder.RecordParams) (err error)
	}
)

var (
	svc seederService
)

func Seeder(ctx context.Context, app serviceInitializer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "seeder",
		Short: "Seeds fake data",
		Long:  "Generates fake data for user and records",

		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			if err = app.InitServices(cli.Context()); err != nil {
				return err
			}

			if err = cService.Activate(ctx); err != nil {
				return err
			}

			svc = seeder.Seeder(ctx, seeder.DefaultStore, dal.Service(), seeder.Faker())
			return nil
		},
	}

	// @todo improve cli command with sub-command
	cmd.AddCommand(
		cmdUsers(),
		cmdRecords(),
		cmdDelete(),
	)

	return cmd
}

func cmdUsers() (cmd *cobra.Command) {
	var (
		limit int
	)
	cmd = &cobra.Command{
		Use:   "users",
		Short: "Seed users",
		Args:  cobra.MaximumNArgs(0),
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

func cmdRecords() (cmd *cobra.Command) {
	var (
		namespaceID int
		moduleID    int

		params = seeder.RecordParams{}
	)
	cmd = &cobra.Command{
		Use:       "records",
		Short:     "Seed records",
		Args:      cobra.MinimumNArgs(1),
		ValidArgs: []string{"create", "delete"},
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "create",
			Short: "Create records",
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
				cli.HandleError(svc.DeleteAllRecord(&cTypes.Module{}))

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

func cmdDelete() (cmd *cobra.Command) {
	var (
		namespaceID int
		moduleID    int

		params = seeder.RecordParams{}
	)
	cmd = &cobra.Command{
		Use:   "delete",
		Short: "delete all seeded data",
		Run: func(cmd *cobra.Command, args []string) {
			params.NamespaceID = uint64(namespaceID)
			params.ModuleID = uint64(moduleID)

			cli.HandleError(svc.DeleteAll(&params))

			cmd.Println("                     Deleted    all    fake    data")
		},
	}

	cmd.Flags().IntVarP(&params.Limit, "limit", "l", 1, "How many users to be created")
	// @todo: Can be improved by adding one flag which accept string(handle) or id(int) for namespace and module
	cmd.Flags().IntVarP(&namespaceID, "namespace-id", "n", 0, "namespace id for recode creation")
	cmd.Flags().StringVarP(&params.NamespaceHandle, "namespace-handle", "a", "", "namespace handle for recode creation")
	cmd.Flags().IntVarP(&moduleID, "module-id", "m", 0, "module id for recode creation")
	cmd.Flags().StringVarP(&params.ModuleHandle, "module-handle", "b", "", "module handle for recode creation")

	return cmd
}
