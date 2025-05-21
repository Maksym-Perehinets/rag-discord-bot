package message

import (
	"github.com/Maksym-Perehinets/shared/database"
	"github.com/bwmarrin/discordgo"
)

type MessagesService interface {
	// GetMessage by ID
	GetMessage(messagesInput database.Messages) (*discordgo.Message, error)

	// GetMessages by IDs
	GetMessages(messagesInput []database.Messages) ([]*discordgo.Message, error)

	//// GetFollowingMessage provide a message that goes after the provided messageID
	//GetFollowingMessage(messageID string) (*discordgo.Message, error)
}
