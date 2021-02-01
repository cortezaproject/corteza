package automation

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/automation/users_handler.yaml

import (
	"context"
	atypes "github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
	"github.com/cortezaproject/corteza-server/system/types"
)

var _ wfexec.ExecResponse

type (
	usersHandlerRegistry interface {
		AddFunctions(ff ...*atypes.Function)
		Type(ref string) expr.Type
	}
)

func (h usersHandler) register() {
	h.reg.AddFunctions(
		h.Lookup(),
		h.Search(),
		h.Each(),
		h.Create(),
		h.Update(),
		h.Delete(),
		h.Recover(),
		h.Suspend(),
		h.Unsuspend(),
	)
}

type (
	usersLookupArgs struct {
		hasLookup    bool
		Lookup       interface{}
		lookupID     uint64
		lookupHandle string
		lookupEmail  string
		lookupRes    *types.User
	}

	usersLookupResults struct {
		User *types.User
	}
)

func (a usersLookupArgs) GetLookup() (bool, uint64, string, string, *types.User) {
	return a.hasLookup, a.lookupID, a.lookupHandle, a.lookupEmail, a.lookupRes
}

// Lookup function Looks-up for user by ID
//
// expects implementation of lookup function:
// func (h usersHandler) lookup(ctx context.Context, args *usersLookupArgs) (results *usersLookupResults, err error) {
//    return
// }
func (h usersHandler) Lookup() *atypes.Function {
	return &atypes.Function{
		Ref:  "usersLookup",
		Kind: "function",
		Meta: &atypes.FunctionMeta{
			Short: "Looks-up for user by ID",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "lookup",
				Types: []string{"ID", "Handle", "String", "User"}, Required: true,
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "user",
				Types: []string{"User"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &usersLookupArgs{
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
				case h.reg.Type("String").Type():
					args.lookupEmail = aux.Get().(string)
				case h.reg.Type("User").Type():
					args.lookupRes = aux.Get().(*types.User)
				}
			}

			var results *usersLookupResults
			if results, err = h.lookup(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.User (*types.User) to User
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("User").Cast(results.User); err != nil {
					return
				} else if err = expr.Assign(out, "user", tval); err != nil {
					return
				}
			}

			return
		},
	}
}

