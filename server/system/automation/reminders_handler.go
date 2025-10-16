package automation

import (
	"context"
	"fmt"
	"time"
	intAuth "github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	. "github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/wfexec"
	"github.com/cortezaproject/corteza/server/system/types"
	sqlxtypes "github.com/jmoiron/sqlx/types"
)


type(
	reminderService interface {
		FindByID(ctx context.Context,ID uint64) (*types.Reminder, error)
		Find(ctx context.Context, filter types.ReminderFilter) (types.ReminderSet,types.ReminderFilter,error)
		Create(ctx context.Context, reminder *types.Reminder) (*types.Reminder,error)
		Update(ctx context.Context, reminder *types.Reminder) (*types.Reminder,error)
		Delete(ctx context.Context, ID uint64) error
		Dismiss(ctx context.Context, ID uint64) error
		Snooze(ctx context.Context, ID uint64,remindAt *time.Time) error
	}

	remindersHandler struct {
		reg remindersHandlerRegistry
		rSvc reminderService
	}

	reminderSetIterator struct {
		// Item buffer, current item pointer, and total items traversed
		ptr    uint
		buffer types.ReminderSet
		total  uint

		// When filter limit is set, this constraints it
		iterLimit    uint
		useIterLimit bool

		// Item loader for additional chunks
		filter types.ReminderFilter
		loader func() error
	}
	reminderLookup interface {
		GetLookup() (bool, uint64, *types.Reminder)
	}

)
func  RemindersHandler(reg remindersHandlerRegistry,rSvc reminderService) *remindersHandler {
	h := &remindersHandler{
		reg : reg,
		rSvc: rSvc,
	}
	h.register()
	return h

}

func (h remindersHandler) lookup(ctx context.Context, args *remindersLookupArgs) (results *remindersLookupResults, err error){
	results = &remindersLookupResults{}
	results.Reminder, err = lookupReminder(ctx, h.rSvc, args)
	return
}

