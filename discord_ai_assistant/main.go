package main

import (
	"context"
	"github.com/Maksym-Perehinets/discord_ai_assistant/ai_client"
	"github.com/Maksym-Perehinets/discord_ai_assistant/bot"
	"os"
)

func main() {
	AIAssistantAPIURL := os.Getenv("AI_ASSISTANT_API_URL")
	b := bot.StartBot()
	exitHandler := b.Run()
	aiClient := ai_client.NewClient(AIAssistantAPIURL, 600) // Initialize AI client with a timeout of 5000 milliseconds

	commandHandlers := bot.AddCommandHandlerForQuery(context.TODO(), "what", aiClient.ProcessQuery)

	b.SetUpSlashCommands("1350120716497846374", commands, commandHandlers)

	exitHandler()
}
