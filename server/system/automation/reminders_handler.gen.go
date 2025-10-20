package automation

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/automation/reminders_handler.yaml

import (
	"context"
	atypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/wfexec"
	"github.com/cortezaproject/corteza/server/system/types"
	"time"
)

var _ wfexec.ExecResponse

type (
	remindersHandlerRegistry interface {
		AddFunctions(ff ...*atypes.Function)
		Type(ref string) expr.Type
	}
)

func (h remindersHandler) register() {
	h.reg.AddFunctions(
		h.Lookup(),
		h.Search(),
		h.Each(),
		h.Create(),
		h.Update(),
		h.Dismiss(),
		h.Snooze(),
		h.Delete(),
	)
}

type (
	remindersLookupArgs struct {
		hasLookup bool
		Lookup    interface{}
		lookupID  uint64
		lookupRes *types.Reminder
	}

	remindersLookupResults struct {
		Reminder *types.Reminder
	}
)

func (a remindersLookupArgs) GetLookup() (bool, uint64, *types.Reminder) {
	return a.hasLookup, a.lookupID, a.lookupRes
}

// Lookup function Reminder lookup
//
// expects implementation of lookup function:
//
//	func (h remindersHandler) lookup(ctx context.Context, args *remindersLookupArgs) (results *remindersLookupResults, err error) {
//	   return
//	}
func (h remindersHandler) Lookup() *atypes.Function {
	return &atypes.Function{
		Ref:    "remindersLookup",
		Kind:   "function",
		Labels: map[string]string{"reminders": "step,workflow"},
		Meta: &atypes.FunctionMeta{
			Short:       "Reminder lookup",
			Description: "Find specific reminder by ID",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "lookup",
				Types: []string{"ID", "Reminder"}, Required: true,
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "reminder",
				Types: []string{"Reminder"},
			},
		},

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
				case h.reg.Type("ID").Type():
					args.lookupID = aux.Get().(uint64)
				case h.reg.Type("Reminder").Type():
					args.lookupRes = aux.Get().(*types.Reminder)
				}
			}

			var results *remindersLookupResults
			if results, err = h.lookup(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Reminder (*types.Reminder) to Reminder
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("Reminder").Cast(results.Reminder); err != nil {
					return
				} else if err = expr.Assign(out, "reminder", tval); err != nil {
					return
				}
			}

			return
		},
	}
}

type (
	remindersSearchArgs struct {
		hasResource bool
		Resource    string

		hasAssignedTo bool
		AssignedTo    uint64

		hasExcludeDismissed bool
		ExcludeDismissed    bool

		hasScheduledOnly bool
		ScheduledOnly    bool

		hasSort bool
		Sort    string

		hasLimit bool
		Limit    uint64

		hasPageCursor bool
		PageCursor    string
	}

	remindersSearchResults struct {
		Reminders []*types.Reminder
		Total     uint64
	}
)

