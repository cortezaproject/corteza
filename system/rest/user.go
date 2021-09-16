package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/corredor"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/service/event"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

var _ = errors.Wrap

type (
	User struct {
		user service.UserService
		role service.RoleService
	}

	userSetPayload struct {
		Filter types.UserFilter `json:"filter"`
		Set    types.UserSet    `json:"set"`
	}
)

func (User) New() *User {
	ctrl := &User{}
	ctrl.user = service.DefaultUser
	ctrl.role = service.DefaultRole
	return ctrl
}

func (ctrl User) List(ctx context.Context, r *request.UserList) (interface{}, error) {
	var (
		err error
		set types.UserSet
		f   = types.UserFilter{
			UserID:    payload.ParseUint64s(r.UserID),
			RoleID:    payload.ParseUint64s(r.RoleID),
			Query:     r.Query,
			Email:     r.Email,
			Username:  r.Username,
			Handle:    r.Handle,
			Kind:      r.Kind,
			Labels:    r.Labels,
			Suspended: filter.State(r.Suspended),
			Deleted:   filter.State(r.Deleted),
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	if r.IncSuspended && f.Suspended == 0 {
		f.Suspended = filter.StateInclusive
	}

	if r.IncDeleted && f.Deleted == 0 {
		f.Deleted = filter.StateInclusive
	}

	set, f, err = ctrl.user.Find(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, f, err)
}

func (ctrl User) Create(ctx context.Context, r *request.UserCreate) (interface{}, error) {
	user := &types.User{
		Email:  r.Email,
		Name:   r.Name,
		Handle: r.Handle,
		Kind:   r.Kind,
		Labels: r.Labels,
	}

	return ctrl.user.Create(ctx, user)
}

func (ctrl User) Update(ctx context.Context, r *request.UserUpdate) (interface{}, error) {
	user := &types.User{
		ID:     r.UserID,
		Email:  r.Email,
		Name:   r.Name,
		Handle: r.Handle,
		Kind:   r.Kind,
		Labels: r.Labels,
	}

	return ctrl.user.Update(ctx, user)
}

type (
	patchOp struct {
		Operation string          `json:"op"`
		Path      string          `json:"path"`
		Value     json.RawMessage `json:"value"`
	}
)

// PartialUpdate
//
// experimental resource management with partial updates (patching) using
// JavaScript Object Notation (JSON) Patch standard (RFC 6902)
//
// If this proves useful, we'll use it on other resources & fields
func (ctrl User) PartialUpdate(ctx context.Context, r *request.UserPartialUpdate) (interface{}, error) {
	u, err := ctrl.user.FindByID(ctx, r.UserID)
	if err != nil {
		return nil, err
	}

	return func(w http.ResponseWriter, r *http.Request) {
		err = func() (err error) {
			var (
				ops = make([]*patchOp, 0)
			)

			if err = json.NewDecoder(r.Body).Decode(&ops); err != nil {
				return err
			}

			for _, p := range ops {
				if p.Operation != "replace" {
					return fmt.Errorf("unsupported operation '%s'", p.Operation)
				}

				var aux interface{}
				err = json.Unmarshal(p.Value, &aux)

				switch p.Path {
				case "/meta/preferredLanguage":
					u.Meta.PreferredLanguage, err = cast.ToStringE(aux)

				case "/meta/securityPolicy/mfa/enforcedEmailOTP":
					u.Meta.SecurityPolicy.MFA.EnforcedEmailOTP, err = cast.ToBoolE(aux)

				case "/meta/securityPolicy/mfa/enforcedTOTP":
					u.Meta.SecurityPolicy.MFA.EnforcedTOTP, err = cast.ToBoolE(aux)

				case "/emailConfirmed":
					// unfortunately, this cannot be passed to update right now
					// internal limitations
					u.EmailConfirmed, err = cast.ToBoolE(aux)
					err = ctrl.user.ToggleEmailConfirmation(ctx, u.ID, u.EmailConfirmed)

				default:
					return fmt.Errorf("unknown path: %s", p.Path)
				}

				if err != nil {
					return fmt.Errorf("could not replace falue of %s: %w", p.Path, err)
				}
			}

			u, err = ctrl.user.Update(ctx, u)
			return err
		}()

		if err != nil {
			api.Send(w, r, err)
		}

		api.Send(w, r, u)
	}, nil
}

func (ctrl User) Read(ctx context.Context, r *request.UserRead) (interface{}, error) {
	return ctrl.user.FindByID(ctx, r.UserID)
}

func (ctrl User) Delete(ctx context.Context, r *request.UserDelete) (interface{}, error) {
	return api.OK(), ctrl.user.Delete(ctx, r.UserID)
}

func (ctrl User) Suspend(ctx context.Context, r *request.UserSuspend) (interface{}, error) {
	return api.OK(), ctrl.user.Suspend(ctx, r.UserID)
}

func (ctrl User) Unsuspend(ctx context.Context, r *request.UserUnsuspend) (interface{}, error) {
	return api.OK(), ctrl.user.Unsuspend(ctx, r.UserID)
}

func (ctrl User) Undelete(ctx context.Context, r *request.UserUndelete) (interface{}, error) {
	return api.OK(), ctrl.user.Undelete(ctx, r.UserID)
}

func (ctrl User) SetPassword(ctx context.Context, r *request.UserSetPassword) (interface{}, error) {
	return api.OK(), ctrl.user.SetPassword(ctx, r.UserID, r.Password)
}

func (ctrl User) MembershipList(ctx context.Context, r *request.UserMembershipList) (interface{}, error) {
	if mm, err := ctrl.role.Membership(ctx, r.UserID); err != nil {
		return nil, err
	} else {
		rval := make([]string, len(mm))
		for i := range mm {
			rval[i] = payload.Uint64toa(mm[i].RoleID)
		}
		return rval, nil
	}
}

func (ctrl User) MembershipAdd(ctx context.Context, r *request.UserMembershipAdd) (interface{}, error) {
	return api.OK(), ctrl.role.MemberAdd(ctx, r.RoleID, r.UserID)
}

func (ctrl User) MembershipRemove(ctx context.Context, r *request.UserMembershipRemove) (interface{}, error) {
	return api.OK(), ctrl.role.MemberRemove(ctx, r.RoleID, r.UserID)
}

func (ctrl *User) TriggerScript(ctx context.Context, r *request.UserTriggerScript) (rsp interface{}, err error) {
	var (
		user *types.User
	)

	if user, err = ctrl.user.FindByID(ctx, r.UserID); err != nil {
		return
	}

	// @todo implement same behaviour as we have on record - user+oldUser
	err = corredor.Service().Exec(ctx, r.Script, event.UserOnManual(user, user))
	return user, err

}

func (ctrl *User) SessionsRemove(ctx context.Context, r *request.UserSessionsRemove) (rsp interface{}, err error) {
	var (
		user *types.User
	)

	if user, err = ctrl.user.FindByID(ctx, r.UserID); err != nil {
		return
	}

	if err = ctrl.user.DeleteAuthSessionsByUserID(ctx, user.ID); err != nil {
		return
	}

	if err = ctrl.user.DeleteAuthTokensByUserID(ctx, user.ID); err != nil {
		return
	}

	return
}

func (ctrl User) makeFilterPayload(ctx context.Context, uu types.UserSet, f types.UserFilter, err error) (*userSetPayload, error) {
	if err != nil {
		return nil, err
	}

	if len(uu) == 0 {
		uu = make([]*types.User, 0)
	}

	return &userSetPayload{Filter: f, Set: uu}, nil
}
