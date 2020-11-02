package rest

import (
	"context"
	"errors"

	"github.com/cortezaproject/corteza-server/federation/rest/request"
	"github.com/cortezaproject/corteza-server/federation/service"
	"github.com/cortezaproject/corteza-server/federation/types"
)

type (
	ManageStructure struct{}
)

func (ManageStructure) New() *ManageStructure {
	return &ManageStructure{}
}

func (ctrl ManageStructure) CreateExposed(ctx context.Context, r *request.ManageStructureCreateExposed) (interface{}, error) {
	var (
		mod = &types.ExposedModule{
			NodeID:             r.NodeID,
			ComposeModuleID:    r.ComposeModuleID,
			ComposeNamespaceID: r.ComposeNamespaceID,
			Fields:             r.Fields,
		}
	)
	return (service.DefaultExposedModule).Create(ctx, mod)
}

func (ctrl ManageStructure) ReadExposed(ctx context.Context, r *request.ManageStructureReadExposed) (interface{}, error) {
	return (service.DefaultExposedModule).FindByID(ctx, r.GetNodeID(), r.GetModuleID())
}

func (ctrl ManageStructure) UpdateExposed(ctx context.Context, r *request.ManageStructureUpdateExposed) (interface{}, error) {
	var (
		em = &types.ExposedModule{
			ID:                 r.ModuleID,
			NodeID:             r.NodeID,
			ComposeModuleID:    r.ComposeModuleID,
			ComposeNamespaceID: r.ComposeNamespaceID,
			Fields:             r.Fields,
		}
	)
	return (service.DefaultExposedModule).Update(ctx, em)
}

func (ctrl ManageStructure) RemoveExposed(ctx context.Context, r *request.ManageStructureRemoveExposed) (interface{}, error) {
	return (service.DefaultExposedModule).DeleteByID(ctx, r.NodeID, r.ModuleID)
}

func (ctrl ManageStructure) ReadShared(ctx context.Context, r *request.ManageStructureReadShared) (interface{}, error) {
	return (service.DefaultSharedModule).FindByID(ctx, r.GetNodeID(), r.GetModuleID())
}

func (ctrl ManageStructure) CreateMappings(ctx context.Context, r *request.ManageStructureCreateMappings) (interface{}, error) {
	mm := &types.ModuleMapping{
		FederationModuleID: r.ModuleID,
		ComposeModuleID:    r.ComposeModuleID,
		ComposeNamespaceID: r.ComposeNamespaceID,
		FieldMapping:       r.Fields,
	}

	return (service.DefaultModuleMapping).Create(ctx, mm)
}

func (ctrl ManageStructure) ReadMappings(ctx context.Context, r *request.ManageStructureReadMappings) (interface{}, error) {
	return (service.DefaultModuleMapping).FindByID(ctx, r.ModuleID)
}

func (ctrl ManageStructure) ListAll(ctx context.Context, r *request.ManageStructureListAll) (interface{}, error) {
	var (
		list interface{}
		err  error
	)

	switch true {
	case r.Exposed:
		list, _, err = (service.DefaultExposedModule).Find(ctx, types.ExposedModuleFilter{
			NodeID: r.NodeID,
		})
		break
	case r.Shared:
		list, _, err = (service.DefaultSharedModule).Find(ctx, types.SharedModuleFilter{
			NodeID: r.NodeID,
		})
		break
	default:
		return nil, errors.New("TODO - http 400 bad request - either use ?exposed or ?shared")
	}

	return list, err
}
