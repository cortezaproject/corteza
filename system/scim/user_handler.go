package scim

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/go-chi/chi/v5"
)

type (
	passwordSetter interface {
		SetPassword(context.Context, uint64, string) error
	}

	usersHandler struct {
		externalIdAsPrimary bool
		externalIdValidator *regexp.Regexp

		svc     service.UserService
		passSvc passwordSetter
		sec     getSecurityContextFn
	}
)

func (h usersHandler) get(w http.ResponseWriter, r *http.Request) {
	var (
		res = h.lookup(h.sec(r), chi.URLParam(r, "id"), w)
	)

	if res == nil {
		return
	}

	send(w, http.StatusOK, newUserResourceResponse(res))
}

func (h usersHandler) create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var (
		ctx      = h.sec(r)
		svc      = h.svc
		payload  = &userResourceRequest{}
		err      error
		existing *types.User
		code     = http.StatusBadRequest
	)

	if err = payload.decodeJSON(r.Body); err != nil {
		sendError(w, newErrorResponse(code, err))
		return
	}

	{
		// do we need to upsert?
		if payload.ExternalId != nil {
			existing, err = h.lookupByExternalId(ctx, *payload.ExternalId)
			if err != nil {
				sendError(w, newErrorResponse(code, err))
				return
			}
		} else if email := payload.Emails.getFirst(); email != "" {
			existing, err = svc.FindByEmail(ctx, email)
			if err != nil && !errors.Is(err, service.UserErrNotFound()) {
				sendError(w, newErrorResponse(http.StatusInternalServerError, err))
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

	send(w, status, newUserResourceResponse(res))
}

func (h usersHandler) replace(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var (
		ctx      = h.sec(r)
		existing = h.lookup(ctx, chi.URLParam(r, "id"), w)
		payload  = &userResourceRequest{}
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

	send(w, status, newUserResourceResponse(res))
}

func (h usersHandler) save(ctx context.Context, req *userResourceRequest, existing *types.User) (res *types.User, err error) {
	var (
		svc = h.svc
	)

	if existing == nil || !existing.Valid() {
		// in case when we did not find a valid user,
		// start from blank
		existing = &types.User{}
	}

	res = existing
	req.applyTo(res)

	if res.ID > 0 {
		res, err = svc.Update(ctx, res)
	} else {
		res, err = svc.Create(ctx, res)
	}

	if err != nil {
		return nil, err
	}

	if req.Password != nil && *req.Password != "" {
		err = h.passSvc.SetPassword(ctx, res.ID, *req.Password)
		if err != nil {
			return
		}
	}

	return res, nil
}

func (h usersHandler) delete(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = h.sec(r)
		svc = h.svc
		res = h.lookup(ctx, chi.URLParam(r, "id"), w)
	)

	if res == nil {
		return
	}

	if err := svc.Delete(ctx, res.ID); err != nil {
		sendError(w, newErrorResponse(http.StatusBadRequest, err))
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

// loads role from request path params
//
// handles errors by writing them to response
func (h usersHandler) lookup(ctx context.Context, id string, w http.ResponseWriter) *types.User {
	var (
		svc = h.svc
	)

	if h.externalIdAsPrimary {
		res, err := h.lookupByExternalId(ctx, id)
		if err != nil {
			sendError(w, err)
			return nil
		}

		if res == nil {
			sendError(w, newErrorResponse(http.StatusNotFound, fmt.Errorf("user not found")))
		}

		return res
	} else {
		groupId, err := strconv.ParseUint(id, 10, 64)
		if err != nil || groupId == 0 {
			sendError(w, newErrorResponse(http.StatusBadRequest, err))
			return nil
		}

		role, err := svc.FindByID(ctx, groupId)
		if err != nil {
			sendError(w, newErrorResponse(http.StatusBadRequest, err))
			return nil
		}

		return role
	}
}

func (h usersHandler) lookupByExternalId(ctx context.Context, id string) (r *types.User, err error) {
	return lookupUserByExternalId(ctx, h.svc, h.externalIdValidator, id)
}

func lookupUserByExternalId(ctx context.Context, svc service.UserService, v *regexp.Regexp, id string) (r *types.User, err error) {
	if v != nil && !v.MatchString(id) {
		return nil, newErrorfResponse(http.StatusBadRequest, "invalid external ID")
	}

	rr, _, err := svc.Find(ctx, types.UserFilter{Labels: map[string]string{userLabel_SCIM_externalId: id}})
	if err != nil {
		return nil, newErrorResponse(http.StatusInternalServerError, err)
	}

	switch len(rr) {
	case 0:
		return nil, nil
	case 1:
		return rr[0], nil
	default:
		return nil, newErrorfResponse(http.StatusPreconditionFailed, "more than one user matches this externalId")
	}
}
