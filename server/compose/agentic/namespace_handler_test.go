package agentic

import (
	"context"
	"testing"

	hmcp "github.com/crusttech/human/server/system/agentic/mcp"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// stubRegistrar collects tool names passed to RegisterTool.
type stubRegistrar struct {
	names []string
}

func (s *stubRegistrar) RegisterTool(tool mcp.Tool, _ string, _ server.ToolHandlerFunc, _ ...hmcp.RegisterOption) {
	s.names = append(s.names, tool.Name)
}

func TestNamespaceHandlerRegistersTools(t *testing.T) {
	reg := &stubRegistrar{}
	NamespaceHandler(reg, nil)

	want := map[string]bool{
		"compose_namespace_lookup": true,
		"compose_namespace_create": true,
		"compose_namespace_update": true,
		"compose_namespace_delete": true,
	}

	for _, name := range reg.names {
		delete(want, name)
	}
	for name := range want {
		t.Errorf("expected tool %q to be registered but it was not", name)
	}
}

func TestParseBoolArg(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		fallback bool
		want     bool
	}{
		{name: "bool true", input: true, fallback: false, want: true},
		{name: "bool false", input: false, fallback: true, want: false},
		{name: "string true", input: "true", fallback: false, want: true},
		{name: "string false", input: "false", fallback: true, want: false},
		{name: "string 1", input: "1", fallback: false, want: true},
		{name: "string 0", input: "0", fallback: true, want: false},
		{name: "invalid string uses fallback true", input: "yes", fallback: true, want: true},
		{name: "invalid string uses fallback false", input: "yes", fallback: false, want: false},
		{name: "nil uses fallback", input: nil, fallback: true, want: true},
		{name: "int uses fallback", input: 1, fallback: false, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseBoolArg(tt.input, tt.fallback)
			if got != tt.want {
				t.Errorf("parseBoolArg(%v, %v) = %v, want %v", tt.input, tt.fallback, got, tt.want)
			}
		})
	}
}

func toolReq(args map[string]any) mcp.CallToolRequest {
	return mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
}

func invalidToolReq() mcp.CallToolRequest {
	return mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: "not a map"}}
}

func TestNamespaceCreate_Validation(t *testing.T) {
	h := &namespaceHandler{}
	ctx := context.Background()

	t.Run("invalid args type", func(t *testing.T) {
		_, err := h.create(ctx, invalidToolReq())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("missing name", func(t *testing.T) {
		_, err := h.create(ctx, toolReq(map[string]any{"slug": "my-ns"}))
		if err == nil {
			t.Fatal("expected error for missing name, got nil")
		}
	})

	t.Run("missing slug", func(t *testing.T) {
		_, err := h.create(ctx, toolReq(map[string]any{"name": "My NS"}))
		if err == nil {
			t.Fatal("expected error for missing slug, got nil")
		}
	})

	t.Run("empty name", func(t *testing.T) {
		_, err := h.create(ctx, toolReq(map[string]any{"name": "", "slug": "my-ns"}))
		if err == nil {
			t.Fatal("expected error for empty name, got nil")
		}
	})

	t.Run("empty slug", func(t *testing.T) {
		_, err := h.create(ctx, toolReq(map[string]any{"name": "My NS", "slug": ""}))
		if err == nil {
			t.Fatal("expected error for empty slug, got nil")
		}
	})
}

func TestNamespaceUpdate_InvalidArgs(t *testing.T) {
	h := &namespaceHandler{}
	ctx := context.Background()

	_, err := h.update(ctx, invalidToolReq())
	if err == nil {
		t.Fatal("expected error for invalid args type, got nil")
	}
}

func TestNamespaceDelete_InvalidArgs(t *testing.T) {
	h := &namespaceHandler{}
	ctx := context.Background()

	_, err := h.del(ctx, invalidToolReq())
	if err == nil {
		t.Fatal("expected error for invalid args type, got nil")
	}
}
