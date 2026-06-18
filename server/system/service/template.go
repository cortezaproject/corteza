package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/pkg/options"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/renderer"
	"github.com/crusttech/human/server/system/types"
)

// The CRUD skeleton (FindByID, Search, Create, Update, DeleteByID,
// UndeleteByID, loadTemplate and toLabeledTemplates) is generated in
// template.gen.go from system/template.cue.
//
// This file owns the struct, access-controller interface, constructor, the
// `validate` hook the generated Create / Update call into, and the custom
// methods (FindByHandle, FindByAny, Render and helpers).

type (
	template struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        templateAccessController

		renderer rendererService
	}

	templateAccessController interface {
		CanCreateTemplate(context.Context) bool
		CanSearchTemplates(context.Context) bool
		CanReadTemplate(context.Context, *types.Template) bool
		CanUpdateTemplate(context.Context, *types.Template) bool
		CanDeleteTemplate(context.Context, *types.Template) bool
		CanRenderTemplate(context.Context, *types.Template) bool
	}

	rendererService interface {
		Render(ctx context.Context, p *renderer.RendererPayload) (io.ReadSeeker, error)
		Drivers() []renderer.DriverDefinition
	}

	TemplateService interface {
		FindByID(ctx context.Context, ID uint64) (*types.Template, error)
		FindByHandle(ct context.Context, handle string) (*types.Template, error)
		FindByAny(ctx context.Context, identifier interface{}) (*types.Template, error)
		Search(context.Context, types.TemplateFilter) (types.TemplateSet, types.TemplateFilter, error)

		Create(ctx context.Context, tpl *types.Template) (*types.Template, error)
		Update(ctx context.Context, tpl *types.Template) (*types.Template, error)

		DeleteByID(ctx context.Context, ID uint64) error
		UndeleteByID(ctx context.Context, ID uint64) error

		Drivers() []renderer.DriverDefinition
		Render(ctx context.Context, templateID uint64, dstType string, variables map[string]interface{}, options map[string]string) (io.ReadSeeker, error)
	}
)

func Renderer(cfg options.TemplateOpt) *template {
	return &template{
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		ac:        DefaultAccessControl,

		renderer: renderer.Renderer(cfg),
	}
}

// validate runs at the top of the generated Create and Update, on the incoming
// resource.
func (svc *template) validate(ctx context.Context, res *types.Template) error {
	if !handle.IsValid(res.Handle) {
		return TemplateErrInvalidHandle()
	}

	if res.Meta.Short == "" {
		return TemplateErrMissingShort()
	}

	return nil
}

func (svc *template) FindByHandle(ctx context.Context, h string) (tpl *types.Template, err error) {
	var (
		tplProps = &templateActionProps{template: &types.Template{Handle: h}}
	)

	err = func() error {
		if h == "" || !handle.IsValid(h) {
			return TemplateErrInvalidHandle()
		}

		if tpl, err = store.LookupTemplateByHandle(ctx, svc.store, h); err != nil {
			return TemplateErrNotFound().Wrap(err)
		}

		tplProps.setTemplate(tpl)

		if !svc.ac.CanReadTemplate(ctx, tpl) {
			return TemplateErrNotAllowedToRead()
		}

		return nil
	}()

	return tpl, svc.recordAction(ctx, tplProps, TemplateActionLookup, err)
}

func (svc *template) FindByAny(ctx context.Context, identifier interface{}) (tpl *types.Template, err error) {
	if ID, ok := identifier.(uint64); ok {
		tpl, err = svc.FindByID(ctx, ID)
	} else if strIdentifier, ok := identifier.(string); ok {
		if ID, _ := strconv.ParseUint(strIdentifier, 10, 64); ID > 0 {
			tpl, err = svc.FindByID(ctx, ID)
		} else {
			tpl, err = svc.FindByHandle(ctx, strIdentifier)
		}
	} else {
		err = TemplateErrInvalidID()
	}

	if err != nil {
		return
	}

	return
}

func (svc *template) Drivers() []renderer.DriverDefinition {
	return svc.renderer.Drivers()
}

func (svc *template) Render(ctx context.Context, templateID uint64, dstType string, variables map[string]interface{}, options map[string]string) (document io.ReadSeeker, err error) {
	var (
		tplProps = &templateActionProps{}
		tpl      *types.Template
	)

	err = func() (err error) {
		tpl, err = svc.FindByID(ctx, templateID)
		if err != nil {
			return err
		}
		if tpl == nil {
			return TemplateErrNotFound()
		}
		if tpl.Partial {
			return TemplateErrCannotRenderPartial()
		}

		tplProps.setTemplate(tpl)

		if !svc.ac.CanRenderTemplate(ctx, tpl) {
			return TemplateErrNotAllowedToRender()
		}

		// Prepare partials
		//
		// @todo Make this more sophisticated by inspecting the template or
		//       by requiring users to "import" (specify) what partials to use.
		pp, err := svc.getPartials(ctx, tpl)
		if err != nil {
			return err
		}

		att, err := svc.getAttachments(ctx, tpl)
		if err != nil {
			return err
		}

		// Prepare payload
		p := &renderer.RendererPayload{
			Template:     svc.getSource(tpl),
			TemplateType: tpl.Type,
			TargetType:   types.DocumentType(dstType),
			Variables:    variables,
			Options:      options,
			Partials:     pp,
			Attachments:  att,
		}

		// Render the doc
		document, err = svc.renderer.Render(ctx, p)
		if err != nil {
			return err
		}
		return nil
	}()

	return document, svc.recordAction(ctx, tplProps, TemplateActionRender, err)
}

// Util things

func (svc *template) getSource(tpl *types.Template) io.Reader {
	return bytes.NewBuffer([]byte(tpl.Template))
}

func (svc *template) getPartials(ctx context.Context, tpl *types.Template) ([]*renderer.TemplatePartial, error) {
	pp := make([]*renderer.TemplatePartial, 0, 20)

	set, _, err := svc.Search(ctx, types.TemplateFilter{
		Partial: true,
	})
	if err != nil {
		return nil, err
	}

	// @todo inspect original template to filter partials
	// @todo do some filtering based on partial type and main template type

	for _, t := range set {
		tpl := t.Template
		if !strings.HasPrefix(tpl, "{{define") {
			tpl = fmt.Sprintf(`{{define "%s"}}%s{{end}}`, t.Handle, tpl)
		}

		pp = append(pp, &renderer.TemplatePartial{
			Handle:       t.Handle,
			Template:     bytes.NewBuffer([]byte(tpl)),
			TemplateType: t.Type,
		})
	}

	return pp, nil
}

// @todo...
func (svc *template) getAttachments(ctx context.Context, tpl *types.Template) (renderer.AttachmentIndex, error) {
	return make(renderer.AttachmentIndex), nil
}
