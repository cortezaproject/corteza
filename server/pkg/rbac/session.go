package rbac

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/auth"
)

type (
	// Security/RBAC session
	Session interface {
		// Identity of the subject
		Identity() uint64

		// Roles returns all subject's roles for the session
		Roles() []uint64

		// Context used for expr evaluation
		Context() context.Context
	}

	session struct {
		// identity
		id uint64

		// roles
		rr []uint64

		// context
		ctx context.Context
	}
)

func (s session) Identity() uint64         { return s.id }
func (s session) Roles() []uint64          { return s.rr }
func (s session) Context() context.Context { return s.ctx }

var _ Session = &session{}

func ContextToSession(ctx context.Context) *session {
	return NewSession(ctx, auth.GetIdentityFromContext(ctx))
}

func ParamsToSession(ctx context.Context, user uint64, roles ...uint64) *session {
	return &session{
		id:  user,
		rr:  roles,
		ctx: ctx,
	}
}

func NewSession(ctx context.Context, i auth.Identifiable) *session {
	return &session{
		id:  i.Identity(),
		rr:  i.Roles(),
		ctx: ctx,
	}
}
