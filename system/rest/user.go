package rest

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/corredor"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	envoyStore "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/pkg/envoy/yaml"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/service/event"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/gabriel-vasile/mimetype"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

var _ = errors.Wrap

type (
	User struct {
		user service.UserService
		role service.RoleService

		userAc userAccessController
		roleAc roleAccessController
	}

	userSetPayload struct {
		Filter types.UserFilter `json:"filter"`
		Set    types.UserSet    `json:"set"`
	}

	userAccessController interface {
		CanCreateUser(context.Context) bool
		CanUpdateUser(context.Context, *types.User) bool
	}
)

func (User) New() *User {
	return &User{
		user: service.DefaultUser,
		role: service.DefaultRole,

		userAc: service.DefaultAccessControl,
		roleAc: service.DefaultAccessControl,
	}
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
	err = corredor.Service().Exec(ctx, r.Script, corredor.ExtendScriptArgs(event.UserOnManual(user, user), r.Args))
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

// Export exports users with optional role membership and related roles
//
// @note this is a temporary implementation; it will be reworked when we rework Envoy and related bits.
func (ctrl *User) Export(ctx context.Context, r *request.UserExport) (rsp interface{}, err error) {
	// Users
	uu, _, err := ctrl.user.Find(ctx, types.UserFilter{})
	if err != nil {
		return
	}

	// Roles
	roleIndex := make(map[uint64]*types.Role)
	roleResIndex := make(map[uint64]resource.Interface)
	rr, _, err := ctrl.role.Find(ctx, types.RoleFilter{Paging: filter.Paging{Limit: 0}})
	if err != nil {
		return
	}
	for _, r := range rr {
		roleIndex[r.ID] = r
	}

	// Membership
	resources := make(resource.InterfaceSet, 0, len(uu))
	var membership types.RoleMemberSet
	for _, u := range uu {
		usrRes := resource.NewUser(u)

		if r.InclRoleMembership {
			membership, err = ctrl.role.Membership(ctx, u.ID)
			if err != nil {
				return
			}

			aux := make(types.RoleSet, 0, 2)

			for _, m := range membership {
				if _, ok := roleResIndex[m.RoleID]; !ok {
					roleResIndex[m.RoleID] = resource.NewRole(roleIndex[m.RoleID])
					if r.InclRoles {
						resources = append(resources, roleResIndex[m.RoleID])
					}
				}
				aux = append(aux, roleIndex[m.RoleID])
			}

			usrRes.AddRoles(aux...)
		}

		resources = append(resources, usrRes)
	}

	// Encode
	ye := yaml.NewYamlEncoder(&yaml.EncoderConfig{})
	bld := envoy.NewBuilder(ye)
	g, err := bld.Build(ctx, resources...)
	if err != nil {
		return nil, err
	}

	err = envoy.Encode(ctx, g, ye)
	if err != nil {
		return
	}

	// make archive
	buf := bytes.NewBuffer(nil)
	w := zip.NewWriter(buf)

	var (
		f  io.Writer
		bb []byte
	)
	for _, s := range ye.Stream() {
		// @todo generalize when needed
		f, err = w.Create(fmt.Sprintf("%s.yaml", s.Resource))
		if err != nil {
			return
		}

		bb, err = ioutil.ReadAll(s.Source)
		if err != nil {
			return
		}

		_, err = f.Write(bb)
		if err != nil {
			return
		}
	}

	err = w.Close()
	if err != nil {
		return
	}
	return ctrl.serve(ctx, fmt.Sprintf("%s.zip", r.Filename), bytes.NewReader(buf.Bytes()), nil)
}

// Import imports users with optional role membership and related roles
//
// @note this is a temporary implementation; it will be reworked when we rework Envoy and related bits.
func (ctrl *User) Import(ctx context.Context, r *request.UserImport) (rsp interface{}, err error) {
	// AC
	// @todo refactor when we refactor this part of the sys
	if !ctrl.userAc.CanCreateUser(ctx) {
		err = fmt.Errorf("cannot import users: not allowed to create users")
		return
	}
	if !ctrl.roleAc.CanCreateRole(ctx) {
		err = fmt.Errorf("cannot import users: not allowed to create roles")
		return
	}

	// Parse inputs
	f, err := r.Upload.Open()
	if err != nil {
		return
	}
	defer f.Close()

	mt, err := mimetype.DetectReader(f)
	if err != nil {
		return
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return
	}

	if !mt.Is("application/zip") {
		err = fmt.Errorf("cannot import users: unsupported file format")
		return
	}

	// un-archive
	archive, err := zip.NewReader(f, r.Upload.Size)
	if err != nil {
		return
	}

	// decode with Envoy
	yd := yaml.Decoder()
	nn := make([]resource.Interface, 0, 10)
	var mm []resource.Interface
	for _, archF := range archive.File {
		if archF.FileInfo().IsDir() {
			continue
		}
		var f io.ReadCloser

		f, err = archF.Open()
		if err != nil {
			return
		}
		defer f.Close()

		mm, err = yd.Decode(ctx, f, nil)
		if err != nil {
			return
		}
		nn = append(nn, mm...)
	}

	// Validate
	for _, n := range nn {
		switch n.ResourceType() {
		case types.UserResourceType,
			types.RoleResourceType:
			continue

		default:
			err = fmt.Errorf("cannot import users: invalid resource provided: %s", n.ResourceType())
			return
		}
	}

	se := envoyStore.NewStoreEncoder(service.DefaultStore, dal.Service(), &envoyStore.EncoderConfig{
		OnExisting: resource.Skip,
	})

	bld := envoy.NewBuilder(se)
	g, err := bld.Build(ctx, nn...)
	if err != nil {
		return
	}

	err = envoy.Encode(ctx, g, se)
	return api.OK(), err
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

func (ctrl User) serve(ctx context.Context, fn string, archive io.ReadSeeker, err error) (interface{}, error) {
	if err != nil {
		return nil, err
	}

	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Disposition", "attachment; filename="+fn)

		http.ServeContent(w, req, fn, time.Now(), archive)
	}, nil
}
