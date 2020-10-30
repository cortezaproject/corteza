package commands

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	cs "github.com/cortezaproject/corteza-server/compose/service"
	ct "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/federation/service"
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/decoder"
	"github.com/spf13/cobra"
)

type (
	dataProcesser struct {
		ID                  uint64
		NodeID              uint64
		ComposeModuleID     uint64
		ComposeNamespaceID  uint64
		ModuleMappingValues *ct.RecordValueSet
	}
)

func commandSyncData(ctx context.Context) func(*cobra.Command, []string) {

	const syncDelay = 10

	return func(_ *cobra.Command, _ []string) {

		// todo - get auth from the node
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		ticker := time.NewTicker(time.Second * syncDelay)

		mapper = &service.Mapper{}

		sync = service.NewSync(
			&service.Syncer{},
			&service.Mapper{},
			service.DefaultSharedModule,
			cs.DefaultRecord)

		queueDataForNodes(ctx, sync)

		for {
			select {
			case <-ctx.Done():
				// do any cleanup here
				service.DefaultLogger.Info(fmt.Sprintf("Stopping sync [processed: %d, persisted: %d]", countProcess, countPersist))
				return
			case <-ticker.C:
				// do the whole process again
				queueDataForNodes(ctx, sync)
			case url := <-surls:
				select {
				case <-ctx.Done():
					// do any cleanup here
					service.DefaultLogger.Info(fmt.Sprintf("Stopping sync [processed: %d, persisted: %d]", countProcess, countPersist))
					return
				default:
				}

				s, err := url.Url.String()

				if err != nil {
					continue
				}

				responseBody, err := sync.FetchUrl(ctx, s)

				if err != nil {
					continue
				}

				spayload := service.Spayload{
					Payload: responseBody,
					Meta:    url.Meta,
				}

				spayloads <- spayload
				// payloads <- responseBody
			case p := <-spayloads:
				body, err := ioutil.ReadAll(p.Payload)

				// handle error
				if err != nil {
					continue
				}

				basePath := fmt.Sprintf("/federation/nodes/%d/modules/%d/records/", p.Meta.(*dataProcesser).NodeID, p.Meta.(*dataProcesser).ID)

				u := types.SyncerURI{
					BaseURL: baseURL,
					Path:    basePath,
					Limit:   limit,
				}

				err = sync.ProcessPayload(ctx, body, surls, u, p.Meta.(*dataProcesser))

				if err != nil {
					// handle error
					service.DefaultLogger.Info(fmt.Sprintf("error on handling payload: %s", err))
				} else {
					n, err := service.DefaultNode.FindBySharedNodeID(ctx, p.Meta.(*dataProcesser).NodeID)

					if err != nil {
						service.DefaultLogger.Info(fmt.Sprintf("could not update sync status: %s", err))
					}

					// add to db - nodes_sync
					new := &types.NodeSync{
						NodeID:       (*n).ID,
						SyncStatus:   types.NodeSyncStatusSuccess,
						SyncType:     types.NodeSyncTypeData,
						TimeOfAction: time.Now().UTC(),
					}

					new, err = service.DefaultNodeSync.Create(ctx, new)

					if err != nil {
						service.DefaultLogger.Info(fmt.Sprintf("could not update sync status: %s", err))
					}
				}
			}
		}
	}
}

func queueDataForNodes(ctx context.Context, sync *service.Sync) {
	nodes, err := sync.GetPairedNodes(ctx)

	if err != nil {
		return
	}

	// get all shared modules and their module mappings
	for _, n := range nodes {
		set, err := sync.GetSharedModules(ctx, n.ID)

		if err != nil {
			service.DefaultLogger.Warn(fmt.Sprintf("could not get shared modules for node: %d, skipping", n.ID))
			continue
		}

		// go through set and prepare module mappings for it
		for _, sm := range set {
			mappings, _ := sync.GetModuleMappings(ctx, sm.ID)

			if mappings == nil {
				service.DefaultLogger.Info(fmt.Sprintf("could not prepare module mappings for shared module: %d, skipping", sm.ID))
				continue
			}

			mappingValues, err := sync.PrepareModuleMappings(ctx, mappings)

			if err != nil || mappingValues == nil {
				service.DefaultLogger.Info(fmt.Sprintf("could not prepare module mappings for shared module: %d, skipping", sm.ID))
				continue
			}

			// get the last sync per-node
			lastSync, _ := sync.GetLastSyncTime(ctx, n.ID, types.NodeSyncTypeData)
			basePath := fmt.Sprintf("/federation/nodes/%d/modules/%d/records/", n.SharedNodeID, sm.ExternalFederationModuleID)

			ex := ""
			if lastSync != nil {
				ex = fmt.Sprintf(" from last sync on (%s)", lastSync.Format(time.RFC3339))
			}

			service.DefaultLogger.Info(fmt.Sprintf("starting structure sync%s for node %d, module %d", ex, n.ID, sm.ID))

			url := types.SyncerURI{
				BaseURL:  baseURL,
				Path:     basePath,
				Limit:    limit,
				LastSync: lastSync,
			}

			processer := &dataProcesser{
				ID:                  sm.ExternalFederationModuleID,
				NodeID:              n.SharedNodeID,
				ComposeModuleID:     mappings.ComposeModuleID,
				ComposeNamespaceID:  mappings.ComposeNamespaceID,
				ModuleMappingValues: &mappingValues,
			}

			go queueUrl(&url, surls, processer)
		}
	}
}

// Process gets the payload from syncer and
// uses the decode package to decode the whole set, depending on
// the filtering that was used (limit)
func (dp *dataProcesser) Process(ctx context.Context, payload []byte) error {
	countProcess = countProcess + 1

	o, err := decoder.DecodeFederationRecordSync([]byte(payload))

	if err != nil {
		return err
	}

	if len(o) == 0 {
		return nil
	}

	service.DefaultLogger.Info(fmt.Sprintf("Adding %d objects", len(o)))

	for _, er := range o {
		mapper.Merge(&er.Values, dp.ModuleMappingValues)

		rec := &ct.Record{
			ModuleID:    dp.ComposeModuleID,
			NamespaceID: dp.ComposeNamespaceID,
			Values:      *dp.ModuleMappingValues,
		}

		sync.CreateRecord(ctx, rec)
	}

	return nil
}
