package main

import (
	"github.com/Maksym-Perehinets/rag-discord-bot/bot"
	database "github.com/Maksym-Perehinets/rag-discord-bot/db"
	"github.com/Maksym-Perehinets/rag-discord-bot/message_parsing"
	"github.com/Maksym-Perehinets/rag-discord-bot/pipeline"
	"log"
)

func main() {
	botD := bot.StartBot()
	dbConn := database.New().DB()
	channels := botD.GetChannels(botD.GetGuilds()[0].ID)
	messageService := database.NewMessageService(dbConn)

	bootstrap(botD.Session(), channels, messageService)

	messageCount := message_parsing.MessageCount(botD.Session(), channels)
	for _, channel := range channels {
		log.Printf("Channel %s has %v messages, ID - %s", channel.Name, messageCount[channel.ID], channel.ID)
	}

	parsingIntents, parsingFunction := pipeline.SetUpParsingPipeLine(dbConn)
	botD.RegisterHandler(parsingFunction, parsingIntents...)

	deleteIntents, deleteFunction := pipeline.SetUpDeletePipeLine(dbConn)
	botD.RegisterHandler(deleteFunction, deleteIntents...)

	editIntents, editFunction := pipeline.SetUpEditPipeLine(dbConn)
	botD.RegisterHandler(editFunction, editIntents...)

	exitHandler := botD.Run()
	exitHandler()
}
