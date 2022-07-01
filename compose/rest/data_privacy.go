package rest

import (
	"context"

	"github.com/cortezaproject/corteza-server/compose/rest/request"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/payload"
)

type (
	sensitiveDataSetPayload struct {
		Set []*sensitiveDataPayload `json:"set"`
	}

	sensitiveDataPayload struct {
		NamespaceID uint64 `json:"namespaceID,string"`
		Namespace   string `json:"namespace"`
		ModuleID    uint64 `json:"moduleID,string"`
		Module      string `json:"module"`

		Records []sensitiveData `json:"records"`
	}

	sensitiveData struct {
		RecordID uint64           `json:"recordID,string"`
		Values   []map[string]any `json:"values"`
	}

	privateDataFinder interface {
		FindSensitive(ctx context.Context, filter types.RecordFilter) (set []types.PrivateDataSet, err error)
	}

	DataPrivacy struct {
		record    privateDataFinder
		module    service.ModuleService
		namespace service.NamespaceService
	}
)

func (DataPrivacy) New() *DataPrivacy {
	return &DataPrivacy{
		record:    service.DefaultRecord,
		module:    service.DefaultModule,
		namespace: service.DefaultNamespace,
	}
}

func (ctrl *DataPrivacy) SensitiveDataList(ctx context.Context, r *request.DataPrivacySensitiveDataList) (out interface{}, err error) {
	outSet := sensitiveDataSetPayload{}

	reqConns := make(map[uint64]bool)
	hasReqConns := len(r.ConnectionID) > 0
	for _, connectionID := range payload.ParseUint64s(r.ConnectionID) {
		reqConns[connectionID] = true
	}

	// All namespaces
	namespaces, _, err := ctrl.namespace.Find(ctx, types.NamespaceFilter{})
	if err != nil {
		return
	}

	outSet.Set = make([]*sensitiveDataPayload, 0, 10)

	for _, n := range namespaces {
		// All modules
		modules, _, err := ctrl.module.Find(ctx, types.ModuleFilter{NamespaceID: n.ID})
		if err != nil {
			return nil, err
		}
		for _, m := range modules {
			conn := m.ModelConfig.ConnectionID
			if hasReqConns && !reqConns[conn] {
				continue
			}

			sData, err := ctrl.record.FindSensitive(ctx, types.RecordFilter{ModuleID: m.ID, NamespaceID: m.NamespaceID})
			if err != nil {
				return nil, err
			}
			if len(sData) == 0 {
				continue
			}

			nsMod := &sensitiveDataPayload{
				NamespaceID: n.ID,
				Namespace:   n.Name,

				ModuleID: m.ID,
				Module:   m.Name,

				Records: make([]sensitiveData, 0, len(sData)),
			}
			for _, a := range sData {
				if len(a.Values) == 0 {
					continue
				}
				nsMod.Records = append(nsMod.Records, sensitiveData{
					RecordID: a.ID,
					Values:   a.Values,
				})
			}

			if len(nsMod.Records) == 0 {
				continue
			}

			outSet.Set = append(outSet.Set, nsMod)
		}
	}

	return outSet, nil
}
