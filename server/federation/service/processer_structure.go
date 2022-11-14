package service

import (
	"context"

	"github.com/cortezaproject/corteza-server/federation/service/decoder"
	"github.com/cortezaproject/corteza-server/federation/types"
	st "github.com/cortezaproject/corteza-server/system/types"
)

type (
	structureProcesser struct {
		SyncService *Sync
		Node        *types.Node
		User        *st.User
	}

	structureProcesserResponse struct {
		ModuleID  uint64
		Processed int
	}
)

// Process gets the payload from syncer and
// uses the decode package to decode the whole set, depending on
// the filtering that was used (limit)
func (dp *structureProcesser) Process(ctx context.Context, payload []byte) (ProcesserResponse, error) {
	processed := 0
	o, err := decoder.DecodeFederationModuleSync([]byte(payload))

	if err != nil {
		return structureProcesserResponse{
			Processed: processed,
		}, err
	}

	if len(o) == 0 {
		return structureProcesserResponse{
			Processed: processed,
		}, nil
	}

	for _, em := range o {
		new := &types.SharedModule{
			NodeID:                     dp.Node.ID,
			ExternalFederationModuleID: em.ID,
			Fields:                     em.Fields,
			Handle:                     em.Handle,
			Name:                       em.Name,
		}

		existing, err := dp.SyncService.LookupSharedModule(ctx, new)

		if err != nil {
			return structureProcesserResponse{
				ModuleID:  new.ExternalFederationModuleID,
				Processed: processed,
			}, err
		}

		// create
		{
			if existing == nil {
				_, err := dp.SyncService.CreateSharedModule(ctx, new)

				if err != nil {
					return structureProcesserResponse{
						ModuleID:  new.ExternalFederationModuleID,
						Processed: processed,
					}, err
				}

				processed++
				continue
			}
		}

		// update
		{
			canUpdate, err := dp.SyncService.CanUpdateSharedModule(ctx, new, existing)

			if err != nil {
				return structureProcesserResponse{
					ModuleID:  existing.ExternalFederationModuleID,
					Processed: processed,
				}, err
			}

			// stop the sync for this module
			if !canUpdate {
				return structureProcesserResponse{
					ModuleID:  new.ExternalFederationModuleID,
					Processed: processed,
				}, SharedModuleErrFederationSyncStructureChanged(&sharedModuleActionProps{module: existing})
			}

			existing.Name = new.Name
			existing.Handle = new.Handle
			existing.Fields = new.Fields

			_, err = dp.SyncService.UpdateSharedModule(ctx, existing)

			if err != nil {
				return structureProcesserResponse{
					ModuleID:  existing.ExternalFederationModuleID,
					Processed: processed,
				}, err
			}

			processed++
			continue
		}
	}

	// use paging, but only one per-time
	return structureProcesserResponse{
		ModuleID:  o[0].ID,
		Processed: processed,
	}, err
}
