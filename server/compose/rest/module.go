package rest

import (
	"context"

	"github.com/crusttech/human/server/compose/rest/request"
	"github.com/crusttech/human/server/compose/service"
	"github.com/crusttech/human/server/compose/service/event"
	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/api"
	"github.com/crusttech/human/server/pkg/corredor"
	"github.com/crusttech/human/server/pkg/filter"
)

type (
	moduleSetPayload struct {
		Filter types.ModuleFilter `json:"filter"`
		Set    []*modulePayload   `json:"set"`
	}

	modulePayload struct {
		*types.Module

		Fields []*moduleFieldPayload `json:"fields"`

		CanGrant             bool `json:"canGrant"`
		CanUpdateModule      bool `json:"canUpdateModule"`
		CanDeleteModule      bool `json:"canDeleteModule"`
		CanCreateRecord      bool `json:"canCreateRecord"`
		CanCreateOwnedRecord bool `json:"canCreateOwnedRecord"`
	}

	moduleFieldPayload struct {
		*types.ModuleField

		CanReadRecordValue   bool `json:"canReadRecordValue"`
		CanUpdateRecordValue bool `json:"canUpdateRecordValue"`
	}

	Module struct {
		module    service.ModuleService
		locale    service.ResourceTranslationsManagerService
		namespace service.NamespaceService
		ac        moduleAccessController
	}

	moduleAccessController interface {
		CanGrant(context.Context) bool

		CanUpdateModule(context.Context, *types.Module) bool
		CanDeleteModule(context.Context, *types.Module) bool
		CanCreateRecordOnModule(context.Context, *types.Module) bool
		CanCreateOwnedRecordOnModule(context.Context, *types.Module) bool
		CanReadRecord(context.Context, *types.Record) bool

		CanReadRecordValueOnModuleField(context.Context, *types.ModuleField) bool
		CanUpdateRecordValueOnModuleField(context.Context, *types.ModuleField) bool
	}
)

func (Module) New() *Module {
	return &Module{
		module:    service.DefaultModule,
		namespace: service.DefaultNamespace,
		ac:        service.DefaultAccessControl,
		locale:    service.DefaultResourceTranslation,
	}
}

// makeFilter builds the search filter for the generated List controller.
func (ctrl Module) makeFilter(ctx context.Context, r *request.ModuleList) (types.ModuleFilter, error) {
	var (
		err error
		f   = types.ModuleFilter{
			NamespaceID: r.NamespaceID,
			ProjectID:   r.ProjectID,
			Query:       r.Query,
			Name:        r.Name,
			Handle:      r.Handle,
			Labels:      r.Labels,
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return f, err
	}

	f.IncTotal = r.IncTotal

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return f, err
	}

	return f, nil
}

// beforeCreate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Create controller
// already mapped the plain-value params and the compound (namespace) id.
func (ctrl Module) beforeCreate(ctx context.Context, res *types.Module, r *request.ModuleCreate) error {
	res.Config = r.Config
	res.Meta = r.Meta
	res.Fields = r.Fields

	return nil
}

// beforeUpdate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Update controller
// already mapped the plain-value params, the compound (namespace) id and
// the resource ID/UpdatedAt.
func (ctrl Module) beforeUpdate(ctx context.Context, res *types.Module, r *request.ModuleUpdate) error {
	res.Config = r.Config
	res.Meta = r.Meta
	res.Fields = r.Fields

	return nil
}

func (ctrl *Module) ListTranslations(ctx context.Context, r *request.ModuleListTranslations) (interface{}, error) {
	return ctrl.locale.Module(ctx, r.NamespaceID, r.ModuleID)
}

func (ctrl *Module) UpdateTranslations(ctx context.Context, r *request.ModuleUpdateTranslations) (interface{}, error) {
	return api.OK(), ctrl.locale.Upsert(ctx, r.Translations)
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
	err = corredor.Service().Exec(ctx, r.Script, corredor.ExtendScriptArgs(event.ModuleOnManual(module, module, namespace), r.Args))
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

		CanCreateRecord:      ctrl.ac.CanCreateRecordOnModule(ctx, m),
		CanCreateOwnedRecord: ctrl.ac.CanCreateOwnedRecordOnModule(ctx, m),
	}, nil
}

func (ctrl Module) makeFieldsPayload(ctx context.Context, m *types.Module) (out []*moduleFieldPayload, err error) {
	out = make([]*moduleFieldPayload, len(m.Fields))

	for i, f := range m.Fields {
		out[i] = &moduleFieldPayload{
			ModuleField: f,

			CanReadRecordValue:   ctrl.ac.CanReadRecordValueOnModuleField(ctx, f),
			CanUpdateRecordValue: ctrl.ac.CanUpdateRecordValueOnModuleField(ctx, f),
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
