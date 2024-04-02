package rest

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/gatekeep"
	"github.com/cortezaproject/corteza/server/system/rest/request"
	"github.com/cortezaproject/corteza/server/system/service"
)

type (
	Gatekeep struct {
		gatekeep gatekeepSvc
	}

	gatekeepSvc interface {
		Lock(ctx context.Context, c gatekeep.Constraint) (l gatekeep.Lock, err error)
		Unlock(ctx context.Context, c gatekeep.Constraint) (err error)
		Check(ctx context.Context, c gatekeep.Constraint, ref uint64) (state gatekeep.LockState, err error)
	}

	gatekeepPayload struct {
		LockID   uint64             `json:"lockID,string"`
		Resource string             `json:"resource,omitempty"`
		State    gatekeep.LockState `json:"state"`
	}
)

func (Gatekeep) New() *Gatekeep {
	return &Gatekeep{
		gatekeep: service.DefaultGatekeep,
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

	l, err := ctrl.gatekeep.Lock(ctx, c)
	if err != nil {
		return
	}

	return gatekeepPayload{
		LockID:   l.ID,
		Resource: r.Resource,
		State:    l.State,
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

func (ctrl Gatekeep) Check(ctx context.Context, l *request.GatekeepCheck) (out interface{}, err error) {
	c := gatekeep.Constraint{
		Resource:  l.Resource,
		Operation: gatekeep.OpWrite,
		UserID:    l.UserID,
	}

	if c.UserID == 0 {
		c.UserID = auth.GetIdentityFromContext(ctx).Identity()
	}

	state, err := ctrl.gatekeep.Check(ctx, c, l.LockID)
	if err != nil {
		return
	}

	return gatekeepPayload{
		LockID:   l.LockID,
		Resource: l.Resource,
		State:    state,
	}, nil
}
