package automation

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// compose/automation/namespaces_handler.yaml

import (
	"context"
	atypes "github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
)

var _ wfexec.ExecResponse

type (
	namespacesHandlerRegistry interface {
		AddFunctions(ff ...*atypes.Function)
		Type(ref string) expr.Type
	}
)

func (h namespacesHandler) register() {
	h.reg.AddFunctions(
		h.Lookup(),
	)
}

type (
	namespacesLookupArgs struct {
		hasModule    bool
		Module       interface{}
		moduleID     uint64
		moduleHandle string

		hasNamespace    bool
		Namespace       interface{}
		namespaceID     uint64
		namespaceHandle string
		namespaceRes    *types.Namespace
	}

	namespacesLookupResults struct {
		Module *types.Module
	}
)

func (a namespacesLookupArgs) GetModule() (bool, uint64, string) {
	return a.hasModule, a.moduleID, a.moduleHandle
}

func (a namespacesLookupArgs) GetNamespace() (bool, uint64, string, *types.Namespace) {
	return a.hasNamespace, a.namespaceID, a.namespaceHandle, a.namespaceRes
}

// Lookup function Lookup for compose module by ID
//
// expects implementation of lookup function:
// func (h namespacesHandler) lookup(ctx context.Context, args *namespacesLookupArgs) (results *namespacesLookupResults, err error) {
//    return
// }
func (h namespacesHandler) Lookup() *atypes.Function {
	return &atypes.Function{
		Ref:  "composeNamespacesLookup",
		Kind: "function",
		Meta: &atypes.FunctionMeta{
			Short: "Lookup for compose module by ID",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "module",
				Types: []string{"ID", "String"}, Required: true,
			},
			{
				Name:  "namespace",
				Types: []string{"ID", "String", "ComposeNamespace"}, Required: true,
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
				args = &namespacesLookupArgs{
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
				case h.reg.Type("String").Type():
					args.moduleHandle = aux.Get().(string)
				}
			}

			// Converting Namespace argument
			if args.hasNamespace {
				aux := expr.Must(expr.Select(in, "namespace"))
				switch aux.Type() {
				case h.reg.Type("ID").Type():
					args.namespaceID = aux.Get().(uint64)
				case h.reg.Type("String").Type():
					args.namespaceHandle = aux.Get().(string)
				case h.reg.Type("ComposeNamespace").Type():
					args.namespaceRes = aux.Get().(*types.Namespace)
				}
			}

			var results *namespacesLookupResults
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
