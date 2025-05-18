package database

import (
	"errors"
	"gorm.io/gorm"
	"log"
)

type messageService struct {
	db *gorm.DB
}

func NewMessageService(db *gorm.DB) MessageService {
	log.Printf("Creating new message service with DB: %v", db)
	return &messageService{
		db: db,
	}
}

// UploadMessage stores a list of messages in the database.
func (s *messageService) UploadMessage(message *Messages) error {
	r := s.db.Create(&message)
	if r.Error != nil {
		log.Printf("Error uploading messages: %v", r.Error)
		return r.Error
	}
	return nil
}

func (s *messageService) BatchUploadMessage(message []*Messages, size int) error {
	r := s.db.CreateInBatches(&message, size)
	if r.Error != nil {
		log.Printf("Error uploading messages: %v", r.Error)
		return r.Error
	}
	return nil
}

// IsAny returns true if there are any messages in db already
func (s *messageService) IsAny() (bool, error) {
	m := Messages{}

	r := s.db.Last(&m)
	if r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		log.Printf("Error getting latest uploaded message: %v", r.Error)
		return false, r.Error
	}

	return true, nil
}

// DeleteMessage removes a message from the database.
func (s *messageService) DeleteMessage(messageID string) error {
	r := s.db.Where("message_id = ?", messageID).Delete(&Messages{})
	if r.Error != nil {
		log.Printf("Error deleting message: %v", r.Error)
		return r.Error
	}
	return nil
}

func (s *messageService) UpdateMessage(message *Messages) error {
	r := s.db.Model(&Messages{}).Where("message_id = ?", message.MessageID).Updates(message)
	if r.Error != nil {
		log.Printf("Error updating message: %v", r.Error)
		return r.Error
	}
	return nil
}
