package mcp

import (
	"github.com/go-chi/chi/v5" 
	"github.com/mark3labs/mcp-go/server"
)

type MCPServer struct {
	server	   *server.MCPServer
	httpServer *server.StreamableHTTPServer
}

func NewMCPServer(reg *Registry, name, version string) *MCPServer {
	s := server.NewMCPServer(
		name,
		version,
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true,false),
	)
	for _,t := range reg.tools {
		s.AddTool(t.Tool, t.Handler)
	}
	for _,r := range reg.resources {
		s.AddResource(r.Resource, r.Handler)
	}

	httpServer := server.NewStreamableHTTPServer(s)

	return &MCPServer{
		server: s,
		httpServer: httpServer,
	}
}

func (m *MCPServer) MountRoutes (r chi.Router) {
	r.Handle("/*",m.httpServer)
}