type (
	usersSearchArgs struct {
		hasQuery bool
		Query    string

		hasEmail bool
		Email    string

		hasHandle bool
		Handle    string

		hasLabels bool
		Labels    map[string]string

		hasDeleted bool
		Deleted    uint64

		hasSuspended bool
		Suspended    uint64

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

	usersSearchResults struct {
		Users      []*types.User
		Total      uint64
		PageCursor string
	}
)

// Search function Searches for users and returns them
//
// expects implementation of search function:
// func (h usersHandler) search(ctx context.Context, args *usersSearchArgs) (results *usersSearchResults, err error) {
//    return
// }
func (h usersHandler) Search() *atypes.Function {
	return &atypes.Function{
		Ref:  "usersSearch",
		Kind: "function",
		Meta: &atypes.FunctionMeta{
			Short: "Searches for users and returns them",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "query",
				Types: []string{"String"},
			},
			{
				Name:  "email",
				Types: []string{"String"},
			},
			{
				Name:  "handle",
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
				Name:  "suspended",
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
				Name:    "users",
				Types:   []string{"User"},
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
				args = &usersSearchArgs{
					hasQuery:             in.Has("query"),
					hasEmail:             in.Has("email"),
					hasHandle:            in.Has("handle"),
					hasLabels:            in.Has("labels"),
					hasDeleted:           in.Has("deleted"),
					hasSuspended:         in.Has("suspended"),
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

			var results *usersSearchResults
			if results, err = h.search(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Users (*types.User) to Array (of User)
				var (
					tval expr.TypedValue
					tarr = make([]expr.TypedValue, len(results.Users))
				)

				for i := range results.Users {
					if tarr[i], err = h.reg.Type("User").Cast(results.Users[i]); err != nil {
						return
					}
				}

				if tval, err = expr.NewArray(tarr); err != nil {
					return
				} else if err = expr.Assign(out, "users", tval); err != nil {
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
	usersEachArgs struct {
		hasQuery bool
		Query    string

		hasEmail bool
		Email    string

		hasHandle bool
		Handle    string

		hasLabels bool
		Labels    map[string]string

		hasDeleted bool
		Deleted    uint64

		hasSuspended bool
		Suspended    uint64

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

	usersEachResults struct {
		User  *types.User
		Total uint64
	}
)

// Each function Searches for users and iterates over results
//
// expects implementation of each function:
// func (h usersHandler) each(ctx context.Context, args *usersEachArgs) (results *usersEachResults, err error) {
//    return
// }
func (h usersHandler) Each() *atypes.Function {
	return &atypes.Function{
		Ref:  "usersEach",
		Kind: "iterator",
		Meta: &atypes.FunctionMeta{
			Short: "Searches for users and iterates over results",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "query",
				Types: []string{"String"},
			},
			{
				Name:  "email",
				Types: []string{"String"},
			},
			{
				Name:  "handle",
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
				Name:  "suspended",
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
				Name:  "user",
				Types: []string{"User"},
			},

			{
				Name:  "total",
				Types: []string{"UnsignedInteger"},
			},
		},

		Iterator: func(ctx context.Context, in *expr.Vars) (out wfexec.IteratorHandler, err error) {
			var (
				args = &usersEachArgs{
					hasQuery:             in.Has("query"),
					hasEmail:             in.Has("email"),
					hasHandle:            in.Has("handle"),
					hasLabels:            in.Has("labels"),
					hasDeleted:           in.Has("deleted"),
					hasSuspended:         in.Has("suspended"),
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
	usersCreateArgs struct {
		hasUser bool
		User    *types.User
	}

	usersCreateResults struct {
		User *types.User
	}
)

// Create function Creates new user
//
// expects implementation of create function:
// func (h usersHandler) create(ctx context.Context, args *usersCreateArgs) (results *usersCreateResults, err error) {
//    return
// }
func (h usersHandler) Create() *atypes.Function {
	return &atypes.Function{
		Ref:  "usersCreate",
		Kind: "function",
		Meta: &atypes.FunctionMeta{
			Short: "Creates new user",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "user",
				Types: []string{"User"}, Required: true,
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "user",
				Types: []string{"User"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &usersCreateArgs{
					hasUser: in.Has("user"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			var results *usersCreateResults
			if results, err = h.create(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.User (*types.User) to User
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("User").Cast(results.User); err != nil {
					return
				} else if err = expr.Assign(out, "user", tval); err != nil {
					return
				}
			}

			return
		},
	}
}

type (
	usersUpdateArgs struct {
		hasUser bool
		User    *types.User
	}

	usersUpdateResults struct {
		User *types.User
	}
)

// Update function Updates exiting user
//
// expects implementation of update function:
// func (h usersHandler) update(ctx context.Context, args *usersUpdateArgs) (results *usersUpdateResults, err error) {
//    return
// }
func (h usersHandler) Update() *atypes.Function {
	return &atypes.Function{
		Ref:  "usersUpdate",
		Kind: "function",
		Meta: &atypes.FunctionMeta{
			Short: "Updates exiting user",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "user",
				Types: []string{"User"}, Required: true,
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "user",
				Types: []string{"User"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &usersUpdateArgs{
					hasUser: in.Has("user"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			var results *usersUpdateResults
			if results, err = h.update(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.User (*types.User) to User
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("User").Cast(results.User); err != nil {
					return
				} else if err = expr.Assign(out, "user", tval); err != nil {
					return
				}
			}

			return
		},
	}
}

type (
	usersDeleteArgs struct {
		hasLookup    bool
		Lookup       interface{}
		lookupID     uint64
		lookupHandle string
		lookupEmail  string
		lookupRes    *types.User
	}
)

func (a usersDeleteArgs) GetLookup() (bool, uint64, string, string, *types.User) {
	return a.hasLookup, a.lookupID, a.lookupHandle, a.lookupEmail, a.lookupRes
}

// Delete function Deletes user
//
// expects implementation of delete function:
// func (h usersHandler) delete(ctx context.Context, args *usersDeleteArgs) (err error) {
//    return
// }
func (h usersHandler) Delete() *atypes.Function {
	return &atypes.Function{
		Ref:  "usersDelete",
		Kind: "function",
		Meta: &atypes.FunctionMeta{
			Short: "Deletes user",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "lookup",
				Types: []string{"ID", "Handle", "String", "User"}, Required: true,
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &usersDeleteArgs{
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
				case h.reg.Type("String").Type():
					args.lookupEmail = aux.Get().(string)
				case h.reg.Type("User").Type():
					args.lookupRes = aux.Get().(*types.User)
				}
			}

			return out, h.delete(ctx, args)
		},
	}
}

type (
	usersRecoverArgs struct {
		hasLookup    bool
		Lookup       interface{}
		lookupID     uint64
		lookupHandle string
		lookupEmail  string
		lookupRes    *types.User
	}
)

func (a usersRecoverArgs) GetLookup() (bool, uint64, string, string, *types.User) {
	return a.hasLookup, a.lookupID, a.lookupHandle, a.lookupEmail, a.lookupRes
}

// Recover function Recovers deleted user
//
// expects implementation of recover function:
// func (h usersHandler) recover(ctx context.Context, args *usersRecoverArgs) (err error) {
//    return
// }
func (h usersHandler) Recover() *atypes.Function {
	return &atypes.Function{
		Ref:  "usersRecover",
		Kind: "function",
		Meta: &atypes.FunctionMeta{
			Short: "Recovers deleted user",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "lookup",
				Types: []string{"ID", "Handle", "String", "User"}, Required: true,
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &usersRecoverArgs{
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
				case h.reg.Type("String").Type():
					args.lookupEmail = aux.Get().(string)
				case h.reg.Type("User").Type():
					args.lookupRes = aux.Get().(*types.User)
				}
			}

			return out, h.recover(ctx, args)
		},
	}
}

type (
	usersSuspendArgs struct {
		hasLookup    bool
		Lookup       interface{}
		lookupID     uint64
		lookupHandle string
		lookupEmail  string
		lookupRes    *types.User
	}
)

func (a usersSuspendArgs) GetLookup() (bool, uint64, string, string, *types.User) {
	return a.hasLookup, a.lookupID, a.lookupHandle, a.lookupEmail, a.lookupRes
}

// Suspend function Suspends user
//
// expects implementation of suspend function:
// func (h usersHandler) suspend(ctx context.Context, args *usersSuspendArgs) (err error) {
//    return
// }
func (h usersHandler) Suspend() *atypes.Function {
	return &atypes.Function{
		Ref:  "usersSuspend",
		Kind: "function",
		Meta: &atypes.FunctionMeta{
			Short: "Suspends user",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "lookup",
				Types: []string{"ID", "Handle", "String", "User"}, Required: true,
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &usersSuspendArgs{
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
				case h.reg.Type("String").Type():
					args.lookupEmail = aux.Get().(string)
				case h.reg.Type("User").Type():
					args.lookupRes = aux.Get().(*types.User)
				}
			}

			return out, h.suspend(ctx, args)
		},
	}
}

type (
	usersUnsuspendArgs struct {
		hasLookup    bool
		Lookup       interface{}
		lookupID     uint64
		lookupHandle string
		lookupEmail  string
		lookupRes    *types.User
	}
)

func (a usersUnsuspendArgs) GetLookup() (bool, uint64, string, string, *types.User) {
	return a.hasLookup, a.lookupID, a.lookupHandle, a.lookupEmail, a.lookupRes
}

// Unsuspend function Unsuspends user
//
// expects implementation of unsuspend function:
// func (h usersHandler) unsuspend(ctx context.Context, args *usersUnsuspendArgs) (err error) {
//    return
// }
func (h usersHandler) Unsuspend() *atypes.Function {
	return &atypes.Function{
		Ref:  "usersUnsuspend",
		Kind: "function",
		Meta: &atypes.FunctionMeta{
			Short: "Unsuspends user",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "lookup",
				Types: []string{"ID", "Handle", "String", "User"}, Required: true,
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &usersUnsuspendArgs{
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
				case h.reg.Type("String").Type():
					args.lookupEmail = aux.Get().(string)
				case h.reg.Type("User").Type():
					args.lookupRes = aux.Get().(*types.User)
				}
			}

			return out, h.unsuspend(ctx, args)
		},
	}
}
