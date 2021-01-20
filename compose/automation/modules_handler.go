package automation

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	recordService interface {
		FindByID(ctx context.Context, namespaceID, moduleID, recordID uint64) (*types.Record, error)
		Find(ctx context.Context, filter types.RecordFilter) (set types.RecordSet, f types.RecordFilter, err error)

		Create(ctx context.Context, record *types.Record) (*types.Record, error)
		Update(ctx context.Context, record *types.Record) (*types.Record, error)
		Bulk(ctx context.Context, oo ...*types.RecordBulkOperation) (types.RecordSet, error)

		DeleteByID(ctx context.Context, namespaceID, moduleID uint64, recordID ...uint64) error
	}

	moduleService interface {
		FindByID(ctx context.Context, namespaceID, moduleID uint64) (*types.Module, error)
		FindByHandle(ctx context.Context, namespaceID uint64, handle string) (*types.Module, error)
	}

	namespaceService interface {
		FindByID(ctx context.Context, namespaceID uint64) (*types.Namespace, error)
		FindByHandle(ctx context.Context, handle string) (*types.Namespace, error)
	}

	recordsHandlers struct {
		reg recordsHandlerRegistry
		ns  namespaceService
		mod moduleService
		rec recordService
	}
)

func RecordsHandlers(reg recordsHandlerRegistry, ns namespaceService, mod moduleService, rec recordService) *recordsHandlers {
	h := &recordsHandlers{
		reg: reg,
		ns:  ns,
		mod: mod,
		rec: rec,
	}

	h.register()
	return h
}

func resolveNamespace(ctx context.Context, svc namespaceService, id *uint64, handle string, res *types.Namespace) (err error) {
	if *id == 0 {
		if len(handle) > 0 {
			if res, err = svc.FindByHandle(ctx, handle); err != nil {
				return
			}
		}

		if res != nil {
			*id = res.ID
		}
	}

	return
}

func resolveModule(ctx context.Context, svc moduleService, namespaceID uint64, id *uint64, handle string, res *types.Module) (err error) {
	if *id == 0 {
		if len(handle) > 0 {
			if res, err = svc.FindByHandle(ctx, namespaceID, handle); err != nil {
				return
			}
		}

		if res != nil {
			*id = res.ID
		}
	}

	return
}

func (h recordsHandlers) lookupByID(ctx context.Context, args *recordsLookupByIDArgs) (results *recordsLookupByIDResults, err error) {
	results = &recordsLookupByIDResults{}

	if err = resolveNamespace(ctx, h.ns, &args.namespaceID, args.namespaceHandle, args.namespaceRes); err != nil {
		return nil, err
	}

	if err = resolveModule(ctx, h.mod, args.namespaceID, &args.moduleID, args.moduleHandle, args.moduleRes); err != nil {
		return nil, err
	}

	results.Record, err = h.rec.FindByID(ctx, args.namespaceID, args.moduleID, args.RecordID)
	return
}

func (h recordsHandlers) create(ctx context.Context, args *recordsCreateArgs) (results *recordsCreateResults, err error) {
	results = &recordsCreateResults{}

	if err = resolveNamespace(ctx, h.ns, &args.namespaceID, args.namespaceHandle, args.namespaceRes); err != nil {
		return nil, err
	}

	if err = resolveModule(ctx, h.mod, args.namespaceID, &args.moduleID, args.moduleHandle, args.moduleRes); err != nil {
		return nil, err
	}

	rec := &types.Record{
		ModuleID:    args.moduleID,
		NamespaceID: args.namespaceID,
		Values:      args.Values,
		Labels:      args.Labels,
		OwnedBy:     args.OwnedBy,
	}

	results.Record, err = h.rec.Create(ctx, rec)
	return
}

func (h recordsHandlers) save(ctx context.Context, args *recordsSaveArgs) (results *recordsSaveResults, err error) {
	results = &recordsSaveResults{}
	results.Record, err = h.rec.Update(ctx, args.Record)
	return
}

func (h recordsHandlers) update(ctx context.Context, args *recordsUpdateArgs) (results *recordsUpdateResults, err error) {
	results = &recordsUpdateResults{}
	if err = resolveNamespace(ctx, h.ns, &args.namespaceID, args.namespaceHandle, args.namespaceRes); err != nil {
		return nil, err
	}

	if err = resolveModule(ctx, h.mod, args.namespaceID, &args.moduleID, args.moduleHandle, args.moduleRes); err != nil {
		return nil, err
	}

	rec := &types.Record{
		ModuleID:    args.moduleID,
		NamespaceID: args.namespaceID,
		Values:      args.Values,
		Labels:      args.Labels,
		OwnedBy:     args.OwnedBy,
	}

	results.Record, err = h.rec.Update(ctx, rec)
	return
}

func (h recordsHandlers) delete(ctx context.Context, args *recordsDeleteArgs) (err error) {
	if err = resolveNamespace(ctx, h.ns, &args.namespaceID, args.namespaceHandle, args.namespaceRes); err != nil {
		return err
	}

	if err = resolveModule(ctx, h.mod, args.namespaceID, &args.moduleID, args.moduleHandle, args.moduleRes); err != nil {
		return err
	}

	return h.rec.DeleteByID(ctx, args.namespaceID, args.moduleID, args.RecordID)
}
