package rest

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/renderer"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

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

	driverSetPayload struct {
		Set []*driverPayload `json:"set"`
	}

	driverPayload struct {
		renderer.DriverDefinition
	}

	templateAccessController interface {
		CanGrant(context.Context) bool
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

// makeFilter builds the search filter for the generated List controller.
func (ctrl *Template) makeFilter(ctx context.Context, r *request.TemplateList) (types.TemplateFilter, error) {
	var (
		err error
		f   = types.TemplateFilter{
			Query:   r.Query,
			Handle:  r.Handle,
			Type:    r.Type,
			OwnerID: r.OwnerID,
			Partial: r.Partial,
			Deleted: filter.State(r.Deleted),
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
// already mapped the plain-value params.
func (ctrl *Template) beforeCreate(ctx context.Context, res *types.Template, r *request.TemplateCreate) error {
	res.Type = types.DocumentType(r.Type)
	res.Meta = r.Meta
	return nil
}

// beforeUpdate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Update controller
// already mapped the plain-value params (and ID/UpdatedAt).
func (ctrl *Template) beforeUpdate(ctx context.Context, res *types.Template, r *request.TemplateUpdate) error {
	res.Type = types.DocumentType(r.Type)
	res.Meta = r.Meta
	return nil
}

func (ctrl *Template) RenderDrivers(ctx context.Context, r *request.TemplateRenderDrivers) (interface{}, error) {
	return ctrl.makeSetRenderDriverPayload(ctx, ctrl.renderer.Drivers()), nil
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

func (ctrl Template) makeSetRenderDriverPayload(ctx context.Context, nn []renderer.DriverDefinition) *driverSetPayload {
	msp := &driverSetPayload{Set: make([]*driverPayload, len(nn))}

	for i := range nn {
		msp.Set[i] = &driverPayload{
			DriverDefinition: nn[i],
		}
	}

	return msp
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
