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
	limit = 100

	httpTimeout time.Duration = 10
	baseURL     string        = "http://localhost:8084"
)

var (
	sync   *service.Sync
	mapper *service.Mapper

	surls     = make(chan service.Surl, 1)
	spayloads = make(chan service.Spayload, 1)

	countProcess = 0
	countPersist = 0
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

func queueUrl(url *types.SyncerURI, urls chan service.Surl, meta service.Processer) {
	s, _ := url.String()

	service.DefaultLogger.Info(fmt.Sprintf("Adding %s to queue", s))

	t := service.Surl{
		Url:  *url,
		Meta: meta,
	}

	sync.QueueUrl(t, urls)
}

func commandPreRunInitService(ctx context.Context, app serviceInitializer) func(*cobra.Command, []string) error {
	return func(_ *cobra.Command, _ []string) error {
		return app.InitServices(ctx)
	}
}
