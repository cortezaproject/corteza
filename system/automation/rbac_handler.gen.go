package automation

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/automation/rbac_handler.yaml

import (
	"context"
	atypes "github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
	"github.com/cortezaproject/corteza-server/system/types"
)

var _ wfexec.ExecResponse

type (
	rbacHandlerRegistry interface {
		AddFunctions(ff ...*atypes.Function)
		Type(ref string) expr.Type
	}
)

func (h rbacHandler) register() {
	h.reg.AddFunctions(
		h.Allow(),
		h.Deny(),
		h.Inherit(),
		h.Check(),
	)
}

type (
	rbacAllowArgs struct {
		hasResource bool
		Resource    rbac.Resource

		hasRole    bool
		Role       interface{}
		roleID     uint64
		roleHandle string
		roleRes    *types.Role

		hasOperation bool
		Operation    string
	}
)

func (a rbacAllowArgs) GetRole() (bool, uint64, string, *types.Role) {
	return a.hasRole, a.roleID, a.roleHandle, a.roleRes
}

// Allow function RBAC: Allow operation on resource to a role
//
// expects implementation of allow function:
// func (h rbacHandler) allow(ctx context.Context, args *rbacAllowArgs) (err error) {
//    return
// }
func (h rbacHandler) Allow() *atypes.Function {
	return &atypes.Function{
		Ref:    "rbacAllow",
		Kind:   "function",
		Labels: map[string]string{"users": "rbac"},
		Meta: &atypes.FunctionMeta{
			Short: "RBAC: Allow operation on resource to a role",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "resource",
				Types: []string{"RbacResource"}, Required: true,
			},
			{
				Name:  "role",
				Types: []string{"ID", "Handle", "Role"}, Required: true,
			},
			{
				Name:  "operation",
				Types: []string{"String"}, Required: true,
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &rbacAllowArgs{
					hasResource:  in.Has("resource"),
					hasRole:      in.Has("role"),
					hasOperation: in.Has("operation"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			// Converting Role argument
			if args.hasRole {
				aux := expr.Must(expr.Select(in, "role"))
				switch aux.Type() {
				case h.reg.Type("ID").Type():
					args.roleID = aux.Get().(uint64)
				case h.reg.Type("Handle").Type():
					args.roleHandle = aux.Get().(string)
				case h.reg.Type("Role").Type():
					args.roleRes = aux.Get().(*types.Role)
				}
			}

			return out, h.allow(ctx, args)
		},
	}
}

type (
	rbacDenyArgs struct {
		hasResource bool
		Resource    rbac.Resource

		hasRole    bool
		Role       interface{}
		roleID     uint64
		roleHandle string
		roleRes    *types.Role

		hasOperation bool
		Operation    string
	}
)

func (a rbacDenyArgs) GetRole() (bool, uint64, string, *types.Role) {
	return a.hasRole, a.roleID, a.roleHandle, a.roleRes
}

// Deny function RBAC: Deny operation on resource to a role
//
// expects implementation of deny function:
// func (h rbacHandler) deny(ctx context.Context, args *rbacDenyArgs) (err error) {
//    return
// }
func (h rbacHandler) Deny() *atypes.Function {
	return &atypes.Function{
		Ref:    "rbacDeny",
		Kind:   "function",
		Labels: map[string]string{"users": "rbac"},
		Meta: &atypes.FunctionMeta{
			Short: "RBAC: Deny operation on resource to a role",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "resource",
				Types: []string{"RbacResource"}, Required: true,
			},
			{
				Name:  "role",
				Types: []string{"ID", "Handle", "Role"}, Required: true,
			},
			{
				Name:  "operation",
				Types: []string{"String"}, Required: true,
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &rbacDenyArgs{
					hasResource:  in.Has("resource"),
					hasRole:      in.Has("role"),
					hasOperation: in.Has("operation"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			// Converting Role argument
			if args.hasRole {
				aux := expr.Must(expr.Select(in, "role"))
				switch aux.Type() {
				case h.reg.Type("ID").Type():
					args.roleID = aux.Get().(uint64)
				case h.reg.Type("Handle").Type():
					args.roleHandle = aux.Get().(string)
				case h.reg.Type("Role").Type():
					args.roleRes = aux.Get().(*types.Role)
				}
			}

			return out, h.deny(ctx, args)
		},
	}
}

type (
	rbacInheritArgs struct {
		hasResource bool
		Resource    rbac.Resource

		hasRole    bool
		Role       interface{}
		roleID     uint64
		roleHandle string
		roleRes    *types.Role

		hasOperation bool
		Operation    string
	}
)

func (a rbacInheritArgs) GetRole() (bool, uint64, string, *types.Role) {
	return a.hasRole, a.roleID, a.roleHandle, a.roleRes
}

// Inherit function RBAC: Remove allow/deny operation of a role from resource
//
// expects implementation of inherit function:
// func (h rbacHandler) inherit(ctx context.Context, args *rbacInheritArgs) (err error) {
//    return
// }
func (h rbacHandler) Inherit() *atypes.Function {
	return &atypes.Function{
		Ref:    "rbacInherit",
		Kind:   "function",
		Labels: map[string]string{"users": "rbac"},
		Meta: &atypes.FunctionMeta{
			Short: "RBAC: Remove allow/deny operation of a role from resource",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "resource",
				Types: []string{"RbacResource"}, Required: true,
			},
			{
				Name:  "role",
				Types: []string{"ID", "Handle", "Role"}, Required: true,
			},
			{
				Name:  "operation",
				Types: []string{"String"}, Required: true,
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &rbacInheritArgs{
					hasResource:  in.Has("resource"),
					hasRole:      in.Has("role"),
					hasOperation: in.Has("operation"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			// Converting Role argument
			if args.hasRole {
				aux := expr.Must(expr.Select(in, "role"))
				switch aux.Type() {
				case h.reg.Type("ID").Type():
					args.roleID = aux.Get().(uint64)
				case h.reg.Type("Handle").Type():
					args.roleHandle = aux.Get().(string)
				case h.reg.Type("Role").Type():
					args.roleRes = aux.Get().(*types.Role)
				}
			}

			return out, h.inherit(ctx, args)
		},
	}
}

type (
	rbacCheckArgs struct {
		hasResource bool
		Resource    rbac.Resource

		hasOperation bool
		Operation    string

		hasUser bool
		User    *types.User
	}

	rbacCheckResults struct {
		Can bool
	}
)

// Check function RBAC: Can user perform an operation on a resource
//
// expects implementation of check function:
// func (h rbacHandler) check(ctx context.Context, args *rbacCheckArgs) (results *rbacCheckResults, err error) {
//    return
// }
func (h rbacHandler) Check() *atypes.Function {
	return &atypes.Function{
		Ref:    "rbacCheck",
		Kind:   "function",
		Labels: map[string]string{"users": "rbac"},
		Meta: &atypes.FunctionMeta{
			Short: "RBAC: Can user perform an operation on a resource",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "resource",
				Types: []string{"RbacResource"}, Required: true,
			},
			{
				Name:  "operation",
				Types: []string{"String"}, Required: true,
			},
			{
				Name:  "user",
				Types: []string{"User"},
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "can",
				Types: []string{"Boolean"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &rbacCheckArgs{
					hasResource:  in.Has("resource"),
					hasOperation: in.Has("operation"),
					hasUser:      in.Has("user"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			var results *rbacCheckResults
			if results, err = h.check(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Can (bool) to Boolean
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("Boolean").Cast(results.Can); err != nil {
					return
				} else if err = expr.Assign(out, "can", tval); err != nil {
					return
				}
			}

			return
		},
	}
}
