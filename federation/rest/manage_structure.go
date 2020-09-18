package rest

import (
	"context"
	"errors"

	"github.com/cortezaproject/corteza-server/federation/rest/request"
	"github.com/cortezaproject/corteza-server/federation/service"
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/davecgh/go-spew/spew"
)

type (
	ManageStructure struct{}
)

func (ManageStructure) New() *ManageStructure {
	return &ManageStructure{}
}

func (ctrl ManageStructure) RemoveExposed(ctx context.Context, r *request.ManageStructureRemoveExposed) (interface{}, error) {
	return nil, (service.ExposedModule()).DeleteByID(ctx, r.NodeID, r.ModuleID)
}

func (ctrl ManageStructure) ReadExposed(ctx context.Context, r *request.ManageStructureReadExposed) (interface{}, error) {
	return (service.ExposedModule()).FindByID(context.Background(), r.GetNodeID(), r.GetModuleID())
}

func (ctrl ManageStructure) CreateExposed(ctx context.Context, r *request.ManageStructureCreateExposed) (interface{}, error) {
	spew.Dump("RECEIVED", r)
	// create new type, add to Create()
	// return (service.ExposedModule()).Create(context.Background(), r.GetNodeID(), r.GetModuleID())
	var (
		err error
		mod = &types.ExposedModule{
			NodeID:          r.NodeID,
			ComposeModuleID: r.ComposeModuleID,
			Fields:          r.Fields,
		}
	)

	spew.Dump("MOD", mod)

	mod, err = (service.ExposedModule()).Create(context.Background(), mod)
	spew.Dump(mod, err)
	// return ctrl.makePayload(ctx, mod, err)
	return nil, nil
}

func (ctrl ManageStructure) ReadShared(ctx context.Context, r *request.ManageStructureReadShared) (interface{}, error) {
	return (service.SharedModule()).FindByID(context.Background(), r.GetNodeID(), r.GetModuleID())
}

func (ctrl ManageStructure) ListAll(ctx context.Context, r *request.ManageStructureListAll) (interface{}, error) {
	var (
		list interface{}
		err  error
	)

	switch true {
	case r.Exposed:
		list, _, err = (service.ExposedModule()).Find(context.Background(), types.ExposedModuleFilter{
			NodeID: r.NodeID,
		})
		break
	case r.Shared:
		list, _, err = (service.SharedModule()).Find(context.Background(), types.SharedModuleFilter{
			NodeID: r.NodeID,
		})
		break
	default:
		return nil, errors.New("TODO - http 400 bad request - either use ?exposed or ?shared")
	}

	return list, err
}
