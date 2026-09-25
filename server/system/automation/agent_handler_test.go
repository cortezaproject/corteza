package automation

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	agenticRuntime "github.com/crusttech/human/server/system/agentic/runtime"
	"github.com/crusttech/human/server/system/types"
)

type mockRuntime struct {
	called bool
	req    *agenticRuntime.AgentRequest
}

func (m *mockRuntime) Run(_ context.Context, req *agenticRuntime.AgentRequest) (*agenticRuntime.AgentResponse, error) {
	m.called = true
	m.req = req
	return &agenticRuntime.AgentResponse{Output: "ok", ConversationID: 1}, nil
}

type mockAgentService struct {
	agent    *types.Agent
	err      error
	searched types.AgentFilter
	found    types.AgentSet
}

func (m *mockAgentService) FindByID(_ context.Context, _ uint64) (*types.Agent, error) {
	return m.agent, m.err
}

func (m *mockAgentService) Search(_ context.Context, f types.AgentFilter) (types.AgentSet, types.AgentFilter, error) {
	m.searched = f
	return m.found, f, m.err
}

func TestAgentHandler_AgentNotFound(t *testing.T) {
	req := require.New(t)
	rt := &mockRuntime{}
	h := &agentHandler{
		runtime: rt,
		agents:  &mockAgentService{err: fmt.Errorf("not found")},
	}

	_, err := h.run(context.Background(), &agentRunArgs{AgentID: 999, Input: "hello"})
	req.Error(err)
	req.False(rt.called)
}

func TestAgentHandler_RunsAnyAgentUnattended(t *testing.T) {
	req := require.New(t)
	rt := &mockRuntime{}
	h := &agentHandler{
		runtime: rt,
		agents:  &mockAgentService{agent: &types.Agent{ID: 1}},
	}

	results, err := h.run(context.Background(), &agentRunArgs{AgentID: 1, Input: "hello", ConversationID: 7})
	req.NoError(err)
	req.True(rt.called)
	req.True(rt.req.Unattended)
	req.Equal(uint64(1), rt.req.AgentID)
	req.Equal(uint64(7), rt.req.ConversationID)
	req.Equal("ok", results.Output)
	req.Equal(uint64(1), results.ConversationID)
}

func TestAgentHandler_LookupByID(t *testing.T) {
	req := require.New(t)
	svc := &mockAgentService{agent: &types.Agent{ID: 5, Handle: "five"}}
	h := &agentHandler{agents: svc}

	results, err := h.lookup(context.Background(), &agentLookupArgs{hasLookup: true, lookupID: 5})
	req.NoError(err)
	req.Equal("five", results.Agent.Handle)
	req.Empty(svc.searched.Handle)
}

func TestAgentHandler_LookupByHandleTakesTheFirstMatch(t *testing.T) {
	req := require.New(t)
	svc := &mockAgentService{found: types.AgentSet{{ID: 8, Handle: "twin"}, {ID: 9, Handle: "twin"}}}
	h := &agentHandler{agents: svc}

	results, err := h.lookup(context.Background(), &agentLookupArgs{hasLookup: true, lookupHandle: "twin"})
	req.NoError(err)
	req.Equal(uint64(8), results.Agent.ID)
	req.Equal("twin", svc.searched.Handle)
}

func TestAgentHandler_LookupByHandleMissing(t *testing.T) {
	req := require.New(t)
	h := &agentHandler{agents: &mockAgentService{}}

	_, err := h.lookup(context.Background(), &agentLookupArgs{hasLookup: true, lookupHandle: "nobody"})
	req.ErrorContains(err, `"nobody"`)
}

func TestAgentHandler_LookupPassesAResourceThrough(t *testing.T) {
	req := require.New(t)
	svc := &mockAgentService{err: fmt.Errorf("must not be called")}
	h := &agentHandler{agents: svc}

	results, err := h.lookup(context.Background(), &agentLookupArgs{hasLookup: true, lookupRes: &types.Agent{ID: 3}})
	req.NoError(err)
	req.Equal(uint64(3), results.Agent.ID)
}

func TestAgentHandler_LookupEmpty(t *testing.T) {
	h := &agentHandler{agents: &mockAgentService{}}
	_, err := h.lookup(context.Background(), &agentLookupArgs{})
	require.ErrorContains(t, err, "empty lookup")
}
