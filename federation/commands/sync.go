package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/cortezaproject/corteza-server/federation/service"
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/cli"

	"github.com/spf13/cobra"
)

const (
	httpTimeout time.Duration = 10
	baseURL     string        = "http://localhost:8084"
)

type (
	serviceInitializer interface {
		InitServices(ctx context.Context) error
	}
)

func Sync(app serviceInitializer) *cobra.Command {

	ctx := cli.Context()

	// Sync commands.
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync commands",
	}

	// Sync structure.
	syncStructureCmd := &cobra.Command{
		Use:   "structure",
		Short: "Sync structure",

		PreRunE: commandPreRunInitService(ctx, app),
		Run:     commandSyncStructure(ctx),
	}

	syncDataCmd := &cobra.Command{
		Use:   "data",
		Short: "Sync data",

		PreRunE: commandPreRunInitService(ctx, app),
		Run:     commandSyncData(ctx),
	}

	cmd.AddCommand(
		syncStructureCmd,
		syncDataCmd,
	)

	return cmd
}

func queueUrl(url *types.SyncerURI, urls chan types.SyncerURI, handler *service.Syncer) {
	s, _ := url.String()

	service.DefaultLogger.Info(fmt.Sprintf("Adding %s to queue", s))
	handler.Queue(*url, urls)
}

func getLastSyncTime(ctx context.Context, nodeID uint64, syncType string) *time.Time {
	ns, _ := service.DefaultNodeSync.LookupLastSuccessfulSync(ctx, nodeID, syncType)

	if ns != nil {
		return &ns.TimeOfAction
	}

	return nil
}

func commandPreRunInitService(ctx context.Context, app serviceInitializer) func(*cobra.Command, []string) error {
	return func(_ *cobra.Command, _ []string) error {
		return app.InitServices(ctx)
	}
}
