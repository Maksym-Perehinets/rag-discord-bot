package main

import (
	"context"
	"github.com/Maksym-Perehinets/discord_ai_assistant/ai_client"
	"github.com/Maksym-Perehinets/discord_ai_assistant/bot"
)

func main() {
	b := bot.StartBot()
	exitHandler := b.Run()
	aiClient := ai_client.NewClient("http://localhost:8000", 600) // Initialize AI client with a timeout of 5000 milliseconds

	commandHandlers := bot.AddCommandHandlerForQuery(context.TODO(), "what", aiClient.ProcessQuery)

	//resp, err := aiClient.ProcessQuery(context.TODO(), ai_client.ChatRequest{
	//	UserID:    "1234567890", // Example user ID
	//	MessageID: "0987654321", // Example message ID
	//	Query: []ai_client.ChatMessage{
	//		{
	//			Role:    "user",
	//			Content: "discussion regarding tasty snacks",
	//		},
	//	},
	//})

	//if err != nil {
	//	log.Fatal(err) // Handle error appropriately in production code
	//}
	b.SetUpSlashCommands("1350120716497846374", commands, commandHandlers)

	exitHandler()
}
