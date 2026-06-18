package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"

	"github.com/crusttech/human/server/pkg/actionlog"
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
		res = new
		return svc.onCreate(ctx, new)
	}()

	return res, svc.recordAction(ctx, aProps, ChatbotActionCreate, err)
}

func (svc *chatbot) Update(ctx context.Context, upd *types.Chatbot) (res *types.Chatbot, err error) {
	var (
		aProps = &chatbotActionProps{update: upd}
	)

	err = func() (err error) {
		res, err = svc.onUpdate(ctx, upd, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, ChatbotActionUpdate, err)
}

func (svc *chatbot) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &chatbotActionProps{}
	)

	err = func() (err error) {
		return svc.onDelete(ctx, ID, aProps)
	}()

	return svc.recordAction(ctx, aProps, ChatbotActionDelete, err)
}

func (svc *chatbot) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &chatbotActionProps{}
	)

	err = func() (err error) {
		return svc.onUndelete(ctx, ID, aProps)
	}()

	return svc.recordAction(ctx, aProps, ChatbotActionUndelete, err)
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
