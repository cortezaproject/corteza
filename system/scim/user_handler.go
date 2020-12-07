package scim

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-chi/chi"
	"io"
	"net/http"
	"regexp"
	"strconv"
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
		ctx = h.sec(r)
	)

	if u, code, err := h.createFromJSON(ctx, r.Body); err != nil {
		sendError(w, newErrorResonse(code, err))
	} else {
		send(w, http.StatusCreated, newUserResourceResponse(u))
	}
}

func (h usersHandler) createFromJSON(ctx context.Context, j io.Reader) (res *types.User, code int, err error) {
	var (
		svc = h.svc.With(ctx)
		//roles   = h.rleSvc.With(ctx)
		payload = &userResourceRequest{}
	)

	code = http.StatusBadRequest
	if err = payload.decodeJSON(j); err != nil {
		return
	}

	// do we need to upsert?
	if payload.ExternalId != nil {
		res, code, err = h.lookupByExternalId(ctx, *payload.ExternalId)
		if err != nil && code != http.StatusNotFound {
			return
		}
	} else if email := payload.Emails.getFirst(); email != "" {
		res, err = svc.FindByEmail(email)
		if err != nil && !errors.Is(err, service.UserErrNotFound()) {
			return nil, http.StatusInternalServerError, err
		}
	}

	if res == nil || !res.Valid() {
		// in case when we did not find a valid user,
		// start from blank
		res = &types.User{}
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

	if payload.Password != nil && *payload.Password != "" {
		err = h.passSvc.SetPassword(ctx, res.ID, *payload.Password)
		if err != nil {
			return
		}
	}

	return res, 0, nil
}

func (h usersHandler) replace(w http.ResponseWriter, r *http.Request) {
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
		send(w, http.StatusOK, newUserResourceResponse(res))
	}
}

func (h usersHandler) updateFromJSON(ctx context.Context, res *types.User, j io.Reader) (*types.User, error) {
	var (
		payload = &userResourceRequest{}
	)

	if err := payload.decodeJSON(j); err != nil {
		return nil, err
	}

	payload.applyTo(res)

	return h.svc.With(ctx).Update(res)
}

func (h usersHandler) delete(w http.ResponseWriter, r *http.Request) {
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
func (h usersHandler) lookup(ctx context.Context, id string, w http.ResponseWriter) *types.User {
	var (
		svc = h.svc.With(ctx)
	)
	spew.Dump(h.externalIdAsPrimary)
	if h.externalIdAsPrimary {
		role, code, err := h.lookupByExternalId(ctx, id)
		if err != nil {
			sendError(w, newErrorResonse(code, err))
			return nil
		}

		return role
	} else {
		groupId, err := strconv.ParseUint(id, 10, 64)
		if err != nil || groupId == 0 {
			sendError(w, newErrorResonse(http.StatusBadRequest, err))
			return nil
		}

		role, err := svc.FindByID(groupId)
		if err != nil {
			sendError(w, newErrorResonse(http.StatusBadRequest, err))
			return nil
		}

		return role
	}
}

func (h usersHandler) lookupByExternalId(ctx context.Context, id string) (r *types.User, code int, err error) {
	spew.Dump(id)
	if h.externalIdValidator != nil && !h.externalIdValidator.MatchString(id) {
		return nil, http.StatusBadRequest, fmt.Errorf("invalid external ID")
	}

	rr, _, err := h.svc.With(ctx).Find(types.UserFilter{Labels: map[string]string{groupLabel_SCIM_externalId: id}})
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	switch len(rr) {
	case 0:
		return nil, http.StatusNotFound, fmt.Errorf("user not found")
	case 1:
		return rr[0], 0, nil
	default:
		return nil, http.StatusPreconditionFailed, fmt.Errorf("more than one user matches this externalId")
	}
}
