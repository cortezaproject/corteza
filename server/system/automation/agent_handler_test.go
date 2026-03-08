package automation

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	agenticRuntime "github.com/cortezaproject/corteza/server/system/agentic/runtime"
	"github.com/cortezaproject/corteza/server/system/types"
)

type mockRuntime struct {
	called bool
}

func (m *mockRuntime) Run(_ context.Context, _ *agenticRuntime.AgentRequest) (*agenticRuntime.AgentResponse, error) {
	m.called = true
	return &agenticRuntime.AgentResponse{Output: "ok", ConversationID: 1}, nil
}

type mockAgentService struct {
	agent *types.Agent
	err   error
}

func (m *mockAgentService) FindByID(_ context.Context, _ uint64) (*types.Agent, error) {
	return m.agent, m.err
}

type mockUserService struct {
	user *types.User
	err  error
}

func (m *mockUserService) FindByID(_ context.Context, _ uint64) (*types.User, error) {
	return m.user, m.err
}

func TestAgentHandler_AgentNotFound(t *testing.T) {
	req := require.New(t)
	rt := &mockRuntime{}
	h := &agentHandler{
		runtime: rt,
		agents:  &mockAgentService{err: fmt.Errorf("not found")},
		users:   &mockUserService{},
	}

	_, err := h.run(context.Background(), &agentRunArgs{AgentID: 999, Input: "hello"})
	req.Error(err)
	req.False(rt.called)
}

func TestAgentHandler_SystemInvocationDisabled(t *testing.T) {
	req := require.New(t)
	rt := &mockRuntime{}
	h := &agentHandler{
		runtime: rt,
		agents: &mockAgentService{agent: &types.Agent{
			Invocation: types.AgentInvocation{
				System: types.AgentInvocationSystem{Enabled: false},
			},
		}},
		users: &mockUserService{},
	}

	_, err := h.run(context.Background(), &agentRunArgs{AgentID: 1, Input: "hello"})
	req.Error(err)
	req.Contains(err.Error(), "not available for system invocation")
	req.False(rt.called)
}

func TestAgentHandler_SystemInvocationEnabled(t *testing.T) {
	req := require.New(t)
	rt := &mockRuntime{}
	h := &agentHandler{
		runtime: rt,
		agents: &mockAgentService{agent: &types.Agent{
			Invocation: types.AgentInvocation{
				System: types.AgentInvocationSystem{Enabled: true},
			},
		}},
		users: &mockUserService{},
	}

	res, err := h.run(context.Background(), &agentRunArgs{AgentID: 1, Input: "hello"})
	req.NoError(err)
	req.Equal("ok", res.Output)
	req.True(rt.called)
}

func TestAgentHandler_ServiceAccountResolved(t *testing.T) {
	req := require.New(t)
	rt := &mockRuntime{}
	h := &agentHandler{
		runtime: rt,
		agents: &mockAgentService{agent: &types.Agent{
			Invocation: types.AgentInvocation{
				System: types.AgentInvocationSystem{
					Enabled:        true,
					ServiceAccount: 42,
				},
			},
		}},
		users: &mockUserService{user: &types.User{ID: 42}},
	}

	res, err := h.run(context.Background(), &agentRunArgs{AgentID: 1, Input: "hello"})
	req.NoError(err)
	req.Equal("ok", res.Output)
	req.True(rt.called)
}

func TestAgentHandler_ServiceAccountNotFound(t *testing.T) {
	req := require.New(t)
	rt := &mockRuntime{}
	h := &agentHandler{
		runtime: rt,
		agents: &mockAgentService{agent: &types.Agent{
			Invocation: types.AgentInvocation{
				System: types.AgentInvocationSystem{
					Enabled:        true,
					ServiceAccount: 99,
				},
			},
		}},
		users: &mockUserService{err: fmt.Errorf("not found")},
	}

	_, err := h.run(context.Background(), &agentRunArgs{AgentID: 1, Input: "hello"})
	req.Error(err)
	req.Contains(err.Error(), "could not resolve service account")
	req.False(rt.called)
}
