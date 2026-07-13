package service

import (
	"context"

	"github.com/crusttech/human/server/pkg/actionlog"
	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

type (
	projectPrivacy struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        projectPrivacyAccessController
	}

	projectPrivacyAccessController interface {
		CanSearchProjectPrivacys(context.Context) bool
		CanCreateProjectPrivacy(context.Context) bool
		CanReadProjectPrivacy(context.Context, *types.ProjectPrivacy) bool
		CanUpdateProjectPrivacy(context.Context, *types.ProjectPrivacy) bool
		CanDeleteProjectPrivacy(context.Context, *types.ProjectPrivacy) bool
	}
)

func ProjectPrivacy() *projectPrivacy {
	return &projectPrivacy{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	}
}

func (svc projectPrivacy) FindByID(ctx context.Context, ID uint64) (r *types.ProjectPrivacy, err error) {
	var (
		raProps = &projectPrivacyActionProps{projectPrivacy: &types.ProjectPrivacy{ID: ID}}
	)

	err = func() error {
		if r, err = loadProjectPrivacy(ctx, svc.store, ID); err != nil {
			return err
		}

		raProps.setProjectPrivacy(r)

		if !svc.ac.CanReadProjectPrivacy(ctx, r) {
			return ProjectPrivacyErrNotAllowedToRead()
		}

		return nil
	}()

	return r, svc.recordAction(ctx, raProps, ProjectPrivacyActionLookup, err)
}

func (svc projectPrivacy) Search(ctx context.Context, filter types.ProjectPrivacyFilter) (rr types.ProjectPrivacySet, f types.ProjectPrivacyFilter, err error) {
	var (
		raProps = &projectPrivacyActionProps{filter: &filter}
	)

	filter.Check = func(res *types.ProjectPrivacy) (bool, error) {
		if !svc.ac.CanReadProjectPrivacy(ctx, res) {
			return false, nil
		}
		return true, nil
	}

	err = func() error {
		if !svc.ac.CanSearchProjectPrivacys(ctx) {
			return ProjectPrivacyErrNotAllowedToSearch()
		}

		if rr, f, err = store.SearchProjectPrivacys(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return rr, f, svc.recordAction(ctx, raProps, ProjectPrivacyActionSearch, err)
}

func (svc projectPrivacy) Create(ctx context.Context, new *types.ProjectPrivacy) (r *types.ProjectPrivacy, err error) {
	var (
		raProps = &projectPrivacyActionProps{projectPrivacy: new, new: new}
	)

	err = func() (err error) {
		if !svc.ac.CanCreateProjectPrivacy(ctx) {
			return ProjectPrivacyErrNotAllowedToCreate()
		}

		if len(new.Title) == 0 {
			new.Title = firstNonEmpty(new.Description, "(untitled)")
		}
		if len(new.Status) == 0 {
			new.Status = "Open"
		}

		new.ID = nextID()
		new.CreatedAt = *now()
		new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.CreateProjectPrivacy(ctx, svc.store, new); err != nil {
			return
		}

		r = new
		return nil
	}()

	return r, svc.recordAction(ctx, raProps, ProjectPrivacyActionCreate, err)
}

func (svc projectPrivacy) Update(ctx context.Context, upd *types.ProjectPrivacy) (r *types.ProjectPrivacy, err error) {
	var (
		raProps = &projectPrivacyActionProps{projectPrivacy: upd, update: upd}
	)

	err = func() (err error) {
		if r, err = loadProjectPrivacy(ctx, svc.store, upd.ID); err != nil {
			return
		}

		raProps.setProjectPrivacy(r)

		if !svc.ac.CanUpdateProjectPrivacy(ctx, r) {
			return ProjectPrivacyErrNotAllowedToUpdate()
		}

		r.Title = upd.Title
		r.Description = upd.Description
		r.RequestType = upd.RequestType
		r.Status = upd.Status
		r.Severity = upd.Severity
		r.Risk = upd.Risk
		r.RequestOwner = upd.RequestOwner
		r.ChangeOwner = upd.ChangeOwner
		r.ChangeApprovedBy = upd.ChangeApprovedBy
		r.RiskAssessment = upd.RiskAssessment
		r.ChangeRequired = upd.ChangeRequired
		r.RiskChange = upd.RiskChange
		r.Backlog = upd.Backlog
		r.DateDue = upd.DateDue
		r.UpdatedAt = now()
		r.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateProjectPrivacy(ctx, svc.store, r); err != nil {
			return err
		}

		return nil
	}()

	return r, svc.recordAction(ctx, raProps, ProjectPrivacyActionUpdate, err)
}

func (svc projectPrivacy) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		raProps = &projectPrivacyActionProps{}
		r       *types.ProjectPrivacy
	)

	err = func() (err error) {
		if r, err = loadProjectPrivacy(ctx, svc.store, ID); err != nil {
			return
		}

		raProps.setProjectPrivacy(r)

		if !svc.ac.CanDeleteProjectPrivacy(ctx, r) {
			return ProjectPrivacyErrNotAllowedToDelete()
		}

		r.DeletedAt = now()
		if err = store.UpdateProjectPrivacy(ctx, svc.store, r); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, raProps, ProjectPrivacyActionDelete, err)
}

func loadProjectPrivacy(ctx context.Context, s store.ProjectPrivacys, ID uint64) (res *types.ProjectPrivacy, err error) {
	if ID == 0 {
		return nil, ProjectPrivacyErrInvalidID()
	}

	if res, err = store.LookupProjectPrivacyByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, ProjectPrivacyErrNotFound()
	}

	return
}
