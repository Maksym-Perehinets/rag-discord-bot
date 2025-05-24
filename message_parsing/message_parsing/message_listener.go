package message_parsing

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

func CreateMessageListener(message chan *discordgo.MessageCreate) (
	intents []discordgo.Intent,
	messageListenerFunc func(s *discordgo.Session, m *discordgo.MessageCreate),
) {
	intents = []discordgo.Intent{
		discordgo.IntentGuildMessages,
		discordgo.IntentsMessageContent,
	}

	messageListenerFunc = func(
		s *discordgo.Session,
		m *discordgo.MessageCreate,
	) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		select {
		case message <- m:
			log.Printf("Message from %s: %s", m.Author.Username, m.Content)
		}
	}

	return intents, messageListenerFunc
}

// DeleteMessageListener detects and returns id of deleted messages
func DeleteMessageListener(message chan string) (
	intents []discordgo.Intent,
	messageListenerFunc func(s *discordgo.Session, m *discordgo.MessageDelete),
) {
	intents = []discordgo.Intent{
		discordgo.IntentGuildMessages,
	}

	messageListenerFunc = func(
		s *discordgo.Session,
		m *discordgo.MessageDelete,
	) {
		select {
		case message <- m.ID:
			log.Printf("Message with id: %s was deleted", m.ID)
		}
	}

	return intents, messageListenerFunc
}

func EditMessageListener(message chan *discordgo.MessageUpdate) (
	intents []discordgo.Intent,
	messageListenerFunc func(s *discordgo.Session, m *discordgo.MessageUpdate),
) {
	intents = []discordgo.Intent{
		discordgo.IntentGuildMessages,
	}

	messageListenerFunc = func(
		s *discordgo.Session,
		m *discordgo.MessageUpdate,
	) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		select {
		case message <- m:
			log.Printf("Message from %s: %s", m.Author.Username, m.Content)
		}
	}

	return intents, messageListenerFunc
}
