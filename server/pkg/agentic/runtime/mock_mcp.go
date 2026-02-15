package runtime

import (
	"context"
)

// mockMCP is a no-op MCP client for development/testing.
type mockMCP struct{}

func MockMCP() *mockMCP { return &mockMCP{} }

func (m *mockMCP) GetTools(_ context.Context, _ uint64) ([]Tool, error) {
	return nil, nil
}

func (m *mockMCP) ExecuteTool(_ context.Context, _ string, _ map[string]any) (any, error) {
	return nil, nil
}
