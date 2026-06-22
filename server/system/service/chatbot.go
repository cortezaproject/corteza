package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"strconv"

	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

func (svc *chatbot) Get(ctx context.Context, ID uint64) (*types.Chatbot, error) {
	return svc.FindByID(ctx, ID)
}

func (svc *chatbot) onLookup(ctx context.Context, ID uint64, aProps *chatbotActionProps) (c *types.Chatbot, err error) {
	if c, err = loadChatbot(ctx, svc.store, ID); err != nil {
		return nil, err
	}

	aProps.setChatbot(c)

	if !svc.ac.CanReadChatbot(ctx, c) {
		return nil, ChatbotErrNotAllowedToRead()
	}

	if err = label.Load(ctx, svc.store, c); err != nil {
		return nil, err
	}

	return c, nil
}

func (svc *chatbot) onCreate(ctx context.Context, new *types.Chatbot) (err error) {
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

	return nil
}

func (svc *chatbot) onUpdate(ctx context.Context, s store.Storer, upd, c *types.Chatbot, aProps *chatbotActionProps, _ func() error, _ func() error) error {
	if !svc.ac.CanUpdateChatbot(ctx, upd) {
		return ChatbotErrNotAllowedToUpdate()
	}

	aProps.setChatbot(c)

	upd.UpdatedAt = now()
	upd.CreatedAt = c.CreatedAt
	upd.DeletedAt = c.DeletedAt

	if err := prepareChatbotOnUpdate(upd, c); err != nil {
		return err
	}

	if err := store.UpdateChatbot(ctx, s, upd); err != nil {
		return err
	}

	// reflect updated record back into res
	*c = *upd

	return label.Update(ctx, s, upd)
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

func (svc *chatbot) onDelete(ctx context.Context, s store.Storer, c *types.Chatbot, aProps *chatbotActionProps) error {
	aProps.setChatbot(c)

	if !svc.ac.CanDeleteChatbot(ctx, c) {
		return ChatbotErrNotAllowedToDelete()
	}

	c.DeletedAt = now()
	if err := store.UpdateChatbot(ctx, s, c); err != nil {
		return err
	}

	svc.softDeleteChatbotAttachments(ctx, c.ID)

	return nil
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

func (svc *chatbot) onUndelete(ctx context.Context, s store.Storer, c *types.Chatbot, aProps *chatbotActionProps) error {
	aProps.setChatbot(c)

	if !svc.ac.CanDeleteChatbot(ctx, c) {
		return ChatbotErrNotAllowedToDelete()
	}

	c.DeletedAt = nil
	return store.UpdateChatbot(ctx, s, c)
}

func (svc *chatbot) onSearch(ctx context.Context, filter types.ChatbotFilter, aProps *chatbotActionProps) (set types.ChatbotSet, f types.ChatbotFilter, err error) {
	filter.Check = func(res *types.Chatbot) (bool, error) {
		if !svc.ac.CanReadChatbot(ctx, res) {
			return false, nil
		}
		return true, nil
	}

	if len(filter.Labels) > 0 {
		filter.LabeledIDs, err = label.Search(
			ctx,
			svc.store,
			types.Chatbot{}.LabelResourceKind(),
			filter.Labels,
		)

		if err != nil {
			return set, f, err
		}

		if len(filter.LabeledIDs) == 0 {
			return set, f, nil
		}
	}

	if set, f, err = store.SearchChatbots(ctx, svc.store, filter); err != nil {
		return set, f, err
	}

	if err = label.Load(ctx, svc.store, toLabeledChatbots(set)...); err != nil {
		return set, f, err
	}

	return set, f, nil
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

