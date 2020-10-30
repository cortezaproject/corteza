package commands

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	cs "github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/federation/service"
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/decoder"
	"github.com/spf13/cobra"
)

type (
	structureProcesser struct {
		NodeID uint64
	}
)

func commandSyncStructure(ctx context.Context) func(*cobra.Command, []string) {

	const syncDelay = 10

	return func(_ *cobra.Command, _ []string) {

		// todo - get auth from the node
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		ticker := time.NewTicker(time.Second * syncDelay)

		sync = service.NewSync(
			&service.Syncer{},
			&service.Mapper{},
			service.DefaultSharedModule,
			cs.DefaultRecord)

		queueStructureForNodes(ctx, sync)

		for {
			select {
			case <-ctx.Done():
				// do any cleanup here
				service.DefaultLogger.Info(fmt.Sprintf("Stopping sync [processed: %d, persisted: %d]", countProcess, countPersist))
				return
			case <-ticker.C:
				// do the whole process again
				queueStructureForNodes(ctx, sync)
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
			case p := <-spayloads:
				body, err := ioutil.ReadAll(p.Payload)

				// handle error
				if err != nil {
					continue
				}

				basePath := fmt.Sprintf("/federation/nodes/%d/modules/exposed/", p.Meta.(*structureProcesser).NodeID)

				u := types.SyncerURI{
					BaseURL: baseURL,
					Path:    basePath,
					Limit:   limit,
				}

				err = sync.ProcessPayload(ctx, body, surls, u, p.Meta.(*structureProcesser))

				if err != nil {
					// handle error
					service.DefaultLogger.Info(fmt.Sprintf("error on handling payload: %s", err))
				} else {
					n, err := service.DefaultNode.FindBySharedNodeID(ctx, p.Meta.(*structureProcesser).NodeID)

					if err != nil {
						service.DefaultLogger.Info(fmt.Sprintf("could not update sync status: %s", err))
					}

					// add to db - nodes_sync
					new := &types.NodeSync{
						NodeID:       (*n).ID,
						SyncStatus:   types.NodeSyncStatusSuccess,
						SyncType:     types.NodeSyncTypeStructure,
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

func queueStructureForNodes(ctx context.Context, sync *service.Sync) {
	nodes, err := sync.GetPairedNodes(ctx)

	if err != nil {
		return
	}

	// get all shared modules and their module mappings
	for _, n := range nodes {
		// get the last sync per-node
		lastSync, _ := sync.GetLastSyncTime(ctx, n.ID, types.NodeSyncTypeStructure)
		basePath := fmt.Sprintf("/federation/nodes/%d/modules/exposed/", n.SharedNodeID)

		ex := ""
		if lastSync != nil {
			ex = fmt.Sprintf(" from last sync on (%s)", lastSync.Format(time.RFC3339))
		}

		service.DefaultLogger.Info(fmt.Sprintf("starting structure sync%s for node %d", ex, n.ID))

		url := types.SyncerURI{
			BaseURL:  baseURL,
			Path:     basePath,
			Limit:    limit,
			LastSync: lastSync,
		}

		processer := &structureProcesser{
			NodeID: n.SharedNodeID,
		}

		go queueUrl(&url, surls, processer)
	}
}

// Process gets the payload from syncer and
// uses the decode package to decode the whole set, depending on
// the filtering that was used (limit)
func (dp *structureProcesser) Process(ctx context.Context, payload []byte) error {
	countProcess = countProcess + 1

	now := time.Now()
	o, err := decoder.DecodeFederationModuleSync([]byte(payload))

	if err != nil {
		return err
	}

	if len(o) == 0 {
		return nil
	}

	service.DefaultLogger.Info(fmt.Sprintf("Adding %d objects", len(o)))

	for i, em := range o {
		n := &types.SharedModule{
			NodeID:                     dp.NodeID,
			ExternalFederationModuleID: em.ID,
			Fields:                     em.Fields,
			Handle:                     fmt.Sprintf("Handle %d %d", i, now.Unix()),
			Name:                       fmt.Sprintf("Name %d %d", i, now.Unix()),
		}

		n, err := sync.CreateSharedModule(ctx, n)

		service.DefaultLogger.Info(fmt.Sprintf("Added shared module: %d", n.ID))

		if err != nil {
			service.DefaultLogger.Error(fmt.Sprintf("could not create shared module: %s", err))
			continue
		}

		countPersist = countPersist + 1
	}

	return nil
}
