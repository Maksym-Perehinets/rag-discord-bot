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
	//
	//// Initialize the database connection
	defer dbService.Close()
	//
	searchService := search.NewSearchService(dbService.DB(), discordBot.Session())
	//m := searchService.Search("sudo", 10, 0.80)
	//
	//for _, match := range m {
	//	log.Printf("Message ID: %s", match.MessageID)
	//	log.Printf("Content: %s", match.Content)
	//}

	s := mcp.NewMCPService()

	s.RegisterTool(mcp.AddSearchTool(searchService))

	s.Start()
}
