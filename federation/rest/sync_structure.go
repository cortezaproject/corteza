package rest

import (
	"context"

	"github.com/cortezaproject/corteza-server/federation/rest/request"
	"github.com/cortezaproject/corteza-server/federation/service"
)

type (
	SyncStructure struct{}
)

func (SyncStructure) New() *SyncStructure {
	return &SyncStructure{}
}

func (ctrl SyncStructure) Remove(ctx context.Context, r *request.SyncStructureRemove) (interface{}, error) {
	return nil, (service.ExposedModule()).DeleteByID(ctx, r.NodeID, r.ModuleID)
}

func (ctrl SyncStructure) ReadExposed(ctx context.Context, r *request.SyncStructureReadExposed) (interface{}, error) {
	return (service.ExposedModule()).FindByID(context.Background(), r.GetNodeID(), r.GetModuleID())
}
