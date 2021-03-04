package rest

import (
	"context"

	"github.com/cortezaproject/corteza-server/federation/rest/request"
	"github.com/cortezaproject/corteza-server/federation/service"
	"github.com/cortezaproject/corteza-server/federation/types"
)

type (
	ManageStructure struct{}

	exposedModulePayload struct {
		*types.ExposedModule

		CanManageModule bool `json:"canManageModule"`
	}

	sharedModulePayload struct {
		*types.SharedModule

		CanMapModule bool `json:"canMapModule"`
	}

	moduleMappingPayload struct {
		*types.ModuleMapping
	}
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
			Name:               r.Name,
			Handle:             r.Handle,
			Fields:             r.Fields,
		}
	)

	em, err := (service.DefaultExposedModule).Create(ctx, mod)
	return ctrl.makePayload(ctx, em, err)
}

func (ctrl ManageStructure) ReadExposed(ctx context.Context, r *request.ManageStructureReadExposed) (interface{}, error) {
	em, err := (service.DefaultExposedModule).FindByID(ctx, r.GetNodeID(), r.GetModuleID())
	return ctrl.makePayload(ctx, em, err)
}

func (ctrl ManageStructure) UpdateExposed(ctx context.Context, r *request.ManageStructureUpdateExposed) (interface{}, error) {
	var (
		em = &types.ExposedModule{
			ID:                 r.ModuleID,
			NodeID:             r.NodeID,
			ComposeModuleID:    r.ComposeModuleID,
			ComposeNamespaceID: r.ComposeNamespaceID,
			Name:               r.Name,
			Handle:             r.Handle,
			Fields:             r.Fields,
		}
	)

	em, err := (service.DefaultExposedModule).Update(ctx, em)
	return ctrl.makePayload(ctx, em, err)
}

func (ctrl ManageStructure) RemoveExposed(ctx context.Context, r *request.ManageStructureRemoveExposed) (interface{}, error) {
	em, err := (service.DefaultExposedModule).DeleteByID(ctx, r.NodeID, r.ModuleID)
	return ctrl.makePayload(ctx, em, err)
}

func (ctrl ManageStructure) ReadShared(ctx context.Context, r *request.ManageStructureReadShared) (interface{}, error) {
	list, err := (service.DefaultSharedModule).FindByID(ctx, r.GetNodeID(), r.GetModuleID())
	return ctrl.makePayload(ctx, list, err)
}

func (ctrl ManageStructure) CreateMappings(ctx context.Context, r *request.ManageStructureCreateMappings) (interface{}, error) {
	mm := &types.ModuleMapping{
		NodeID:             r.NodeID,
		FederationModuleID: r.ModuleID,
		ComposeModuleID:    r.ComposeModuleID,
		ComposeNamespaceID: r.ComposeNamespaceID,
		FieldMapping:       r.Fields,
	}

	// check if it exists, do an upsert
	existing, _, err := service.DefaultModuleMapping.Find(ctx, types.ModuleMappingFilter{
		NodeID:             r.NodeID,
		ComposeModuleID:    r.ComposeModuleID,
		ComposeNamespaceID: r.ComposeNamespaceID,
		FederationModuleID: r.ModuleID,
	})

	if err != nil {
		return nil, err
	}

	if len(existing) > 0 {
		// do an update
		return service.DefaultModuleMapping.Update(ctx, mm)
	}

	return service.DefaultModuleMapping.Create(ctx, mm)
}

// ReadMappings outputs an object with the module mappings per each federation module
func (ctrl ManageStructure) ReadMappings(ctx context.Context, r *request.ManageStructureReadMappings) (interface{}, error) {
	f := types.ModuleMappingFilter{
		NodeID:             r.NodeID,
		FederationModuleID: r.ModuleID,
	}

	if r.ComposeModuleID > 0 {
		f.ComposeModuleID = r.ComposeModuleID
	}

	set, _, err := service.DefaultModuleMapping.Find(ctx, f)

	if err != nil {
		return nil, err
	}

	if len(set) != 1 {
		return nil, service.ModuleMappingErrNotFound()
	}

	return set[0], nil
}

// ListAll show the list of exposed / shared / mapped modules
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
	case r.Mapped:
		list, _, err = (service.DefaultModuleMapping).Find(ctx, types.ModuleMappingFilter{
			NodeID: r.NodeID,
		})
		break
	default:
		return nil, service.ExposedModuleErrRequestParametersInvalid()
	}

	return ctrl.makePayload(ctx, list, err)
}

func (ctrl ManageStructure) makeSharedModulePayload(ctx context.Context, sm *types.SharedModule, err error) (*sharedModulePayload, error) {
	if err != nil || sm == nil {
		return nil, err
	}

	return &sharedModulePayload{
		SharedModule: sm,

		CanMapModule: service.DefaultAccessControl.CanMapModule(ctx, sm),
	}, nil
}

func (ctrl ManageStructure) makeExposedModulePayload(ctx context.Context, em *types.ExposedModule, err error) (*exposedModulePayload, error) {
	if err != nil || em == nil {
		return nil, err
	}

	return &exposedModulePayload{
		ExposedModule: em,

		CanManageModule: service.DefaultAccessControl.CanManageModule(ctx, em),
	}, nil
}

func (ctrl ManageStructure) makeModuleMappingPayload(ctx context.Context, em *types.ModuleMapping, err error) (*moduleMappingPayload, error) {
	if err != nil || em == nil {
		return nil, err
	}

	return &moduleMappingPayload{
		ModuleMapping: em,
	}, nil
}

func (ctrl ManageStructure) makePayload(ctx context.Context, payload interface{}, err error) (interface{}, error) {
	if err != nil || payload == nil {
		return nil, err
	}

	switch payload.(type) {
	case types.SharedModuleSet:
		set := make([]sharedModulePayload, len(payload.(types.SharedModuleSet)))

		for i, d := range payload.(types.SharedModuleSet) {
			mp, _ := ctrl.makeSharedModulePayload(ctx, d, err)
			set[i] = *mp
		}

		return set, err
	case types.ExposedModuleSet:
		set := make([]exposedModulePayload, len(payload.(types.ExposedModuleSet)))

		for i, d := range payload.(types.ExposedModuleSet) {
			mp, _ := ctrl.makeExposedModulePayload(ctx, d, err)
			set[i] = *mp
		}

		return set, err
	case types.ModuleMappingSet:
		set := make([]moduleMappingPayload, len(payload.(types.ModuleMappingSet)))

		for i, d := range payload.(types.ModuleMappingSet) {
			mp, _ := ctrl.makeModuleMappingPayload(ctx, d, err)
			set[i] = *mp
		}

		return set, err
	case *types.ExposedModule:
		return ctrl.makeExposedModulePayload(ctx, payload.(*types.ExposedModule), err)
	case *types.SharedModule:
		return ctrl.makeSharedModulePayload(ctx, payload.(*types.SharedModule), err)
	default:
		return nil, nil
	}
}
