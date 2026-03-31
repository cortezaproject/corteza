package automation

import (
	"context"

	atypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	ngRolesHandler struct {
		reg  constructSvc
		tReg typeRegistry
		h    *rolesHandler
	}
)

func NgRolesHandler(reg constructSvc, tReg typeRegistry, roleSvc roleService, userSvc userService) *ngRolesHandler {
	ngh := &ngRolesHandler{
		reg:  reg,
		tReg: tReg,
		h: &rolesHandler{
			rSvc: roleSvc,
			uSvc: userSvc,
		},
	}

	ngh.register()
	return ngh
}

func (h ngRolesHandler) register() {}

func (h ngRolesHandler) Lookup() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "rolesLookup",
		Kind:   "function",
		Groups: []string{"Roles"},
		Labels: map[string]string{"users": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Find Role",
			Description: "Look up a role by ID or handle",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "users"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "lookup", Types: []string{"ID", "Handle", "Role"}, Required: true},
		},

		Results: []*atypes.Param{
			{ArgumentName: "role", Types: []string{"Role"}},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Lookup (ID/Handle/Role)", Argument: "lookup"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var args = &rolesLookupArgs{
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
				case h.tReg.Type("Role").Type():
					args.lookupRes = aux.Get().(*types.Role)
				}
			}

			var results *rolesLookupResults
			if results, err = h.h.lookup(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}
			if tval, err := h.tReg.Type("Role").Cast(results.Role); err == nil {
				_ = expr.Assign(out, "role", tval)
			}
			return
		},
	}
}

func (h ngRolesHandler) AddMember() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "rolesAddMember",
		Kind:   "function",
		Groups: []string{"Roles"},
		Labels: map[string]string{"users": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Add Role Member",
			Description: "Add a user to a role",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "users"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "role", Types: []string{"ID", "Handle", "Role"}, Required: true},
			{ArgumentName: "user", Types: []string{"ID", "Handle", "String", "User"}, Required: true},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Role (ID/Handle/Role)", Argument: "role"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "User (ID/Handle/Email/User)", Argument: "user"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var args = &rolesAddMemberArgs{
				hasRole: in.Has("role"),
				hasUser: in.Has("user"),
			}

			if err = in.Decode(args); err != nil {
				return
			}

			if args.hasRole {
				aux := expr.Must(expr.Select(in, "role"))
				switch aux.Type() {
				case h.tReg.Type("ID").Type():
					args.roleID = aux.Get().(uint64)
				case h.tReg.Type("Handle").Type():
					args.roleHandle = aux.Get().(string)
				case h.tReg.Type("Role").Type():
					args.roleRes = aux.Get().(*types.Role)
				}
			}

			if args.hasUser {
				aux := expr.Must(expr.Select(in, "user"))
				switch aux.Type() {
				case h.tReg.Type("ID").Type():
					args.userID = aux.Get().(uint64)
				case h.tReg.Type("Handle").Type():
					args.userHandle = aux.Get().(string)
				case h.tReg.Type("String").Type():
					args.userEmail = aux.Get().(string)
				case h.tReg.Type("User").Type():
					args.userRes = aux.Get().(*types.User)
				}
			}

			return out, h.h.addMember(ctx, args)
		},
	}
}

func (h ngRolesHandler) RemoveMember() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "rolesRemoveMember",
		Kind:   "function",
		Groups: []string{"Roles"},
		Labels: map[string]string{"users": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Remove Role Member",
			Description: "Remove a user from a role",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "users"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "role", Types: []string{"ID", "Handle", "Role"}, Required: true},
			{ArgumentName: "user", Types: []string{"ID", "Handle", "String", "User"}, Required: true},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Role (ID/Handle/Role)", Argument: "role"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "User (ID/Handle/Email/User)", Argument: "user"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var args = &rolesRemoveMemberArgs{
				hasRole: in.Has("role"),
				hasUser: in.Has("user"),
			}

			if err = in.Decode(args); err != nil {
				return
			}

			if args.hasRole {
				aux := expr.Must(expr.Select(in, "role"))
				switch aux.Type() {
				case h.tReg.Type("ID").Type():
					args.roleID = aux.Get().(uint64)
				case h.tReg.Type("Handle").Type():
					args.roleHandle = aux.Get().(string)
				case h.tReg.Type("Role").Type():
					args.roleRes = aux.Get().(*types.Role)
				}
			}

			if args.hasUser {
				aux := expr.Must(expr.Select(in, "user"))
				switch aux.Type() {
				case h.tReg.Type("ID").Type():
					args.userID = aux.Get().(uint64)
				case h.tReg.Type("Handle").Type():
					args.userHandle = aux.Get().(string)
				case h.tReg.Type("String").Type():
					args.userEmail = aux.Get().(string)
				case h.tReg.Type("User").Type():
					args.userRes = aux.Get().(*types.User)
				}
			}

			return out, h.h.removeMember(ctx, args)
		},
	}
}

