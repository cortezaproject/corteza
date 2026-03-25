package automation

import (
	"context"

	atypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	ngUsersHandler struct {
		reg  constructSvc
		tReg typeRegistry
		h    *usersHandler
	}
)

func NgUsersHandler(reg constructSvc, tReg typeRegistry, uSvc userService, rSvc roleService) *ngUsersHandler {
	ngh := &ngUsersHandler{
		reg:  reg,
		tReg: tReg,
		h: &usersHandler{
			uSvc: uSvc,
			rSvc: rSvc,
		},
	}

	ngh.register()
	return ngh
}

func (h ngUsersHandler) register() {
	h.reg.AddFunctions(
		h.Lookup(),
		h.Create(),
		h.Update(),
		h.Delete(),
		h.Suspend(),
		h.Unsuspend(),
	)
}

func (h ngUsersHandler) Lookup() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "usersLookup",
		Kind:   "function",
		Groups: []string{"System User"},
		Labels: map[string]string{"users": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "User lookup",
			Description: "Find specific user by ID, handle or string",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "user"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "lookup", Types: []string{"ID", "Handle", "String", "User"}, Required: true},
		},

		Results: []*atypes.Param{
			{ArgumentName: "user", Types: []string{"User"}},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Lookup (ID/Handle/Email/User)", Argument: "lookup"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var args = &usersLookupArgs{
				hasLookup: in.Has("lookup"),
			}

			if err = in.Decode(args); err != nil {
				return
			}

			if args.hasLookup {
				aux := expr.Must(expr.Select(in, "lookup"))
				switch aux.Type() {
				case h.tReg.Type("ID").Type():
					args.lookupID = aux.Get().(uint64)
				case h.tReg.Type("Handle").Type():
					args.lookupHandle = aux.Get().(string)
				case h.tReg.Type("String").Type():
					args.lookupEmail = aux.Get().(string)
				case h.tReg.Type("User").Type():
					args.lookupRes = aux.Get().(*types.User)
				}
			}

			var results *usersLookupResults
			if results, err = h.h.lookup(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}
			if tval, err := h.tReg.Type("User").Cast(results.User); err == nil {
				_ = expr.Assign(out, "user", tval)
			}
			return
		},
	}
}

func (h ngUsersHandler) Create() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "usersCreate",
		Kind:   "function",
		Groups: []string{"System User"},
		Labels: map[string]string{"users": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short: "User create",
			Icon:  &atypes.NgAutomationIcon{Type: "name", Value: "user-plus"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "email", Types: []string{"String"}, Required: true},
			{ArgumentName: "name", Types: []string{"String"}},
			{ArgumentName: "handle", Types: []string{"String"}},
			{ArgumentName: "username", Types: []string{"String"}},
			{ArgumentName: "kind", Types: []string{"String"}},
		},

		Results: []*atypes.Param{
			{ArgumentName: "user", Types: []string{"User"}},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Email", Argument: "email"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Name", Argument: "name"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Handle", Argument: "handle"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Username", Argument: "username"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Kind", Argument: "kind"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var args = &usersCreateArgs{
				hasUser: true,
				User:    &types.User{},
			}

			if in.Has("email") {
				aux := expr.Must(expr.Select(in, "email"))
				args.User.Email, _ = expr.CastToString(aux.Get())
			}
			if in.Has("name") {
				aux := expr.Must(expr.Select(in, "name"))
				args.User.Name, _ = expr.CastToString(aux.Get())
			}
			if in.Has("handle") {
				aux := expr.Must(expr.Select(in, "handle"))
				args.User.Handle, _ = expr.CastToString(aux.Get())
			}
			if in.Has("username") {
				aux := expr.Must(expr.Select(in, "username"))
				args.User.Username, _ = expr.CastToString(aux.Get())
			}
			if in.Has("kind") {
				aux := expr.Must(expr.Select(in, "kind"))
				if kindStr, castErr := expr.CastToString(aux.Get()); castErr == nil {
					args.User.Kind = types.UserKind(kindStr)
				}
			}

			var results *usersCreateResults
			if results, err = h.h.create(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}
			if tval, err := h.tReg.Type("User").Cast(results.User); err == nil {
				_ = expr.Assign(out, "user", tval)
			}
			return
		},
	}
}

