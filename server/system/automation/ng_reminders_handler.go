package automation

import (
	"context"

	atypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	ngRemindersHandler struct {
		reg  constructSvc
		tReg typeRegistry
		h    *remindersHandler
	}
)

func NgRemindersHandler(reg constructSvc, tReg typeRegistry, svc reminderService) *ngRemindersHandler {
	ngh := &ngRemindersHandler{
		reg:  reg,
		tReg: tReg,
		h:    &remindersHandler{rSvc: svc}, // internal handler only uses rSvc
	}

	ngh.register()
	return ngh
}

func (h ngRemindersHandler) register() {
	h.reg.AddFunctions(
		h.Lookup(),
		h.Search(),
		h.Create(),
		h.Update(),
		h.Dismiss(),
		h.Snooze(),
		h.Delete(),
	)
}

func (h ngRemindersHandler) Lookup() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "remindersLookup",
		Kind:   "function",
		Groups: []string{"Reminders"},
		Labels: map[string]string{"reminders": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Find Reminder",
			Description: "Look up a reminder by ID",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "clock"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "lookup", Types: []string{"ID", "Reminder"}, Required: true},
		},

		Results: []*atypes.Param{
			{ArgumentName: "reminder", Types: []string{"Reminder"}},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Lookup (ID/Reminder)", Argument: "lookup"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &remindersLookupArgs{
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
				case h.tReg.Type("ID").Type():
					args.lookupID = aux.Get().(uint64)
				case h.tReg.Type("Reminder").Type():
					args.lookupRes = aux.Get().(*types.Reminder)
				}
			}

			var results *remindersLookupResults
			if results, err = h.h.lookup(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			if tval, err := h.tReg.Type("Reminder").Cast(results.Reminder); err == nil {
				_ = expr.Assign(out, "reminder", tval)
			}

			return
		},
	}
}

func (h ngRemindersHandler) Search() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "remindersSearch",
		Kind:   "function",
		Groups: []string{"Reminders"},
		Labels: map[string]string{"reminders": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Search Reminders",
			Description: "Search reminders with filters",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "clock"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "resource", Types: []string{"String"}},
			{ArgumentName: "assignedTo", Types: []string{"ID"}},
			{ArgumentName: "excludeDismissed", Types: []string{"Boolean"}},
			{ArgumentName: "scheduledOnly", Types: []string{"Boolean"}},
			{ArgumentName: "sort", Types: []string{"String"}},
			{ArgumentName: "limit", Types: []string{"UnsignedInteger"}},
			{ArgumentName: "pageCursor", Types: []string{"String"}},
		},

		Results: []*atypes.Param{
			{ArgumentName: "reminders", Types: []string{"Reminder"}, IsArray: true},
			{ArgumentName: "total", Types: []string{"UnsignedInteger"}},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Resource", Argument: "resource"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Assigned To", Argument: "assignedTo"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Exclude Dismissed", Argument: "excludeDismissed"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Scheduled Only", Argument: "scheduledOnly"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Sort", Argument: "sort"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Limit", Argument: "limit"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Page Cursor", Argument: "pageCursor"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var args = &remindersSearchArgs{
				hasResource:         in.Has("resource"),
				hasAssignedTo:       in.Has("assignedTo"),
				hasExcludeDismissed: in.Has("excludeDismissed"),
				hasScheduledOnly:    in.Has("scheduledOnly"),
				hasSort:             in.Has("sort"),
				hasLimit:            in.Has("limit"),
				hasPageCursor:       in.Has("pageCursor"),
			}

			if err = in.Decode(args); err != nil {
				return
			}

			var results *remindersSearchResults
			if results, err = h.h.search(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				var (
					tval expr.TypedValue
					tarr = make([]expr.TypedValue, len(results.Reminders))
				)
				for i := range results.Reminders {
					if tarr[i], err = h.tReg.Type("Reminder").Cast(results.Reminders[i]); err != nil {
						return
					}
				}
				if tval, err = expr.NewArray(tarr); err != nil {
					return
				} else if err = expr.Assign(out, "reminders", tval); err != nil {
					return
				}
			}

			if tval, err := h.tReg.Type("UnsignedInteger").Cast(results.Total); err == nil {
				_ = expr.Assign(out, "total", tval)
			}

			return
		},
	}
}

func (h ngRemindersHandler) Create() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "remindersCreate",
		Kind:   "function",
		Groups: []string{"Reminders"},
		Labels: map[string]string{"create": "step", "reminders": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short: "Create Reminder",
			Icon:  &atypes.NgAutomationIcon{Type: "name", Value: "clock"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "resource", Types: []string{"String"}},
			{ArgumentName: "assignedTo", Types: []string{"ID"}},
			{ArgumentName: "assignedBy", Types: []string{"ID"}},
			{ArgumentName: "remindAt", Types: []string{"DateTime"}, Required: true},
			{ArgumentName: "payload", Types: []string{"Bytes"}},
		},

		Results: []*atypes.Param{
			{ArgumentName: "reminder", Types: []string{"Reminder"}},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Resource", Argument: "resource"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Assigned To", Argument: "assignedTo"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Assigned By", Argument: "assignedBy"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Remind At (Date)", Argument: "remindAt"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Payload (Bytes)", Argument: "payload"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var args = &remindersCreateArgs{
				hasResource:   in.Has("resource"),
				hasAssignedTo: in.Has("assignedTo"),
				hasAssignedBy: in.Has("assignedBy"),
				hasRemindAt:   in.Has("remindAt"),
				hasPayload:    in.Has("payload"),
			}

			if err = in.Decode(args); err != nil {
				return
			}

			var results *remindersCreateResults
			if results, err = h.h.create(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}
			if tval, err := h.tReg.Type("Reminder").Cast(results.Reminder); err == nil {
				_ = expr.Assign(out, "reminder", tval)
			}
			return
		},
	}
}

