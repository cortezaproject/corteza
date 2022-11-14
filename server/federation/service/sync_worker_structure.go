package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"go.uber.org/zap"
)

type (
	syncWorkerStructure struct {
		syncService *Sync
		logger      *zap.Logger
		delay       time.Duration
		limit       int
	}
)

func WorkerStructure(sync *Sync, logger *zap.Logger) *syncWorkerStructure {
	return &syncWorkerStructure{
		syncService: sync,
		logger:      logger,
		limit:       defaultPage,
		delay:       defaultDelay,
	}
}

func (w *syncWorkerStructure) queueUrl(url *types.SyncerURI, urls chan Url, meta Processer) {
	t := Url{
		Url:  *url,
		Meta: meta,
	}

	w.syncService.QueueUrl(t, urls)
}

func (w *syncWorkerStructure) PrepareForNodes(ctx context.Context, urls chan Url) {
	nodes, err := w.syncService.GetPairedNodes(ctx)

	if err != nil {
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

		// get the last sync per-node
		lastSync, _ := w.syncService.GetLastSyncTime(ctx, n.ID, types.NodeSyncTypeStructure)
		basePath := fmt.Sprintf("/nodes/%d/modules/exposed/", n.SharedNodeID)

		z := []zap.Field{
			zap.Uint64("nodeID", n.ID),
			zap.String("host", n.BaseURL),
			zap.String("delay", w.delay.String()),
			zap.Int("pagesize", w.limit),
		}

		if lastSync != nil {
			z = append(z, zap.Time("lastSync", *lastSync))
		}

		w.logger.Info("starting structure sync", z...)

		url := types.SyncerURI{
			BaseURL:  n.BaseURL,
			Path:     basePath,
			Limit:    w.limit,
			LastSync: lastSync,
		}

		processer := &structureProcesser{
			SyncService: w.syncService,
			User:        u,
			Node:        n,
		}

		go w.queueUrl(&url, urls, processer)
	}
}

func (w *syncWorkerStructure) Watch(ctx context.Context, delay time.Duration, limit int) {
	var (
		urls     = make(chan Url, 100)
		payloads = make(chan Payload, 10)

		countProcess = 0
	)

	w.delay = delay
	w.limit = limit

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ctx = auth.SetIdentityToContext(ctx, auth.FederationUser())

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
			meta := url.Meta.(*structureProcesser)

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
			meta := p.Meta.(*structureProcesser)
			body, err := ioutil.ReadAll(p.Payload)

			// handle error
			if err != nil {
				w.logger.Error("could not read data from synced url, skipping",
					zap.Error(err),
					zap.Uint64("nodeID", meta.Node.ID),
					zap.String("host", meta.Node.BaseURL))

				continue
			}

			basePath := fmt.Sprintf("/nodes/%d/modules/exposed/", meta.Node.SharedNodeID)

			u := types.SyncerURI{
				BaseURL: p.Meta.(*structureProcesser).Node.BaseURL,
				Path:    basePath,
				Limit:   w.limit,
			}

			processed, errProcess := w.syncService.ProcessPayload(ctx, body, urls, u, meta)
			countProcess += processed.(structureProcesserResponse).Processed

			// error raised before the actual persist process
			// ignore
			if moduleID := processed.(structureProcesserResponse).ModuleID; moduleID == 0 {
				w.logger.Info("no structure change from last change, skipping",
					zap.Uint64("nodeID", meta.Node.ID),
					zap.String("host", meta.Node.BaseURL))

				continue
			}

			syncStatus := types.NodeSyncStatusSuccess
			if errProcess != nil {
				syncStatus = types.NodeSyncStatusError
			}

			new := &types.NodeSync{
				NodeID:       meta.Node.ID,
				ModuleID:     processed.(structureProcesserResponse).ModuleID,
				SyncStatus:   syncStatus,
				SyncType:     types.NodeSyncTypeStructure,
				TimeOfAction: time.Now().UTC(),
			}

			new, err = DefaultNodeSync.Create(ctx, new)

			if err != nil {
				w.logger.Info("could not update sync status", zap.Error(err))
			}

			if errProcess != nil {
				w.logger.Info("error on handling payload", zap.Error(errProcess))
			} else {
				w.logger.Info("processed objects",
					zap.Int("processed", processed.(structureProcesserResponse).Processed),
					zap.Uint64("nodeID", meta.Node.ID),
					zap.String("host", meta.Node.BaseURL))
			}

		}
	}
}