func (h ngUsersHandler) Update() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "usersUpdate",
		Kind:   "function",
		Groups: []string{"System User"},
		Labels: map[string]string{"users": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short: "User update",
			Icon:  &atypes.NgAutomationIcon{Type: "name", Value: "pencil"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "user", Types: []string{"User"}, Required: true},
		},

		Results: []*atypes.Param{
			{ArgumentName: "user", Types: []string{"User"}},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "User", Argument: "user"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var args = &usersUpdateArgs{
				hasUser: in.Has("user"),
			}

			if err = in.Decode(args); err != nil {
				return
			}

			var results *usersUpdateResults
			if results, err = h.h.update(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}
			if tval, err := h.tReg.Type("User").Cast(results.User); err == nil {
				_ = expr.Assign(out, "user", tval)
			}
			return
		},
	}
}

func (h ngUsersHandler) Delete() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "usersDelete",
		Kind:   "function",
		Groups: []string{"System User"},
		Labels: map[string]string{"delete": "step", "users": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short: "User delete",
			Icon:  &atypes.NgAutomationIcon{Type: "name", Value: "trash"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "lookup", Types: []string{"ID", "Handle", "String", "User"}, Required: true},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Lookup (ID/Handle/Email/User)", Argument: "lookup"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var args = &usersDeleteArgs{
				hasLookup: in.Has("lookup"),
			}

			if err = in.Decode(args); err != nil {
				return
			}

			if args.hasLookup {
				aux := expr.Must(expr.Select(in, "lookup"))
				switch aux.Type() {
				case h.tReg.Type("ID").Type():
					args.lookupID = aux.Get().(uint64)
				case h.tReg.Type("Handle").Type():
					args.lookupHandle = aux.Get().(string)
				case h.tReg.Type("String").Type():
					args.lookupEmail = aux.Get().(string)
				case h.tReg.Type("User").Type():
					args.lookupRes = aux.Get().(*types.User)
				}
			}

			return out, h.h.delete(ctx, args)
		},
	}
}

func (h ngUsersHandler) Suspend() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "usersSuspend",
		Kind:   "function",
		Groups: []string{"System User"},
		Labels: map[string]string{"suspend": "step", "users": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short: "User suspend",
			Icon:  &atypes.NgAutomationIcon{Type: "name", Value: "lock"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "lookup", Types: []string{"ID", "Handle", "String", "User"}, Required: true},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Lookup (ID/Handle/Email/User)", Argument: "lookup"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var args = &usersSuspendArgs{
				hasLookup: in.Has("lookup"),
			}

			if err = in.Decode(args); err != nil {
				return
			}

			if args.hasLookup {
				aux := expr.Must(expr.Select(in, "lookup"))
				switch aux.Type() {
				case h.tReg.Type("ID").Type():
					args.lookupID = aux.Get().(uint64)
				case h.tReg.Type("Handle").Type():
					args.lookupHandle = aux.Get().(string)
				case h.tReg.Type("String").Type():
					args.lookupEmail = aux.Get().(string)
				case h.tReg.Type("User").Type():
					args.lookupRes = aux.Get().(*types.User)
				}
			}

			return out, h.h.suspend(ctx, args)
		},
	}
}

func (h ngUsersHandler) Unsuspend() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "usersUnsuspend",
		Kind:   "function",
		Groups: []string{"System User"},
		Labels: map[string]string{"unsuspend": "step", "users": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short: "User unsuspend",
			Icon:  &atypes.NgAutomationIcon{Type: "name", Value: "unlock"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "lookup", Types: []string{"ID", "Handle", "String", "User"}, Required: true},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Lookup (ID/Handle/Email/User)", Argument: "lookup"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var args = &usersUnsuspendArgs{
				hasLookup: in.Has("lookup"),
			}

			if err = in.Decode(args); err != nil {
				return
			}

			if args.hasLookup {
				aux := expr.Must(expr.Select(in, "lookup"))
				switch aux.Type() {
				case h.tReg.Type("ID").Type():
					args.lookupID = aux.Get().(uint64)
				case h.tReg.Type("Handle").Type():
					args.lookupHandle = aux.Get().(string)
				case h.tReg.Type("String").Type():
					args.lookupEmail = aux.Get().(string)
				case h.tReg.Type("User").Type():
					args.lookupRes = aux.Get().(*types.User)
				}
			}

			return out, h.h.unsuspend(ctx, args)
		},
	}
}