func (h remindersHandler) search(ctx context.Context,args *remindersSearchArgs) (results *remindersSearchResults,err error){
	results = &remindersSearchResults{}
	var (
		f=types.ReminderFilter{
			Resource: args.Resource,
			AssignedTo: args.AssignedTo,
			ExcludeDismissed: args.ExcludeDismissed,
			ScheduledOnly: args.ScheduledOnly,
		}
	)

	if args.hasSort {
		if err = f.Sort.Set(args.Sort); err != nil {
			return
		}
	}

	if args.hasPageCursor {
		if err = f.PageCursor.Decode(args.PageCursor); err != nil {
			return
		}
	}
	
	if args.hasLimit {
		f.Limit = uint(args.Limit)
	}
	var auxf types.ReminderFilter
	results.Reminders, auxf, err = h.rSvc.Find(ctx, f)
	results.Total=uint64(auxf.Total)
	return
}
func (h remindersHandler) each(ctx context.Context, args *remindersEachArgs) (out wfexec.IteratorHandler, err error){
	var(
		i = &reminderSetIterator{}
		f=types.ReminderFilter{
			Resource: args.Resource,
			AssignedTo: args.AssignedTo,
			ExcludeDismissed: args.ExcludeDismissed,
			ScheduledOnly: args.ScheduledOnly,
		}
	)
	if args.hasSort {
		if err = f.Sort.Set(args.Sort); err != nil {
			return
		}
	}
	if args.hasPageCursor {
		if err = f.PageCursor.Decode(args.PageCursor); err != nil {
			return
		}
	}   
	if args.hasLimit {
		i.useIterLimit = true
		i.iterLimit = uint(args.Limit)
		f.Limit = uint(args.Limit)
		if args.Limit > uint64(wfexec.MaxIteratorBufferSize) {
			f.Limit = wfexec.MaxIteratorBufferSize
		}
		i.iterLimit = uint(args.Limit)
	} else {
		f.Limit = wfexec.MaxIteratorBufferSize
	}
	i.filter = f
	i.loader = func() (err error) {
		if i.filter.PageCursor != nil && i.filter.NextPage == nil {
			return
		}
		i.total += i.ptr
		i.ptr = 0
		i.filter.PageCursor = i.filter.NextPage
		i.filter.NextPage = nil
		i.buffer, i.filter, err = h.rSvc.Find(ctx, i.filter)
		return
	}
	return i, i.loader()
}

 func (h remindersHandler) create(ctx context.Context, args *remindersCreateArgs) (results *remindersCreateResults, err error) {
	
  	results = &remindersCreateResults{}
	currentUser := intAuth.GetIdentityFromContext(ctx).Identity()

	reminder := &types.Reminder{
		Resource:args.Resource,
		AssignedTo: args.AssignedTo,
		AssignedBy: args.AssignedBy,
		RemindAt: args.RemindAt,
		Payload: args.Payload,
	}
	if reminder.AssignedTo == 0 {
		reminder.AssignedTo = currentUser
	}
	if reminder.AssignedBy == 0 {
		reminder.AssignedBy = currentUser
	}
	if len(reminder.Payload) == 0 {
      reminder.Payload = sqlxtypes.JSONText(`{"title": "", "note": ""}`)
  	}
	if reminder.RemindAt == nil {
      return nil, fmt.Errorf("remindAt is required")
  	}
  	results.Reminder, err = h.rSvc.Create(ctx, reminder)
  	return
  }

  func (h remindersHandler) update(ctx context.Context, args *remindersUpdateArgs) (results *remindersUpdateResults, err error) {
  	results = &remindersUpdateResults{}
	reminder, err := lookupReminder(ctx, h.rSvc, args)
	if err != nil {
		return nil, err
	}
	if args.hasResource {
		reminder.Resource = args.Resource
	}
	if args.hasAssignedTo {
		reminder.AssignedTo = args.AssignedTo
	}
	if args.hasAssignedBy {
		reminder.AssignedBy = args.AssignedBy
	}
	if args.hasRemindAt {
		reminder.RemindAt = args.RemindAt
	}
	if args.hasPayload {
		reminder.Payload = args.Payload
	}
  	results.Reminder, err = h.rSvc.Update(ctx, reminder)
  	return
  }
  func (h remindersHandler) dismiss(ctx context.Context, args *remindersDismissArgs) error {
  	if id, err := getReminderID(ctx, h.rSvc, args); err != nil {
  		return err
  	} else {
  		return h.rSvc.Dismiss(ctx, id)
  	}
  }
  func (h remindersHandler) snooze(ctx context.Context, args *remindersSnoozeArgs) error {
  	if id, err := getReminderID(ctx, h.rSvc, args); err != nil {
  		return err
  	} else {
  		return h.rSvc.Snooze(ctx, id, args.RemindAt)
  	}
  }
  func (h remindersHandler) delete(ctx context.Context, args *remindersDeleteArgs) error {
  	if id, err := getReminderID(ctx, h.rSvc, args); err != nil {
  		return err
  	} else {
  		return h.rSvc.Delete(ctx, id)
  	}
  }
  func getReminderID(ctx context.Context, svc reminderService, args reminderLookup) (uint64, error) {
	_, ID, reminder := args.GetLookup()

	switch {
	case reminder != nil:
		return reminder.ID, nil
	case ID > 0:
		return ID, nil
	}

	lookupResult, err := lookupReminder(ctx, svc, args)
	if err != nil {
		return 0, err
	}
	return lookupResult.ID, nil
}

func lookupReminder(ctx context.Context, svc reminderService, args reminderLookup) (*types.Reminder, error) {
	_, ID, reminder := args.GetLookup()
	switch {
	case reminder != nil:
		return reminder, nil
	case ID > 0:
		return svc.FindByID(ctx, ID)
	}
	return nil, fmt.Errorf("empty lookup params")
}

func (i *reminderSetIterator) More(context.Context, *expr.Vars) (bool, error) {
	a := wfexec.GenericResourceNextCheck(i.useIterLimit, i.ptr, uint(len(i.buffer)), i.total, i.iterLimit, i.filter.NextPage != nil)
	return a, nil
}
func (i *reminderSetIterator) Start(context.Context, *expr.Vars) error {
	i.ptr = 0
	return nil
}
func (i *reminderSetIterator) Next(context.Context, *expr.Vars) (out *expr.Vars, err error) {
	if len(i.buffer)-int(i.ptr) <= 0 {
		if err = i.loader(); err != nil {
			panic(err)
		}
	}

	out = &expr.Vars{}
	reminder := *i.buffer[i.ptr]  // Make a copy
  	reminder.Payload = sqlxtypes.JSONText(string(i.buffer[i.ptr].Payload))
  	out.Set("reminder", Must(NewReminder(&reminder)))

	i.ptr++
	return out, nil
}
