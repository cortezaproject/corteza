package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	cs "github.com/cortezaproject/corteza-server/compose/service"
	ct "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/system/service"
	ss "github.com/cortezaproject/corteza-server/system/service"
	st "github.com/cortezaproject/corteza-server/system/types"
)

type (
	Sync struct {
		syncer               *Syncer
		mapper               *Mapper
		sharedModuleService  SharedModuleService
		composeRecordService cs.RecordService
		systemUserService    ss.UserService
		systemRoleService    ss.RoleService
	}
)

func NewSync(s *Syncer, m *Mapper, sm SharedModuleService, cs cs.RecordService, us ss.UserService, rs ss.RoleService) *Sync {
	return &Sync{
		syncer:               s,
		mapper:               m,
		sharedModuleService:  sm,
		composeRecordService: cs,
		systemUserService:    us,
		systemRoleService:    rs,
	}
}

// CanUpdateSharedModule checks the origin and destination module
// compatibility of fields
// It currently checks if all of the fields are exactly the same
// TODO - check if any of the newly missing fields are actually being used so a safe update
// is possible
func (s *Sync) CanUpdateSharedModule(ctx context.Context, new *types.SharedModule, existing *types.SharedModule) (bool, error) {

	// check for mapped fields
	fstr, err := json.Marshal(new.Fields)
	f2str, err := json.Marshal(existing.Fields)

	if err != nil {
		return false, err
	}

	if string(fstr) == string(f2str) {
		return true, nil
	}

	return false, nil
}

// ProcessPayload passes the payload to the syncer lib
func (s *Sync) ProcessPayload(ctx context.Context, payload []byte, out chan Url, url types.SyncerURI, processer Processer) (int, error) {
	return s.syncer.Process(ctx, payload, out, url, processer)
}

// QueueUrl passes the url to the syncer
func (s *Sync) QueueUrl(url Url, out chan Url) {
	s.syncer.Queue(url, out)
}

// FetchUrl passes the url to be fetched to the syncer
func (s *Sync) FetchUrl(ctx context.Context, url string) (io.Reader, error) {
	return s.syncer.Fetch(ctx, url)
}

// CreateRecord wraps the compose Record service Create
func (s *Sync) CreateRecord(ctx context.Context, rec *ct.Record) (*ct.Record, error) {
	return s.composeRecordService.With(ctx).Create(rec)
}

// LookupSharedModule find the shared module if exists
func (s *Sync) LookupSharedModule(ctx context.Context, new *types.SharedModule) (*types.SharedModule, error) {
	var sm *types.SharedModule

	list, _, err := service.DefaultStore.SearchFederationSharedModules(ctx, types.SharedModuleFilter{
		NodeID:                     new.NodeID,
		ExternalFederationModuleID: new.ExternalFederationModuleID})

	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	sm = list[0]

	return sm, nil
}

func (s *Sync) UpdateSharedModule(ctx context.Context, updated *types.SharedModule) (*types.SharedModule, error) {
	return s.sharedModuleService.Update(ctx, updated)
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

// LoadUserWithRoles gets the federation user, that was
// created at node pairing process
func (s *Sync) LoadUserWithRoles(ctx context.Context, nodeID uint64) (*st.User, error) {
	var (
		u   *st.User
		err error
	)

	// get the federated user, associated for this node
	if u, err = s.systemUserService.With(ctx).FindByHandle(fmt.Sprintf("federation_%d", nodeID)); err != nil {
		return nil, err
	}

	// attach the roles
	rr, _, err := s.systemRoleService.With(ctx).Find(st.RoleFilter{MemberID: u.ID})

	if err != nil {
		return nil, err
	}

	u.SetRoles(rr.IDs())

	return u, nil
}
