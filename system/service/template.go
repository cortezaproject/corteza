package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/label"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/renderer"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	template struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        templateAccessController

		renderer rendererService
	}

	templateAccessController interface {
		CanCreateTemplate(context.Context) bool
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

func Renderer(cfg options.TemplateOpt) TemplateService {
	return (&template{
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		ac:        DefaultAccessControl,

		renderer: renderer.Renderer(cfg),
	})
}

func (svc template) FindByID(ctx context.Context, ID uint64) (tpl *types.Template, err error) {
	var (
		tplProps = &templateActionProps{template: &types.Template{ID: ID}}
	)

	err = func() error {
		if ID == 0 {
			return TemplateErrInvalidID()
		}

		if tpl, err = store.LookupTemplateByID(ctx, svc.store, ID); err != nil {
			return TemplateErrInvalidID().Wrap(err)
		}

		tplProps.setTemplate(tpl)

		if !svc.ac.CanReadTemplate(ctx, tpl) {
			return TemplateErrNotAllowedToRead()
		}

		return nil
	}()

	return tpl, svc.recordAction(ctx, tplProps, TemplateActionLookup, err)
}

func (svc template) FindByHandle(ctx context.Context, h string) (tpl *types.Template, err error) {
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

func (svc template) FindByAny(ctx context.Context, identifier interface{}) (tpl *types.Template, err error) {
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

func (svc template) Search(ctx context.Context, filter types.TemplateFilter) (set types.TemplateSet, f types.TemplateFilter, err error) {
	var (
		aProps = &templateActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.Template) (bool, error) {
		if !svc.ac.CanReadTemplate(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if len(filter.Labels) > 0 {
			filter.LabeledIDs, err = label.Search(
				ctx,
				svc.store,
				types.Template{}.LabelResourceKind(),
				filter.Labels,
			)

			if err != nil {
				return err
			}

			// labels specified but no labeled resources found
			if len(filter.LabeledIDs) == 0 {
				return nil
			}
		}

		if set, f, err = store.SearchTemplates(ctx, svc.store, filter); err != nil {
			return err
		}

		if err = label.Load(ctx, svc.store, toLabeledTemplates(set)...); err != nil {
			return err
		}

		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, TemplateActionSearch, err)
}

func (svc template) Create(ctx context.Context, new *types.Template) (tpl *types.Template, err error) {
	var (
		tplProps = &templateActionProps{new: new}
	)

	err = func() (err error) {
		if !svc.ac.CanCreateTemplate(ctx) {
			return TemplateErrNotAllowedToCreate()
		}

		// @todo corredor?

		// Set new values after beforeCreate events are emitted
		new.ID = nextID()
		new.CreatedAt = *now()

		if err = store.CreateTemplate(ctx, svc.store, new); err != nil {
			return
		}

		if err = label.Create(ctx, svc.store, new); err != nil {
			return
		}

		tpl = new

		return nil
	}()

	return tpl, svc.recordAction(ctx, tplProps, TemplateActionCreate, err)
}

func (svc template) Update(ctx context.Context, upd *types.Template) (tpl *types.Template, err error) {
	var (
		tplProps = &templateActionProps{update: upd}
	)

	err = func() (err error) {
		if upd.ID == 0 {
			return TemplateErrInvalidID()
		}

		if tpl, err = store.LookupTemplateByID(ctx, svc.store, upd.ID); err != nil {
			return
		}

		tplProps.setTemplate(tpl)

		if !svc.ac.CanUpdateTemplate(ctx, tpl) {
			return TemplateErrNotAllowedToUpdate()
		}

		// @todo corredor?

		tpl.Handle = upd.Handle
		tpl.Language = upd.Language
		tpl.Type = upd.Type
		tpl.Partial = upd.Partial
		tpl.Meta = upd.Meta
		tpl.Template = upd.Template
		tpl.OwnerID = upd.OwnerID
		tpl.UpdatedAt = now()

		if err = store.UpdateTemplate(ctx, svc.store, tpl); err != nil {
			return err
		}

		if label.Changed(tpl.Labels, upd.Labels) {
			if err = label.Update(ctx, svc.store, upd); err != nil {
				return
			}
			tpl.Labels = upd.Labels
		}

		return nil
	}()

	return tpl, svc.recordAction(ctx, tplProps, TemplateActionUpdate, err)
}

func (svc template) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		tplProps = &templateActionProps{}
		tpl      *types.Template
	)

	err = func() (err error) {
		if ID == 0 {
			return TemplateErrInvalidID()
		}

		if tpl, err = store.LookupTemplateByID(ctx, svc.store, ID); err != nil {
			return
		}

		tplProps.setTemplate(tpl)

		if !svc.ac.CanDeleteTemplate(ctx, tpl) {
			return TemplateErrNotAllowedToDelete()
		}

		// @todo corredor?

		tpl.DeletedAt = now()
		if err = store.UpdateTemplate(ctx, svc.store, tpl); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, tplProps, TemplateActionDelete, err)
}

func (svc template) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		tplProps = &templateActionProps{}
		tpl      *types.Template
	)

	err = func() (err error) {
		if ID == 0 {
			return TemplateErrInvalidID()
		}

		if tpl, err = store.LookupTemplateByID(ctx, svc.store, ID); err != nil {
			return
		}

		tplProps.setTemplate(tpl)

		if !svc.ac.CanDeleteTemplate(ctx, tpl) {
			return TemplateErrNotAllowedToUndelete()
		}

		// @todo corredor?
		tpl.DeletedAt = nil
		if err = store.UpdateTemplate(ctx, svc.store, tpl); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, tplProps, TemplateActionUndelete, err)
}

func (svc template) Drivers() []renderer.DriverDefinition {
	return svc.renderer.Drivers()
}

func (svc template) Render(ctx context.Context, templateID uint64, dstType string, variables map[string]interface{}, options map[string]string) (document io.ReadSeeker, err error) {
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

func (svc template) getSource(tpl *types.Template) io.Reader {
	return bytes.NewBuffer([]byte(tpl.Template))
}

func (svc template) getPartials(ctx context.Context, tpl *types.Template) ([]*renderer.TemplatePartial, error) {
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
func (svc template) getAttachments(ctx context.Context, tpl *types.Template) (renderer.AttachmentIndex, error) {
	return make(renderer.AttachmentIndex), nil
	// fpath := "..."

	// att := make(renderer.AttachmentIndex)
	// return att, filepath.Walk(fpath, func(fpath string, info os.FileInfo, err error) error {
	// 	if err != nil {
	// 		return err
	// 	}

	// 	if info.IsDir() {
	// 		return nil
	// 	}

	// 	f, err := os.Open(fpath)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	defer f.Close()

	// 	bb := make([]byte, info.Size())
	// 	_, err = f.Read(bb)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	att[info.Name()] = &renderer.Attachment{
	// 		Source: bytes.NewBuffer(bb),
	// 		// @todo proper implementation!!!
	// 		Mime: "image/png",
	// 		Name: info.Name(),
	// 	}

	// 	return nil
	// })
}

// toLabeledTemplates converts to []label.LabeledResource
func toLabeledTemplates(set []*types.Template) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}
