package database

import "gorm.io/gorm"

type Service interface {
	// Health checks the health of the database connection by pinging it.
	// It returns a map with keys indicating various health statistics.
	Health() map[string]string
	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	// DB exposes the GORM DB instance for application use.
	DB() *gorm.DB
}

type MessageService interface {
	// IsAny returns true if there are any messages in db already
	IsAny() (bool, error)

	// UploadMessage stores a list of messages in the database.
	UploadMessage(messages *Messages) error

	// BatchUploadMessage stores a batch of messages in the database.
	BatchUploadMessage(message []*Messages, size int) error

	// DeleteMessage removes a message from the database.
	DeleteMessage(messageID string) error

	// UpdateMessage updates a message in the database.
	UpdateMessage(message *Messages) error
}
