package database_service

import (
	"github.com/Maksym-Perehinets/shared/database"
	"github.com/pgvector/pgvector-go"
)

type MessageService interface {
	// Search performs a semantic search using the provided query and returns the results.
	Search(vector pgvector.Vector, limit int, bottomLine float64) ([]database.Messages, error)
}
