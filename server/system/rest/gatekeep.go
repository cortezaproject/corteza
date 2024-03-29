package rest

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/gatekeep"
	"github.com/cortezaproject/corteza/server/system/rest/request"
)

type (
	Gatekeep struct {
		gatekeep gatekeepSvc
	}

	gatekeepSvc interface {
		Lock(ctx context.Context, c gatekeep.Constraint) (ref uint64, state gatekeep.LockState, err error)
		Unlock(ctx context.Context, c gatekeep.Constraint) (err error)
		ProbeLock(ctx context.Context, c gatekeep.Constraint, ref uint64) (state gatekeep.LockState, err error)
	}

	gatekeepPayload struct {
		LockID   uint64             `json:"lockID,string"`
		Resource string             `json:"resource,omitempty"`
		State    gatekeep.LockState `json:"state"`
	}
)

func (Gatekeep) New() *Gatekeep {
	return &Gatekeep{
		gatekeep: gatekeep.Service(),
	}
}

func (ctrl Gatekeep) Lock(ctx context.Context, r *request.GatekeepLock) (out interface{}, err error) {
	c := gatekeep.Constraint{
		Resource:  r.Resource,
		Operation: gatekeep.OpWrite,
		UserID:    r.UserID,
	}

	if r.ExpiresIn != "" {
		c.ExpiresIn, err = time.ParseDuration(r.ExpiresIn)
		if err != nil {
			return
		}
	}

	if c.UserID == 0 {
		c.UserID = auth.GetIdentityFromContext(ctx).Identity()
	}

	ref, state, err := ctrl.gatekeep.Lock(ctx, c)
	if err != nil {
		return
	}

	return gatekeepPayload{
		LockID:   ref,
		Resource: r.Resource,
		State:    state,
	}, nil
}

func (ctrl Gatekeep) Unlock(ctx context.Context, r *request.GatekeepUnlock) (out interface{}, err error) {
	c := gatekeep.Constraint{
		Resource:  r.Resource,
		Operation: gatekeep.OpWrite,
		UserID:    r.UserID,
	}

	if c.UserID == 0 {
		c.UserID = auth.GetIdentityFromContext(ctx).Identity()
	}

	return api.OK(), ctrl.gatekeep.Unlock(ctx, c)
}

func (ctrl Gatekeep) Check(ctx context.Context, r *request.GatekeepCheck) (out interface{}, err error) {
	c := gatekeep.Constraint{
		Resource:  r.Resource,
		Operation: gatekeep.OpWrite,
		UserID:    r.UserID,
	}

	if c.UserID == 0 {
		c.UserID = auth.GetIdentityFromContext(ctx).Identity()
	}

	state, err := ctrl.gatekeep.ProbeLock(ctx, c, r.LockID)
	if err != nil {
		return
	}

	return gatekeepPayload{
		LockID:   r.LockID,
		Resource: r.Resource,
		State:    state,
	}, nil
}
