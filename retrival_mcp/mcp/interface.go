package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Service interface {
	Start()

	RegisterTool(t mcp.Tool, h server.ToolHandlerFunc)
}
