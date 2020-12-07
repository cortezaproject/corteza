package scim

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/go-chi/chi"
	"io"
	"net/http"
	"regexp"
	"strconv"
)

type (
	groupsHandler struct {
		externalIdAsPrimary bool
		externalIdValidator *regexp.Regexp

		svc service.RoleService
		sec getSecurityContextFn
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
		ctx = h.sec(r)
	)

	if u, code, err := h.createFromJSON(ctx, r.Body); err != nil {
		sendError(w, newErrorResonse(code, err))
	} else {
		send(w, http.StatusCreated, newGroupResourceResponse(u))
	}
}

func (h groupsHandler) createFromJSON(ctx context.Context, j io.Reader) (res *types.Role, code int, err error) {
	var (
		svc     = h.svc.With(ctx)
		payload = &groupResourceRequest{}
	)

	code = http.StatusBadRequest
	if err = payload.decodeJSON(j); err != nil {
	}

	// do we need to upsert?
	if payload.ExternalId != nil {
		res, code, err = h.lookupByExternalId(ctx, *payload.ExternalId)
		if err != nil && code != http.StatusNotFound {
			return
		}
	} else if payload.Name != nil {
		res, err = svc.FindByName(*payload.Name)
		if err != nil && !errors.Is(err, service.RoleErrNotFound()) {
			return nil, http.StatusInternalServerError, err
		}
	}

	if res == nil || res.ID == 0 {
		// in case when we did not find a valid group,
		// start from blank
		res = &types.Role{}
	}

	payload.applyTo(res)

	if res.ID > 0 {
		res, err = svc.Update(res)
	} else {
		res, err = svc.Create(res)
	}

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return res, 0, nil
}

func (h groupsHandler) replace(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var (
		ctx      = h.sec(r)
		existing = h.lookup(ctx, chi.URLParam(r, "id"), w)
	)

	if existing == nil {
		return
	}

	if res, err := h.updateFromJSON(ctx, existing, r.Body); err != nil {
		sendError(w, newErrorResonse(http.StatusBadRequest, err))
	} else {
		send(w, http.StatusOK, newGroupResourceResponse(res))
	}
}

func (h groupsHandler) updateFromJSON(ctx context.Context, res *types.Role, j io.Reader) (*types.Role, error) {
	var (
		payload = &groupResourceRequest{}
	)

	if err := payload.decodeJSON(j); err != nil {
		return nil, err
	}

	payload.applyTo(res)

	return h.svc.With(ctx).Update(res)
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
		sendError(w, newErrorResonse(http.StatusBadRequest, err))
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
		role, code, err := h.lookupByExternalId(ctx, id)
		if err != nil {
			sendError(w, newErrorResonse(code, err))
			return nil
		}

		return role
	} else {
		resId, err := strconv.ParseUint(id, 10, 64)
		if err != nil || resId == 0 {
			sendError(w, newErrorResonse(http.StatusBadRequest, err))
			return nil
		}

		role, err := svc.FindByID(resId)
		if err != nil {
			sendError(w, newErrorResonse(http.StatusBadRequest, err))
			return nil
		}

		return role
	}
}

func (h groupsHandler) lookupByExternalId(ctx context.Context, id string) (r *types.Role, code int, err error) {
	if h.externalIdValidator != nil && !h.externalIdValidator.MatchString(id) {
		return nil, http.StatusBadRequest, fmt.Errorf("invalid external ID")
	}

	rr, _, err := h.svc.With(ctx).Find(types.RoleFilter{Labels: map[string]string{groupLabel_SCIM_externalId: id}})
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	switch len(rr) {
	case 0:
		return nil, http.StatusNotFound, fmt.Errorf("role not found")
	case 1:
		return rr[0], 0, nil
	default:
		return nil, http.StatusPreconditionFailed, fmt.Errorf("more than one role matches this externalId")
	}
}
