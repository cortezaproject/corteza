package automation

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// compose/automation/records_handler.yaml

import (
	"context"
	atypes "github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
)

var _ wfexec.ExecResponse

type (
	recordsHandlerRegistry interface {
		AddFunctions(ff ...*atypes.Function)
		Type(ref string) expr.Type
	}
)

func (h recordsHandler) register() {
	h.reg.AddFunctions(
		h.Lookup(),
		h.Search(),
		h.Each(),
		h.New(),
		h.Validate(),
		h.Create(),
		h.Update(),
		h.Delete(),
		h.Report(),
	)
}

type (
	recordsLookupArgs struct {
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

		hasRecord bool
		Record    interface{}
		recordID  uint64
		recordRes *types.Record
	}

	recordsLookupResults struct {
		Record *types.Record
	}
)

func (a recordsLookupArgs) GetModule() (bool, uint64, string, *types.Module) {
	return a.hasModule, a.moduleID, a.moduleHandle, a.moduleRes
}

func (a recordsLookupArgs) GetNamespace() (bool, uint64, string, *types.Namespace) {
	return a.hasNamespace, a.namespaceID, a.namespaceHandle, a.namespaceRes
}

func (a recordsLookupArgs) GetRecord() (bool, uint64, *types.Record) {
	return a.hasRecord, a.recordID, a.recordRes
}

