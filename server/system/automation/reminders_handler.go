package automation

import (
	"context"
	"fmt"
	. "github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/wfexec"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/spf13/cast"
)


type(
	reminderService interface {
		FindByID(ctx context.Context,ID uint64) (*types.Reminder, error)
		Find(ctx context.Context, filter types.ReminderFilter) (types.ReminderSet,types.ReminderFilter,error)
		Create(ctx context.Context, reminder *types.Reminder) (*types.Reminder,error)
		Update(ctx context.Context, reminder *types.Reminder) (*types.Reminder,error)
		Delete(ctx context.Context, ID uint64) error
		Dismiss(ctx context.Context, ID uint64) error
		Snooze(ctx context.Context, ID uint64) error
	}

	remindersHandler struct {
		reg remindersHandlerRegistry
		rSvc reminderService
	}

	reminderSetIterator struct {
		// Item buffer, current item pointer, and total items traversed
		ptr    uint
		buffer types.RoleSet
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

func lookupReminder(ctx context.Context, svc reminderService,args reminderLookup) (*types.Reminder,error){
	_, ID, reminder := args.GetLookup()
	switch {
	case reminder != nil:
		return reminder,nil 
	case ID>0:
		return svc.FindByID(ctx,ID)
	}
	return nil,fmt.Errorf("empty lookup params")
}