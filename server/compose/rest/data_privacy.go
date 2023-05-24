package rest

import (
	"context"

	"github.com/cortezaproject/corteza/server/compose/rest/request"
	"github.com/cortezaproject/corteza/server/compose/service"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/payload"
)

type (
	sensitiveRecordsSetPayload struct {
		Set []*sensitiveDataPayload `json:"set"`
	}

	sensitiveDataPayload struct {
		NamespaceID uint64 `json:"namespaceID,string"`
		Namespace   string `json:"namespace"`
		ModuleID    uint64 `json:"moduleID,string"`
		Module      string `json:"module"`

		Records []sensitiveRecords `json:"records"`
	}

	sensitiveRecords struct {
		RecordID uint64           `json:"recordID,string"`
		Values   []map[string]any `json:"values"`
	}

	privacyModuleSetPayload struct {
		Filter types.PrivacyModuleFilter `json:"filter"`
		Set    []*types.PrivacyModule    `json:"set"`
	}

	privateDataFinder interface {
		SearchSensitive(ctx context.Context) (set []types.SensitiveRecordSet, err error)
	}

	DataPrivacy struct {
		record    privateDataFinder
		module    service.ModuleService
		namespace service.NamespaceService
		privacy   service.DataPrivacyService
	}
)

func (DataPrivacy) New() *DataPrivacy {
	return &DataPrivacy{
		record:    service.DefaultRecord,
		module:    service.DefaultModule,
		namespace: service.DefaultNamespace,
		privacy:   service.DefaultDataPrivacy,
	}
}

func (ctrl *DataPrivacy) RecordList(ctx context.Context, r *request.DataPrivacyRecordList) (out interface{}, err error) {
	// If we're requesting only specific connections, prepare filter params here
	reqConns := make(map[uint64]bool)
	hasReqConns := len(r.ConnectionID) > 0
	for _, connectionID := range payload.ParseUint64s(r.ConnectionID) {
		reqConns[connectionID] = true
	}

	// Collect sensitive records
	ss, err := ctrl.record.SearchSensitive(ctx)
	if err != nil {
		return
	}

	outSet := sensitiveRecordsSetPayload{
		Set: make([]*sensitiveDataPayload, 0, 10),
	}

	for _, s := range ss {
		// Skip the ones we don't want
		if hasReqConns && !reqConns[s.Module.Config.DAL.ConnectionID] {
			continue
		}

		// Build the payload
		payload := &sensitiveDataPayload{
			NamespaceID: s.Namespace.ID,
			Namespace:   s.Namespace.Name,

			ModuleID: s.Module.ID,
			Module:   s.Module.Name,

			Records: make([]sensitiveRecords, 0, len(s.Records)),
		}
		for _, a := range s.Records {
			if len(a.Values) == 0 {
				continue
			}
			payload.Records = append(payload.Records, sensitiveRecords{
				RecordID: a.RecordID,
				Values:   a.Values,
			})
		}

		outSet.Set = append(outSet.Set, payload)
	}

	return outSet, nil
}

func (ctrl *DataPrivacy) ModuleList(ctx context.Context, r *request.DataPrivacyModuleList) (out interface{}, err error) {
	var (
		f = types.PrivacyModuleFilter{
			ConnectionID: r.ConnectionID,
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	set, f, err := ctrl.privacy.FindModules(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, f, err)
}

func (ctrl DataPrivacy) makeFilterPayload(_ context.Context, mm types.PrivacyModuleSet, f types.PrivacyModuleFilter, err error) (*privacyModuleSetPayload, error) {
	if err != nil {
		return nil, err
	}

	if len(mm) == 0 {
		mm = make([]*types.PrivacyModule, 0)
	}

	return &privacyModuleSetPayload{Filter: f, Set: mm}, nil
}
