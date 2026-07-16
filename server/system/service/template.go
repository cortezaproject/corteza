package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/pkg/options"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/renderer"
	"github.com/crusttech/human/server/system/types"
)

type (
	templateServices struct {
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
		services:  &templateServices{renderer: renderer.Renderer(cfg)},
	}
}

func (svc *template) validate(ctx context.Context, t *types.Template) error {
	if !handle.IsValid(t.Handle) {
		return TemplateErrInvalidHandle()
	}
	if t.Meta.Short == "" {
		return TemplateErrMissingShort()
	}
	return nil
}

func (svc *template) onRender(ctx context.Context, aProps *templateActionProps, templateID uint64, dstType string, variables map[string]interface{}, options map[string]string) (io.ReadSeeker, error) {
	tpl, err := svc.FindByID(ctx, templateID)
	if err != nil {
		return nil, err
	}
	if tpl == nil {
		return nil, TemplateErrNotFound()
	}
	if tpl.Partial {
		return nil, TemplateErrCannotRenderPartial()
	}

	aProps.setTemplate(tpl)

	if !svc.ac.CanRenderTemplate(ctx, tpl) {
		return nil, TemplateErrNotAllowedToRender()
	}

	pp, err := svc.getPartials(ctx, tpl)
	if err != nil {
		return nil, err
	}

	att, err := svc.getAttachments(ctx, tpl)
	if err != nil {
		return nil, err
	}

	// Optional header/footer templates, used by drivers that support them
	header, err := svc.getAuxTemplateSource(ctx, tpl.Meta.HeaderTemplateID)
	if err != nil {
		return nil, err
	}

	footer, err := svc.getAuxTemplateSource(ctx, tpl.Meta.FooterTemplateID)
	if err != nil {
		return nil, err
	}

	p := &renderer.RendererPayload{
		Template:     svc.getSource(tpl),
		TemplateType: tpl.Type,
		TargetType:   types.DocumentType(dstType),
		Variables:    variables,
		Options:      options,
		Partials:     pp,
		Attachments:  att,
		Header:       header,
		Footer:       footer,
	}

	return svc.services.renderer.Render(ctx, p)
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

	return
}

func (svc *template) Drivers() []renderer.DriverDefinition {
	return svc.services.renderer.Drivers()
}

func (svc *template) getSource(tpl *types.Template) io.Reader {
	return bytes.NewBuffer([]byte(tpl.Template))
}

// getAuxTemplateSource loads raw source of the referenced header/footer
// template; (nil, nil) when ID is unset
func (svc *template) getAuxTemplateSource(ctx context.Context, ID uint64) (io.Reader, error) {
	if ID == 0 {
		return nil, nil
	}

	tpl, err := svc.FindByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	return svc.getSource(tpl), nil
}

func (svc *template) getPartials(ctx context.Context, tpl *types.Template) ([]*renderer.TemplatePartial, error) {
	pp := make([]*renderer.TemplatePartial, 0, 20)

	set, _, err := svc.Search(ctx, types.TemplateFilter{
		Partial: true,
	})
	if err != nil {
		return nil, err
	}

	for _, t := range set {
		src := t.Template
		if !strings.HasPrefix(src, "{{define") {
			src = fmt.Sprintf(`{{define "%s"}}%s{{end}}`, t.Handle, src)
		}

		pp = append(pp, &renderer.TemplatePartial{
			Handle:       t.Handle,
			Template:     bytes.NewBuffer([]byte(src)),
			TemplateType: t.Type,
		})
	}

	return pp, nil
}

func (svc *template) getAttachments(_ context.Context, _ *types.Template) (renderer.AttachmentIndex, error) {
	return make(renderer.AttachmentIndex), nil
}

