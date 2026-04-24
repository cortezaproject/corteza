package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"strconv"

	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

type (
	chatbot struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        chatbotAccessController
	}

	chatbotAccessController interface {
		CanCreateChatbot(ctx context.Context) bool
		CanSearchChatbots(ctx context.Context) bool
		CanReadChatbot(ctx context.Context, c *types.Chatbot) bool
		CanUpdateChatbot(ctx context.Context, c *types.Chatbot) bool
		CanDeleteChatbot(ctx context.Context, c *types.Chatbot) bool
	}
)

func Chatbot() *chatbot {
	return &chatbot{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	}
}

func (svc *chatbot) Get(ctx context.Context, ID uint64) (*types.Chatbot, error) {
	return svc.FindByID(ctx, ID)
}

func (svc *chatbot) FindByID(ctx context.Context, ID uint64) (c *types.Chatbot, err error) {
	err = func() error {
		if c, err = loadChatbot(ctx, svc.store, ID); err != nil {
			return err
		}

		if !svc.ac.CanReadChatbot(ctx, c) {
			return ChatbotErrNotAllowedToRead()
		}

		if err = label.Load(ctx, svc.store, c); err != nil {
			return err
		}

		return nil
	}()

	return c, err
}

func (svc *chatbot) Create(ctx context.Context, new *types.Chatbot) (c *types.Chatbot, err error) {
	err = func() (err error) {
		if !svc.ac.CanCreateChatbot(ctx) {
			return ChatbotErrNotAllowedToCreate()
		}

		new.ID = nextID()
		new.CreatedAt = *now()

		if err = prepareChatbotOnCreate(new); err != nil {
			return
		}

		if err = store.CreateChatbot(ctx, svc.store, new); err != nil {
			return
		}

		if err = label.Create(ctx, svc.store, new); err != nil {
			return
		}

		c = new
		return nil
	}()

	return c, err
}

func (svc *chatbot) Update(ctx context.Context, upd *types.Chatbot) (c *types.Chatbot, err error) {
	err = func() (err error) {
		if !svc.ac.CanUpdateChatbot(ctx, upd) {
			return ChatbotErrNotAllowedToUpdate()
		}

		var existing *types.Chatbot
		if existing, err = store.LookupChatbotByID(ctx, svc.store, upd.ID); err != nil {
			return ChatbotErrNotFound()
		}

		if isStale(upd.UpdatedAt, existing.UpdatedAt, existing.CreatedAt) {
			return ChatbotErrStaleData()
		}

		upd.UpdatedAt = now()
		upd.CreatedAt = existing.CreatedAt
		upd.DeletedAt = existing.DeletedAt

		if err = prepareChatbotOnUpdate(upd, existing); err != nil {
			return
		}

		if err = store.UpdateChatbot(ctx, svc.store, upd); err != nil {
			return
		}

		if err = label.Update(ctx, svc.store, upd); err != nil {
			return
		}

		c = upd
		return nil
	}()

	return c, err
}

func (svc *chatbot) RegenerateWidgetKey(ctx context.Context, ID uint64) (c *types.Chatbot, err error) {
	err = func() (err error) {
		var existing *types.Chatbot
		if existing, err = loadChatbot(ctx, svc.store, ID); err != nil {
			return
		}

		if !svc.ac.CanUpdateChatbot(ctx, existing) {
			return ChatbotErrNotAllowedToUpdate()
		}

		key, err := generateChatbotWidgetKey()
		if err != nil {
			return err
		}
		existing.WidgetKey = key
		existing.UpdatedAt = now()

		if err = store.UpdateChatbot(ctx, svc.store, existing); err != nil {
			return
		}

		c = existing
		return nil
	}()

	return c, err
}

