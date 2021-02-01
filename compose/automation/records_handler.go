package automation

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/compose/types"
	. "github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
)

type (
	recordService interface {
		FindByID(ctx context.Context, namespaceID, moduleID, recordID uint64) (*types.Record, error)
		Find(ctx context.Context, filter types.RecordFilter) (set types.RecordSet, f types.RecordFilter, err error)

		Create(ctx context.Context, record *types.Record) (*types.Record, error)
		Update(ctx context.Context, record *types.Record) (*types.Record, error)
		Bulk(ctx context.Context, oo ...*types.RecordBulkOperation) (types.RecordSet, error)

		Validate(ctx context.Context, rec *types.Record) error

		DeleteByID(ctx context.Context, namespaceID, moduleID uint64, recordID ...uint64) error
	}

	recordsHandler struct {
		reg recordsHandlerRegistry
		ns  namespaceService
		mod moduleService
		rec recordService
	}

	recordSetIterator struct {
		ptr    int
		set    types.RecordSet
		filter types.RecordFilter
	}

	recordLookup interface {
		GetRecord() (bool, uint64, *types.Record)
	}
)

func RecordsHandler(reg recordsHandlerRegistry, ns namespaceService, mod moduleService, rec recordService) *recordsHandler {
	h := &recordsHandler{
		reg: reg,
		ns:  ns,
		mod: mod,
		rec: rec,
	}

	h.register()
	return h
}

func (h recordsHandler) lookup(ctx context.Context, args *recordsLookupArgs) (results *recordsLookupResults, err error) {
	results = &recordsLookupResults{}
	results.Record, err = h.lookupRecord(ctx, args)
	return
}

func (h recordsHandler) search(ctx context.Context, args *recordsSearchArgs) (results *recordsSearchResults, err error) {
	results = &recordsSearchResults{}

	var (
		f = types.RecordFilter{
			Query:   args.Query,
			Labels:  args.Labels,
			Deleted: filter.State(args.Deleted),
		}
	)

	if ns, mod, err := h.loadCombo(ctx, args); err != nil {
		return nil, err
	} else {
		f.ModuleID = mod.ID
		f.NamespaceID = ns.ID
	}

	if args.hasSort {
		if err = f.Sort.Set(args.Sort); err != nil {
			return
		}
	}

	if args.hasPageCursor {
		if err = f.PageCursor.Decode(args.PageCursor); err != nil {
			return
		}
	}

	f.IncTotal = args.IncTotal
	f.IncPageNavigation = args.IncPageNavigation

	if args.hasLabels {
		f.Labels = args.Labels
	}

	if args.hasLimit {
		f.Limit = uint(args.Limit)
	}

	results.Records, _, err = h.rec.Find(ctx, f)
	return
}

func (h recordsHandler) each(ctx context.Context, args *recordsEachArgs) (out wfexec.IteratorHandler, err error) {
	var (
		i = &recordSetIterator{}
		f = types.RecordFilter{
			Query:   args.Query,
			Labels:  args.Labels,
			Deleted: filter.State(args.Deleted),
		}
	)

	if ns, mod, err := h.loadCombo(ctx, args); err != nil {
		return nil, err
	} else {
		f.ModuleID = mod.ID
		f.NamespaceID = ns.ID
	}

	if args.hasSort {
		if err = f.Sort.Set(args.Sort); err != nil {
			return nil, err
		}
	}

	if args.hasPageCursor {
		if err = f.PageCursor.Decode(args.PageCursor); err != nil {
			return nil, err
		}
	}

	f.IncTotal = args.IncTotal
	f.IncPageNavigation = args.IncPageNavigation

	if args.hasLabels {
		f.Labels = args.Labels
	}

	if args.hasLimit {
		f.Limit = uint(args.Limit)
	}

	i.set, i.filter, err = h.rec.Find(ctx, f)
	return i, err
}

func (h recordsHandler) validate(ctx context.Context, args *recordsValidateArgs) (*recordsValidateResults, error) {
	results := &recordsValidateResults{Valid: true}
	if err := h.rec.Validate(ctx, args.Record); err != nil {
		results.Valid = false

		//if rves, is := err.(*types.RecordValueErrorSet); is {
		//	results.Errors = rves
		//} else {
		//	return nil, err
		//}
	}

	return results, nil
}

func (h recordsHandler) new(ctx context.Context, args *recordsNewArgs) (*recordsNewResults, error) {
	results := &recordsNewResults{}
	namespace, module, err := h.loadCombo(ctx, args)
	if err != nil {
		return nil, err
	}

	results.Record = &types.Record{
		ModuleID:    module.ID,
		NamespaceID: namespace.ID,
	}

	results.Record.SetModule(module)
	return results, nil
}

func (h recordsHandler) create(ctx context.Context, args *recordsCreateArgs) (results *recordsCreateResults, err error) {
	results = &recordsCreateResults{}
	results.Record, err = h.rec.Create(ctx, args.Record)
	return
}

func (h recordsHandler) update(ctx context.Context, args *recordsUpdateArgs) (results *recordsUpdateResults, err error) {
	results = &recordsUpdateResults{}
	results.Record, err = h.rec.Update(ctx, args.Record)
	return
}

func (h recordsHandler) delete(ctx context.Context, args *recordsDeleteArgs) error {
	if rec, err := h.lookupRecord(ctx, args); err != nil {
		return err
	} else {
		return h.rec.DeleteByID(ctx, rec.NamespaceID, rec.ModuleID, rec.ID)
	}
}

func (h recordsHandler) lookupRecord(ctx context.Context, args recordLookup) (record *types.Record, err error) {
	var (
		namespace *types.Namespace
		module    *types.Module
		recordID  uint64
	)

	if _, recordID, record = args.GetRecord(); record != nil {
		return
	}

	if namespace, module, err = h.loadCombo(ctx, args); err != nil {
		return
	}

	return h.rec.FindByID(ctx, namespace.ID, module.ID, recordID)
}

func (h recordsHandler) loadCombo(ctx context.Context, args interface{}) (namespace *types.Namespace, module *types.Module, err error) {
	if lkp, is := args.(namespaceLookup); is {
		if namespace, err = lookupNamespace(ctx, h.ns, lkp); err != nil {
			err = fmt.Errorf("could not load namespace: %w", err)
			return
		}
	} else {
		err = fmt.Errorf("could not extract namespace lookup arguments")
		return
	}

	if lkp, is := args.(moduleLookup); is {
		if module, err = lookupModule(ctx, h.ns, h.mod, lkp); err != nil {
			err = fmt.Errorf("could not load module: %w", err)
			return
		}
	} else {
		err = fmt.Errorf("could not extract module lookup arguments")
		return
	}

	return
}

func (i *recordSetIterator) More(context.Context, *Vars) (bool, error) {
	return i.ptr < len(i.set), nil
}

func (i *recordSetIterator) Start(context.Context, *Vars) error { i.ptr = 0; return nil }

func (i *recordSetIterator) Next(context.Context, *Vars) (*Vars, error) {
	out := RVars{
		"record": Must(NewComposeRecord(i.set[i.ptr])),
		"index":  Must(NewUnsignedInteger(i.ptr)),
		"total":  Must(NewUnsignedInteger(i.filter.Total)),
	}

	i.ptr++
	return out.Vars(), nil
}
