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
	"github.com/cortezaproject/corteza-server/pkg/label"
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
		h.LookupByID(),
		h.Save(),
		h.Create(),
		h.Update(),
		h.Delete(),
	)
}

type (
	modulesLookupByIDArgs struct {
		hasRecordID bool
		RecordID    uint64

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

	modulesLookupByIDResults struct {
		Record *types.Record
	}
)

func (a modulesLookupByIDArgs) GetModule() (bool, uint64, string, *types.Module) {
	return a.hasModule, a.moduleID, a.moduleHandle, a.moduleRes
}

func (a modulesLookupByIDArgs) GetNamespace() (bool, uint64, string, *types.Namespace) {
	return a.hasNamespace, a.namespaceID, a.namespaceHandle, a.namespaceRes
}

// LookupByID function Lookup for compose record by ID
//
// expects implementation of lookupByID function:
// func (h modulesHandler) lookupByID(ctx context.Context, args *modulesLookupByIDArgs) (results *modulesLookupByIDResults, err error) {
//    return
// }
func (h modulesHandler) LookupByID() *atypes.Function {
	return &atypes.Function{
		Ref:  "composeModulesLookupByID",
		Kind: "function",
		Meta: &atypes.FunctionMeta{
			Short: "Lookup for compose record by ID",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "recordID",
				Types: []string{"ID"}, Required: true,
			},
			{
				Name:  "module",
				Types: []string{"ID", "String", "ComposeModule"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label:       "Module to set record type",
					Description: "Even with unique record ID across all modules, module needs to be known\nbefore doing any record operations. Mainly because records of different\nmodules can be located in different stores.",
				},
			},
			{
				Name:  "namespace",
				Types: []string{"ID", "String", "ComposeNamespace"}, Required: true,
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "record",
				Types: []string{"ComposeRecord"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &modulesLookupByIDArgs{
					hasRecordID:  in.Has("recordID"),
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
				case h.reg.Type("String").Type():
					args.namespaceHandle = aux.Get().(string)
				case h.reg.Type("ComposeNamespace").Type():
					args.namespaceRes = aux.Get().(*types.Namespace)
				}
			}

			var results *modulesLookupByIDResults
			if results, err = h.lookupByID(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Record (*types.Record) to ComposeRecord
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("ComposeRecord").Cast(results.Record); err != nil {
					return
				} else if err = expr.Assign(out, "record", tval); err != nil {
					return
				}
			}

			return
		},
	}
}

type (
	modulesSaveArgs struct {
		hasRecord bool
		Record    *types.Record
	}

	modulesSaveResults struct {
		Record *types.Record
	}
)

// Save function Save record
//
// expects implementation of save function:
// func (h modulesHandler) save(ctx context.Context, args *modulesSaveArgs) (results *modulesSaveResults, err error) {
//    return
// }
func (h modulesHandler) Save() *atypes.Function {
	return &atypes.Function{
		Ref:  "composeModulesSave",
		Kind: "function",
		Meta: &atypes.FunctionMeta{
			Short: "Save record",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "record",
				Types: []string{"ComposeRecord"}, Required: true,
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "record",
				Types: []string{"ComposeRecord"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &modulesSaveArgs{
					hasRecord: in.Has("record"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			var results *modulesSaveResults
			if results, err = h.save(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Record (*types.Record) to ComposeRecord
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("ComposeRecord").Cast(results.Record); err != nil {
					return
				} else if err = expr.Assign(out, "record", tval); err != nil {
					return
				}
			}

			return
		},
	}
}

type (
	modulesCreateArgs struct {
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

		hasValues bool
		Values    types.RecordValueSet

		hasLabels bool
		Labels    label.Labels

		hasOwnedBy bool
		OwnedBy    uint64
	}

	modulesCreateResults struct {
		Record *types.Record
	}
)

func (a modulesCreateArgs) GetModule() (bool, uint64, string, *types.Module) {
	return a.hasModule, a.moduleID, a.moduleHandle, a.moduleRes
}

func (a modulesCreateArgs) GetNamespace() (bool, uint64, string, *types.Namespace) {
	return a.hasNamespace, a.namespaceID, a.namespaceHandle, a.namespaceRes
}

// Create function Creates and stores a new record
//
// expects implementation of create function:
// func (h modulesHandler) create(ctx context.Context, args *modulesCreateArgs) (results *modulesCreateResults, err error) {
//    return
// }
func (h modulesHandler) Create() *atypes.Function {
	return &atypes.Function{
		Ref:  "composeModulesCreate",
		Kind: "function",
		Meta: &atypes.FunctionMeta{
			Short: "Creates and stores a new record",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "module",
				Types: []string{"ID", "String", "ComposeModule"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label:       "Module to set record type",
					Description: "Even with unique record ID across all modules, module needs to be known\nbefore doing any record operations. Mainly because records of different\nmodules can be located in different stores.",
				},
			},
			{
				Name:  "namespace",
				Types: []string{"ID", "String", "ComposeNamespace"}, Required: true,
			},
			{
				Name:  "values",
				Types: []string{"KV"},
			},
			{
				Name:  "labels",
				Types: []string{"KV"},
			},
			{
				Name:  "ownedBy",
				Types: []string{"ID"},
				Meta: &atypes.ParamMeta{
					Label:  "Record owner",
					Visual: map[string]interface{}{"ref": "users"},
				},
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "record",
				Types: []string{"ComposeRecord"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &modulesCreateArgs{
					hasModule:    in.Has("module"),
					hasNamespace: in.Has("namespace"),
					hasValues:    in.Has("values"),
					hasLabels:    in.Has("labels"),
					hasOwnedBy:   in.Has("ownedBy"),
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
				case h.reg.Type("String").Type():
					args.namespaceHandle = aux.Get().(string)
				case h.reg.Type("ComposeNamespace").Type():
					args.namespaceRes = aux.Get().(*types.Namespace)
				}
			}

			var results *modulesCreateResults
			if results, err = h.create(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Record (*types.Record) to ComposeRecord
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("ComposeRecord").Cast(results.Record); err != nil {
					return
				} else if err = expr.Assign(out, "record", tval); err != nil {
					return
				}
			}

			return
		},
	}
}

type (
	modulesUpdateArgs struct {
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

		hasValues bool
		Values    types.RecordValueSet

		hasLabels bool
		Labels    label.Labels

		hasOwnedBy bool
		OwnedBy    uint64
	}

	modulesUpdateResults struct {
		Record *types.Record
	}
)

func (a modulesUpdateArgs) GetModule() (bool, uint64, string, *types.Module) {
	return a.hasModule, a.moduleID, a.moduleHandle, a.moduleRes
}

func (a modulesUpdateArgs) GetNamespace() (bool, uint64, string, *types.Namespace) {
	return a.hasNamespace, a.namespaceID, a.namespaceHandle, a.namespaceRes
}

// Update function Updates an existing record
//
// expects implementation of update function:
// func (h modulesHandler) update(ctx context.Context, args *modulesUpdateArgs) (results *modulesUpdateResults, err error) {
//    return
// }
func (h modulesHandler) Update() *atypes.Function {
	return &atypes.Function{
		Ref:  "composeModulesUpdate",
		Kind: "function",
		Meta: &atypes.FunctionMeta{
			Short: "Updates an existing record",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "module",
				Types: []string{"ID", "String", "ComposeModule"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label:       "Module to set record type",
					Description: "Even with unique record ID across all modules, module needs to be known\nbefore doing any record operations. Mainly because records of different\nmodules can be located in different stores.",
				},
			},
			{
				Name:  "namespace",
				Types: []string{"ID", "String", "ComposeNamespace"}, Required: true,
			},
			{
				Name:  "values",
				Types: []string{"KV"},
			},
			{
				Name:  "labels",
				Types: []string{"KV"},
			},
			{
				Name:  "ownedBy",
				Types: []string{"ID"},
				Meta: &atypes.ParamMeta{
					Label:  "Record owner",
					Visual: map[string]interface{}{"ref": "users"},
				},
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "record",
				Types: []string{"ComposeRecord"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &modulesUpdateArgs{
					hasModule:    in.Has("module"),
					hasNamespace: in.Has("namespace"),
					hasValues:    in.Has("values"),
					hasLabels:    in.Has("labels"),
					hasOwnedBy:   in.Has("ownedBy"),
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
				case h.reg.Type("String").Type():
					args.namespaceHandle = aux.Get().(string)
				case h.reg.Type("ComposeNamespace").Type():
					args.namespaceRes = aux.Get().(*types.Namespace)
				}
			}

			var results *modulesUpdateResults
			if results, err = h.update(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Record (*types.Record) to ComposeRecord
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("ComposeRecord").Cast(results.Record); err != nil {
					return
				} else if err = expr.Assign(out, "record", tval); err != nil {
					return
				}
			}

			return
		},
	}
}

type (
	modulesDeleteArgs struct {
		hasRecordID bool
		RecordID    uint64

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
)

func (a modulesDeleteArgs) GetModule() (bool, uint64, string, *types.Module) {
	return a.hasModule, a.moduleID, a.moduleHandle, a.moduleRes
}

func (a modulesDeleteArgs) GetNamespace() (bool, uint64, string, *types.Namespace) {
	return a.hasNamespace, a.namespaceID, a.namespaceHandle, a.namespaceRes
}

// Delete function Soft deletes compose record by ID
//
// expects implementation of delete function:
// func (h modulesHandler) delete(ctx context.Context, args *modulesDeleteArgs) (err error) {
//    return
// }
func (h modulesHandler) Delete() *atypes.Function {
	return &atypes.Function{
		Ref:  "composeModulesDelete",
		Kind: "function",
		Meta: &atypes.FunctionMeta{
			Short: "Soft deletes compose record by ID",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "recordID",
				Types: []string{"ID"}, Required: true,
			},
			{
				Name:  "module",
				Types: []string{"ID", "String", "ComposeModule"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label:       "Module to set record type",
					Description: "Even with unique record ID across all modules, module needs to be known\nbefore doing any record operations. Mainly because records of different\nmodules can be located in different stores.",
				},
			},
			{
				Name:  "namespace",
				Types: []string{"ID", "String", "ComposeNamespace"}, Required: true,
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &modulesDeleteArgs{
					hasRecordID:  in.Has("recordID"),
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
				case h.reg.Type("String").Type():
					args.namespaceHandle = aux.Get().(string)
				case h.reg.Type("ComposeNamespace").Type():
					args.namespaceRes = aux.Get().(*types.Namespace)
				}
			}

			return out, h.delete(ctx, args)
		},
	}
}
