package mcp

import (
	"context"
	"encoding/json"
	"github.com/Maksym-Perehinets/retrival_mcp/search"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"log"
)

const (
	limit = 20
)

func AddSearchTool(search search.SemanticSearch) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("semantic_search",
		mcp.WithDescription("Semantic search for similar messages in the database"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Formated query from the user request. This is the text to search for in the database."),
		),
	)

	searchToolHandler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("searchTool called with request: %v", request)
		query, err := request.RequireString("query")

		topMatch := search.Search(query, limit, 0.55)

		r, err := json.Marshal(topMatch)
		if err != nil {
			log.Printf("Error marshaling top matches: %v", err)
			return mcp.NewToolResultError("Failed to process the search results"), err
		}

		if err != nil {
			return mcp.NewToolResultError("Invalid query parameter The 'query' parameter must be a string"), err
		}
		return mcp.NewToolResultText(string(r)), nil
	}

	return tool, searchToolHandler
}
