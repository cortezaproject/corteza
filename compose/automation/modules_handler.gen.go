package automation

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// compose/automation/modules_handler.yaml

import (
	"context"
	atypes "github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
)

var _ wfexec.ExecResponse

type (
	modulesHandlerRegistry interface {
		AddFunctions(ff ...*atypes.Function)
		Type(ref string) expr.Type
	}
)

func (h modulesHandler) register() {
	h.reg.AddFunctions(
		h.Lookup(),
	)
}

type (
	modulesLookupArgs struct {
		hasModule    bool
		Module       interface{}
		moduleID     uint64
		moduleHandle string
		moduleRes    *types.Module

		hasNamespace    bool
		Namespace       interface{}
		namespaceID     uint64
		namespaceHandle string
		namespaceRes    *types.Namespace
	}

	modulesLookupResults struct {
		Module *types.Module
	}
)

func (a modulesLookupArgs) GetModule() (bool, uint64, string, *types.Module) {
	return a.hasModule, a.moduleID, a.moduleHandle, a.moduleRes
}

func (a modulesLookupArgs) GetNamespace() (bool, uint64, string, *types.Namespace) {
	return a.hasNamespace, a.namespaceID, a.namespaceHandle, a.namespaceRes
}

// Lookup function Lookup for compose Lookup by ID
//
// expects implementation of lookup function:
// func (h modulesHandler) lookup(ctx context.Context, args *modulesLookupArgs) (results *modulesLookupResults, err error) {
//    return
// }
func (h modulesHandler) Lookup() *atypes.Function {
	return &atypes.Function{
		Ref:    "composeModulesLookup",
		Kind:   "function",
		Labels: map[string]string{"compose": "step,workflow", "lookup": "step", "module": "step,workflow"},
		Meta: &atypes.FunctionMeta{
			Short: "Lookup for compose Lookup by ID",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "module",
				Types: []string{"ID", "Handle", "ComposeModule"}, Required: true,
			},
			{
				Name:  "namespace",
				Types: []string{"ID", "Handle", "ComposeNamespace"}, Required: true,
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "module",
				Types: []string{"ComposeModule"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &modulesLookupArgs{
					hasModule:    in.Has("module"),
					hasNamespace: in.Has("namespace"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			// Converting Module argument
			if args.hasModule {
				aux := expr.Must(expr.Select(in, "module"))
				switch aux.Type() {
				case h.reg.Type("ID").Type():
					args.moduleID = aux.Get().(uint64)
				case h.reg.Type("Handle").Type():
					args.moduleHandle = aux.Get().(string)
				case h.reg.Type("ComposeModule").Type():
					args.moduleRes = aux.Get().(*types.Module)
				}
			}

			// Converting Namespace argument
			if args.hasNamespace {
				aux := expr.Must(expr.Select(in, "namespace"))
				switch aux.Type() {
				case h.reg.Type("ID").Type():
					args.namespaceID = aux.Get().(uint64)
				case h.reg.Type("Handle").Type():
					args.namespaceHandle = aux.Get().(string)
				case h.reg.Type("ComposeNamespace").Type():
					args.namespaceRes = aux.Get().(*types.Namespace)
				}
			}

			var results *modulesLookupResults
			if results, err = h.lookup(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Module (*types.Module) to ComposeModule
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("ComposeModule").Cast(results.Module); err != nil {
					return
				} else if err = expr.Assign(out, "module", tval); err != nil {
					return
				}
			}

			return
		},
	}
}
