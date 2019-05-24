package rest

import (
	"context"

	"github.com/titpetric/factory/resputil"

	"github.com/cortezaproject/corteza-server/compose/internal/service"
	"github.com/cortezaproject/corteza-server/compose/rest/request"
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	modulePayload struct {
		*types.Module

		Fields []*moduleFieldPayload `json:"fields"`

		CanGrant        bool `json:"canGrant"`
		CanUpdateModule bool `json:"canUpdateModule"`
		CanDeleteModule bool `json:"canDeleteModule"`
		CanCreateRecord bool `json:"canCreateRecord"`
		CanReadRecord   bool `json:"canReadRecord"`
		CanUpdateRecord bool `json:"canUpdateRecord"`
		CanDeleteRecord bool `json:"canDeleteRecord"`
	}

	moduleFieldPayload struct {
		*types.ModuleField

		CanReadRecordValue   bool `json:"canReadRecordValue"`
		CanUpdateRecordValue bool `json:"canUpdateRecordValue"`
	}

	moduleSetPayload struct {
		Filter types.ModuleFilter `json:"filter"`
		Set    []*modulePayload   `json:"set"`
	}

	Module struct {
		module service.ModuleService
		record service.RecordService
		ac     moduleAccessController
	}

	moduleAccessController interface {
		CanGrant(context.Context) bool

		CanUpdateModule(context.Context, *types.Module) bool
		CanDeleteModule(context.Context, *types.Module) bool
		CanCreateRecord(context.Context, *types.Module) bool
		CanReadRecord(context.Context, *types.Module) bool
		CanUpdateRecord(context.Context, *types.Module) bool
		CanDeleteRecord(context.Context, *types.Module) bool

		CanReadRecordValue(context.Context, *types.ModuleField) bool
		CanUpdateRecordValue(context.Context, *types.ModuleField) bool
	}
)

func (Module) New() *Module {
	return &Module{
		module: service.DefaultModule,
		record: service.DefaultRecord,
		ac:     service.DefaultAccessControl,
	}
}

func (ctrl *Module) List(ctx context.Context, r *request.ModuleList) (interface{}, error) {
	f := types.ModuleFilter{
		NamespaceID: r.NamespaceID,
		Query:       r.Query,
		PerPage:     r.PerPage,
		Page:        r.Page,
	}

	set, filter, err := ctrl.module.With(ctx).Find(f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *Module) Read(ctx context.Context, r *request.ModuleRead) (interface{}, error) {
	mod, err := ctrl.module.With(ctx).FindByID(r.NamespaceID, r.ModuleID)
	return ctrl.makePayload(ctx, mod, err)
}

func (ctrl *Module) Create(ctx context.Context, r *request.ModuleCreate) (interface{}, error) {
	var (
		err error
		mod = &types.Module{
			NamespaceID: r.NamespaceID,
			Name:        r.Name,
			Fields:      r.Fields,
			Meta:        r.Meta,
		}
	)

	mod, err = ctrl.module.With(ctx).Create(mod)
	return ctrl.makePayload(ctx, mod, err)
}

func (ctrl *Module) Update(ctx context.Context, r *request.ModuleUpdate) (interface{}, error) {
	var (
		err error
		mod = &types.Module{
			ID:          r.ModuleID,
			NamespaceID: r.NamespaceID,
			Name:        r.Name,
			Fields:      r.Fields,
			Meta:        r.Meta,
			UpdatedAt:   r.UpdatedAt,
		}
	)

	mod, err = ctrl.module.With(ctx).Update(mod)
	return ctrl.makePayload(ctx, mod, err)
}

func (ctrl *Module) Delete(ctx context.Context, r *request.ModuleDelete) (interface{}, error) {
	_, err := ctrl.module.With(ctx).FindByID(r.NamespaceID, r.ModuleID)
	if err != nil {
		return nil, err
	}

	return resputil.OK(), ctrl.module.With(ctx).DeleteByID(r.NamespaceID, r.ModuleID)
}

func (ctrl Module) makePayload(ctx context.Context, m *types.Module, err error) (*modulePayload, error) {
	if err != nil || m == nil {
		return nil, err
	}

	mfp, err := ctrl.makeFieldsPayload(ctx, m)
	if err != nil {
		return nil, err
	}

	return &modulePayload{
		Module: m,

		Fields: mfp,

		CanGrant: ctrl.ac.CanGrant(ctx),

		CanUpdateModule: ctrl.ac.CanUpdateModule(ctx, m),
		CanDeleteModule: ctrl.ac.CanDeleteModule(ctx, m),
		CanCreateRecord: ctrl.ac.CanCreateRecord(ctx, m),
		CanReadRecord:   ctrl.ac.CanReadRecord(ctx, m),
		CanUpdateRecord: ctrl.ac.CanUpdateRecord(ctx, m),
		CanDeleteRecord: ctrl.ac.CanDeleteRecord(ctx, m),
	}, nil
}

func (ctrl Module) makeFieldsPayload(ctx context.Context, m *types.Module) (out []*moduleFieldPayload, err error) {
	out = make([]*moduleFieldPayload, len(m.Fields))

	for i, f := range m.Fields {
		out[i] = &moduleFieldPayload{
			ModuleField: f,

			CanReadRecordValue:   ctrl.ac.CanReadRecordValue(ctx, f),
			CanUpdateRecordValue: ctrl.ac.CanUpdateRecordValue(ctx, f),
		}
	}

	return
}

func (ctrl Module) makeFilterPayload(ctx context.Context, nn types.ModuleSet, f types.ModuleFilter, err error) (*moduleSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &moduleSetPayload{Filter: f, Set: make([]*modulePayload, len(nn))}

	for i := range nn {
		msp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return msp, nil
}
