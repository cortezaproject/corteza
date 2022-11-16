package commands

import (
	"context"

	cs "github.com/cortezaproject/corteza/server/compose/service"
	"github.com/cortezaproject/corteza/server/federation/service"
	ss "github.com/cortezaproject/corteza/server/system/service"
	"github.com/spf13/cobra"
)

func commandSyncStructure(ctx context.Context) func(*cobra.Command, []string) {
	return func(_ *cobra.Command, _ []string) {
		syncService := service.NewSync(
			&service.Syncer{},
			&service.Mapper{},
			service.DefaultSharedModule,
			cs.DefaultRecord,
			ss.DefaultUser,
			ss.DefaultRole)

		syncStructure := service.WorkerStructure(syncService, service.DefaultLogger)
		syncStructure.Watch(ctx, service.DefaultOptions.StructureMonitorInterval, service.DefaultOptions.StructurePageSize)
	}
}
