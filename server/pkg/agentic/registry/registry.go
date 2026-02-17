package registry

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	registry struct {
		mu     sync.RWMutex
		agents map[uint64]*types.Agent
	}

	ValidationError struct {
		Field string
		Error string
	}
)

func Registry() *registry {
	return &registry{
		agents: make(map[uint64]*types.Agent),
	}
}

func (r *registry) Add(ctx context.Context, agent *types.Agent) (*types.Agent, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if agent.ID == 0 {
		return nil, fmt.Errorf("agent ID is required")
	}

	if _, exists := r.agents[agent.ID]; exists {
		return nil, fmt.Errorf("agent with ID %d already exists", agent.ID)
	}

	// Store copy
	stored := *agent
	r.agents[agent.ID] = &stored

	return agent, nil
}

func (r *registry) Remove(ctx context.Context, id uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	nid := id
	existing, exists := r.agents[nid]
	if !exists {
		return fmt.Errorf("agent not found: %d", id)
	}

	if existing.Status == "active" {
		return fmt.Errorf("cannot remove active agent, disable it first")
	}

	delete(r.agents, nid)
	return nil
}

func (r *registry) Get(ctx context.Context, id uint64) (*types.Agent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	agent, exists := r.agents[id]
	if !exists {
		return nil, fmt.Errorf("agent not found: %d", id)
	}
	return agent, nil
}

func (r *registry) Activate(ctx context.Context, id uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	agent, exists := r.agents[id]
	if !exists {
		return fmt.Errorf("agent not found: %d", id)
	}

	if agent.Status == "active" {
		return nil
	}

	agent.Status = "active"
	now := time.Now()
	agent.UpdatedAt = &now
	return nil
}

func (r *registry) Disable(ctx context.Context, id uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	agent, exists := r.agents[id]
	if !exists {
		return fmt.Errorf("agent not found: %d", id)
	}

	if agent.Status == "disabled" {
		return nil
	}

	agent.Status = "disabled"
	now := time.Now()
	agent.UpdatedAt = &now
	return nil
}
