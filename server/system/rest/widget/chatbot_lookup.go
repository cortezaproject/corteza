package widget

import (
	"context"
	"fmt"
	"sync"

	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

// ChatbotByKeyLookup resolves a Chatbot from its widget key.
type ChatbotByKeyLookup struct {
	mu    sync.RWMutex
	cache map[string]uint64
	store store.Chatbots
}

func NewChatbotByKeyLookup(s store.Chatbots) *ChatbotByKeyLookup {
	return &ChatbotByKeyLookup{cache: map[string]uint64{}, store: s}
}

func (l *ChatbotByKeyLookup) Invalidate(widgetKey string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.cache, widgetKey)
}

func (l *ChatbotByKeyLookup) Find(ctx context.Context, widgetKey string) (*types.Chatbot, error) {
	if widgetKey == "" {
		return nil, fmt.Errorf("widget: missing key")
	}

	l.mu.RLock()
	id, hit := l.cache[widgetKey]
	l.mu.RUnlock()

	ctx = auth.SetIdentityToContext(ctx, auth.ServiceUser())

	if hit {
		c, err := store.LookupChatbotByID(ctx, l.store, id)
		if err == nil && c != nil && c.DeletedAt == nil && c.WidgetKey == widgetKey {
			return c, nil
		}
		l.Invalidate(widgetKey)
	}

	c, err := store.LookupChatbotByWidgetKey(ctx, l.store, widgetKey)
	if err == nil && c != nil && c.DeletedAt == nil {
		l.mu.Lock()
		l.cache[widgetKey] = c.ID
		l.mu.Unlock()
		return c, nil
	}

	set, _, err := store.SearchChatbots(ctx, l.store, types.ChatbotFilter{Deleted: filter.StateExcluded})
	if err != nil {
		return nil, err
	}
	for _, c := range set {
		if c.WidgetKey == widgetKey {
			l.mu.Lock()
			l.cache[widgetKey] = c.ID
			l.mu.Unlock()
			return c, nil
		}
	}
	return nil, fmt.Errorf("widget: unknown key")
}
