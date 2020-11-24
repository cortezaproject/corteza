package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/cortezaproject/corteza-server/federation/types"
	"go.uber.org/zap"
)

type (
	syncWorkerData struct {
		syncService *Sync
		logger      *zap.Logger
		delay       time.Duration
		limit       int
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
		// get the user, associated for this node
		u, err := w.syncService.LoadUserWithRoles(ctx, n.ID)

		if err != nil {
			w.logger.Info("could not preload federation user, skipping",
				zap.Uint64("nodeID", n.ID),
				zap.Error(err))

			continue
		}

		set, err := w.syncService.GetSharedModules(ctx, n.ID)

		if err != nil {
			w.logger.Info("could not get shared modules, skipping",
				zap.Uint64("nodeID", n.ID),
				zap.Error(err))

			continue
		}

		if len(set) == 0 {
			w.logger.Info("there are no valid shared modules, skipping",
				zap.Uint64("nodeID", n.ID),
				zap.Error(err))

			continue
		}

		// go through set and prepare module mappings for it
		for _, sm := range set {
			z := []zap.Field{
				zap.Uint64("nodeID", n.ID),
				zap.Uint64("moduleID", sm.ID),
				zap.String("host", n.BaseURL),
				zap.String("delay", w.delay.String()),
				zap.Int("pagesize", w.limit),
			}

			// if the last sync was error'd, log and skip
			lastSyncStatus, err := w.syncService.GetLastStructureSyncStatus(ctx, n.ID, sm.ExternalFederationModuleID)

			if err != nil {
				w.logger.Info("could not get last sync status, skipping", z...)
				continue
			}

			if lastSyncStatus == types.NodeSyncStatusError {
				w.logger.Info("last structure sync was not complete, admin resolve needed, skipping", z...)
				continue
			}

			mappings, _ := w.syncService.GetModuleMappings(ctx, sm.ID)

			if mappings == nil {
				w.logger.Info("could not prepare module mappings for shared module, skipping", z...)
				continue
			}

			mappingValues, err := w.syncService.PrepareModuleMappings(ctx, mappings)

			if err != nil || mappingValues == nil {
				w.logger.Info("could not prepare module mappings for shared module, skipping", z...)
				continue
			}

			// get the last sync per-node
			lastSync, _ := w.syncService.GetLastSyncTime(ctx, n.ID, types.NodeSyncTypeData)
			basePath := fmt.Sprintf("/nodes/%d/modules/%d/records/", n.SharedNodeID, sm.ExternalFederationModuleID)

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
				ComposeModuleID:     mappings.ComposeModuleID,
				ComposeNamespaceID:  mappings.ComposeNamespaceID,
				ModuleMappings:      &mappings.FieldMapping,
				ModuleMappingValues: &mappingValues,
				SyncService:         w.syncService,
				User:                u,
				Node:                n,
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

			s, _ := url.Url.String()
			meta := url.Meta.(*dataProcesser)

			// use the authToken from node pairing
			ctx = context.WithValue(ctx, FederationUserToken, meta.Node.AuthToken)

			responseBody, err := w.syncService.FetchUrl(ctx, s)

			if err != nil {
				w.logger.Error("could not fetch data from url, skipping",
					zap.Error(err),
					zap.String("url", s),
					zap.Uint64("nodeID", meta.Node.ID),
					zap.String("host", meta.Node.BaseURL))

				continue
			}

			spayload := Payload{
				Payload: responseBody,
				Meta:    url.Meta,
			}

			payloads <- spayload
		case p := <-payloads:
			meta := p.Meta.(*dataProcesser)
			body, err := ioutil.ReadAll(p.Payload)

			// handle error
			if err != nil {
				w.logger.Error("could not read data from synced url, skipping",
					zap.Error(err),
					zap.Uint64("nodeID", meta.Node.ID),
					zap.String("host", meta.Node.BaseURL))

				continue
			}

			basePath := fmt.Sprintf("/nodes/%d/modules/%d/records/", meta.Node.SharedNodeID, meta.ID)

			u := types.SyncerURI{
				BaseURL: meta.NodeBaseURL,
				Path:    basePath,
				Limit:   w.limit,
			}

			processed, errProcess := w.syncService.ProcessPayload(ctx, body, urls, u, meta)
			countProcess += processed.(dataProcesserResponse).Processed

			// error raised before the actual persist process
			// ignore
			syncStatus := types.NodeSyncStatusSuccess
			if errProcess != nil {
				syncStatus = types.NodeSyncStatusError
			}

			new := &types.NodeSync{
				NodeID:       meta.Node.ID,
				ModuleID:     meta.ID,
				SyncStatus:   syncStatus,
				SyncType:     types.NodeSyncTypeStructure,
				TimeOfAction: time.Now().UTC(),
			}

			new, err = DefaultNodeSync.Create(ctx, new)

			if err != nil {
				w.logger.Info("could not update sync status", zap.Error(err))
			}

			if errProcess != nil {
				w.logger.Info("error on persisting structure", zap.Error(errProcess))
			} else {
				w.logger.Info("processed objects",
					zap.Int("processed", processed.(dataProcesserResponse).Processed),
					zap.Uint64("nodeID", meta.Node.ID))
			}
		}
	}
}
