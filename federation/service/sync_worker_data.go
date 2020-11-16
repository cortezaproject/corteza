package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	ct "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/decoder"
	"go.uber.org/zap"
)

type (
	syncWorkerData struct {
		syncService *Sync
		logger      *zap.Logger
		delay       time.Duration
		limit       int
	}

	dataProcesser struct {
		ID                  uint64
		NodeID              uint64
		ComposeModuleID     uint64
		ComposeNamespaceID  uint64
		NodeBaseURL         string
		ModuleMappings      *types.ModuleFieldMappingSet
		ModuleMappingValues *ct.RecordValueSet
		SyncService         *Sync
	}
)

func WorkerData(sync *Sync, logger *zap.Logger) *syncWorkerData {
	return &syncWorkerData{
		syncService: sync,
		logger:      logger,
		limit:       defaultPage,
		delay:       defaultDelay,
	}
}

func (w *syncWorkerData) queueUrl(url *types.SyncerURI, urls chan Url, meta Processer) {
	t := Url{
		Url:  *url,
		Meta: meta,
	}

	w.syncService.QueueUrl(t, urls)
}

func (w *syncWorkerData) PrepareForNodes(ctx context.Context, urls chan Url) {

	nodes, err := w.syncService.GetPairedNodes(ctx)

	if err != nil {
		w.logger.Info("could not get paired nodes", zap.Error(err))
		return
	}

	// get all shared modules and their module mappings
	for _, n := range nodes {
		set, err := w.syncService.GetSharedModules(ctx, n.ID)

		if err != nil {
			w.logger.Info("could not get shared modules, skipping", zap.Uint64("nodeID", n.ID), zap.Error(err))
			continue
		}

		if len(set) == 0 {
			w.logger.Info("there are no valid shared modules, skipping", zap.Uint64("nodeID", n.ID), zap.Error(err))
			continue
		}

		// go through set and prepare module mappings for it
		for _, sm := range set {
			mappings, _ := w.syncService.GetModuleMappings(ctx, sm.ID)

			if mappings == nil {
				w.logger.Info("could not prepare module mappings for shared module, skipping", zap.Uint64("nodeID", n.ID), zap.Uint64("moduleID", sm.ID))
				continue
			}

			mappingValues, err := w.syncService.PrepareModuleMappings(ctx, mappings)

			if err != nil || mappingValues == nil {
				w.logger.Info("could not prepare module mappings for shared module, skipping", zap.Uint64("nodeID", n.ID), zap.Uint64("moduleID", sm.ID))
				continue
			}

			// get the last sync per-node
			lastSync, _ := w.syncService.GetLastSyncTime(ctx, n.ID, types.NodeSyncTypeData)
			basePath := fmt.Sprintf("/federation/nodes/%d/modules/%d/records/", n.SharedNodeID, sm.ExternalFederationModuleID)

			z := []zap.Field{zap.Uint64("nodeID", n.ID)}

			if lastSync != nil {
				z = append(z, zap.Time("lastSync", *lastSync))
			}

			w.logger.Info("starting data sync", z...)

			url := types.SyncerURI{
				BaseURL:  n.BaseURL,
				Path:     basePath,
				Limit:    w.limit,
				LastSync: lastSync,
			}

			processer := &dataProcesser{
				ID:                  sm.ExternalFederationModuleID,
				NodeID:              n.SharedNodeID,
				ComposeModuleID:     mappings.ComposeModuleID,
				ComposeNamespaceID:  mappings.ComposeNamespaceID,
				ModuleMappings:      &mappings.FieldMapping,
				ModuleMappingValues: &mappingValues,
				NodeBaseURL:         n.BaseURL,
				SyncService:         w.syncService,
			}

			go w.queueUrl(&url, urls, processer)
		}
	}

}

func (w *syncWorkerData) Watch(ctx context.Context, delay time.Duration, limit int) {
	var (
		urls     = make(chan Url, 100)
		payloads = make(chan Payload, 10)

		countProcess = 0
	)

	w.delay = delay
	w.limit = limit

	// todo - get auth from the node
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ticker := time.NewTicker(delay)

	w.PrepareForNodes(ctx, urls)

	for {
		select {
		case <-ctx.Done():
			// do any cleanup here
			w.logger.Info("stopping sync", zap.Int("processed", countProcess))
			return
		case <-ticker.C:
			// do the whole process again
			w.PrepareForNodes(ctx, urls)
		case url := <-urls:
			select {
			case <-ctx.Done():
				// do any cleanup here
				w.logger.Info("stopping sync", zap.Int("processed", countProcess))
				return
			default:
			}

			s, err := url.Url.String()

			if err != nil {
				continue
			}

			responseBody, err := w.syncService.FetchUrl(ctx, s)

			if err != nil {
				continue
			}

			spayload := Payload{
				Payload: responseBody,
				Meta:    url.Meta,
			}

			payloads <- spayload
		case p := <-payloads:

			body, err := ioutil.ReadAll(p.Payload)

			// handle error
			if err != nil {
				continue
			}

			basePath := fmt.Sprintf("/federation/nodes/%d/modules/%d/records/", p.Meta.(*dataProcesser).NodeID, p.Meta.(*dataProcesser).ID)

			u := types.SyncerURI{
				BaseURL: p.Meta.(*dataProcesser).NodeBaseURL,
				Path:    basePath,
				Limit:   w.limit,
			}

			processed, err := w.syncService.ProcessPayload(ctx, body, urls, u, p.Meta.(*dataProcesser))
			countProcess += processed

			if err != nil {
				// handle error
				w.logger.Info("error on handling payload", zap.Error(err))
			} else {
				// TODO
				n, err := DefaultNode.FindBySharedNodeID(ctx, p.Meta.(*dataProcesser).NodeID)

				if err != nil {
					w.logger.Info("could not update sync status", zap.Error(err))
					continue
				}

				// add to db - nodes_sync
				new := &types.NodeSync{
					NodeID:       (*n).ID,
					SyncStatus:   types.NodeSyncStatusSuccess,
					SyncType:     types.NodeSyncTypeData,
					TimeOfAction: time.Now().UTC(),
				}

				// todo
				new, err = DefaultNodeSync.Create(ctx, new)

				if err != nil {
					w.logger.Info("could not update sync status", zap.Error(err))
				}

				w.logger.Info("processed objects", zap.Int("processed", processed), zap.Uint64("nodeID", n.ID))
			}
		}
	}
}

// Process gets the payload from syncer and
// uses the decode package to decode the whole set, depending on
// the filtering that was used (limit)
func (dp *dataProcesser) Process(ctx context.Context, payload []byte) (int, error) {

	processed := 0
	o, err := decoder.DecodeFederationRecordSync([]byte(payload))

	if err != nil {
		return processed, err
	}

	if len(o) == 0 {
		return processed, nil
	}

	for _, er := range o {
		dp.SyncService.mapper.Merge(&er.Values, dp.ModuleMappingValues, dp.ModuleMappings)

		rec := &ct.Record{
			ModuleID:    dp.ComposeModuleID,
			NamespaceID: dp.ComposeNamespaceID,
			Values:      *dp.ModuleMappingValues,
		}

		AddFederationLabel(rec, dp.NodeBaseURL)

		_, err := dp.SyncService.CreateRecord(ctx, rec)

		if err != nil {
			continue
		}

		processed++
	}

	return processed, nil
}
