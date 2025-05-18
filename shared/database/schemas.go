package database

import "github.com/pgvector/pgvector-go"

type Messages struct {
	ID                int             `json:"id" gorm:"primaryKey"`
	ChannelID         string          `json:"channel_id" gorm:"column:channel_id"`
	MessageID         string          `json:"message_id" gorm:"column:message_id"`
	AuthorID          string          `json:"author_id" gorm:"column:author_id"`
	VectorizedMessage pgvector.Vector `json:"vectorized_message" gorm:"column:vectorized_message;type:vector(1024)"` // TODO handle different vector size
}