func (h ngRolesHandler) Create() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "rolesCreate",
		Kind:   "function",
		Groups: []string{"Roles"},
		Labels: map[string]string{"create": "step", "users": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Create Role",
			Description: "Create a new role",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "users"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "name", Types: []string{"String"}, Required: true},
			{ArgumentName: "handle", Types: []string{"String"}},
		},

		Results: []*atypes.Param{
			{ArgumentName: "role", Types: []string{"Role"}},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Name", Argument: "name"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Handle", Argument: "handle"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var args = &rolesCreateArgs{
				hasRole: true,
				Role:    &types.Role{},
			}

			if in.Has("name") {
				aux := expr.Must(expr.Select(in, "name"))
				args.Role.Name, _ = expr.CastToString(aux.Get())
			}
			if in.Has("handle") {
				aux := expr.Must(expr.Select(in, "handle"))
				args.Role.Handle, _ = expr.CastToString(aux.Get())
			}

			var results *rolesCreateResults
			if results, err = h.h.create(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}
			if tval, err := h.tReg.Type("Role").Cast(results.Role); err == nil {
				_ = expr.Assign(out, "role", tval)
			}
			return
		},
	}
}

func (h ngRolesHandler) Update() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "rolesUpdate",
		Kind:   "function",
		Groups: []string{"Roles"},
		Labels: map[string]string{"update": "step", "users": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Update Role",
			Description: "Update an existing role",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "users"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "role", Types: []string{"Role"}, Required: true},
		},

		Results: []*atypes.Param{
			{ArgumentName: "role", Types: []string{"Role"}},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Role", Argument: "role"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var args = &rolesUpdateArgs{
				hasRole: in.Has("role"),
			}

			if err = in.Decode(args); err != nil {
				return
			}

			var results *rolesUpdateResults
			if results, err = h.h.update(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}
			if tval, err := h.tReg.Type("Role").Cast(results.Role); err == nil {
				_ = expr.Assign(out, "role", tval)
			}
			return
		},
	}
}

func (h ngRolesHandler) Delete() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "rolesDelete",
		Kind:   "function",
		Groups: []string{"Roles"},
		Labels: map[string]string{"delete": "step", "users": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Delete Role",
			Description: "Delete a role",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "users"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "lookup", Types: []string{"ID", "Handle", "Role"}, Required: true},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Lookup (ID/Handle/Role)", Argument: "lookup"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var args = &rolesDeleteArgs{
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
				case h.tReg.Type("Role").Type():
					args.lookupRes = aux.Get().(*types.Role)
				}
			}

			return out, h.h.delete(ctx, args)
		},
	}
}

func (h ngRolesHandler) Archive() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "rolesArchive",
		Kind:   "function",
		Groups: []string{"Roles"},
		Labels: map[string]string{"archive": "step", "users": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Archive Role",
			Description: "Archive a role",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "users"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "lookup", Types: []string{"ID", "Handle", "Role"}, Required: true},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Lookup (ID/Handle/Role)", Argument: "lookup"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var args = &rolesArchiveArgs{
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
				case h.tReg.Type("Role").Type():
					args.lookupRes = aux.Get().(*types.Role)
				}
			}

			return out, h.h.archive(ctx, args)
		},
	}
}

func (h ngRolesHandler) Unarchive() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "rolesUnarchive",
		Kind:   "function",
		Groups: []string{"Roles"},
		Labels: map[string]string{"unarchive": "step", "users": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Unarchive Role",
			Description: "Unarchive a role",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "users"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "lookup", Types: []string{"ID", "Handle", "Role"}, Required: true},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Lookup (ID/Handle/Role)", Argument: "lookup"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var args = &rolesUnarchiveArgs{
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
				case h.tReg.Type("Role").Type():
					args.lookupRes = aux.Get().(*types.Role)
				}
			}

			return out, h.h.unarchive(ctx, args)
		},
	}
}
