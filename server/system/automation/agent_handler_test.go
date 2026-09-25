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
	agent *types.Agent
	err   error
}

func (m *mockAgentService) FindByID(_ context.Context, _ uint64) (*types.Agent, error) {
	return m.agent, m.err
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
