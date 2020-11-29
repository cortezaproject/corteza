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
	usersHandler struct {
		svc service.UserService
	}
)

func (h usersHandler) get(w http.ResponseWriter, r *http.Request) {
	var (
		id, _ = strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
		ctx   = auth.SetSuperUserContext(r.Context())
		svc   = h.svc.With(ctx)
	)

	if id == 0 {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	if u, err := svc.FindByID(id); err != nil {
		sendError(w, newErrorResonse(http.StatusBadRequest, err))
		return
	} else {
		send(w, http.StatusOK, newUserResourceResponse(u))
	}

	w.WriteHeader(http.StatusOK)
}

func (h usersHandler) create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var (
		ctx = auth.SetSuperUserContext(r.Context())
	)

	if u, err := h.createFromJSON(ctx, r.Body); err != nil {
		sendError(w, newErrorResonse(http.StatusBadRequest, err))
	} else {
		send(w, http.StatusCreated, newUserResourceResponse(u))
	}
}

func (h usersHandler) createFromJSON(ctx context.Context, j io.Reader) (u *types.User, err error) {
	var (
		svc     = h.svc.With(ctx)
		payload = &userResourceRequest{}
	)

	if err = payload.decodeJSON(j); err != nil {
		return
	}

	// do we need to upsert?
	if email := payload.Emails.getFirst(); email != "" {
		u, err = svc.FindByEmail(email)
		if err != nil && !errors.Is(err, service.UserErrNotFound()) {
			return
		}
	}

	if u == nil || !u.Valid() {
		// in case when we did not find a valid user,
		// start from blank
		u = &types.User{}
	}

	payload.applyTo(u)

	if u.ID > 0 {
		return svc.Update(u)
	} else {
		return svc.Create(u)
	}
}

func (h usersHandler) replace(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var (
		ctx       = auth.SetSuperUserContext(r.Context())
		userID, _ = strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	)

	if u, err := h.updateFromJSON(ctx, userID, r.Body); err != nil {
		sendError(w, newErrorResonse(http.StatusBadRequest, err))
	} else {
		send(w, http.StatusOK, newUserResourceResponse(u))
	}
}

func (h usersHandler) updateFromJSON(ctx context.Context, id uint64, j io.Reader) (u *types.User, err error) {
	var (
		svc     = h.svc.With(ctx)
		payload = &userResourceRequest{}
	)

	if u, err = svc.FindByID(id); err != nil {
		return
	}

	if u == nil || !u.Valid() {
		return nil, fmt.Errorf("refusing to update invalid user")
	}

	if err = payload.decodeJSON(j); err != nil {
		return
	}

	payload.applyTo(u)

	return h.svc.With(ctx).Update(u)
}

func (h usersHandler) delete(w http.ResponseWriter, r *http.Request) {
	var (
		ctx       = auth.SetSuperUserContext(r.Context())
		userID, _ = strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
		svc       = h.svc.With(ctx)
	)

	if err := svc.Delete(userID); err != nil {
		sendError(w, newErrorResonse(http.StatusBadRequest, err))
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