func (svc *chatbot) DeleteByID(ctx context.Context, ID uint64) (err error) {
	err = func() (err error) {
		var c *types.Chatbot
		if c, err = loadChatbot(ctx, svc.store, ID); err != nil {
			return
		}

		if !svc.ac.CanDeleteChatbot(ctx, c) {
			return ChatbotErrNotAllowedToDelete()
		}

		c.DeletedAt = now()
		if err = store.UpdateChatbot(ctx, svc.store, c); err != nil {
			return
		}

		svc.softDeleteChatbotAttachments(ctx, ID)

		return nil
	}()

	return err
}

// softDeleteChatbotAttachments marks every chatbot-kind attachment whose
// Meta.Labels["chatbotID"] matches ID as deleted. Errors are logged but don't
// fail the chatbot delete — the attachments are orphaned at worst.
func (svc *chatbot) softDeleteChatbotAttachments(ctx context.Context, ID uint64) {
	idStr := strconv.FormatUint(ID, 10)
	aa, _, err := store.SearchAttachments(ctx, svc.store, types.AttachmentFilter{
		Kind: types.AttachmentKindChatbot,
		Check: func(a *types.Attachment) (bool, error) {
			return a.Meta.Labels["chatbotID"] == idStr, nil
		},
	})
	if err != nil {
		return
	}
	for _, a := range aa {
		a.DeletedAt = now()
		_ = store.UpdateAttachment(ctx, svc.store, a)
	}
}

func (svc *chatbot) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	err = func() (err error) {
		var c *types.Chatbot
		if c, err = loadChatbot(ctx, svc.store, ID); err != nil {
			return
		}

		if !svc.ac.CanDeleteChatbot(ctx, c) {
			return ChatbotErrNotAllowedToDelete()
		}

		c.DeletedAt = nil
		if err = store.UpdateChatbot(ctx, svc.store, c); err != nil {
			return
		}

		return nil
	}()

	return err
}

func (svc *chatbot) Search(ctx context.Context, filter types.ChatbotFilter) (set types.ChatbotSet, f types.ChatbotFilter, err error) {
	filter.Check = func(res *types.Chatbot) (bool, error) {
		if !svc.ac.CanReadChatbot(ctx, res) {
			return false, nil
		}
		return true, nil
	}

	err = func() error {
		if !svc.ac.CanSearchChatbots(ctx) {
			return ChatbotErrNotAllowedToSearch()
		}

		if len(filter.Labels) > 0 {
			filter.LabeledIDs, err = label.Search(
				ctx,
				svc.store,
				types.Chatbot{}.LabelResourceKind(),
				filter.Labels,
			)

			if err != nil {
				return err
			}

			if len(filter.LabeledIDs) == 0 {
				return nil
			}
		}

		if set, f, err = store.SearchChatbots(ctx, svc.store, filter); err != nil {
			return err
		}

		if err = label.Load(ctx, svc.store, toLabeledChatbots(set)...); err != nil {
			return err
		}

		return nil
	}()

	return set, f, err
}

func toLabeledChatbots(set types.ChatbotSet) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}

func prepareChatbotOnCreate(c *types.Chatbot) error {
	if err := validateChatbotScenarios(c.Scenarios); err != nil {
		return err
	}
	if c.WidgetKey == "" {
		k, err := generateChatbotWidgetKey()
		if err != nil {
			return err
		}
		c.WidgetKey = k
	}
	return nil
}

func prepareChatbotOnUpdate(upd, existing *types.Chatbot) error {
	if err := validateChatbotScenarios(upd.Scenarios); err != nil {
		return err
	}
	if upd.WidgetKey == "" {
		upd.WidgetKey = existing.WidgetKey
	}
	if upd.WidgetKey == "" {
		k, err := generateChatbotWidgetKey()
		if err != nil {
			return err
		}
		upd.WidgetKey = k
	}
	return nil
}

func validateChatbotScenarios(ss types.ChatbotScenarios) error {
	for _, s := range ss {
		if s.Type == "conversation" && s.AgentID == 0 {
			return ChatbotErrConversationScenarioMissingAgent()
		}
	}
	return nil
}

func generateChatbotWidgetKey() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
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
