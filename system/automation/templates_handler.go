package automation

import (
	"context"
	"fmt"
	"io"

	. "github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/spf13/cast"
)

type (
	templateService interface {
		FindByID(ctx context.Context, ID uint64) (*types.Template, error)
		FindByHandle(ct context.Context, handle string) (*types.Template, error)
		FindByAny(ctx context.Context, identifier interface{}) (*types.Template, error)
		Search(context.Context, types.TemplateFilter) (types.TemplateSet, types.TemplateFilter, error)

		Create(ctx context.Context, tpl *types.Template) (*types.Template, error)
		Update(ctx context.Context, tpl *types.Template) (*types.Template, error)

		DeleteByID(ctx context.Context, ID uint64) error
		UndeleteByID(ctx context.Context, ID uint64) error

		Render(ctx context.Context, templateID uint64, dstType string, variables map[string]interface{}, options map[string]string) (io.ReadSeeker, error)
	}

	templatesHandler struct {
		reg  templatesHandlerRegistry
		tSvc templateService
	}

	templateSetIterator struct {
		// Item buffer, current item pointer, and total items traversed
		ptr    uint
		buffer types.TemplateSet
		total  uint

		// When filter limit is set, this constraints it
		iterLimit    uint
		useIterLimit bool

		// Item loader for additional chunks
		filter types.TemplateFilter
		loader func() error
	}

	templateLookup interface {
		GetLookup() (bool, uint64, string, *types.Template)
	}
)

func TemplatesHandler(reg templatesHandlerRegistry, tSvc templateService) *templatesHandler {
	h := &templatesHandler{
		reg:  reg,
		tSvc: tSvc,
	}

	h.register()
	return h
}

func (h templatesHandler) lookup(ctx context.Context, args *templatesLookupArgs) (results *templatesLookupResults, err error) {
	results = &templatesLookupResults{}
	results.Template, err = lookupTemplate(ctx, h.tSvc, args)
	return
}

func (h templatesHandler) search(ctx context.Context, args *templatesSearchArgs) (results *templatesSearchResults, err error) {
	results = &templatesSearchResults{}

	var (
		f = types.TemplateFilter{
			Handle:  args.Handle,
			Type:    args.Type,
			OwnerID: args.OwnerID,
			Partial: args.Partial,
			Labels:  args.Labels,
		}
	)

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

	if args.hasLabels {
		f.Labels = args.Labels
	}

	if args.hasLimit {
		f.Limit = uint(args.Limit)
	}

	results.Templates, _, err = h.tSvc.Search(ctx, f)
	return
}

func (h templatesHandler) each(ctx context.Context, args *templatesEachArgs) (out wfexec.IteratorHandler, err error) {
	var (
		i = &templateSetIterator{}
		f = types.TemplateFilter{
			Handle:  args.Handle,
			Type:    args.Type,
			OwnerID: args.OwnerID,
			Partial: args.Partial,
			Labels:  args.Labels,
		}
	)

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

	if args.hasLabels {
		f.Labels = args.Labels
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
		i.buffer, i.filter, err = h.tSvc.Search(ctx, i.filter)

		return
	}

	// Initial load
	return i, i.loader()
}

func (h templatesHandler) create(ctx context.Context, args *templatesCreateArgs) (results *templatesCreateResults, err error) {
	results = &templatesCreateResults{}
	results.Template, err = h.tSvc.Create(ctx, args.Template)
	return
}

func (h templatesHandler) update(ctx context.Context, args *templatesUpdateArgs) (results *templatesUpdateResults, err error) {
	results = &templatesUpdateResults{}
	results.Template, err = h.tSvc.Update(ctx, args.Template)
	return
}

func (h templatesHandler) delete(ctx context.Context, args *templatesDeleteArgs) error {
	if id, err := getTemplateID(ctx, h.tSvc, args); err != nil {
		return err
	} else {
		return h.tSvc.DeleteByID(ctx, id)
	}
}

func (h templatesHandler) recover(ctx context.Context, args *templatesRecoverArgs) error {
	if id, err := getTemplateID(ctx, h.tSvc, args); err != nil {
		return err
	} else {
		return h.tSvc.UndeleteByID(ctx, id)
	}
}

func (h templatesHandler) render(ctx context.Context, args *templatesRenderArgs) (*templatesRenderResults, error) {
	var err error

	opts := make(map[string]string)
	if args.hasOptions {
		opts, err = cast.ToStringMapStringE(args.Options)
		if err != nil {
			return nil, err
		}
	}

	tplID, err := getTemplateID(ctx, h.tSvc, args)
	if err != nil {
		return nil, err
	}

	var vars map[string]interface{}
	if args.hasVariables {
		vars = args.Variables.Dict()
	}

	doc, err := h.tSvc.Render(ctx, tplID, args.DocumentType, vars, opts)
	if err != nil {
		return nil, err
	}

	rr := &templatesRenderResults{
		Document: &renderedDocument{
			Document: doc,
			Name:     args.DocumentName,
			Type:     args.DocumentType,
		},
	}

	return rr, nil
}

func (i *templateSetIterator) More(context.Context, *Vars) (bool, error) {
	return wfexec.GenericResourceNextCheck(i.useIterLimit, i.ptr, uint(len(i.buffer)), i.total, i.iterLimit, i.filter.NextPage != nil), nil
}

func (i *templateSetIterator) Start(context.Context, *Vars) error { i.ptr = 0; return nil }

func (i *templateSetIterator) Next(context.Context, *Vars) (out *Vars, err error) {
	if len(i.buffer)-int(i.ptr) <= 0 {
		if err = i.loader(); err != nil {
			panic(err)
		}
	}

	out = &Vars{}
	out.Set("template", Must(NewTemplate(i.buffer[i.ptr])))
	out.Set("index", Must(NewInteger(i.total+i.ptr)))
	out.Set("total", Must(NewInteger(i.filter.Total)))

	i.ptr++
	return out, nil
}

func lookupTemplate(ctx context.Context, svc templateService, args templateLookup) (*types.Template, error) {
	_, ID, handle, template := args.GetLookup()

	switch {
	case template != nil:
		return template, nil
	case ID > 0:
		return svc.FindByID(ctx, ID)
	case len(handle) > 0:
		return svc.FindByHandle(ctx, handle)
	}

	return nil, fmt.Errorf("empty lookup params")
}

func getTemplateID(ctx context.Context, svc templateService, args templateLookup) (uint64, error) {
	_, ID, _, _ := args.GetLookup()

	if ID > 0 {
		return ID, nil
	}

	tpl, err := lookupTemplate(ctx, svc, args)
	if err != nil {
		return 0, err
	}

	return tpl.ID, nil
}
