package pipeline

import (
	database "github.com/Maksym-Perehinets/rag-discord-bot/db"
	"github.com/Maksym-Perehinets/rag-discord-bot/message_parsing"
	"github.com/Maksym-Perehinets/rag-discord-bot/vectorizer"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
	"log"
)

func SetUpParsingPipeLine(db *gorm.DB) ([]discordgo.Intent, func(s *discordgo.Session, m *discordgo.MessageCreate)) {
	v := vectorizer.NewVectorizer()

	// recive message step
	originalMessage := make(chan *discordgo.MessageCreate, 100)
	processedMessage := make(chan *database.Messages, 100)

	messageService := database.NewMessageService(db)

	i, f := message_parsing.CreateMessageListener(originalMessage)

	// vectorize step
	go func() {
		for {
			select {
			case message, ok := <-originalMessage:
				if !ok {
					log.Printf("Chanel is closed")
				}
				r, err := v.VectorizeMessage(message.Content)
				if err != nil {
					log.Printf("Could not vectorize message %s: %v", message.ID, err)
				}

				log.Printf("Vectorized message %s", message.Content)

				if message.Content == "" {
					log.Printf("Message %s is empty, skipping", message.ID)
					continue
				}

				processedMessage <- &database.Messages{
					MessageID:         message.ID,
					ChannelID:         message.ChannelID,
					AuthorID:          message.Author.ID,
					VectorizedMessage: database.ToPgVector(r),
				}
			}
		}
	}()

	// save to db step
	go func() {
		for {
			select {
			case message, ok := <-processedMessage:
				if !ok {
					log.Printf("Chanel is closed")
				}
				r := messageService.UploadMessage(message)
				if r != nil {
					log.Printf("Could not save message %s: %v", message.MessageID, r)
				}
			}
		}
	}()

	return i, f
}

func SetUpDeletePipeLine(db *gorm.DB) ([]discordgo.Intent, func(s *discordgo.Session, m *discordgo.MessageDelete)) {
	deleteMessageID := make(chan string, 100)

	messageService := database.NewMessageService(db)

	i, f := message_parsing.DeleteMessageListener(deleteMessageID)

	go func() {
		for {
			select {
			case messageID, ok := <-deleteMessageID:
				if !ok {
					panic("Chanel is closed")
				}
				err := messageService.DeleteMessage(messageID)
				if err != nil {
					log.Printf("Could not delete message %s: %v", messageID, err)

				}
			}
		}
	}()

	return i, f
}

func SetUpEditPipeLine(db *gorm.DB) ([]discordgo.Intent, func(s *discordgo.Session, m *discordgo.MessageUpdate)) {
	editMessage := make(chan *discordgo.MessageUpdate, 100)
	processedMessage := make(chan *database.Messages, 100)

	messageService := database.NewMessageService(db)
	vectorizerService := vectorizer.NewVectorizer()

	i, f := message_parsing.EditMessageListener(editMessage)

	go func() {
		for {
			select {
			case message, ok := <-editMessage:
				if !ok {
					panic("Chanel is closed")
				}
				vectorizedMessage, err := vectorizerService.VectorizeMessage(message.Content)
				if err != nil {
					log.Printf("Could not vectorize message %s: %v", message.ID, err)
				}
				processedMessage <- &database.Messages{
					MessageID:         message.ID,
					ChannelID:         message.ChannelID,
					AuthorID:          message.Author.ID,
					VectorizedMessage: database.ToPgVector(vectorizedMessage),
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case message, ok := <-processedMessage:
				if !ok {
					panic("Chanel is closed")
				}
				err := messageService.UpdateMessage(message)
				if err != nil {
					log.Printf("Could not update message %s: %v", message.ID, err)
				}
			}
		}
	}()

	return i, f
}
