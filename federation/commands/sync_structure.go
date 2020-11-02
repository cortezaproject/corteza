package commands

import (
	"context"
	"time"

	cs "github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/federation/service"
	"github.com/spf13/cobra"
)

func commandSyncStructure(ctx context.Context) func(*cobra.Command, []string) {
	return func(_ *cobra.Command, _ []string) {
		syncService := service.NewSync(
			&service.Syncer{},
			&service.Mapper{},
			service.DefaultSharedModule,
			cs.DefaultRecord)

		syncStructure := service.WorkerStructure(syncService, service.DefaultLogger)
		syncStructure.Watch(ctx, time.Second*10, 10)
	}
}
