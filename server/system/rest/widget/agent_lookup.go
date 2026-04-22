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

// AgentByKeyLookup resolves an Agent from its widget key. The agent set is
// small and mutations go through the admin API, so a plain mutex-guarded map
// with invalidation on miss is enough.
type AgentByKeyLookup struct {
	mu    sync.RWMutex
	cache map[string]uint64 // widgetKey -> agentID
	store store.Agents
}

func NewAgentByKeyLookup(s store.Agents) *AgentByKeyLookup {
	return &AgentByKeyLookup{cache: map[string]uint64{}, store: s}
}

func (l *AgentByKeyLookup) Invalidate(widgetKey string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.cache, widgetKey)
}

func (l *AgentByKeyLookup) Find(ctx context.Context, widgetKey string) (*types.Agent, error) {
	if widgetKey == "" {
		return nil, fmt.Errorf("widget: missing key")
	}

	l.mu.RLock()
	id, hit := l.cache[widgetKey]
	l.mu.RUnlock()

	// query under service identity so visibility isn't governed by caller's RBAC
	ctx = auth.SetIdentityToContext(ctx, auth.ServiceUser())

	if hit {
		a, err := store.LookupAgentByID(ctx, l.store, id)
		if err == nil && a != nil && a.DeletedAt == nil && a.Chatbot.WidgetKey == widgetKey {
			return a, nil
		}
		l.Invalidate(widgetKey)
	}

	set, _, err := store.SearchAgents(ctx, l.store, types.AgentFilter{Deleted: filter.StateExcluded})
	if err != nil {
		return nil, err
	}
	for _, a := range set {
		if a.Chatbot.WidgetKey == widgetKey {
			l.mu.Lock()
			l.cache[widgetKey] = a.ID
			l.mu.Unlock()
			return a, nil
		}
	}
	return nil, fmt.Errorf("widget: unknown key")
}
