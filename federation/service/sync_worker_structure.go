package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/decoder"
	"github.com/davecgh/go-spew/spew"
	"go.uber.org/zap"
)

type (
	syncWorkerStructure struct {
		syncService *Sync
		logger      *zap.Logger
		delay       time.Duration
		limit       int
	}

	structureProcesser struct {
		NodeID       uint64
		SharedNodeID uint64
		NodeBaseURL  string
		SyncService  *Sync
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
	// s, _ := url.String()

	// w.logger.Debug(fmt.Sprintf("adding %s to queue", s))

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
		// get the last sync per-node
		lastSync, _ := w.syncService.GetLastSyncTime(ctx, n.ID, types.NodeSyncTypeStructure)
		basePath := fmt.Sprintf("/federation/nodes/%d/modules/exposed/", n.SharedNodeID)

		z := []zap.Field{zap.Uint64("nodeID", n.ID), zap.String("host", n.BaseURL)}

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
			NodeID:       n.ID,
			SharedNodeID: n.SharedNodeID,
			NodeBaseURL:  n.BaseURL,
			SyncService:  w.syncService,
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
				spew.Dump("ERR", err)
				continue
			}

			responseBody, err := w.syncService.FetchUrl(ctx, s)

			if err != nil {
				spew.Dump("ERR", err)
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
				spew.Dump("ERR", err)
				continue
			}

			basePath := fmt.Sprintf("/federation/nodes/%d/modules/exposed/", p.Meta.(*structureProcesser).SharedNodeID)

			u := types.SyncerURI{
				BaseURL: p.Meta.(*structureProcesser).NodeBaseURL,
				Path:    basePath,
				Limit:   w.limit,
			}

			processed, err := w.syncService.ProcessPayload(ctx, body, urls, u, p.Meta.(*structureProcesser))
			countProcess += processed

			if err != nil {
				w.logger.Info("error on handling payload", zap.Error(err))
			} else {
				n, err := DefaultNode.FindBySharedNodeID(ctx, p.Meta.(*structureProcesser).SharedNodeID)

				if err != nil {
					w.logger.Info("could not update sync status", zap.Error(err))
					continue
				}

				// add to db - nodes_sync
				new := &types.NodeSync{
					NodeID:       (*n).ID,
					SyncStatus:   types.NodeSyncStatusSuccess,
					SyncType:     types.NodeSyncTypeStructure,
					TimeOfAction: time.Now().UTC(),
				}

				new, err = DefaultNodeSync.Create(ctx, new)

				if err != nil {
					w.logger.Info("could not update sync status", zap.Error(err))
				}

				w.logger.Info("processed objects", zap.Int("processed", processed), zap.Uint64("nodeID", n.ID), zap.String("host", n.BaseURL))
			}
		}
	}
}

// Process gets the payload from syncer and
// uses the decode package to decode the whole set, depending on
// the filtering that was used (limit)
func (dp *structureProcesser) Process(ctx context.Context, payload []byte) (int, error) {
	processed := 0

	o, err := decoder.DecodeFederationModuleSync([]byte(payload))

	if err != nil {
		return processed, err
	}

	if len(o) == 0 {
		return processed, nil
	}

	for _, em := range o {
		new := &types.SharedModule{
			NodeID:                     dp.NodeID,
			ExternalFederationModuleID: em.ID,
			Fields:                     em.Fields,
			Handle:                     em.Handle,
			Name:                       em.Name,
		}

		existing, err := dp.SyncService.LookupSharedModule(ctx, new)

		if err != nil {
			continue
		}

		if existing == nil {
			_, err := dp.SyncService.CreateSharedModule(ctx, new)

			if err != nil {
				continue
			}

			processed++
			continue
		}

		canUpdate, err := dp.SyncService.CanUpdateSharedModule(ctx, new, existing)

		if err != nil {
			continue
		}

		if canUpdate {
			_, err := dp.SyncService.UpdateSharedModule(ctx, existing)

			if err != nil {
				continue
			}

			processed++
			continue
		}
	}

	return processed, nil
}
