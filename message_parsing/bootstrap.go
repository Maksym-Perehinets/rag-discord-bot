package main

import (
	database "github.com/Maksym-Perehinets/rag-discord-bot/db"
	"github.com/Maksym-Perehinets/rag-discord-bot/message_parsing"
	"github.com/Maksym-Perehinets/rag-discord-bot/vectorizer"
	"github.com/bwmarrin/discordgo"
	"log"
	"sync"
)

// bootstrap check if initial db bootstrap required if so perform it
func bootstrap(s *discordgo.Session, channels []*discordgo.Channel, messageService database.MessageService) {
	resp, err := messageService.IsAny()
	if err != nil {
		log.Printf("Error checking if any messages exist: %v", err)
		return
	}

	if resp {
		log.Printf("Database already initialized, skipping bootstrap")
		return
	}

	log.Printf("Database not initialized, performing bootstrap")
	messagesMap := <-message_parsing.GetChannelMessages(s, channels)
	// Vectorize the messages and store them in the database
	v := vectorizer.NewVectorizer()
	var wg sync.WaitGroup
	for channelID, messages := range messagesMap {
		messageCh := make(chan *database.Messages, len(messages))
		for _, message := range messages {
			wg.Add(1)
			log.Printf("Processing message %s from channel %s", message.ID, channelID)
			go func() {
				defer wg.Done()

				vectorizedMessage, err := v.VectorizeMessage(message.Content)
				if err != nil {
					log.Printf("Error vectorizing message: %v", err)
					return
				}
				messageCh <- &database.Messages{
					ChannelID:         channelID,
					MessageID:         message.ID,
					AuthorID:          message.Author.ID,
					VectorizedMessage: database.ToPgVector(vectorizedMessage),
				}
				log.Printf("message %s from channel %s vectorized", message.ID, channelID)
			}()
		}
		wg.Wait()
		close(messageCh)
		var messagesToUpload []*database.Messages
		for msg := range messageCh {
			messagesToUpload = append(messagesToUpload, msg)
		}
		err = messageService.BatchUploadMessage(messagesToUpload, len(messagesToUpload)/20)
		if err != nil {
			log.Printf("Error uploading messages: %v", err)
		}
	}

}
