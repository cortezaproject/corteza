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

	if r.ComposeModuleID == 0 {
		return nil, errors.New("TODO - http 400 bad request - use compose module id in request")
	}

	if r.ComposeNamespaceID == 0 {
		return nil, errors.New("TODO - http 400 bad request - use compose namespace id in request")
	}

	return (service.ExposedModule()).Create(ctx, mod)
}

func (ctrl ManageStructure) ReadExposed(ctx context.Context, r *request.ManageStructureReadExposed) (interface{}, error) {
	return (service.ExposedModule()).FindByID(ctx, r.GetNodeID(), r.GetModuleID())
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
	return (service.ExposedModule()).Update(ctx, em)
}

func (ctrl ManageStructure) RemoveExposed(ctx context.Context, r *request.ManageStructureRemoveExposed) (interface{}, error) {
	return (service.ExposedModule()).DeleteByID(ctx, r.NodeID, r.ModuleID)
}

func (ctrl ManageStructure) ReadShared(ctx context.Context, r *request.ManageStructureReadShared) (interface{}, error) {
	return (service.SharedModule()).FindByID(ctx, r.GetNodeID(), r.GetModuleID())
}

func (ctrl ManageStructure) CreateMappings(ctx context.Context, r *request.ManageStructureCreateMappings) (interface{}, error) {
	mm := &types.ModuleMapping{
		FederationModuleID: r.ModuleID,
		ComposeModuleID:    r.ComposeModuleID,
		ComposeNamespaceID: r.ComposeNamespaceID,
		FieldMapping:       r.Fields,
	}

	return (service.ModuleMapping()).Create(ctx, mm)
}

func (ctrl ManageStructure) ReadMappings(ctx context.Context, r *request.ManageStructureReadMappings) (interface{}, error) {
	return (service.ModuleMapping()).FindByID(ctx, r.ModuleID)
}

func (ctrl ManageStructure) ListAll(ctx context.Context, r *request.ManageStructureListAll) (interface{}, error) {
	var (
		list interface{}
		err  error
	)

	switch true {
	case r.Exposed:
		list, _, err = (service.ExposedModule()).Find(ctx, types.ExposedModuleFilter{
			NodeID: r.NodeID,
		})
		break
	case r.Shared:
		list, _, err = (service.SharedModule()).Find(ctx, types.SharedModuleFilter{
			NodeID: r.NodeID,
		})
		break
	default:
		return nil, errors.New("TODO - http 400 bad request - either use ?exposed or ?shared")
	}

	return list, err
}
