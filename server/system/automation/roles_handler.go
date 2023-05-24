package automation

import (
	"context"
	"fmt"

	. "github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/wfexec"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/spf13/cast"
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

		Membership(ctx context.Context, userID uint64) (types.RoleMemberSet, error)
		MemberList(ctx context.Context, roleID uint64) (types.RoleMemberSet, error)
		MemberAdd(ctx context.Context, roleID, userID uint64) error
		MemberRemove(ctx context.Context, roleID, userID uint64) error
	}

	rolesHandler struct {
		reg  rolesHandlerRegistry
		rSvc roleService
		uSvc userService
	}

	roleSetIterator struct {
		// Item buffer, current item pointer, and total items traversed
		ptr    uint
		buffer types.RoleSet
		total  uint

		// When filter limit is set, this constraints it
		iterLimit    uint
		useIterLimit bool

		// Item loader for additional chunks
		filter types.RoleFilter
		loader func() error
	}

	roleLookup interface {
		GetLookup() (bool, uint64, string, *types.Role)
	}
)

func RolesHandler(reg rolesHandlerRegistry, rSvc roleService, uSvc userService) *rolesHandler {
	h := &rolesHandler{
		reg:  reg,
		rSvc: rSvc,
		uSvc: uSvc,
	}

	h.register()
	return h
}

func (h rolesHandler) lookup(ctx context.Context, args *rolesLookupArgs) (results *rolesLookupResults, err error) {
	results = &rolesLookupResults{}
	results.Role, err = lookupRole(ctx, h.rSvc, args)
	return
}

func (h rolesHandler) searchMembers(ctx context.Context, args *rolesSearchMembersArgs) (results *rolesSearchMembersResults, err error) {
	results = &rolesSearchMembersResults{}

	rl, err := lookupRole(ctx, h.rSvc, args)
	if err != nil {
		return
	}
	if rl == nil {
		return nil, fmt.Errorf("role not found")
	}

	// Get membership info
	mm, err := h.rSvc.MemberList(ctx, rl.ID)
	if err != nil {
		return
	}

	if len(mm) == 0 {
		results.Users = []*types.User{}
		return
	}

	// Get actual users
	uu := make([]string, len(mm))
	for i, m := range mm {
		uu[i] = cast.ToString(m.UserID)
	}
	results.Users, _, err = h.uSvc.Find(ctx, types.UserFilter{
		UserID: uu,
	})
	results.Total = uint64(len(results.Users))

	return
}

func (h rolesHandler) eachMember(ctx context.Context, args *rolesEachMemberArgs) (out wfexec.IteratorHandler, err error) {
	var (
		i = &userSetIterator{}
	)

	rl, err := lookupRole(ctx, h.rSvc, args)
	if err != nil {
		return
	}
	if rl == nil {
		return nil, fmt.Errorf("role not found")
	}

	// Get membership info
	mm, err := h.rSvc.MemberList(ctx, rl.ID)
	if err != nil {
		return
	}

	if len(mm) == 0 {
		i.buffer = []*types.User{}
		i.filter = types.UserFilter{}
		return i, nil
	}

	// Get actual users
	uu := make([]string, len(mm))
	for i, m := range mm {
		uu[i] = cast.ToString(m.UserID)
	}
	i.buffer, i.filter, err = h.uSvc.Find(ctx, types.UserFilter{
		UserID: uu,
	})
	return i, err
}

func (h rolesHandler) addMember(ctx context.Context, args *rolesAddMemberArgs) (err error) {
	role, err := lookupRole(ctx, h.rSvc, &rolesLookupArgs{
		hasLookup:    args.hasRole,
		Lookup:       args.Role,
		lookupID:     args.roleID,
		lookupHandle: args.roleHandle,
		lookupRes:    args.roleRes,
	})
	if err != nil {
		return
	}
	if role == nil {
		return fmt.Errorf("role not found")
	}

	user, err := lookupUser(ctx, h.uSvc, &usersLookupArgs{
		hasLookup:    args.hasUser,
		Lookup:       args.User,
		lookupID:     args.userID,
		lookupHandle: args.userHandle,
		lookupEmail:  args.userEmail,
		lookupRes:    args.userRes,
	})
	if err != nil {
		return
	}
	if role == nil {
		return fmt.Errorf("user not found")
	}

	return h.rSvc.MemberAdd(ctx, role.ID, user.ID)
}

func (h rolesHandler) removeMember(ctx context.Context, args *rolesRemoveMemberArgs) (err error) {
	role, err := lookupRole(ctx, h.rSvc, &rolesLookupArgs{
		hasLookup:    args.hasRole,
		Lookup:       args.Role,
		lookupID:     args.roleID,
		lookupHandle: args.roleHandle,
		lookupRes:    args.roleRes,
	})
	if err != nil {
		return
	}
	if role == nil {
		return fmt.Errorf("role not found")
	}

	user, err := lookupUser(ctx, h.uSvc, &usersLookupArgs{
		hasLookup:    args.hasUser,
		Lookup:       args.User,
		lookupID:     args.userID,
		lookupHandle: args.userHandle,
		lookupEmail:  args.userEmail,
		lookupRes:    args.userRes,
	})
	if err != nil {
		return
	}
	if role == nil {
		return fmt.Errorf("user not found")
	}

	return h.rSvc.MemberRemove(ctx, role.ID, user.ID)
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

	var auxf types.RoleFilter
	results.Roles, auxf, err = h.rSvc.Find(ctx, f)
	results.Total = uint64(auxf.Total)
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
		i.buffer, i.filter, err = h.rSvc.Find(ctx, i.filter)

		return
	}

	// Initial load
	return i, i.loader()
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
	a := wfexec.GenericResourceNextCheck(i.useIterLimit, i.ptr, uint(len(i.buffer)), i.total, i.iterLimit, i.filter.NextPage != nil)
	return a, nil
}

func (i *roleSetIterator) Start(context.Context, *Vars) error { i.ptr = 0; return nil }

func (i *roleSetIterator) Next(context.Context, *Vars) (out *Vars, err error) {
	if len(i.buffer)-int(i.ptr) <= 0 {
		if err = i.loader(); err != nil {
			panic(err)
		}
	}

	out = &Vars{}
	out.Set("role", Must(NewRole(i.buffer[i.ptr])))
	out.Set("index", Must(NewInteger(i.total+i.ptr)))
	out.Set("total", Must(NewInteger(i.filter.Total)))

	i.ptr++
	return out, nil
}
