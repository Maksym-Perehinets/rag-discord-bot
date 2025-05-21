package message

import (
	"github.com/Maksym-Perehinets/shared/database"
	"github.com/bwmarrin/discordgo"
	"log"
	"sync"
)

type messageService struct {
	sessions *discordgo.Session
}

func NewMessageService(session *discordgo.Session) MessagesService {
	return &messageService{
		sessions: session,
	}
}

func (m *messageService) GetMessage(messagesInput database.Messages) (*discordgo.Message, error) {
	message, err := m.sessions.ChannelMessage(messagesInput.ChannelID, messagesInput.MessageID)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (m *messageService) GetMessages(messagesInput []database.Messages) ([]*discordgo.Message, error) {
	var resp []*discordgo.Message
	messageAmount := len(messagesInput)
	wg := sync.WaitGroup{}
	messages := make(chan *discordgo.Message, messageAmount)

	wg.Add(messageAmount)
	for _, mess := range messagesInput {
		go func(i database.Messages) {
			defer wg.Done()
			message, err := m.sessions.ChannelMessage(i.ChannelID, i.MessageID)
			if err != nil {
				log.Printf("Error getting message: %v", err)
				//panic(err)
				return
			}
			messages <- message
			return
		}(mess)
	}

	wg.Wait()
	close(messages)
	for message := range messages {
		resp = append(resp, message)
	}

	return resp, nil
}
