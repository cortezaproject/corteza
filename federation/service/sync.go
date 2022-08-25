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
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/filter"
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
	var (
		fstr, f2str []byte
		err         error
	)

	if fstr, err = json.Marshal(new.Fields); err != nil {
		return false, err
	}

	if f2str, err = json.Marshal(existing.Fields); err != nil {
		return false, err
	}

	return string(fstr) == string(f2str), nil
}

// ProcessPayload passes the payload to the syncer lib
func (s *Sync) ProcessPayload(ctx context.Context, payload []byte, out chan Url, url types.SyncerURI, processer Processer) (ProcesserResponse, error) {
	ctx = auth.SetIdentityToContext(ctx, auth.FederationUser())

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
func (s *Sync) CreateRecord(ctx context.Context, rec *ct.Record) (out *ct.Record, err error) {
	out, _, err = s.composeRecordService.Create(ctx, rec)
	return
}

// UpdateRecord wraps the compose Record service Update
func (s *Sync) UpdateRecord(ctx context.Context, rec *ct.Record) (out *ct.Record, err error) {
	out, _, err = s.composeRecordService.Update(ctx, rec)
	return
}

// DeleteRecord wraps the compose Record service Update
func (s *Sync) DeleteRecord(ctx context.Context, rec *ct.Record) error {
	return s.composeRecordService.DeleteByID(ctx, rec.NamespaceID, rec.ModuleID, rec.ID)
}

// FindRecord find the record via federation label
func (s *Sync) FindRecords(ctx context.Context, filter ct.RecordFilter) (set ct.RecordSet, err error) {
	set, _, err = s.composeRecordService.Find(ctx, filter)
	return
}

// LookupSharedModule find the shared module if exists
func (s *Sync) LookupSharedModule(ctx context.Context, new *types.SharedModule) (*types.SharedModule, error) {
	var sm *types.SharedModule

	list, _, err := s.sharedModuleService.Find(ctx, types.SharedModuleFilter{
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

// UpdateSharedModule wraps the federation SharedModule service Update
func (s *Sync) UpdateSharedModule(ctx context.Context, updated *types.SharedModule) (*types.SharedModule, error) {
	return s.sharedModuleService.Update(ctx, updated)
}

// CreateSharedModule wraps the federation SharedModule service Create
func (s *Sync) CreateSharedModule(ctx context.Context, new *types.SharedModule) (*types.SharedModule, error) {
	return s.sharedModuleService.Create(ctx, new)
}

// GetPairedNodes finds successfuly paired nodes
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

func (s *Sync) GetLastStructureSyncStatus(ctx context.Context, nodeID, externalFederationModuleID uint64) (syncStatus string, err error) {
	var list types.NodeSyncSet

	list, _, err = DefaultNodeSync.Search(ctx, types.NodeSyncFilter{
		NodeID:   nodeID,
		ModuleID: externalFederationModuleID,
		SyncType: types.NodeSyncTypeStructure,
		Sorting: filter.Sorting{
			Sort: filter.SortExprSet{
				&filter.SortExpr{Column: "time_action", Descending: true},
			},
		},
		Paging: filter.Paging{Limit: 1},
	})

	if err != nil {
		syncStatus = types.NodeSyncStatusSuccess
		return
	}

	if len(list) == 0 {
		syncStatus = types.NodeSyncStatusSuccess
		return
	}

	syncStatus = list[0].SyncStatus

	return
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
	if u, err = s.systemUserService.FindByHandle(ctx, fmt.Sprintf("federation_%d", nodeID)); err != nil {
		return nil, err
	}

	// attach the roles
	rr, _, err := s.systemRoleService.Find(ctx, st.RoleFilter{MemberID: u.ID})

	if err != nil {
		return nil, err
	}

	u.SetRoles(rr.IDs()...)

	return u, nil
}
