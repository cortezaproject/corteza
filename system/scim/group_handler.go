package scim

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/go-chi/chi"
	"io"
	"net/http"
	"strconv"
)

type (
	groupsHandler struct {
		svc service.RoleService
	}
)

func (h groupsHandler) get(w http.ResponseWriter, r *http.Request) {
	var (
		id, _ = strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
		ctx   = auth.SetSuperUserContext(r.Context())
		svc   = h.svc.With(ctx)
	)

	if id == 0 {
		http.Error(w, "invalid group id", http.StatusBadRequest)
		return
	}

	if u, err := svc.FindByID(id); err != nil {
		sendError(w, newErrorResonse(http.StatusBadRequest, err))
	} else {
		send(w, http.StatusOK, newGroupResourceResponse(u))
	}
}

func (h groupsHandler) create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var (
		ctx = auth.SetSuperUserContext(r.Context())
	)

	if u, err := h.createFromJSON(ctx, r.Body); err != nil {
		sendError(w, newErrorResonse(http.StatusBadRequest, err))
	} else {
		send(w, http.StatusCreated, newGroupResourceResponse(u))
	}
}

func (h groupsHandler) createFromJSON(ctx context.Context, j io.Reader) (r *types.Role, err error) {
	var (
		svc     = h.svc.With(ctx)
		payload = &groupResourceRequest{}
	)

	if err = payload.decodeJSON(j); err != nil {
		return
	}

	// do we need to upsert?
	if *payload.Name != "" {
		r, err = svc.FindByName(*payload.Name)
		if err != nil && !errors.Is(err, service.RoleErrNotFound()) {
			return
		}
	}

	if r == nil || r.ID == 0 {
		// in case when we did not find a valid group,
		// start from blank
		r = &types.Role{}
	}

	payload.applyTo(r)

	if r.ID > 0 {
		return svc.Update(r)
	} else {
		return svc.Create(r)
	}
}

func (h groupsHandler) replace(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var (
		ctx        = auth.SetSuperUserContext(r.Context())
		groupID, _ = strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	)

	if u, err := h.updateFromJSON(ctx, groupID, r.Body); err != nil {
		sendError(w, newErrorResonse(http.StatusBadRequest, err))
	} else {
		send(w, http.StatusOK, newGroupResourceResponse(u))
	}
}

func (h groupsHandler) updateFromJSON(ctx context.Context, id uint64, j io.Reader) (r *types.Role, err error) {
	var (
		svc     = h.svc.With(ctx)
		payload = &groupResourceRequest{}
	)

	if r, err = svc.FindByID(id); err != nil {
		return
	}

	if r == nil {
		return nil, fmt.Errorf("refusing to update invalid group")
	}

	if err = payload.decodeJSON(j); err != nil {
		return
	}

	payload.applyTo(r)

	return h.svc.With(ctx).Update(r)
}

func (h groupsHandler) delete(w http.ResponseWriter, r *http.Request) {
	var (
		ctx        = auth.SetSuperUserContext(r.Context())
		groupID, _ = strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
		svc        = h.svc.With(ctx)
	)

	if err := svc.Delete(groupID); err != nil {
		sendError(w, newErrorResonse(http.StatusBadRequest, err))
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
