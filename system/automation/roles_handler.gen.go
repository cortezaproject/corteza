package automation

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/automation/roles_handler.yaml

import (
	"context"
	atypes "github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
	"github.com/cortezaproject/corteza-server/system/types"
)

var _ wfexec.ExecResponse

type (
	rolesHandlerRegistry interface {
		AddFunctions(ff ...*atypes.Function)
		Type(ref string) expr.Type
	}
)

func (h rolesHandler) register() {
	h.reg.AddFunctions(
		h.Lookup(),
		h.Search(),
		h.Each(),
		h.Create(),
		h.Update(),
		h.Delete(),
		h.Recover(),
		h.Archive(),
		h.Unarchive(),
	)
}

type (
	rolesLookupArgs struct {
		hasLookup    bool
		Lookup       interface{}
		lookupID     uint64
		lookupHandle string
		lookupRes    *types.Role
	}

	rolesLookupResults struct {
		Role *types.Role
	}
)

func (a rolesLookupArgs) GetLookup() (bool, uint64, string, *types.Role) {
	return a.hasLookup, a.lookupID, a.lookupHandle, a.lookupRes
}

// Lookup function Looks-up for role by ID
//
// expects implementation of lookup function:
// func (h rolesHandler) lookup(ctx context.Context, args *rolesLookupArgs) (results *rolesLookupResults, err error) {
//    return
// }
func (h rolesHandler) Lookup() *atypes.Function {
	return &atypes.Function{
		Ref:  "rolesLookup",
		Kind: "function",
		Meta: &atypes.FunctionMeta{
			Short: "Looks-up for role by ID",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "lookup",
				Types: []string{"ID", "Handle", "Role"}, Required: true,
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "role",
				Types: []string{"Role"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &rolesLookupArgs{
					hasLookup: in.Has("lookup"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			// Converting Lookup argument
			if args.hasLookup {
				aux := expr.Must(expr.Select(in, "lookup"))
				switch aux.Type() {
				case h.reg.Type("ID").Type():
					args.lookupID = aux.Get().(uint64)
				case h.reg.Type("Handle").Type():
					args.lookupHandle = aux.Get().(string)
				case h.reg.Type("Role").Type():
					args.lookupRes = aux.Get().(*types.Role)
				}
			}

			var results *rolesLookupResults
			if results, err = h.lookup(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Role (*types.Role) to Role
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("Role").Cast(results.Role); err != nil {
					return
				} else if err = expr.Assign(out, "role", tval); err != nil {
					return
				}
			}

			return
		},
	}
}

type (
	rolesSearchArgs struct {
		hasQuery bool
		Query    string

		hasMemberID bool
		MemberID    uint64

		hasHandle bool
		Handle    string

		hasName bool
		Name    string

		hasLabels bool
		Labels    map[string]string

		hasDeleted bool
		Deleted    uint64

		hasArchived bool
		Archived    uint64

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

	rolesSearchResults struct {
		Roles      []*types.Role
		Total      uint64
		PageCursor string
	}
)

// Search function Searches for roles and returns them
//
// expects implementation of search function:
// func (h rolesHandler) search(ctx context.Context, args *rolesSearchArgs) (results *rolesSearchResults, err error) {
//    return
// }
func (h rolesHandler) Search() *atypes.Function {
	return &atypes.Function{
		Ref:  "rolesSearch",
		Kind: "function",
		Meta: &atypes.FunctionMeta{
			Short: "Searches for roles and returns them",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "query",
				Types: []string{"String"},
			},
			{
				Name:  "memberID",
				Types: []string{"ID"},
			},
			{
				Name:  "handle",
				Types: []string{"String"},
			},
			{
				Name:  "name",
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
				Name:  "archived",
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
				Name:    "roles",
				Types:   []string{"Role"},
				IsArray: true,
			},

			{
				Name:  "total",
				Types: []string{"UnsignedInteger"},
			},

			{
				Name:  "pageCursor",
				Types: []string{"String"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &rolesSearchArgs{
					hasQuery:             in.Has("query"),
					hasMemberID:          in.Has("memberID"),
					hasHandle:            in.Has("handle"),
					hasName:              in.Has("name"),
					hasLabels:            in.Has("labels"),
					hasDeleted:           in.Has("deleted"),
					hasArchived:          in.Has("archived"),
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

			var results *rolesSearchResults
			if results, err = h.search(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Roles (*types.Role) to Array (of Role)
				var (
					tval expr.TypedValue
					tarr = make([]expr.TypedValue, len(results.Roles))
				)

				for i := range results.Roles {
					if tarr[i], err = h.reg.Type("Role").Cast(results.Roles[i]); err != nil {
						return
					}
				}

				if tval, err = expr.NewArray(tarr); err != nil {
					return
				} else if err = expr.Assign(out, "roles", tval); err != nil {
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
	rolesEachArgs struct {
		hasQuery bool
		Query    string

		hasMemberID bool
		MemberID    uint64

		hasHandle bool
		Handle    string

		hasName bool
		Name    string

		hasLabels bool
		Labels    map[string]string

		hasDeleted bool
		Deleted    uint64

		hasArchived bool
		Archived    uint64

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

	rolesEachResults struct {
		Role  *types.Role
		Total uint64
	}
)

// Each function Searches for roles and iterates over results
//
// expects implementation of each function:
// func (h rolesHandler) each(ctx context.Context, args *rolesEachArgs) (results *rolesEachResults, err error) {
//    return
// }
func (h rolesHandler) Each() *atypes.Function {
	return &atypes.Function{
		Ref:  "rolesEach",
		Kind: "iterator",
		Meta: &atypes.FunctionMeta{
			Short: "Searches for roles and iterates over results",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "query",
				Types: []string{"String"},
			},
			{
				Name:  "memberID",
				Types: []string{"ID"},
			},
			{
				Name:  "handle",
				Types: []string{"String"},
			},
			{
				Name:  "name",
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
				Name:  "archived",
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
				Name:  "role",
				Types: []string{"Role"},
			},

			{
				Name:  "total",
				Types: []string{"UnsignedInteger"},
			},
		},

		Iterator: func(ctx context.Context, in *expr.Vars) (out wfexec.IteratorHandler, err error) {
			var (
				args = &rolesEachArgs{
					hasQuery:             in.Has("query"),
					hasMemberID:          in.Has("memberID"),
					hasHandle:            in.Has("handle"),
					hasName:              in.Has("name"),
					hasLabels:            in.Has("labels"),
					hasDeleted:           in.Has("deleted"),
					hasArchived:          in.Has("archived"),
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

			return h.each(ctx, args)
		},
	}
}

type (
	rolesCreateArgs struct {
		hasRole bool
		Role    *types.Role
	}

	rolesCreateResults struct {
		Role *types.Role
	}
)

// Create function Creates new role
//
// expects implementation of create function:
// func (h rolesHandler) create(ctx context.Context, args *rolesCreateArgs) (results *rolesCreateResults, err error) {
//    return
// }
func (h rolesHandler) Create() *atypes.Function {
	return &atypes.Function{
		Ref:  "rolesCreate",
		Kind: "function",
		Meta: &atypes.FunctionMeta{
			Short: "Creates new role",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "role",
				Types: []string{"Role"}, Required: true,
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "role",
				Types: []string{"Role"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &rolesCreateArgs{
					hasRole: in.Has("role"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			var results *rolesCreateResults
			if results, err = h.create(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Role (*types.Role) to Role
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("Role").Cast(results.Role); err != nil {
					return
				} else if err = expr.Assign(out, "role", tval); err != nil {
					return
				}
			}

			return
		},
	}
}

type (
	rolesUpdateArgs struct {
		hasRole bool
		Role    *types.Role
	}

	rolesUpdateResults struct {
		Role *types.Role
	}
)

// Update function Updates exiting role
//
// expects implementation of update function:
// func (h rolesHandler) update(ctx context.Context, args *rolesUpdateArgs) (results *rolesUpdateResults, err error) {
//    return
// }
func (h rolesHandler) Update() *atypes.Function {
	return &atypes.Function{
		Ref:  "rolesUpdate",
		Kind: "function",
		Meta: &atypes.FunctionMeta{
			Short: "Updates exiting role",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "role",
				Types: []string{"Role"}, Required: true,
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "role",
				Types: []string{"Role"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &rolesUpdateArgs{
					hasRole: in.Has("role"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			var results *rolesUpdateResults
			if results, err = h.update(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Role (*types.Role) to Role
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("Role").Cast(results.Role); err != nil {
					return
				} else if err = expr.Assign(out, "role", tval); err != nil {
					return
				}
			}

			return
		},
	}
}

type (
	rolesDeleteArgs struct {
		hasLookup    bool
		Lookup       interface{}
		lookupID     uint64
		lookupHandle string
		lookupRes    *types.Role
	}
)

func (a rolesDeleteArgs) GetLookup() (bool, uint64, string, *types.Role) {
	return a.hasLookup, a.lookupID, a.lookupHandle, a.lookupRes
}

// Delete function Deletes the role
//
// expects implementation of delete function:
// func (h rolesHandler) delete(ctx context.Context, args *rolesDeleteArgs) (err error) {
//    return
// }
func (h rolesHandler) Delete() *atypes.Function {
	return &atypes.Function{
		Ref:  "rolesDelete",
		Kind: "function",
		Meta: &atypes.FunctionMeta{
			Short: "Deletes the role",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "lookup",
				Types: []string{"ID", "Handle", "Role"}, Required: true,
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &rolesDeleteArgs{
					hasLookup: in.Has("lookup"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			// Converting Lookup argument
			if args.hasLookup {
				aux := expr.Must(expr.Select(in, "lookup"))
				switch aux.Type() {
				case h.reg.Type("ID").Type():
					args.lookupID = aux.Get().(uint64)
				case h.reg.Type("Handle").Type():
					args.lookupHandle = aux.Get().(string)
				case h.reg.Type("Role").Type():
					args.lookupRes = aux.Get().(*types.Role)
				}
			}

			return out, h.delete(ctx, args)
		},
	}
}

type (
	rolesRecoverArgs struct {
		hasLookup    bool
		Lookup       interface{}
		lookupID     uint64
		lookupHandle string
		lookupRes    *types.Role
	}
)

func (a rolesRecoverArgs) GetLookup() (bool, uint64, string, *types.Role) {
	return a.hasLookup, a.lookupID, a.lookupHandle, a.lookupRes
}

// Recover function Recovers deleted role
//
// expects implementation of recover function:
// func (h rolesHandler) recover(ctx context.Context, args *rolesRecoverArgs) (err error) {
//    return
// }
func (h rolesHandler) Recover() *atypes.Function {
	return &atypes.Function{
		Ref:  "rolesRecover",
		Kind: "function",
		Meta: &atypes.FunctionMeta{
			Short: "Recovers deleted role",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "lookup",
				Types: []string{"ID", "Handle", "Role"}, Required: true,
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &rolesRecoverArgs{
					hasLookup: in.Has("lookup"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			// Converting Lookup argument
			if args.hasLookup {
				aux := expr.Must(expr.Select(in, "lookup"))
				switch aux.Type() {
				case h.reg.Type("ID").Type():
					args.lookupID = aux.Get().(uint64)
				case h.reg.Type("Handle").Type():
					args.lookupHandle = aux.Get().(string)
				case h.reg.Type("Role").Type():
					args.lookupRes = aux.Get().(*types.Role)
				}
			}

			return out, h.recover(ctx, args)
		},
	}
}

type (
	rolesArchiveArgs struct {
		hasLookup    bool
		Lookup       interface{}
		lookupID     uint64
		lookupHandle string
		lookupRes    *types.Role
	}
)

func (a rolesArchiveArgs) GetLookup() (bool, uint64, string, *types.Role) {
	return a.hasLookup, a.lookupID, a.lookupHandle, a.lookupRes
}

// Archive function Archives the role
//
// expects implementation of archive function:
// func (h rolesHandler) archive(ctx context.Context, args *rolesArchiveArgs) (err error) {
//    return
// }
func (h rolesHandler) Archive() *atypes.Function {
	return &atypes.Function{
		Ref:  "rolesArchive",
		Kind: "function",
		Meta: &atypes.FunctionMeta{
			Short: "Archives the role",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "lookup",
				Types: []string{"ID", "Handle", "Role"}, Required: true,
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &rolesArchiveArgs{
					hasLookup: in.Has("lookup"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			// Converting Lookup argument
			if args.hasLookup {
				aux := expr.Must(expr.Select(in, "lookup"))
				switch aux.Type() {
				case h.reg.Type("ID").Type():
					args.lookupID = aux.Get().(uint64)
				case h.reg.Type("Handle").Type():
					args.lookupHandle = aux.Get().(string)
				case h.reg.Type("Role").Type():
					args.lookupRes = aux.Get().(*types.Role)
				}
			}

			return out, h.archive(ctx, args)
		},
	}
}

type (
	rolesUnarchiveArgs struct {
		hasLookup    bool
		Lookup       interface{}
		lookupID     uint64
		lookupHandle string
		lookupRes    *types.Role
	}
)

func (a rolesUnarchiveArgs) GetLookup() (bool, uint64, string, *types.Role) {
	return a.hasLookup, a.lookupID, a.lookupHandle, a.lookupRes
}

// Unarchive function Unarchives the role
//
// expects implementation of unarchive function:
// func (h rolesHandler) unarchive(ctx context.Context, args *rolesUnarchiveArgs) (err error) {
//    return
// }
func (h rolesHandler) Unarchive() *atypes.Function {
	return &atypes.Function{
		Ref:  "rolesUnarchive",
		Kind: "function",
		Meta: &atypes.FunctionMeta{
			Short: "Unarchives the role",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "lookup",
				Types: []string{"ID", "Handle", "Role"}, Required: true,
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &rolesUnarchiveArgs{
					hasLookup: in.Has("lookup"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			// Converting Lookup argument
			if args.hasLookup {
				aux := expr.Must(expr.Select(in, "lookup"))
				switch aux.Type() {
				case h.reg.Type("ID").Type():
					args.lookupID = aux.Get().(uint64)
				case h.reg.Type("Handle").Type():
					args.lookupHandle = aux.Get().(string)
				case h.reg.Type("Role").Type():
					args.lookupRes = aux.Get().(*types.Role)
				}
			}

			return out, h.unarchive(ctx, args)
		},
	}
}
