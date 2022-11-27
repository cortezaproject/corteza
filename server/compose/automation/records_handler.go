package automation

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza/server/compose/types"
	. "github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/wfexec"
)

type (
	recordService interface {
		FindByID(ctx context.Context, namespaceID, moduleID, recordID uint64) (*types.Record, *types.RecordValueErrorSet, error)
		Find(ctx context.Context, filter types.RecordFilter) (set types.RecordSet, f types.RecordFilter, err error)

		Create(ctx context.Context, record *types.Record) (*types.Record, *types.RecordValueErrorSet, error)
		Update(ctx context.Context, record *types.Record) (*types.Record, *types.RecordValueErrorSet, error)
		Bulk(ctx context.Context, oo ...*types.RecordBulkOperation) (types.RecordSet, *types.RecordValueErrorSet, error)
		Report(ctx context.Context, namespaceID, moduleID uint64, metrics, dimensions, filter string) (out any, err error)

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
		// Item buffer, current item pointer, and total items traversed
		ptr    uint
		buffer types.RecordSet
		total  uint

		// When filter limit is set, this constraints it
		iterLimit    uint
		useIterLimit bool

		// Item loader for additional chunks
		filter types.RecordFilter
		loader func() error
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
			Meta:    args.Meta,
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
		if args.PageCursor != "" {
			f.PageCursor = &filter.PagingCursor{}
			if err = f.PageCursor.UnmarshalJSON([]byte(args.PageCursor)); err != nil {
				return
			}
		}
	}

	f.IncTotal = args.IncTotal
	f.IncPageNavigation = args.IncPageNavigation

	if args.hasMeta {
		f.Meta = args.Meta
	}

	if args.hasLimit {
		f.Limit = uint(args.Limit)
	}

	var auxf types.RecordFilter
	results.Records, auxf, err = h.rec.Find(ctx, f)
	results.Total = uint64(auxf.Total)
	if auxf.NextPage != nil {
		results.NextPage = auxf.NextPage.Encode()
	}
	if auxf.PrevPage != nil {
		results.PrevPage = auxf.PrevPage.Encode()
	}

	// Always assure at least empty cursor when requesting page nav.
	if auxf.PageNavigation == nil && auxf.IncPageNavigation {
		results.PageNavigation = append(results.PageNavigation, nil)
	}

	for _, pn := range auxf.PageNavigation {
		if pn == nil {
			continue
		}
		// first page is null -- no cursor
		if pn.Cursor == nil {
			results.PageNavigation = append(results.PageNavigation, nil)
		} else {
			results.PageNavigation = append(results.PageNavigation, Must(NewString(pn.Cursor.Encode())))
		}
	}
	return
}

func (h recordsHandler) first(ctx context.Context, args *recordsFirstArgs) (results *recordsFirstResults, err error) {
	r, err := h.fetchEdge(ctx, args, true)
	if err != nil {
		return nil, err
	}

	return &recordsFirstResults{
		Record: r,
	}, nil
}

func (h recordsHandler) last(ctx context.Context, args *recordsLastArgs) (results *recordsLastResults, err error) {
	r, err := h.fetchEdge(ctx, args, false)
	if err != nil {
		return nil, err
	}

	return &recordsLastResults{
		Record: r,
	}, nil
}

func (h recordsHandler) each(ctx context.Context, args *recordsEachArgs) (out wfexec.IteratorHandler, err error) {
	var (
		i = &recordSetIterator{}
		f = types.RecordFilter{
			Query:   args.Query,
			Meta:    args.Meta,
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
		if args.PageCursor != "" {
			f.NextPage = &filter.PagingCursor{}
			if err = f.NextPage.UnmarshalJSON([]byte(args.PageCursor)); err != nil {
				return
			}
		}
	}

	f.IncTotal = args.IncTotal
	f.IncPageNavigation = args.IncPageNavigation

	if args.hasMeta {
		f.Meta = args.Meta
	}

	if args.hasLimit {
		i.useIterLimit = true
		i.iterLimit = uint(args.Limit)

		f.Limit = uint(args.Limit)

		if args.Limit > uint64(wfexec.MaxIteratorBufferSize) {
			f.Limit = wfexec.MaxIteratorBufferSize
		}
		i.iterLimit = uint(args.Limit)
	} else {
		f.Limit = wfexec.MaxIteratorBufferSize
	}

	i.filter = f
	i.loader = func() (err error) {
		// Edgecase
		if i.filter.PageCursor != nil && i.filter.NextPage == nil {
			return
		}

		i.total += i.ptr
		i.ptr = 0

		i.filter.PageCursor = i.filter.NextPage
		i.filter.NextPage = nil
		i.buffer, i.filter, err = h.rec.Find(ctx, i.filter)

		return
	}

	// Initial load
	return i, i.loader()
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
	results.Record, _, err = h.rec.Create(ctx, args.Record)
	return
}

func (h recordsHandler) update(ctx context.Context, args *recordsUpdateArgs) (results *recordsUpdateResults, err error) {
	results = &recordsUpdateResults{}
	results.Record, _, err = h.rec.Update(ctx, args.Record)
	return
}

func (h recordsHandler) delete(ctx context.Context, args *recordsDeleteArgs) error {
	if rec, err := h.lookupRecord(ctx, args); err != nil {
		return err
	} else {
		return h.rec.DeleteByID(ctx, rec.NamespaceID, rec.ModuleID, rec.ID)
	}
}

func (h recordsHandler) report(ctx context.Context, args *recordsReportArgs) (*recordsReportResults, error) {
	r := &recordsReportResults{}

	ns, mod, err := h.loadCombo(ctx, args)
	if err != nil {
		return nil, err
	}

	r.Report, err = h.rec.Report(ctx, ns.ID, mod.ID, args.Metrics, args.Dimensons, args.Filter)
	if err != nil {
		return nil, err
	}

	return r, nil
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

	record, _, err = h.rec.FindByID(ctx, namespace.ID, module.ID, recordID)
	return
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

func (h recordsHandler) fetchEdge(ctx context.Context, args interface{}, first bool) (*types.Record, error) {
	f := types.RecordFilter{}

	if first {
		f.Sort.Set("createdAt DESC")
	} else {
		f.Sort.Set("createdAt ASC")
	}

	f.Limit = 1

	if ns, mod, err := h.loadCombo(ctx, args); err != nil {
		return nil, err
	} else {
		f.ModuleID = mod.ID
		f.NamespaceID = ns.ID
	}

	rr, _, err := h.rec.Find(ctx, f)
	if err != nil {
		return nil, err
	}
	if len(rr) == 0 {
		return nil, fmt.Errorf("could not fetch records: no records found")
	}
	return rr[0], nil
}

func (i *recordSetIterator) More(context.Context, *Vars) (bool, error) {
	return wfexec.GenericResourceNextCheck(i.useIterLimit, i.ptr, uint(len(i.buffer)), i.total, i.iterLimit, i.filter.NextPage != nil), nil
}

func (i *recordSetIterator) Start(context.Context, *Vars) error { i.ptr = 0; return nil }

func (i *recordSetIterator) Next(context.Context, *Vars) (out *Vars, err error) {
	if len(i.buffer)-int(i.ptr) <= 0 {
		if err = i.loader(); err != nil {
			panic(err)
		}
	}

	out = &Vars{}
	out.Set("record", Must(NewComposeRecord(i.buffer[i.ptr])))
	out.Set("index", Must(NewInteger(i.total+i.ptr)))
	out.Set("total", Must(NewInteger(i.filter.Total)))

	i.ptr++
	return out, nil
}
