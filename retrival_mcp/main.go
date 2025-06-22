package main

import (
	"github.com/Maksym-Perehinets/retrival_mcp/mcp"
	"github.com/Maksym-Perehinets/retrival_mcp/search"
	"github.com/Maksym-Perehinets/shared/database"
	"github.com/Maksym-Perehinets/shared/discord"
)

func main() {
	dbService := database.New()
	discordBot := discord.StartBot()

	// Initialize the database connection
	defer dbService.Close()

	searchService := search.NewSearchService(dbService.DB(), discordBot.Session())

	s := mcp.NewMCPService()

	s.RegisterTool(mcp.AddSearchTool(searchService))

	s.Start()
}
