package automation

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/rbac"
)

type (
	rbacService interface {
		Can(ses rbac.Session, op string, res rbac.Resource) bool
		Grant(ctx context.Context, rules ...*rbac.Rule) (err error)
	}

	rbacHandler struct {
		reg rbacHandlerRegistry
		svc rbacService
	}

	auxResource struct{ res string }
)

var _ rbac.Resource = &auxResource{}

func (r auxResource) RbacResource() string {
	return r.res
}

func RbacHandler(reg rbacHandlerRegistry, svc rbacService) *rbacHandler {
	h := &rbacHandler{
		reg: reg,
		svc: svc,
	}

	h.register()
	return h
}

func (h rbacHandler) allow(ctx context.Context, args *rbacAllowArgs) (err error) {
	return h.grant(ctx, rbac.AllowRule(args.roleID, args.Operation, args.Resource.RbacResource()))
}

func (h rbacHandler) deny(ctx context.Context, args *rbacDenyArgs) (err error) {
	return h.grant(ctx, rbac.DenyRule(args.roleID, args.Operation, args.Resource.RbacResource()))
}

func (h rbacHandler) inherit(ctx context.Context, args *rbacInheritArgs) (err error) {
	return h.grant(ctx, rbac.InheritRule(args.roleID, args.Operation, args.Resource.RbacResource()))
}

// verifies grant op (granter needs to be allowed to grant on the component

func (h rbacHandler) grant(ctx context.Context, r *rbac.Rule) (err error) {
	// we can safely create a fake component resource since there are
	// no extra metadata hidden under it,
	cmpRes := &auxResource{rbac.ResourceComponent(r.Resource)}

	if !h.svc.Can(rbac.ContextToSession(ctx), "grant", cmpRes) {
		return fmt.Errorf("not allowed to grant %s", r.String())
	}

	return h.svc.Grant(ctx, r)
}

func (h rbacHandler) check(ctx context.Context, args *rbacCheckArgs) (b *rbacCheckResults, err error) {
	b = &rbacCheckResults{}

	b.Can = h.svc.Can(rbac.ContextToSession(ctx), args.Operation, args.Resource)
	return
}
