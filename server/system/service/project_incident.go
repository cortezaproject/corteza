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
	projectIncident struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        projectIncidentAccessController
	}

	projectIncidentAccessController interface {
		CanSearchProjectIncidents(context.Context) bool
		CanCreateProjectIncident(context.Context) bool
		CanReadProjectIncident(context.Context, *types.ProjectIncident) bool
		CanUpdateProjectIncident(context.Context, *types.ProjectIncident) bool
		CanDeleteProjectIncident(context.Context, *types.ProjectIncident) bool
	}
)

func ProjectIncident() *projectIncident {
	return &projectIncident{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	}
}

func (svc projectIncident) FindByID(ctx context.Context, ID uint64) (r *types.ProjectIncident, err error) {
	var (
		raProps = &projectIncidentActionProps{projectIncident: &types.ProjectIncident{ID: ID}}
	)

	err = func() error {
		if r, err = loadProjectIncident(ctx, svc.store, ID); err != nil {
			return err
		}

		raProps.setProjectIncident(r)

		if !svc.ac.CanReadProjectIncident(ctx, r) {
			return ProjectIncidentErrNotAllowedToRead()
		}

		return nil
	}()

	return r, svc.recordAction(ctx, raProps, ProjectIncidentActionLookup, err)
}

func (svc projectIncident) Search(ctx context.Context, filter types.ProjectIncidentFilter) (rr types.ProjectIncidentSet, f types.ProjectIncidentFilter, err error) {
	var (
		raProps = &projectIncidentActionProps{filter: &filter}
	)

	filter.Check = func(res *types.ProjectIncident) (bool, error) {
		if !svc.ac.CanReadProjectIncident(ctx, res) {
			return false, nil
		}
		return true, nil
	}

	err = func() error {
		if !svc.ac.CanSearchProjectIncidents(ctx) {
			return ProjectIncidentErrNotAllowedToSearch()
		}

		if rr, f, err = store.SearchProjectIncidents(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return rr, f, svc.recordAction(ctx, raProps, ProjectIncidentActionSearch, err)
}

func (svc projectIncident) Create(ctx context.Context, new *types.ProjectIncident) (r *types.ProjectIncident, err error) {
	var (
		raProps = &projectIncidentActionProps{projectIncident: new, new: new}
	)

	err = func() (err error) {
		if !svc.ac.CanCreateProjectIncident(ctx) {
			return ProjectIncidentErrNotAllowedToCreate()
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

		if err = store.CreateProjectIncident(ctx, svc.store, new); err != nil {
			return
		}

		r = new
		return nil
	}()

	return r, svc.recordAction(ctx, raProps, ProjectIncidentActionCreate, err)
}

func (svc projectIncident) Update(ctx context.Context, upd *types.ProjectIncident) (r *types.ProjectIncident, err error) {
	var (
		raProps = &projectIncidentActionProps{projectIncident: upd, update: upd}
	)

	err = func() (err error) {
		if r, err = loadProjectIncident(ctx, svc.store, upd.ID); err != nil {
			return
		}

		raProps.setProjectIncident(r)

		if !svc.ac.CanUpdateProjectIncident(ctx, r) {
			return ProjectIncidentErrNotAllowedToUpdate()
		}

		r.Title = upd.Title
		r.Description = upd.Description
		r.IncidentType = upd.IncidentType
		r.GroupSystem = upd.GroupSystem
		r.Status = upd.Status
		r.Severity = upd.Severity
		r.Risk = upd.Risk
		r.IssueOwner = upd.IssueOwner
		r.ChangeOwner = upd.ChangeOwner
		r.ChangeApprovedBy = upd.ChangeApprovedBy
		r.RiskIssue = upd.RiskIssue
		r.ChangeRequired = upd.ChangeRequired
		r.RiskChange = upd.RiskChange
		r.Backlog = upd.Backlog
		r.DateDue = upd.DateDue
		r.CompletedDate = upd.CompletedDate
		r.UpdatedAt = now()
		r.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateProjectIncident(ctx, svc.store, r); err != nil {
			return err
		}

		return nil
	}()

	return r, svc.recordAction(ctx, raProps, ProjectIncidentActionUpdate, err)
}

func (svc projectIncident) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		raProps = &projectIncidentActionProps{}
		r       *types.ProjectIncident
	)

	err = func() (err error) {
		if r, err = loadProjectIncident(ctx, svc.store, ID); err != nil {
			return
		}

		raProps.setProjectIncident(r)

		if !svc.ac.CanDeleteProjectIncident(ctx, r) {
			return ProjectIncidentErrNotAllowedToDelete()
		}

		r.DeletedAt = now()
		if err = store.UpdateProjectIncident(ctx, svc.store, r); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, raProps, ProjectIncidentActionDelete, err)
}

func loadProjectIncident(ctx context.Context, s store.ProjectIncidents, ID uint64) (res *types.ProjectIncident, err error) {
	if ID == 0 {
		return nil, ProjectIncidentErrInvalidID()
	}

	if res, err = store.LookupProjectIncidentByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, ProjectIncidentErrNotFound()
	}

	return
}
