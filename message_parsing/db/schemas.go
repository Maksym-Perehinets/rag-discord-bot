package database

import "github.com/pgvector/pgvector-go"

type Messages struct {
	ID                int             `json:"id" gorm:"primaryKey"`
	ChannelID         string          `json:"channel_id" gorm:"column:channel_id"`
	MessageID         string          `json:"message_id" gorm:"column:message_id"`
	AuthorID          string          `json:"author_id" gorm:"column:author_id"`
	VectorizedMessage pgvector.Vector `json:"vectorized_message" gorm:"column:vectorized_message;type:vector(1024)"` // TODO handle different vector size
}

// TODO redo to following tables
//
//type User struct {
//	UserID int    `json:"user_id" gorm:"primaryKey;autoIncrement:false"` // Primary Key
//	Name   string `json:"name" gorm:"column:name"`
//}
//
//// Channel represents the channels where messages are posted.
//type Channel struct {
//	ChannelID int    `json:"channel_id" gorm:"primaryKey;autoIncrement:false"` // Primary Key
//	Name      string `json:"name" gorm:"column:name;unique"`
//}
//
//// Message represents individual messages.
//type Message struct {
//	MessageID         int             `json:"message_id" gorm:"primaryKey;autoIncrement:false"`        // Primary Key for Message
//	ChannelID         int             `json:"channel_id" gorm:"column:channel_id"`                     // Foreign Key for Channel table
//	AuthorID          int             `json:"author_id" gorm:"column:author_id"`                       // Foreign Key for User table
//	VectorizedMessage pgvector.Vector `json:"vectorized_message" gorm:"column:vectorized_message;type:vector(1024)"`
//
//	User    User    `gorm:"foreignKey:AuthorID;references:UserID"`
//	Channel Channel `gorm:"foreignKey:ChannelID;references:ChannelID"`
//}
