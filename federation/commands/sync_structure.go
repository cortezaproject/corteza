package commands

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/cortezaproject/corteza-server/federation/service"
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/decoder"
	"github.com/spf13/cobra"
)

var (
	urls            = make(chan types.SyncerURI, 1)
	payloads        = make(chan io.Reader, 1)
	countProcess    = 0
	countPersist    = 0
	structureSyncer *service.Syncer
)

func commandSyncStructure(ctx context.Context) func(*cobra.Command, []string) {

	return func(_ *cobra.Command, _ []string) {

		// todo - get node by id
		const (
			nodeID    = 276342359342989444
			limit     = 5
			syncDelay = 10
		)

		node, err := service.DefaultNode.FindByID(ctx, nodeID)

		if err != nil {
			service.DefaultLogger.Info("could not find any nodes to sync")
			return
		}

		ticker := time.NewTicker(time.Second * syncDelay)
		basePath := fmt.Sprintf("/federation/nodes/%d/modules/exposed/", node.SharedNodeID)

		structureSyncer = service.NewSyncer()

		// todo - get auth from the node
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		// get the last sync per-node
		lastSync := getLastSyncTime(ctx, nodeID, types.NodeSyncTypeStructure)

		ex := ""
		if lastSync != nil {
			ex = fmt.Sprintf(" from last sync on (%s)", lastSync.Format(time.RFC3339))
		}

		service.DefaultLogger.Info(fmt.Sprintf("Starting structure sync%s", ex))

		url := types.SyncerURI{
			BaseURL:  baseURL,
			Path:     basePath,
			Limit:    limit,
			LastSync: lastSync,
		}

		go queueUrl(&url, urls, structureSyncer)

		for {
			select {
			case <-ctx.Done():
				// do any cleanup here
				service.DefaultLogger.Info(fmt.Sprintf("Stopping sync [processed: %d, persisted: %d]", countProcess, countPersist))
				return
			case <-ticker.C:
				lastSync := getLastSyncTime(ctx, nodeID, types.NodeSyncTypeStructure)

				url := types.SyncerURI{
					BaseURL:  baseURL,
					Path:     basePath,
					Limit:    limit,
					LastSync: lastSync,
				}

				if lastSync == nil {
					service.DefaultLogger.Info(fmt.Sprintf("Start fetching modules from beginning of time, since lastSync == nil"))
				} else {
					service.DefaultLogger.Info(fmt.Sprintf("Start fetching modules, start with time: %s", lastSync.Format(time.RFC3339)))
				}

				go queueUrl(&url, urls, structureSyncer)
			case url := <-urls:
				select {
				case <-ctx.Done():
					// do any cleanup here
					service.DefaultLogger.Info(fmt.Sprintf("Stopping sync [processed: %d, persisted: %d]", countProcess, countPersist))
					return
				default:
				}

				s, err := url.String()

				if err != nil {
					continue
				}

				responseBody, err := structureSyncer.Fetch(ctx, s)

				if err != nil {
					continue
				}

				payloads <- responseBody
			case p := <-payloads:
				body, err := ioutil.ReadAll(p)

				// handle error
				if err != nil {
					continue
				}

				u := types.SyncerURI{
					BaseURL: baseURL,
					Path:    basePath,
					Limit:   limit,
				}

				err = structureSyncer.Process(ctx, body, node.ID, urls, u, structureSyncer, persistExposedModuleStructure)

				if err != nil {
					// handle error
					service.DefaultLogger.Error(fmt.Sprintf("error on handling payload: %s", err))
				} else {
					n := time.Now().UTC()

					// add to db - nodes_sync
					new := &types.NodeSync{
						NodeID:       node.ID,
						SyncStatus:   types.NodeSyncStatusSuccess,
						SyncType:     types.NodeSyncTypeStructure,
						TimeOfAction: n,
					}

					new, err := service.DefaultNodeSync.Create(ctx, new)

					if err != nil {
						service.DefaultLogger.Warn(fmt.Sprintf("could not update sync status: %s", err))
					}
				}
			}
		}
	}
}

// persistExposedModuleStructure gets the payload from syncer and
// uses the decode package to decode the whole set, depending on
// the filtering that was used (limit)
func persistExposedModuleStructure(ctx context.Context, payload []byte, nodeID uint64, structureSyncer *service.Syncer) error {
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
			NodeID:                     nodeID,
			ExternalFederationModuleID: em.ID,
			Fields:                     em.Fields,
			Handle:                     fmt.Sprintf("Handle %d %d", i, now.Unix()),
			Name:                       fmt.Sprintf("Name %d %d", i, now.Unix()),
		}

		n, err := structureSyncer.Service.Create(ctx, n)

		service.DefaultLogger.Info(fmt.Sprintf("Added shared module: %d", n.ID))

		if err != nil {
			service.DefaultLogger.Error(fmt.Sprintf("could not create shared module: %s", err))
			continue
		}

		countPersist = countPersist + 1
	}

	return nil
}
