package service

import (
	"context"
	"io"
	"time"

	cs "github.com/cortezaproject/corteza-server/compose/service"
	ct "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/federation/types"
)

type (
	Sync struct {
		syncer               *Syncer
		mapper               *Mapper
		sharedModuleService  SharedModuleService
		composeRecordService cs.RecordService
	}
)

func NewSync(s *Syncer, m *Mapper, sm SharedModuleService, cs cs.RecordService) *Sync {
	return &Sync{
		syncer:               s,
		mapper:               m,
		sharedModuleService:  sm,
		composeRecordService: cs,
	}
}

func (s *Sync) ProcessPayload(ctx context.Context, payload []byte, out chan Surl, url types.SyncerURI, processer Processer) error {
	return s.syncer.Process(ctx, payload, out, url, processer)
}

func (s *Sync) QueueUrl(url Surl, out chan Surl) {
	s.syncer.Queue(url, out)
}

func (s *Sync) FetchUrl(ctx context.Context, url string) (io.Reader, error) {
	return s.syncer.Fetch(ctx, url)
}

func (s *Sync) CreateRecord(ctx context.Context, rec *ct.Record) (*ct.Record, error) {
	return s.composeRecordService.With(ctx).Create(rec)
}

func (s *Sync) CreateSharedModule(ctx context.Context, new *types.SharedModule) (*types.SharedModule, error) {
	return s.sharedModuleService.Create(ctx, new)
}

func (s *Sync) GetPairedNodes(ctx context.Context) (types.NodeSet, error) {
	set, _, err := DefaultNode.Search(ctx, types.NodeFilter{Status: types.NodeStatusPaired})

	if err != nil {
		return nil, err
	}

	return set, nil
}

func (s *Sync) GetSharedModules(ctx context.Context, nodeID uint64) (types.SharedModuleSet, error) {
	set, _, err := DefaultSharedModule.Find(ctx, types.SharedModuleFilter{NodeID: nodeID})

	if err != nil {
		return nil, err
	}

	set, _ = set.Filter(func(sm *types.SharedModule) (bool, error) {
		return len(sm.Fields) > 0, nil
	})

	return set, nil
}

func (s *Sync) GetModuleMappings(ctx context.Context, moduleID uint64) (out *types.ModuleMapping, err error) {
	out, err = DefaultModuleMapping.FindByID(ctx, moduleID)
	return
}

func (s *Sync) PrepareModuleMappings(ctx context.Context, mappings *types.ModuleMapping) (ct.RecordValueSet, error) {
	return s.mapper.Prepare((*mappings).FieldMapping), nil
}

func (s *Sync) GetLastSyncTime(ctx context.Context, nodeID uint64, syncType string) (*time.Time, error) {
	ns, err := DefaultNodeSync.LookupLastSuccessfulSync(ctx, nodeID, syncType)

	if err != nil || ns == nil {
		return nil, err
	}

	return &ns.TimeOfAction, nil
}