func (h ngRemindersHandler) Update() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "remindersUpdate",
		Kind:   "function",
		Groups: []string{"Reminders"},
		Labels: map[string]string{"reminders": "step,workflow", "update": "step"},
		Meta: &atypes.ConstructFunctionMeta{
			Short: "Update Reminder",
			Icon:  &atypes.NgAutomationIcon{Type: "name", Value: "clock"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "lookup", Types: []string{"ID", "Reminder"}, Required: true},
			{ArgumentName: "resource", Types: []string{"String"}},
			{ArgumentName: "assignedTo", Types: []string{"ID"}},
			{ArgumentName: "assignedBy", Types: []string{"ID"}},
			{ArgumentName: "remindAt", Types: []string{"DateTime"}},
			{ArgumentName: "payload", Types: []string{"Bytes"}},
		},

		Results: []*atypes.Param{
			{ArgumentName: "reminder", Types: []string{"Reminder"}},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Lookup (ID/Reminder)", Argument: "lookup"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Resource", Argument: "resource"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Assigned To", Argument: "assignedTo"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Assigned By", Argument: "assignedBy"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Remind At (Date)", Argument: "remindAt"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Payload (Bytes)", Argument: "payload"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var args = &remindersUpdateArgs{
				hasLookup:     in.Has("lookup"),
				hasResource:   in.Has("resource"),
				hasAssignedTo: in.Has("assignedTo"),
				hasAssignedBy: in.Has("assignedBy"),
				hasRemindAt:   in.Has("remindAt"),
				hasPayload:    in.Has("payload"),
			}

			if err = in.Decode(args); err != nil {
				return
			}

			if args.hasLookup {
				aux := expr.Must(expr.Select(in, "lookup"))
				switch aux.Type() {
				case h.tReg.Type("ID").Type():
					args.lookupID = aux.Get().(uint64)
				case h.tReg.Type("Reminder").Type():
					args.lookupRes = aux.Get().(*types.Reminder)
				}
			}

			var results *remindersUpdateResults
			if results, err = h.h.update(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}
			if tval, err := h.tReg.Type("Reminder").Cast(results.Reminder); err == nil {
				_ = expr.Assign(out, "reminder", tval)
			}
			return
		},
	}
}

func (h ngRemindersHandler) Dismiss() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "remindersDismiss",
		Kind:   "function",
		Groups: []string{"Reminders"},
		Labels: map[string]string{"dismiss": "step", "reminders": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short: "Dismiss Reminder",
			Icon:  &atypes.NgAutomationIcon{Type: "name", Value: "clock"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "lookup", Types: []string{"ID", "Reminder"}, Required: true},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Lookup (ID/Reminder)", Argument: "lookup"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var args = &remindersDismissArgs{
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
				case h.tReg.Type("Reminder").Type():
					args.lookupRes = aux.Get().(*types.Reminder)
				}
			}

			return out, h.h.dismiss(ctx, args)
		},
	}
}

func (h ngRemindersHandler) Snooze() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "remindersSnooze",
		Kind:   "function",
		Groups: []string{"Reminders"},
		Labels: map[string]string{"reminders": "step,workflow", "snooze": "step"},
		Meta: &atypes.ConstructFunctionMeta{
			Short: "Snooze Reminder",
			Icon:  &atypes.NgAutomationIcon{Type: "name", Value: "clock"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "lookup", Types: []string{"ID", "Reminder"}, Required: true},
			{ArgumentName: "remindAt", Types: []string{"DateTime"}, Required: true},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Lookup (ID/Reminder)", Argument: "lookup"}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Remind At (Date)", Argument: "remindAt"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var args = &remindersSnoozeArgs{
				hasLookup:   in.Has("lookup"),
				hasRemindAt: in.Has("remindAt"),
			}

			if err = in.Decode(args); err != nil {
				return
			}

			if args.hasLookup {
				aux := expr.Must(expr.Select(in, "lookup"))
				switch aux.Type() {
				case h.tReg.Type("ID").Type():
					args.lookupID = aux.Get().(uint64)
				case h.tReg.Type("Reminder").Type():
					args.lookupRes = aux.Get().(*types.Reminder)
				}
			}

			return out, h.h.snooze(ctx, args)
		},
	}
}

func (h ngRemindersHandler) Delete() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "remindersDelete",
		Kind:   "function",
		Groups: []string{"Reminders"},
		Labels: map[string]string{"delete": "step", "reminders": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short: "Delete Reminder",
			Icon:  &atypes.NgAutomationIcon{Type: "name", Value: "clock"},
		},

		Parameters: []*atypes.Param{
			{ArgumentName: "lookup", Types: []string{"ID", "Reminder"}, Required: true},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Lookup (ID/Reminder)", Argument: "lookup"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var args = &remindersDeleteArgs{
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
				case h.tReg.Type("Reminder").Type():
					args.lookupRes = aux.Get().(*types.Reminder)
				}
			}

			return out, h.h.delete(ctx, args)
		},
	}
}
