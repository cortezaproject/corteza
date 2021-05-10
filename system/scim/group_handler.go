package scim

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/go-chi/chi"
	"net/http"
	"regexp"
	"strconv"
)

type (
	groupsHandler struct {
		externalIdAsPrimary bool
		externalIdValidator *regexp.Regexp

		svc     service.RoleService
		userSvc service.UserService
		sec     getSecurityContextFn
	}
)

func (h groupsHandler) get(w http.ResponseWriter, r *http.Request) {
	var (
		res = h.lookup(h.sec(r), chi.URLParam(r, "id"), w)
	)

	if res == nil {
		return
	}

	send(w, http.StatusOK, newGroupResourceResponse(res))
}

func (h groupsHandler) create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var (
		ctx      = h.sec(r)
		svc      = h.svc.With(ctx)
		payload  = &groupResourceRequest{}
		err      error
		existing *types.Role
	)

	if err = payload.decodeJSON(r.Body); err != nil {
		sendError(w, newErrorResponse(http.StatusBadRequest, err))
		return
	}

	{
		// do we need to upsert?
		if payload.ExternalId != nil {
			existing, err = h.lookupByExternalId(ctx, *payload.ExternalId)
			if err != nil {
				sendError(w, err)
				return
			}
		} else if *payload.Name != "" {
			existing, err = svc.FindByName(*payload.Name)
			if err != nil && !errors.Is(err, service.RoleErrNotFound()) {
				sendError(w, err)
				return
			}
		}
	}

	res, err := h.save(ctx, payload, existing)
	if err != nil {
		sendError(w, err)
		return
	}

	status := http.StatusOK
	if res.UpdatedAt == nil {
		status = http.StatusCreated
	}

	send(w, status, newGroupResourceResponse(res))
}

func (h groupsHandler) replace(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var (
		ctx      = h.sec(r)
		existing = h.lookup(ctx, chi.URLParam(r, "id"), w)
		payload  = &groupResourceRequest{}
	)

	if err := payload.decodeJSON(r.Body); err != nil {
		sendError(w, newErrorResponse(http.StatusBadRequest, err))
		return
	}

	res, err := h.save(ctx, payload, existing)
	if err != nil {
		sendError(w, err)
		return
	}

	status := http.StatusOK
	if res.UpdatedAt == nil {
		status = http.StatusCreated
	}

	send(w, status, newGroupResourceResponse(res))
}

// patches group
//
// only supports adding and removing members
func (h groupsHandler) patch(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var (
		ctx     = h.sec(r)
		svc     = h.svc.With(ctx)
		res     = h.lookup(ctx, chi.URLParam(r, "id"), w)
		payload = &operationsRequest{}
	)

	if res == nil {
		return
	}

	if err := payload.decodeJSON(r.Body); err != nil {
		sendError(w, newErrorResponse(http.StatusBadRequest, err))
		return
	}

	var (
		u           *types.User
		ops         = make([]func() error, 0, len(payload.Operations))
		err         error
		memberships = make(map[uint64]bool)
	)

	{
		// collect all existing memberships into a simple map
		// to ensure we dont step on our feet (too much)
		//
		// this is not 100% bulletproof for concurrent modifications
		mm, _, err := store.SearchRoleMembers(ctx, service.DefaultStore, types.RoleMemberFilter{RoleID: res.ID})
		if err != nil {
			sendError(w, err)
			return
		}

		for _, m := range mm {
			memberships[m.UserID] = true
		}
	}

	// validate and collect operations
	for _, op := range payload.Operations {
		// iterate through operation's values, load user and schedule op
		for path, userExternalId := range op.Value {
			if path != "members" {
				// allow only "members" path
				sendError(w, newErrorfResponse(http.StatusBadRequest, "unsupported path: %q", op.Path))
				return
			}

			u, err = lookupUserByExternalId(ctx, h.userSvc, h.externalIdValidator, userExternalId)
			if err != nil {
				sendError(w, err)
				return
			}

			if u == nil {
				sendError(w, newErrorfResponse(http.StatusBadRequest, "no such user: %q", userExternalId))
				return
			}

			// making sure u is not overwritten
			// in the next iteration
			memberId := u.ID

			switch op.Operation {
			case patchOpAdd:
				// support for add operation,
				// check if there members already exist
				ops = append(ops, func() error {
					if memberships[memberId] {
						// already added
						return nil
					}

					memberships[memberId] = true
					return svc.MemberAdd(res.ID, memberId)
				})

			case patchOpRemove:
				// support for remove operation,
				// check if there members are missing
				ops = append(ops, func() error {
					if !memberships[memberId] {
						// already removed
						return nil
					}

					delete(memberships, memberId)
					return svc.MemberRemove(res.ID, memberId)
				})

			default:
				sendError(w, newErrorfResponse(http.StatusBadRequest, "unsupported operation: %q", op.Operation))
				return
			}
		}
	}

	// run all scheduled ops
	for _, op := range ops {
		if err = op(); err != nil {
			sendError(w, err)
			return
		}
	}

	send(w, http.StatusNoContent, nil)
}

func (h groupsHandler) save(ctx context.Context, req *groupResourceRequest, existing *types.Role) (res *types.Role, err error) {
	var (
		svc = h.svc.With(ctx)
	)

	if existing == nil {
		// in case when we did not find a valid group,
		// start from blank
		existing = &types.Role{}
	}

	res = existing
	req.applyTo(res)

	if res.ID > 0 {
		res, err = svc.Update(res)
	} else {
		res, err = svc.Create(res)
	}

	if err != nil {
		return nil, newErrorResponse(http.StatusInternalServerError, err)
	}

	return res, nil
}

func (h groupsHandler) delete(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = h.sec(r)
		svc = h.svc.With(ctx)
		res = h.lookup(ctx, chi.URLParam(r, "id"), w)
	)

	if res == nil {
		return
	}

	if err := svc.Delete(res.ID); err != nil {
		sendError(w, newErrorResponse(http.StatusBadRequest, err))
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

// loads role from request path params
//
// handles errors by writing them to response
func (h groupsHandler) lookup(ctx context.Context, id string, w http.ResponseWriter) *types.Role {
	var (
		svc = h.svc.With(ctx)
	)

	if h.externalIdAsPrimary {
		res, err := h.lookupByExternalId(ctx, id)
		if err != nil {
			sendError(w, err)
			return nil
		}

		if res == nil {
			sendError(w, newErrorResponse(http.StatusNotFound, fmt.Errorf("group not found")))
		}

		return res
	} else {
		id, err := strconv.ParseUint(id, 10, 64)
		if err != nil || id == 0 {
			sendError(w, newErrorResponse(http.StatusBadRequest, err))
			return nil
		}

		role, err := svc.FindByID(id)
		if err != nil {
			sendError(w, newErrorResponse(http.StatusBadRequest, err))
			return nil
		}

		return role
	}
}

func (h groupsHandler) lookupByExternalId(ctx context.Context, id string) (r *types.Role, err error) {
	if h.externalIdValidator != nil && !h.externalIdValidator.MatchString(id) {
		return nil, newErrorfResponse(http.StatusBadRequest, "invalid external ID")
	}

	rr, _, err := h.svc.With(ctx).Find(types.RoleFilter{Labels: map[string]string{groupLabel_SCIM_externalId: id}})
	if err != nil {
		return nil, newErrorResponse(http.StatusInternalServerError, err)
	}

	switch len(rr) {
	case 0:
		return nil, nil
	case 1:
		return rr[0], nil
	default:
		return nil, newErrorfResponse(http.StatusPreconditionFailed, "more than one group matches this externalId")
	}
}