// Lookup function Lookup for compose record by ID
//
// expects implementation of lookup function:
// func (h recordsHandler) lookup(ctx context.Context, args *recordsLookupArgs) (results *recordsLookupResults, err error) {
//    return
// }
func (h recordsHandler) Lookup() *atypes.Function {
	return &atypes.Function{
		Ref:    "composeRecordsLookup",
		Kind:   "function",
		Labels: map[string]string{"compose": "step,workflow", "record": "step,workflow"},
		Meta: &atypes.FunctionMeta{
			Short: "Lookup for compose record by ID",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "module",
				Types: []string{"ID", "Handle", "ComposeModule"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label:       "Module to set record type",
					Description: "Even with unique record ID across all modules, module needs to be known\nbefore doing any record operations. Mainly because records of different\nmodules can be located in different stores.",
				},
			},
			{
				Name:  "namespace",
				Types: []string{"ID", "Handle", "ComposeNamespace"}, Required: true,
			},
			{
				Name:  "record",
				Types: []string{"ID", "ComposeRecord"}, Required: true,
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
				args = &recordsLookupArgs{
					hasModule:    in.Has("module"),
					hasNamespace: in.Has("namespace"),
					hasRecord:    in.Has("record"),
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

			// Converting Record argument
			if args.hasRecord {
				aux := expr.Must(expr.Select(in, "record"))
				switch aux.Type() {
				case h.reg.Type("ID").Type():
					args.recordID = aux.Get().(uint64)
				case h.reg.Type("ComposeRecord").Type():
					args.recordRes = aux.Get().(*types.Record)
				}
			}

			var results *recordsLookupResults
			if results, err = h.lookup(ctx, args); err != nil {
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
	recordsSearchArgs struct {
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

		hasQuery bool
		Query    string

		hasLabels bool
		Labels    map[string]string

		hasDeleted bool
		Deleted    uint64

		hasSort bool
		Sort    string

		hasLimit bool
		Limit    uint64

		hasIncTotal bool
		IncTotal    bool

		hasIncPageNavigation bool
		IncPageNavigation    bool

		hasPageCursor bool
		PageCursor    string
	}

	recordsSearchResults struct {
		Records    []*types.Record
		Total      uint64
		PageCursor string
	}
)

func (a recordsSearchArgs) GetModule() (bool, uint64, string, *types.Module) {
	return a.hasModule, a.moduleID, a.moduleHandle, a.moduleRes
}

func (a recordsSearchArgs) GetNamespace() (bool, uint64, string, *types.Namespace) {
	return a.hasNamespace, a.namespaceID, a.namespaceHandle, a.namespaceRes
}

// Search function Searches for records and returns them
//
// expects implementation of search function:
// func (h recordsHandler) search(ctx context.Context, args *recordsSearchArgs) (results *recordsSearchResults, err error) {
//    return
// }
func (h recordsHandler) Search() *atypes.Function {
	return &atypes.Function{
		Ref:    "composeRecordsSearch",
		Kind:   "function",
		Labels: map[string]string{"compose": "step,workflow", "record": "step,workflow"},
		Meta: &atypes.FunctionMeta{
			Short: "Searches for records and returns them",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "module",
				Types: []string{"ID", "Handle", "ComposeModule"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label:       "Module to set record type",
					Description: "Even with unique record ID across all modules, module needs to be known\nbefore doing any record operations. Mainly because records of different\nmodules can be located in different stores.",
				},
			},
			{
				Name:  "namespace",
				Types: []string{"ID", "Handle", "ComposeNamespace"}, Required: true,
			},
			{
				Name:  "query",
				Types: []string{"String"},
			},
			{
				Name:  "labels",
				Types: []string{"KV"},
			},
			{
				Name:  "deleted",
				Types: []string{"UnsignedInteger"},
			},
			{
				Name:  "sort",
				Types: []string{"String"},
			},
			{
				Name:  "limit",
				Types: []string{"UnsignedInteger"},
			},
			{
				Name:  "incTotal",
				Types: []string{"Boolean"},
			},
			{
				Name:  "incPageNavigation",
				Types: []string{"Boolean"},
			},
			{
				Name:  "pageCursor",
				Types: []string{"String"},
			},
		},

		Results: []*atypes.Param{

			{
				Name:    "records",
				Types:   []string{"ComposeRecord"},
				IsArray: true,
			},

			{
				Name:  "total",
				Types: []string{"UnsignedInteger"},
				Meta: &atypes.ParamMeta{
					Label:       "Total records found",
					Description: "Total items that satisfy given conditions.\n\nNeeds to be explicitly requested with incTotal argument",
				},
			},

			{
				Name:  "pageCursor",
				Types: []string{"String"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &recordsSearchArgs{
					hasModule:            in.Has("module"),
					hasNamespace:         in.Has("namespace"),
					hasQuery:             in.Has("query"),
					hasLabels:            in.Has("labels"),
					hasDeleted:           in.Has("deleted"),
					hasSort:              in.Has("sort"),
					hasLimit:             in.Has("limit"),
					hasIncTotal:          in.Has("incTotal"),
					hasIncPageNavigation: in.Has("incPageNavigation"),
					hasPageCursor:        in.Has("pageCursor"),
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

			var results *recordsSearchResults
			if results, err = h.search(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Records (*types.Record) to Array (of ComposeRecord)
				var (
					tval expr.TypedValue
					tarr = make([]expr.TypedValue, len(results.Records))
				)

				for i := range results.Records {
					if tarr[i], err = h.reg.Type("ComposeRecord").Cast(results.Records[i]); err != nil {
						return
					}
				}

				if tval, err = expr.NewArray(tarr); err != nil {
					return
				} else if err = expr.Assign(out, "records", tval); err != nil {
					return
				}
			}

			{
				// converting results.Total (uint64) to UnsignedInteger
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("UnsignedInteger").Cast(results.Total); err != nil {
					return
				} else if err = expr.Assign(out, "total", tval); err != nil {
					return
				}
			}

			{
				// converting results.PageCursor (string) to String
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("String").Cast(results.PageCursor); err != nil {
					return
				} else if err = expr.Assign(out, "pageCursor", tval); err != nil {
					return
				}
			}

			return
		},
	}
}

type (
	recordsEachArgs struct {
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

		hasQuery bool
		Query    string

		hasLabels bool
		Labels    map[string]string

		hasDeleted bool
		Deleted    uint64

		hasSort bool
		Sort    string

		hasLimit bool
		Limit    uint64

		hasIncTotal bool
		IncTotal    bool

		hasIncPageNavigation bool
		IncPageNavigation    bool

		hasPageCursor bool
		PageCursor    string
	}

	recordsEachResults struct {
		Record *types.Record
		Index  uint64
		Total  uint64
	}
)

func (a recordsEachArgs) GetModule() (bool, uint64, string, *types.Module) {
	return a.hasModule, a.moduleID, a.moduleHandle, a.moduleRes
}

func (a recordsEachArgs) GetNamespace() (bool, uint64, string, *types.Namespace) {
	return a.hasNamespace, a.namespaceID, a.namespaceHandle, a.namespaceRes
}

// Each function Searches for records and iterates over results
//
// expects implementation of each function:
// func (h recordsHandler) each(ctx context.Context, args *recordsEachArgs) (results *recordsEachResults, err error) {
//    return
// }
func (h recordsHandler) Each() *atypes.Function {
	return &atypes.Function{
		Ref:    "composeRecordsEach",
		Kind:   "iterator",
		Labels: map[string]string{"compose": "step,workflow", "record": "step,workflow"},
		Meta: &atypes.FunctionMeta{
			Short: "Searches for records and iterates over results",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "module",
				Types: []string{"ID", "Handle", "ComposeModule"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label:       "Module to set record type",
					Description: "Even with unique record ID across all modules, module needs to be known\nbefore doing any record operations. Mainly because records of different\nmodules can be located in different stores.",
				},
			},
			{
				Name:  "namespace",
				Types: []string{"ID", "Handle", "ComposeNamespace"}, Required: true,
			},
			{
				Name:  "query",
				Types: []string{"String"},
			},
			{
				Name:  "labels",
				Types: []string{"KV"},
			},
			{
				Name:  "deleted",
				Types: []string{"UnsignedInteger"},
			},
			{
				Name:  "sort",
				Types: []string{"String"},
			},
			{
				Name:  "limit",
				Types: []string{"UnsignedInteger"},
			},
			{
				Name:  "incTotal",
				Types: []string{"Boolean"},
			},
			{
				Name:  "incPageNavigation",
				Types: []string{"Boolean"},
			},
			{
				Name:  "pageCursor",
				Types: []string{"String"},
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "record",
				Types: []string{"ComposeRecord"},
			},

			{
				Name:  "index",
				Types: []string{"UnsignedInteger"},
				Meta: &atypes.ParamMeta{
					Label:       "Iteration counter",
					Description: "Zero-based number iteration counter",
				},
			},

			{
				Name:  "total",
				Types: []string{"UnsignedInteger"},
				Meta: &atypes.ParamMeta{
					Label:       "Total records found",
					Description: "Total items that satisfy given conditions.\n\nNeeds to be explicitly requested with incTotal argument",
				},
			},
		},

		Iterator: func(ctx context.Context, in *expr.Vars) (out wfexec.IteratorHandler, err error) {
			var (
				args = &recordsEachArgs{
					hasModule:            in.Has("module"),
					hasNamespace:         in.Has("namespace"),
					hasQuery:             in.Has("query"),
					hasLabels:            in.Has("labels"),
					hasDeleted:           in.Has("deleted"),
					hasSort:              in.Has("sort"),
					hasLimit:             in.Has("limit"),
					hasIncTotal:          in.Has("incTotal"),
					hasIncPageNavigation: in.Has("incPageNavigation"),
					hasPageCursor:        in.Has("pageCursor"),
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

			return h.each(ctx, args)
		},
	}
}

type (
	recordsNewArgs struct {
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

	recordsNewResults struct {
		Record *types.Record
	}
)

func (a recordsNewArgs) GetModule() (bool, uint64, string, *types.Module) {
	return a.hasModule, a.moduleID, a.moduleHandle, a.moduleRes
}

func (a recordsNewArgs) GetNamespace() (bool, uint64, string, *types.Namespace) {
	return a.hasNamespace, a.namespaceID, a.namespaceHandle, a.namespaceRes
}

// New function Make a new record
//
// expects implementation of new function:
// func (h recordsHandler) new(ctx context.Context, args *recordsNewArgs) (results *recordsNewResults, err error) {
//    return
// }
func (h recordsHandler) New() *atypes.Function {
	return &atypes.Function{
		Ref:    "composeRecordsNew",
		Kind:   "function",
		Labels: map[string]string{"compose": "step,workflow", "record": "step,workflow"},
		Meta: &atypes.FunctionMeta{
			Short: "Make a new record",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "module",
				Types: []string{"ID", "Handle", "ComposeModule"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label:       "Module to set record type",
					Description: "Even with unique record ID across all modules, module needs to be known\nbefore doing any record operations. Mainly because records of different\nmodules can be located in different stores.",
				},
			},
			{
				Name:  "namespace",
				Types: []string{"ID", "Handle", "ComposeNamespace"}, Required: true,
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
				args = &recordsNewArgs{
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

			var results *recordsNewResults
			if results, err = h.new(ctx, args); err != nil {
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
	recordsValidateArgs struct {
		hasRecord bool
		Record    *types.Record
	}

	recordsValidateResults struct {
		Valid bool
	}
)

// Validate function Validate record
//
// expects implementation of validate function:
// func (h recordsHandler) validate(ctx context.Context, args *recordsValidateArgs) (results *recordsValidateResults, err error) {
//    return
// }
func (h recordsHandler) Validate() *atypes.Function {
	return &atypes.Function{
		Ref:    "composeRecordsValidate",
		Kind:   "function",
		Labels: map[string]string{"compose": "step,workflow", "record": "step,workflow"},
		Meta: &atypes.FunctionMeta{
			Short: "Validate record",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "record",
				Types: []string{"ComposeRecord"}, Required: true,
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "valid",
				Types: []string{"Boolean"},
				Meta: &atypes.ParamMeta{
					Label: "Set to true when record is valid",
				},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &recordsValidateArgs{
					hasRecord: in.Has("record"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			var results *recordsValidateResults
			if results, err = h.validate(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Valid (bool) to Boolean
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("Boolean").Cast(results.Valid); err != nil {
					return
				} else if err = expr.Assign(out, "valid", tval); err != nil {
					return
				}
			}

			return
		},
	}
}

type (
	recordsCreateArgs struct {
		hasRecord bool
		Record    *types.Record
	}

	recordsCreateResults struct {
		Record *types.Record
	}
)

// Create function Creates and stores a new record
//
// expects implementation of create function:
// func (h recordsHandler) create(ctx context.Context, args *recordsCreateArgs) (results *recordsCreateResults, err error) {
//    return
// }
func (h recordsHandler) Create() *atypes.Function {
	return &atypes.Function{
		Ref:    "composeRecordsCreate",
		Kind:   "function",
		Labels: map[string]string{"compose": "step,workflow", "create": "step", "record": "step,workflow"},
		Meta: &atypes.FunctionMeta{
			Short: "Creates and stores a new record",
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
				args = &recordsCreateArgs{
					hasRecord: in.Has("record"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			var results *recordsCreateResults
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
	recordsUpdateArgs struct {
		hasRecord bool
		Record    *types.Record
	}

	recordsUpdateResults struct {
		Record *types.Record
	}
)

// Update function Updates an existing record
//
// expects implementation of update function:
// func (h recordsHandler) update(ctx context.Context, args *recordsUpdateArgs) (results *recordsUpdateResults, err error) {
//    return
// }
func (h recordsHandler) Update() *atypes.Function {
	return &atypes.Function{
		Ref:    "composeRecordsUpdate",
		Kind:   "function",
		Labels: map[string]string{"compose": "step,workflow", "record": "step,workflow", "update": "step"},
		Meta: &atypes.FunctionMeta{
			Short: "Updates an existing record",
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
				args = &recordsUpdateArgs{
					hasRecord: in.Has("record"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			var results *recordsUpdateResults
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
	recordsDeleteArgs struct {
		hasRecord bool
		Record    interface{}
		recordID  uint64
		recordRes *types.Record
	}
)

func (a recordsDeleteArgs) GetRecord() (bool, uint64, *types.Record) {
	return a.hasRecord, a.recordID, a.recordRes
}

// Delete function Soft deletes compose record by ID
//
// expects implementation of delete function:
// func (h recordsHandler) delete(ctx context.Context, args *recordsDeleteArgs) (err error) {
//    return
// }
func (h recordsHandler) Delete() *atypes.Function {
	return &atypes.Function{
		Ref:    "composeRecordsDelete",
		Kind:   "function",
		Labels: map[string]string{"compose": "step,workflow", "delete": "step", "record": "step,workflow"},
		Meta: &atypes.FunctionMeta{
			Short: "Soft deletes compose record by ID",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "record",
				Types: []string{"ID", "ComposeRecord"}, Required: true,
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &recordsDeleteArgs{
					hasRecord: in.Has("record"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			// Converting Record argument
			if args.hasRecord {
				aux := expr.Must(expr.Select(in, "record"))
				switch aux.Type() {
				case h.reg.Type("ID").Type():
					args.recordID = aux.Get().(uint64)
				case h.reg.Type("ComposeRecord").Type():
					args.recordRes = aux.Get().(*types.Record)
				}
			}

			return out, h.delete(ctx, args)
		},
	}
}

type (
	recordsReportArgs struct {
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

		hasMetrics bool
		Metrics    string

		hasDimensons bool
		Dimensons    string

		hasFilter bool
		Filter    string
	}

	recordsReportResults struct {
		Report interface{}
	}
)

func (a recordsReportArgs) GetModule() (bool, uint64, string, *types.Module) {
	return a.hasModule, a.moduleID, a.moduleHandle, a.moduleRes
}

func (a recordsReportArgs) GetNamespace() (bool, uint64, string, *types.Namespace) {
	return a.hasNamespace, a.namespaceID, a.namespaceHandle, a.namespaceRes
}

// Report function Searches for records and returns them
//
// expects implementation of report function:
// func (h recordsHandler) report(ctx context.Context, args *recordsReportArgs) (results *recordsReportResults, err error) {
//    return
// }
func (h recordsHandler) Report() *atypes.Function {
	return &atypes.Function{
		Ref:    "composeRecordsReport",
		Kind:   "function",
		Labels: map[string]string{"compose": "step,workflow", "record": "step,workflow"},
		Meta: &atypes.FunctionMeta{
			Short: "Searches for records and returns them",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "module",
				Types: []string{"ID", "Handle", "ComposeModule"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label:       "Module to set record type",
					Description: "Even with unique record ID across all modules, module needs to be known\nbefore doing any record operations. Mainly because records of different\nmodules can be located in different stores.",
				},
			},
			{
				Name:  "namespace",
				Types: []string{"ID", "Handle", "ComposeNamespace"}, Required: true,
			},
			{
				Name:  "metrics",
				Types: []string{"String"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label:       "Metrics for records report",
					Description: "List of comma delimited expressions with aggregation functions (count, sum, min, avg)",
				},
			},
			{
				Name:  "dimensons",
				Types: []string{"String"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label:       "Dimensons for records report",
					Description: "List of comma delimited dimension (i.e.: group by) expressions or fields",
				},
			},
			{
				Name:  "filter",
				Types: []string{"String"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label:       "Filter for records report",
					Description: "Filter in CortezaQL format",
				},
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "report",
				Types: []string{"Any"},
				Meta: &atypes.ParamMeta{
					Label:       "Complex structure holding complete records report",
					Description: "Example of a result value\n[]map[string]interface{}{\n\t{\"count\": 3, \"dimension_0\": 1, \"metric_0\": 3},\n\t{\"count\": 2, \"dimension_0\": 2, \"metric_0\": 5},\n\t{\"count\": 1, \"dimension_0\": nil, \"metric_0\": nil},\n}",
				},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &recordsReportArgs{
					hasModule:    in.Has("module"),
					hasNamespace: in.Has("namespace"),
					hasMetrics:   in.Has("metrics"),
					hasDimensons: in.Has("dimensons"),
					hasFilter:    in.Has("filter"),
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

			var results *recordsReportResults
			if results, err = h.report(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Report (interface{}) to Any
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("Any").Cast(results.Report); err != nil {
					return
				} else if err = expr.Assign(out, "report", tval); err != nil {
					return
				}
			}

			return
		},
	}
}