// Search function Reminders search
//
// expects implementation of search function:
//
//	func (h remindersHandler) search(ctx context.Context, args *remindersSearchArgs) (results *remindersSearchResults, err error) {
//	   return
//	}
func (h remindersHandler) Search() *atypes.Function {
	return &atypes.Function{
		Ref:    "remindersSearch",
		Kind:   "function",
		Labels: map[string]string{"reminders": "step,workflow"},
		Meta: &atypes.FunctionMeta{
			Short:       "Reminders search",
			Description: "Search reminders",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "resource",
				Types: []string{"String"},
			},
			{
				Name:  "assignedTo",
				Types: []string{"ID"},
			},
			{
				Name:  "excludeDismissed",
				Types: []string{"Boolean"},
			},
			{
				Name:  "scheduledOnly",
				Types: []string{"Boolean"},
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
				Name:  "pageCursor",
				Types: []string{"String"},
			},
		},

		Results: []*atypes.Param{

			{
				Name:    "reminders",
				Types:   []string{"Reminder"},
				IsArray: true,
			},

			{
				Name:  "total",
				Types: []string{"UnsignedInteger"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &remindersSearchArgs{
					hasResource:         in.Has("resource"),
					hasAssignedTo:       in.Has("assignedTo"),
					hasExcludeDismissed: in.Has("excludeDismissed"),
					hasScheduledOnly:    in.Has("scheduledOnly"),
					hasSort:             in.Has("sort"),
					hasLimit:            in.Has("limit"),
					hasPageCursor:       in.Has("pageCursor"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			var results *remindersSearchResults
			if results, err = h.search(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Reminders (*types.Reminder) to Array (of Reminder)
				var (
					tval expr.TypedValue
					tarr = make([]expr.TypedValue, len(results.Reminders))
				)

				for i := range results.Reminders {
					if tarr[i], err = h.reg.Type("Reminder").Cast(results.Reminders[i]); err != nil {
						return
					}
				}

				if tval, err = expr.NewArray(tarr); err != nil {
					return
				} else if err = expr.Assign(out, "reminders", tval); err != nil {
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

			return
		},
	}
}

type (
	remindersEachArgs struct {
		hasResource bool
		Resource    string

		hasAssignedTo bool
		AssignedTo    uint64

		hasExcludeDismissed bool
		ExcludeDismissed    bool

		hasScheduledOnly bool
		ScheduledOnly    bool

		hasSort bool
		Sort    string

		hasLimit bool
		Limit    uint64

		hasPageCursor bool
		PageCursor    string
	}

	remindersEachResults struct {
		Reminder *types.Reminder
		Total    uint64
	}
)

// Each function Iterate reminders
//
// expects implementation of each function:
//
//	func (h remindersHandler) each(ctx context.Context, args *remindersEachArgs) (results *remindersEachResults, err error) {
//	   return
//	}
func (h remindersHandler) Each() *atypes.Function {
	return &atypes.Function{
		Ref:    "remindersEach",
		Kind:   "iterator",
		Labels: map[string]string(nil),
		Meta: &atypes.FunctionMeta{
			Short: "Iterate reminders",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "resource",
				Types: []string{"String"},
			},
			{
				Name:  "assignedTo",
				Types: []string{"ID"},
			},
			{
				Name:  "excludeDismissed",
				Types: []string{"Boolean"},
			},
			{
				Name:  "scheduledOnly",
				Types: []string{"Boolean"},
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
				Name:  "pageCursor",
				Types: []string{"String"},
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "reminder",
				Types: []string{"Reminder"},
			},

			{
				Name:  "total",
				Types: []string{"UnsignedInteger"},
			},
		},

		Iterator: func(ctx context.Context, in *expr.Vars) (out wfexec.IteratorHandler, err error) {
			var (
				args = &remindersEachArgs{
					hasResource:         in.Has("resource"),
					hasAssignedTo:       in.Has("assignedTo"),
					hasExcludeDismissed: in.Has("excludeDismissed"),
					hasScheduledOnly:    in.Has("scheduledOnly"),
					hasSort:             in.Has("sort"),
					hasLimit:            in.Has("limit"),
					hasPageCursor:       in.Has("pageCursor"),
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
	remindersCreateArgs struct {
		hasResource bool
		Resource    string

		hasAssignedTo bool
		AssignedTo    uint64

		hasAssignedBy bool
		AssignedBy    uint64

		hasRemindAt bool
		RemindAt    *time.Time

		hasPayload bool
		Payload    []byte
	}

	remindersCreateResults struct {
		Reminder *types.Reminder
	}
)

// Create function Reminder create
//
// expects implementation of create function:
//
//	func (h remindersHandler) create(ctx context.Context, args *remindersCreateArgs) (results *remindersCreateResults, err error) {
//	   return
//	}
func (h remindersHandler) Create() *atypes.Function {
	return &atypes.Function{
		Ref:    "remindersCreate",
		Kind:   "function",
		Labels: map[string]string{"create": "step", "reminders": "step,workflow"},
		Meta: &atypes.FunctionMeta{
			Short: "Reminder create",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "resource",
				Types: []string{"String"},
			},
			{
				Name:  "assignedTo",
				Types: []string{"ID"},
			},
			{
				Name:  "assignedBy",
				Types: []string{"ID"},
			},
			{
				Name:  "remindAt",
				Types: []string{"DateTime"}, Required: true,
			},
			{
				Name:  "payload",
				Types: []string{"Bytes"},
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "reminder",
				Types: []string{"Reminder"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &remindersCreateArgs{
					hasResource:   in.Has("resource"),
					hasAssignedTo: in.Has("assignedTo"),
					hasAssignedBy: in.Has("assignedBy"),
					hasRemindAt:   in.Has("remindAt"),
					hasPayload:    in.Has("payload"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			var results *remindersCreateResults
			if results, err = h.create(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Reminder (*types.Reminder) to Reminder
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("Reminder").Cast(results.Reminder); err != nil {
					return
				} else if err = expr.Assign(out, "reminder", tval); err != nil {
					return
				}
			}

			return
		},
	}
}

type (
	remindersUpdateArgs struct {
		hasLookup bool
		Lookup    interface{}
		lookupID  uint64
		lookupRes *types.Reminder

		hasResource bool
		Resource    string

		hasAssignedTo bool
		AssignedTo    uint64

		hasAssignedBy bool
		AssignedBy    uint64

		hasRemindAt bool
		RemindAt    *time.Time

		hasPayload bool
		Payload    []byte
	}

	remindersUpdateResults struct {
		Reminder *types.Reminder
	}
)

func (a remindersUpdateArgs) GetLookup() (bool, uint64, *types.Reminder) {
	return a.hasLookup, a.lookupID, a.lookupRes
}

// Update function Reminder update
//
// expects implementation of update function:
//
//	func (h remindersHandler) update(ctx context.Context, args *remindersUpdateArgs) (results *remindersUpdateResults, err error) {
//	   return
//	}
func (h remindersHandler) Update() *atypes.Function {
	return &atypes.Function{
		Ref:    "remindersUpdate",
		Kind:   "function",
		Labels: map[string]string{"reminders": "step,workflow", "update": "step"},
		Meta: &atypes.FunctionMeta{
			Short: "Reminder update",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "lookup",
				Types: []string{"ID", "Reminder"}, Required: true,
			},
			{
				Name:  "resource",
				Types: []string{"String"},
			},
			{
				Name:  "assignedTo",
				Types: []string{"ID"},
			},
			{
				Name:  "assignedBy",
				Types: []string{"ID"},
			},
			{
				Name:  "remindAt",
				Types: []string{"DateTime"},
			},
			{
				Name:  "payload",
				Types: []string{"Bytes"},
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "reminder",
				Types: []string{"Reminder"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &remindersUpdateArgs{
					hasLookup:     in.Has("lookup"),
					hasResource:   in.Has("resource"),
					hasAssignedTo: in.Has("assignedTo"),
					hasAssignedBy: in.Has("assignedBy"),
					hasRemindAt:   in.Has("remindAt"),
					hasPayload:    in.Has("payload"),
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
				case h.reg.Type("Reminder").Type():
					args.lookupRes = aux.Get().(*types.Reminder)
				}
			}

			var results *remindersUpdateResults
			if results, err = h.update(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Reminder (*types.Reminder) to Reminder
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("Reminder").Cast(results.Reminder); err != nil {
					return
				} else if err = expr.Assign(out, "reminder", tval); err != nil {
					return
				}
			}

			return
		},
	}
}

type (
	remindersDismissArgs struct {
		hasLookup bool
		Lookup    interface{}
		lookupID  uint64
		lookupRes *types.Reminder
	}
)

func (a remindersDismissArgs) GetLookup() (bool, uint64, *types.Reminder) {
	return a.hasLookup, a.lookupID, a.lookupRes
}

// Dismiss function Dismiss reminder
//
// expects implementation of dismiss function:
//
//	func (h remindersHandler) dismiss(ctx context.Context, args *remindersDismissArgs) (err error) {
//	   return
//	}
func (h remindersHandler) Dismiss() *atypes.Function {
	return &atypes.Function{
		Ref:    "remindersDismiss",
		Kind:   "function",
		Labels: map[string]string{"dismiss": "step", "reminders": "step,workflow"},
		Meta: &atypes.FunctionMeta{
			Short: "Dismiss reminder",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "lookup",
				Types: []string{"ID", "Reminder"}, Required: true,
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &remindersDismissArgs{
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
				case h.reg.Type("Reminder").Type():
					args.lookupRes = aux.Get().(*types.Reminder)
				}
			}

			return out, h.dismiss(ctx, args)
		},
	}
}

type (
	remindersSnoozeArgs struct {
		hasLookup bool
		Lookup    interface{}
		lookupID  uint64
		lookupRes *types.Reminder

		hasRemindAt bool
		RemindAt    *time.Time
	}
)

func (a remindersSnoozeArgs) GetLookup() (bool, uint64, *types.Reminder) {
	return a.hasLookup, a.lookupID, a.lookupRes
}

// Snooze function Snooze reminder
//
// expects implementation of snooze function:
//
//	func (h remindersHandler) snooze(ctx context.Context, args *remindersSnoozeArgs) (err error) {
//	   return
//	}
func (h remindersHandler) Snooze() *atypes.Function {
	return &atypes.Function{
		Ref:    "remindersSnooze",
		Kind:   "function",
		Labels: map[string]string{"reminders": "step,workflow", "snooze": "step"},
		Meta: &atypes.FunctionMeta{
			Short: "Snooze reminder",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "lookup",
				Types: []string{"ID", "Reminder"}, Required: true,
			},
			{
				Name:  "remindAt",
				Types: []string{"DateTime"}, Required: true,
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &remindersSnoozeArgs{
					hasLookup:   in.Has("lookup"),
					hasRemindAt: in.Has("remindAt"),
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
				case h.reg.Type("Reminder").Type():
					args.lookupRes = aux.Get().(*types.Reminder)
				}
			}

			return out, h.snooze(ctx, args)
		},
	}
}

type (
	remindersDeleteArgs struct {
		hasLookup bool
		Lookup    interface{}
		lookupID  uint64
		lookupRes *types.Reminder
	}
)

func (a remindersDeleteArgs) GetLookup() (bool, uint64, *types.Reminder) {
	return a.hasLookup, a.lookupID, a.lookupRes
}

// Delete function Reminder delete
//
// expects implementation of delete function:
//
//	func (h remindersHandler) delete(ctx context.Context, args *remindersDeleteArgs) (err error) {
//	   return
//	}
func (h remindersHandler) Delete() *atypes.Function {
	return &atypes.Function{
		Ref:    "remindersDelete",
		Kind:   "function",
		Labels: map[string]string{"delete": "step", "reminders": "step,workflow"},
		Meta: &atypes.FunctionMeta{
			Short: "Reminder delete",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "lookup",
				Types: []string{"ID", "Reminder"}, Required: true,
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &remindersDeleteArgs{
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
				case h.reg.Type("Reminder").Type():
					args.lookupRes = aux.Get().(*types.Reminder)
				}
			}

			return out, h.delete(ctx, args)
		},
	}
}
