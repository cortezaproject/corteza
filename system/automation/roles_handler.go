package automation

import (
	"context"
	"fmt"
	. "github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	roleService interface {
		FindByID(ctx context.Context, roleID uint64) (*types.Role, error)
		FindByHandle(ctx context.Context, handle string) (*types.Role, error)
		Find(ctx context.Context, filter types.RoleFilter) (set types.RoleSet, f types.RoleFilter, err error)

		Create(ctx context.Context, role *types.Role) (*types.Role, error)
		Update(ctx context.Context, role *types.Role) (*types.Role, error)

		Delete(ctx context.Context, id uint64) error
		Archive(ctx context.Context, id uint64) error
		Unarchive(ctx context.Context, id uint64) error
		Undelete(ctx context.Context, id uint64) error
	}

	rolesHandler struct {
		reg  rolesHandlerRegistry
		rSvc roleService
	}

	roleSetIterator struct {
		ptr    int
		set    types.RoleSet
		filter types.RoleFilter
	}

	roleLookup interface {
		GetLookup() (bool, uint64, string, *types.Role)
	}
)

func RolesHandler(reg rolesHandlerRegistry, rSvc roleService) *rolesHandler {
	h := &rolesHandler{
		reg:  reg,
		rSvc: rSvc,
	}

	h.register()
	return h
}

func (h rolesHandler) lookup(ctx context.Context, args *rolesLookupArgs) (results *rolesLookupResults, err error) {
	results = &rolesLookupResults{}
	results.Role, err = lookupRole(ctx, h.rSvc, args)
	return
}

func (h rolesHandler) search(ctx context.Context, args *rolesSearchArgs) (results *rolesSearchResults, err error) {
	results = &rolesSearchResults{}

	var (
		f = types.RoleFilter{
			Query:    args.Query,
			MemberID: args.MemberID,
			Name:     args.Name,
			Handle:   args.Handle,
			Labels:   args.Labels,
			Deleted:  filter.State(args.Deleted),
			Archived: filter.State(args.Archived),
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

	results.Roles, _, err = h.rSvc.Find(ctx, f)
	return
}

func (h rolesHandler) each(ctx context.Context, args *rolesEachArgs) (out wfexec.IteratorHandler, err error) {
	var (
		i = &roleSetIterator{}
		f = types.RoleFilter{
			Query:    args.Query,
			MemberID: args.MemberID,
			Name:     args.Name,
			Handle:   args.Handle,
			Labels:   args.Labels,
			Deleted:  filter.State(args.Deleted),
			Archived: filter.State(args.Archived),
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

	i.set, i.filter, err = h.rSvc.Find(ctx, f)
	return i, err
}

func (h rolesHandler) create(ctx context.Context, args *rolesCreateArgs) (results *rolesCreateResults, err error) {
	results = &rolesCreateResults{}
	results.Role, err = h.rSvc.Create(ctx, args.Role)
	return
}

func (h rolesHandler) update(ctx context.Context, args *rolesUpdateArgs) (results *rolesUpdateResults, err error) {
	results = &rolesUpdateResults{}
	results.Role, err = h.rSvc.Update(ctx, args.Role)
	return
}

func (h rolesHandler) delete(ctx context.Context, args *rolesDeleteArgs) error {
	if id, err := getRoleID(ctx, h.rSvc, args); err != nil {
		return err
	} else {
		return h.rSvc.Delete(ctx, id)
	}
}

func (h rolesHandler) recover(ctx context.Context, args *rolesRecoverArgs) error {
	if id, err := getRoleID(ctx, h.rSvc, args); err != nil {
		return err
	} else {
		return h.rSvc.Undelete(ctx, id)
	}
}

func (h rolesHandler) archive(ctx context.Context, args *rolesArchiveArgs) error {
	if id, err := getRoleID(ctx, h.rSvc, args); err != nil {
		return err
	} else {
		return h.rSvc.Archive(ctx, id)
	}
}

func (h rolesHandler) unarchive(ctx context.Context, args *rolesUnarchiveArgs) error {
	if id, err := getRoleID(ctx, h.rSvc, args); err != nil {
		return err
	} else {
		return h.rSvc.Unarchive(ctx, id)
	}
}

func getRoleID(ctx context.Context, svc roleService, args roleLookup) (uint64, error) {
	_, ID, _, _ := args.GetLookup()

	if ID > 0 {
		return ID, nil
	}

	role, err := lookupRole(ctx, svc, args)
	if err != nil {
		return 0, err
	}

	return role.ID, nil

}

func lookupRole(ctx context.Context, svc roleService, args roleLookup) (*types.Role, error) {
	_, ID, handle, role := args.GetLookup()

	switch {
	case role != nil:
		return role, nil
	case ID > 0:
		return svc.FindByID(ctx, ID)
	case len(handle) > 0:
		return svc.FindByHandle(ctx, handle)
	}

	return nil, fmt.Errorf("empty lookup params")
}

func (i *roleSetIterator) More(context.Context, *Vars) (bool, error) {
	return i.ptr < len(i.set), nil
}

func (i *roleSetIterator) Start(context.Context, *Vars) error { i.ptr = 0; return nil }

func (i *roleSetIterator) Next(context.Context, *Vars) (*Vars, error) {
	out := RVars{
		"role":  Must(NewRole(i.set[i.ptr])),
		"total": Must(NewUnsignedInteger(i.filter.Total)),
	}

	i.ptr++
	return out.Vars(), nil
}
