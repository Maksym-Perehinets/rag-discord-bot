package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"log"
)

const (
	ServerName    = "retrival_mcp"
	ServerVersion = "0.0.1"
)

type mcpService struct {
	s *server.MCPServer
}

func NewMCPService(opts ...server.ServerOption) Service {
	mcpServer := server.NewMCPServer(
		ServerName,
		ServerVersion,
		server.WithLogging(),
		server.WithToolCapabilities(true),
	)

	for _, opt := range opts {
		opt(mcpServer)
	}

	return &mcpService{s: mcpServer}
}

func (m *mcpService) Start() {
	// TODO migrate to Server-Side Streamable HTTP transport as soon as it's released
	sse := server.NewSSEServer(
		m.s,
		server.WithBaseURL("http://localhost:8080"),
		server.WithUseFullURLForMessageEndpoint(true),
	)

	log.Printf("HTTP server listening on :8080")
	if err := sse.Start(":8080"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func (m *mcpService) RegisterTool(t mcp.Tool, h server.ToolHandlerFunc) {
	log.Printf("Registering tool: %s", t.Name)
	m.s.AddTool(t, h)
}
