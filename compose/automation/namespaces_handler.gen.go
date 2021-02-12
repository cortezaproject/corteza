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
		hasNamespace    bool
		Namespace       interface{}
		namespaceID     uint64
		namespaceHandle string
		namespaceRes    *types.Namespace
	}

	namespacesLookupResults struct {
		Namespace *types.Namespace
	}
)

func (a namespacesLookupArgs) GetNamespace() (bool, uint64, string, *types.Namespace) {
	return a.hasNamespace, a.namespaceID, a.namespaceHandle, a.namespaceRes
}

// Lookup function Lookup for compose namespace by ID
//
// expects implementation of lookup function:
// func (h namespacesHandler) lookup(ctx context.Context, args *namespacesLookupArgs) (results *namespacesLookupResults, err error) {
//    return
// }
func (h namespacesHandler) Lookup() *atypes.Function {
	return &atypes.Function{
		Ref:    "composeNamespacesLookup",
		Kind:   "function",
		Labels: map[string]string{"compose": "step,workflow", "lookup": "step", "namespace": "step,workflow"},
		Meta: &atypes.FunctionMeta{
			Short: "Lookup for compose namespace by ID",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "namespace",
				Types: []string{"ID", "Handle", "ComposeNamespace"}, Required: true,
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "namespace",
				Types: []string{"ComposeNamespace"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &namespacesLookupArgs{
					hasNamespace: in.Has("namespace"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
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

			var results *namespacesLookupResults
			if results, err = h.lookup(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Namespace (*types.Namespace) to ComposeNamespace
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("ComposeNamespace").Cast(results.Namespace); err != nil {
					return
				} else if err = expr.Assign(out, "namespace", tval); err != nil {
					return
				}
			}

			return
		},
	}
}
