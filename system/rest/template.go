package rest

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	Template struct {
		renderer service.TemplateService
		ac       templateAccessController
	}

	templateSetPayload struct {
		Filter types.TemplateFilter `json:"filter"`
		Set    []*templatePayload   `json:"set"`
	}

	templatePayload struct {
		*types.Template

		CanGrant          bool `json:"canGrant"`
		CanUpdateTemplate bool `json:"canUpdateTemplate"`
		CanDeleteTemplate bool `json:"canDeleteTemplate"`
	}

	templateAccessController interface {
		CanGrant(context.Context) bool

		CanAccess(context.Context) bool
		CanCreateTemplate(context.Context) bool
		CanReadTemplate(context.Context, *types.Template) bool
		CanUpdateTemplate(context.Context, *types.Template) bool
		CanDeleteTemplate(context.Context, *types.Template) bool
	}
)

func (Template) New() *Template {
	return &Template{
		renderer: service.DefaultRenderer,
		ac:       service.DefaultAccessControl,
	}
}

func (ctrl *Template) Read(ctx context.Context, r *request.TemplateRead) (interface{}, error) {
	tpl, err := ctrl.renderer.FindByID(ctx, r.TemplateID)
	return ctrl.makePayload(ctx, tpl, err)
}

func (ctrl *Template) List(ctx context.Context, r *request.TemplateList) (interface{}, error) {
	var (
		err error
		f   = types.TemplateFilter{
			Handle:  r.Handle,
			Type:    r.Type,
			OwnerID: r.OwnerID,
			Partial: r.Partial,
			Deleted: filter.State(r.Deleted),
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	set, filter, err := ctrl.renderer.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *Template) Create(ctx context.Context, r *request.TemplateCreate) (interface{}, error) {
	var (
		err error
		app = &types.Template{
			Handle:   r.Handle,
			Language: r.Language,
			Type:     types.DocumentType(r.Type),
			Partial:  r.Partial,
			Meta:     r.Meta,
			Template: r.Template,
			OwnerID:  r.OwnerID,
		}
	)

	app, err = ctrl.renderer.Create(ctx, app)
	return ctrl.makePayload(ctx, app, err)
}

func (ctrl *Template) Update(ctx context.Context, r *request.TemplateUpdate) (interface{}, error) {
	var (
		err error
		app = &types.Template{
			ID:       r.TemplateID,
			Handle:   r.Handle,
			Language: r.Language,
			Type:     types.DocumentType(r.Type),
			Partial:  r.Partial,
			Meta:     r.Meta,
			Template: r.Template,
			OwnerID:  r.OwnerID,
		}
	)

	app, err = ctrl.renderer.Update(ctx, app)
	return ctrl.makePayload(ctx, app, err)
}

func (ctrl *Template) Delete(ctx context.Context, r *request.TemplateDelete) (interface{}, error) {
	return api.OK(), ctrl.renderer.DeleteByID(ctx, r.TemplateID)
}

func (ctrl *Template) Undelete(ctx context.Context, r *request.TemplateUndelete) (interface{}, error) {
	return api.OK(), ctrl.renderer.UndeleteByID(ctx, r.TemplateID)
}

func (ctrl *Template) Render(ctx context.Context, r *request.TemplateRender) (interface{}, error) {
	vars := make(map[string]interface{})
	err := json.Unmarshal(r.Variables, &vars)
	if err != nil {
		return nil, err
	}

	opts := make(map[string]string)
	if r.Options != nil {
		err = json.Unmarshal(r.Options, &opts)
		if err != nil {
			return nil, err
		}
	}

	ct := ctrl.getDestinationType(r.Ext)

	doc, err := ctrl.renderer.Render(ctx, r.TemplateID, ct, vars, opts)
	return ctrl.serve(doc, ct, r, err)
}

// Utilities

func (ctrl Template) makeFilterPayload(ctx context.Context, nn types.TemplateSet, f types.TemplateFilter, err error) (*templateSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &templateSetPayload{Filter: f, Set: make([]*templatePayload, len(nn))}

	for i := range nn {
		msp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return msp, nil
}

func (ctrl Template) makePayload(ctx context.Context, tpl *types.Template, err error) (*templatePayload, error) {
	if err != nil || tpl == nil {
		return nil, err
	}

	pl := &templatePayload{
		Template: tpl,

		CanGrant:          ctrl.ac.CanGrant(ctx),
		CanUpdateTemplate: ctrl.ac.CanUpdateTemplate(ctx, tpl),
		CanDeleteTemplate: ctrl.ac.CanDeleteTemplate(ctx, tpl),
	}

	return pl, nil
}

func (ctrl *Template) serve(doc io.ReadSeeker, ct string, r *request.TemplateRender, err error) (interface{}, error) {
	if err != nil {
		return nil, err
	}

	return func(w http.ResponseWriter, req *http.Request) {
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		name := url.QueryEscape(strings.TrimSpace(r.Filename) + "." + strings.TrimSpace(r.Ext))
		w.Header().Add("Content-Disposition", "attachment; filename="+name)
		w.Header().Add("Content-Type", ct+"; charset=utf-8")

		http.ServeContent(w, req, name, time.Now(), doc)
	}, nil
}

func (ctrl *Template) getDestinationType(ext string) string {
	switch ext {
	case "txt":
		return "text/plain"
	case "html":
		return "text/html"
	case "pdf":
		return "application/pdf"
	}

	return "text/plain"
}
