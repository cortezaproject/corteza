package rest

import (
	"context"
	"strconv"

	"github.com/crusttech/human/server/pkg/api"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

type Dml struct{}

func (Dml) New() *Dml { return &Dml{} }

func (ctrl Dml) ConnectionList(ctx context.Context, r *request.DmlConnectionList) (interface{}, error) {
	f := types.DmlConnectionFilter{Handle: r.Handle}
	for _, id := range r.ConnectionID {
		if v, err := strconv.ParseUint(id, 10, 64); err == nil {
			f.ConnectionID = append(f.ConnectionID, v)
		}
	}
	set, err := service.DefaultDmlConnection.Find(ctx, f)
	if err != nil {
		return nil, err
	}
	return &struct {
		Set []*types.DmlConnection `json:"set"`
	}{set}, nil
}

func (ctrl Dml) ConnectionCreate(ctx context.Context, r *request.DmlConnectionCreate) (interface{}, error) {
	return service.DefaultDmlConnection.Create(ctx, types.DmlConnectionInput{
		Handle: r.Handle,
		Label:  r.Label,
		Params: &r.Params,
	})
}

func (ctrl Dml) ConnectionRead(ctx context.Context, r *request.DmlConnectionRead) (interface{}, error) {
	return service.DefaultDmlConnection.FindByID(ctx, r.ConnectionID)
}

func (ctrl Dml) ConnectionModels(ctx context.Context, r *request.DmlConnectionModels) (interface{}, error) {
	set, err := service.DefaultDmlConnection.FindModels(ctx, r.ConnectionID, types.DmlModelFilter{Ident: r.Ident})
	if err != nil {
		return nil, err
	}
	return &struct {
		Set []*types.DmlModel `json:"set"`
	}{set}, nil
}

func (ctrl Dml) MappingCreate(ctx context.Context, r *request.DmlMappingCreate) (interface{}, error) {
	return service.DefaultDmlMapping.Create(ctx, &types.DmlMapping{
		ConnectionID:    r.ConnectionID,
		NamespaceHandle: r.NamespaceHandle,
		SourceIdent:     r.SourceIdent,
		ModuleHandle:    r.ModuleHandle,
		ModuleName:      r.ModuleName,
		Skip:            r.Skip,
		Identifier:      r.Identifier,
		Columns:         r.Columns,
	})
}

func (ctrl Dml) MappingRead(ctx context.Context, r *request.DmlMappingRead) (interface{}, error) {
	return service.DefaultDmlMapping.FindByID(ctx, r.MappingID)
}

func (ctrl Dml) MappingUpdate(ctx context.Context, r *request.DmlMappingUpdate) (interface{}, error) {
	return service.DefaultDmlMapping.Update(ctx, &types.DmlMapping{
		ID:              r.MappingID,
		ConnectionID:    r.ConnectionID,
		NamespaceHandle: r.NamespaceHandle,
		SourceIdent:     r.SourceIdent,
		ModuleHandle:    r.ModuleHandle,
		ModuleName:      r.ModuleName,
		Skip:            r.Skip,
		Identifier:      r.Identifier,
		Columns:         r.Columns,
	})
}

func (ctrl Dml) MappingDelete(ctx context.Context, r *request.DmlMappingDelete) (interface{}, error) {
	return api.OK(), service.DefaultDmlMapping.DeleteByID(ctx, r.MappingID)
}

func (ctrl Dml) ImportRun(ctx context.Context, r *request.DmlImportRun) (interface{}, error) {
	return service.DefaultDmlImporter.RunImport(ctx, r.MappingID, r.Method)
}

func (ctrl Dml) ImportRunRead(ctx context.Context, r *request.DmlImportRunRead) (interface{}, error) {
	return service.DefaultDmlImporter.GetRun(ctx, r.RunID)
}
