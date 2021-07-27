package automation

import (
	"context"
	"errors"
	"fmt"

	. "github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	userService interface {
		FindByID(ctx context.Context, userID uint64) (*types.User, error)
		FindByHandle(ctx context.Context, handle string) (*types.User, error)
		FindByEmail(ctx context.Context, email string) (*types.User, error)
		Find(ctx context.Context, filter types.UserFilter) (set types.UserSet, f types.UserFilter, err error)

		Create(ctx context.Context, user *types.User) (*types.User, error)
		Update(ctx context.Context, user *types.User) (*types.User, error)

		Delete(ctx context.Context, id uint64) error
		Suspend(ctx context.Context, id uint64) error
		Unsuspend(ctx context.Context, id uint64) error
		Undelete(ctx context.Context, id uint64) error
	}

	usersHandler struct {
		reg  usersHandlerRegistry
		uSvc userService
		rSvc roleService
	}

	userSetIterator struct {
		ptr    int
		set    types.UserSet
		filter types.UserFilter
	}

	userLookup interface {
		GetLookup() (bool, uint64, string, string, *types.User)
	}
)

func UsersHandler(reg usersHandlerRegistry, uSvc userService, rSvc roleService) *usersHandler {
	h := &usersHandler{
		reg:  reg,
		uSvc: uSvc,
		rSvc: rSvc,
	}

	h.register()
	return h
}

func (h usersHandler) lookup(ctx context.Context, args *usersLookupArgs) (results *usersLookupResults, err error) {
	results = &usersLookupResults{}
	results.User, err = lookupUser(ctx, h.uSvc, args)
	return
}

func (h usersHandler) searchMembership(ctx context.Context, args *usersSearchMembershipArgs) (results *usersSearchMembershipResults, err error) {
	results = &usersSearchMembershipResults{}
	u, err := lookupUser(ctx, h.uSvc, args)
	if err != nil {
		return
	}
	if u == nil {
		return nil, errors.New("user not found")
	}

	// Get the roles
	mm, err := h.rSvc.Membership(ctx, u.ID)
	if err != nil {
		return
	}
	if len(mm) == 0 {
		results.Roles = []*types.Role{}
		return
	}

	rr := make([]uint64, len(mm))
	for i, r := range mm {
		rr[i] = r.RoleID
	}

	results.Roles, _, err = h.rSvc.Find(ctx, types.RoleFilter{
		RoleID: rr,
	})
	results.Total = uint64(len(results.Roles))
	return
}

func (h usersHandler) checkMembership(ctx context.Context, args *usersCheckMembershipArgs) (results *usersCheckMembershipResults, err error) {
	results = &usersCheckMembershipResults{}

	// Resolve user
	u, err := lookupUser(ctx, h.uSvc, usersLookupArgs{
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
	if u == nil {
		return nil, errors.New("user not found")
	}

	// Get user membershio
	mm, err := h.rSvc.Membership(ctx, u.ID)
	if err != nil {
		return
	}

	if len(mm) == 0 {
		return &usersCheckMembershipResults{
			Member: false,
		}, nil
	}

	// Resolve role
	r, err := lookupRole(ctx, h.rSvc, rolesLookupArgs{
		hasLookup:    args.hasRole,
		Lookup:       args.Role,
		lookupID:     args.roleID,
		lookupHandle: args.roleHandle,
		lookupRes:    args.roleRes,
	})
	if err != nil {
		return
	}
	if r == nil {
		return nil, errors.New("role not found")
	}

	// Check if there
	for _, m := range mm {
		if m.RoleID == r.ID {
			return &usersCheckMembershipResults{
				Member: true,
			}, nil
		}
	}

	return &usersCheckMembershipResults{
		Member: false,
	}, nil
}

func (h usersHandler) search(ctx context.Context, args *usersSearchArgs) (results *usersSearchResults, err error) {
	results = &usersSearchResults{}

	var (
		f = types.UserFilter{
			Query:     args.Query,
			Email:     args.Email,
			Handle:    args.Handle,
			Labels:    args.Labels,
			Deleted:   filter.State(args.Deleted),
			Suspended: filter.State(args.Suspended),
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

	var auxf types.UserFilter
	results.Users, auxf, err = h.uSvc.Find(ctx, f)
	results.Total = uint64(auxf.Total)
	return
}

func (h usersHandler) each(ctx context.Context, args *usersEachArgs) (out wfexec.IteratorHandler, err error) {
	var (
		i = &userSetIterator{}
		f = types.UserFilter{
			Query:     args.Query,
			Email:     args.Email,
			Handle:    args.Handle,
			Labels:    args.Labels,
			Deleted:   filter.State(args.Deleted),
			Suspended: filter.State(args.Suspended),
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

	i.set, i.filter, err = h.uSvc.Find(ctx, f)
	return i, err
}

func (h usersHandler) create(ctx context.Context, args *usersCreateArgs) (results *usersCreateResults, err error) {
	results = &usersCreateResults{}
	results.User, err = h.uSvc.Create(ctx, args.User)
	return
}

func (h usersHandler) update(ctx context.Context, args *usersUpdateArgs) (results *usersUpdateResults, err error) {
	results = &usersUpdateResults{}
	results.User, err = h.uSvc.Update(ctx, args.User)
	return
}

func (h usersHandler) delete(ctx context.Context, args *usersDeleteArgs) error {
	if id, err := getUserID(ctx, h.uSvc, args); err != nil {
		return err
	} else {
		return h.uSvc.Delete(ctx, id)
	}
}

func (h usersHandler) recover(ctx context.Context, args *usersRecoverArgs) error {
	if id, err := getUserID(ctx, h.uSvc, args); err != nil {
		return err
	} else {
		return h.uSvc.Undelete(ctx, id)
	}
}

func (h usersHandler) suspend(ctx context.Context, args *usersSuspendArgs) error {
	if id, err := getUserID(ctx, h.uSvc, args); err != nil {
		return err
	} else {
		return h.uSvc.Suspend(ctx, id)
	}
}

func (h usersHandler) unsuspend(ctx context.Context, args *usersUnsuspendArgs) error {
	if id, err := getUserID(ctx, h.uSvc, args); err != nil {
		return err
	} else {
		return h.uSvc.Unsuspend(ctx, id)
	}
}

func getUserID(ctx context.Context, svc userService, args userLookup) (uint64, error) {
	_, ID, _, _, _ := args.GetLookup()

	if ID > 0 {
		return ID, nil
	}

	user, err := lookupUser(ctx, svc, args)
	if err != nil {
		return 0, err
	}

	return user.ID, nil

}

func lookupUser(ctx context.Context, svc userService, args userLookup) (*types.User, error) {
	_, ID, handle, email, user := args.GetLookup()

	switch {
	case user != nil:
		return user, nil
	case ID > 0:
		return svc.FindByID(ctx, ID)
	case len(email) > 0:
		return svc.FindByEmail(ctx, email)
	case len(handle) > 0:
		return svc.FindByHandle(ctx, handle)
	}

	return nil, fmt.Errorf("empty lookup params")
}

func (i *userSetIterator) More(context.Context, *Vars) (bool, error) {
	return i.ptr < len(i.set), nil
}

func (i *userSetIterator) Start(context.Context, *Vars) error { i.ptr = 0; return nil }

func (i *userSetIterator) Next(context.Context, *Vars) (*Vars, error) {
	out := &Vars{}
	out.Set("user", Must(NewUser(i.set[i.ptr])))
	out.Set("index", i.ptr)
	out.Set("total", i.filter.Total)

	i.ptr++
	return out, nil
}
