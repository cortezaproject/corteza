package rest

import (
	"context"

	"github.com/cortezaproject/corteza-server/compose/rest/request"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/service/event"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/corredor"
	"github.com/cortezaproject/corteza-server/pkg/filter"
)

type (
	moduleSetPayload struct {
		Filter types.ModuleFilter `json:"filter"`
		Set    []*modulePayload   `json:"set"`
	}

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

	Module struct {
		module    service.ModuleService
		namespace service.NamespaceService
		ac        moduleAccessController
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
		module:    service.DefaultModule,
		namespace: service.DefaultNamespace,
		ac:        service.DefaultAccessControl,
	}
}

func (ctrl *Module) List(ctx context.Context, r *request.ModuleList) (interface{}, error) {
	var (
		err error
		f   = types.ModuleFilter{
			NamespaceID: r.NamespaceID,
			Query:       r.Query,
			Name:        r.Name,
			Handle:      r.Handle,
			Labels:      r.Labels,
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	set, filter, err := ctrl.module.Find(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *Module) Read(ctx context.Context, r *request.ModuleRead) (interface{}, error) {
	mod, err := ctrl.module.FindByID(ctx, r.NamespaceID, r.ModuleID)
	return ctrl.makePayload(ctx, mod, err)
}

func (ctrl *Module) Create(ctx context.Context, r *request.ModuleCreate) (interface{}, error) {
	var (
		err error
		mod = &types.Module{
			NamespaceID: r.NamespaceID,
			Name:        r.Name,
			Handle:      r.Handle,
			Fields:      r.Fields,
			Meta:        r.Meta,
			Labels:      r.Labels,
		}
	)

	mod, err = ctrl.module.Create(ctx, mod)
	return ctrl.makePayload(ctx, mod, err)
}

func (ctrl *Module) Update(ctx context.Context, r *request.ModuleUpdate) (interface{}, error) {
	var (
		err error
		mod = &types.Module{
			ID:          r.ModuleID,
			NamespaceID: r.NamespaceID,
			Name:        r.Name,
			Handle:      r.Handle,
			Fields:      r.Fields,
			Meta:        r.Meta,
			Labels:      r.Labels,
			UpdatedAt:   r.UpdatedAt,
		}
	)

	mod, err = ctrl.module.Update(ctx, mod)
	return ctrl.makePayload(ctx, mod, err)
}

func (ctrl *Module) Delete(ctx context.Context, r *request.ModuleDelete) (interface{}, error) {
	_, err := ctrl.module.FindByID(ctx, r.NamespaceID, r.ModuleID)
	if err != nil {
		return nil, err
	}

	return api.OK(), ctrl.module.DeleteByID(ctx, r.NamespaceID, r.ModuleID)
}

func (ctrl *Module) TriggerScript(ctx context.Context, r *request.ModuleTriggerScript) (rsp interface{}, err error) {
	var (
		module    *types.Module
		namespace *types.Namespace
	)

	if module, err = ctrl.module.FindByID(ctx, r.NamespaceID, r.ModuleID); err != nil {
		return
	}

	if namespace, err = ctrl.namespace.FindByID(ctx, r.NamespaceID); err != nil {
		return
	}

	// @todo implement same behaviour as we have on record - module+oldModule
	err = corredor.Service().Exec(ctx, r.Script, event.ModuleOnManual(module, module, namespace))
	return ctrl.makePayload(ctx, module, err)
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
