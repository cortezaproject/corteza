package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"

	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/store"
	types "github.com/crusttech/human/server/system/types"
)

type chatbotAccessController interface {
	CanCreateChatbot(context.Context) bool
	CanSearchChatbots(context.Context) bool
	CanReadChatbot(context.Context, *types.Chatbot) bool
	CanUpdateChatbot(context.Context, *types.Chatbot) bool
	CanDeleteChatbot(context.Context, *types.Chatbot) bool
}

type chatbot struct {
	actionlog actionlog.Recorder
	store     store.Storer
	ac        chatbotAccessController
}

func Chatbot() *chatbot {
	return &chatbot{
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		ac:        DefaultAccessControl,
	}
}

func (svc *chatbot) FindByID(ctx context.Context, ID uint64) (res *types.Chatbot, err error) {
	var (
		aProps = &chatbotActionProps{chatbot: &types.Chatbot{ID: ID}}
	)

	err = func() error {
		res, err = svc.onLookup(ctx, ID, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, ChatbotActionLookup, err)
}

func (svc *chatbot) Search(ctx context.Context, filter types.ChatbotFilter) (set types.ChatbotSet, f types.ChatbotFilter, err error) {
	var (
		aProps = &chatbotActionProps{filter: &filter}
	)

	err = func() error {
		if !svc.ac.CanSearchChatbots(ctx) {
			return ChatbotErrNotAllowedToSearch()
		}
		set, f, err = svc.onSearch(ctx, filter, aProps)
		return err
	}()

	return set, f, svc.recordAction(ctx, aProps, ChatbotActionSearch, err)
}

func (svc *chatbot) Create(ctx context.Context, new *types.Chatbot) (res *types.Chatbot, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &chatbotActionProps{chatbot: new, new: new}
	)

	err = func() (err error) {
		if !svc.ac.CanCreateChatbot(ctx) {
			return ChatbotErrNotAllowedToCreate()
		}
		before := func() error { return nil }
		after := func() error { return nil }
		res = new
		return svc.onCreate(ctx, new, before, after)
	}()

	return res, svc.recordAction(ctx, aProps, ChatbotActionCreate, err)
}

func (svc *chatbot) Update(ctx context.Context, upd *types.Chatbot) (res *types.Chatbot, err error) {
	var (
		aProps = &chatbotActionProps{update: upd}
		old    *types.Chatbot
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadChatbot(ctx, s, upd.ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setChatbot(res)
		aProps.setUpdate(res)
		old = res.Clone()

		if upd.Handle != res.Handle && !handle.IsValid(upd.Handle) {
			return ChatbotErrInvalidHandle()
		}

		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return ChatbotErrStaleData()
		}
		before := func() error { return nil }
		after := func() error { return nil }

		if err = svc.onUpdate(ctx, s, upd, res, aProps, before, after); err != nil {
			return
		}
		res.Handle = upd.Handle
		res.Name = upd.Name
		res.Enabled = upd.Enabled
		res.WidgetKey = upd.WidgetKey
		res.SessionTTL = upd.SessionTTL
		res.CreatedBy = upd.CreatedBy
		res.UpdatedBy = upd.UpdatedBy
		res.DeletedBy = upd.DeletedBy
		res.UpdatedAt = now()

		if err = store.UpdateChatbot(ctx, s, res); err != nil {
			return err
		}

		if label.Changed(res.Labels, upd.Labels) {
			if err = label.Update(ctx, s, upd); err != nil {
				return
			}
			res.Labels = upd.Labels
		}

		return nil
	})

	return res, svc.recordAction(ctx, aProps, ChatbotActionUpdate, err, old, res)
}

func (svc *chatbot) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &chatbotActionProps{}
		res    *types.Chatbot
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadChatbot(ctx, s, ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setChatbot(res)

		return svc.onDelete(ctx, s, res, aProps)
	})

	return svc.recordAction(ctx, aProps, ChatbotActionDelete, err)
}

func (svc *chatbot) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &chatbotActionProps{}
		res    *types.Chatbot
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadChatbot(ctx, s, ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setChatbot(res)

		return svc.onUndelete(ctx, s, res, aProps)
	})

	return svc.recordAction(ctx, aProps, ChatbotActionUndelete, err)
}

func loadChatbot(ctx context.Context, s store.Chatbots, ID uint64) (res *types.Chatbot, err error) {
	if ID == 0 {
		return nil, ChatbotErrInvalidID()
	}

	if res, err = store.LookupChatbotByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, ChatbotErrNotFound()
	}

	return
}

// toLabeledChatbots converts to []label.LabeledResource
func toLabeledChatbots(set []*types.Chatbot) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}
