package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"

	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/store"
	types "github.com/crusttech/human/server/system/types"
)

func (svc *template) FindByID(ctx context.Context, ID uint64) (res *types.Template, err error) {
	var (
		aProps = &templateActionProps{template: &types.Template{ID: ID}}
	)

	err = func() error {
		if res, err = loadTemplate(ctx, svc.store, ID); err != nil {
			return TemplateErrInvalidID().Wrap(err)
		}

		aProps.setTemplate(res)

		if !svc.ac.CanReadTemplate(ctx, res) {
			return TemplateErrNotAllowedToRead()
		}

		return nil
	}()

	return res, svc.recordAction(ctx, aProps, TemplateActionLookup, err)
}

func (svc *template) Search(ctx context.Context, filter types.TemplateFilter) (set types.TemplateSet, f types.TemplateFilter, err error) {
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
		if !svc.ac.CanSearchTemplates(ctx) {
			return TemplateErrNotAllowedToSearch()
		}

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

func (svc *template) Create(ctx context.Context, new *types.Template) (res *types.Template, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &templateActionProps{template: new, new: new}
	)

	err = func() (err error) {
		if err = svc.validate(ctx, new); err != nil {
			return err
		}
		if !svc.ac.CanCreateTemplate(ctx) {
			return TemplateErrNotAllowedToCreate()
		}

		new.ID = nextID()
		new.CreatedAt = *now()

		if err = store.CreateTemplate(ctx, svc.store, new); err != nil {
			return
		}

		if err = label.Create(ctx, svc.store, new); err != nil {
			return
		}

		res = new
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, TemplateActionCreate, err)
}

func (svc *template) Update(ctx context.Context, upd *types.Template) (res *types.Template, err error) {
	var (
		aProps = &templateActionProps{update: upd}
	)

	err = func() (err error) {
		if err = svc.validate(ctx, upd); err != nil {
			return err
		}
		if res, err = loadTemplate(ctx, svc.store, upd.ID); err != nil {
			return
		}

		aProps.setTemplate(res)

		if !svc.ac.CanUpdateTemplate(ctx, res) {
			return TemplateErrNotAllowedToUpdate()
		}

		// Test if stale (update has an older version of data)
		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return TemplateErrStaleData()
		}
		res.Handle = upd.Handle
		res.Language = upd.Language
		res.Type = upd.Type
		res.Partial = upd.Partial
		res.Meta = upd.Meta
		res.Template = upd.Template
		res.OwnerID = upd.OwnerID
		res.UpdatedAt = now()

		if err = store.UpdateTemplate(ctx, svc.store, res); err != nil {
			return err
		}

		if label.Changed(res.Labels, upd.Labels) {
			if err = label.Update(ctx, svc.store, upd); err != nil {
				return
			}
			res.Labels = upd.Labels
		}
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, TemplateActionUpdate, err)
}

func (svc *template) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &templateActionProps{}
		res    *types.Template
	)

	err = func() (err error) {
		if res, err = loadTemplate(ctx, svc.store, ID); err != nil {
			return
		}

		aProps.setTemplate(res)

		if !svc.ac.CanDeleteTemplate(ctx, res) {
			return TemplateErrNotAllowedToDelete()
		}

		res.DeletedAt = now()
		if err = store.UpdateTemplate(ctx, svc.store, res); err != nil {
			return
		}
		return nil
	}()

	return svc.recordAction(ctx, aProps, TemplateActionDelete, err)
}

func (svc *template) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &templateActionProps{}
		res    *types.Template
	)

	err = func() (err error) {
		if res, err = loadTemplate(ctx, svc.store, ID); err != nil {
			return
		}

		aProps.setTemplate(res)

		if !svc.ac.CanDeleteTemplate(ctx, res) {
			return TemplateErrNotAllowedToUndelete()
		}

		res.DeletedAt = nil
		if err = store.UpdateTemplate(ctx, svc.store, res); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, aProps, TemplateActionUndelete, err)
}

func loadTemplate(ctx context.Context, s store.Templates, ID uint64) (res *types.Template, err error) {
	if ID == 0 {
		return nil, TemplateErrInvalidID()
	}

	if res, err = store.LookupTemplateByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, TemplateErrNotFound()
	}

	return
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
