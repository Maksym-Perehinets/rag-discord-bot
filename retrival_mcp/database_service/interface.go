package database_service

import "github.com/Maksym-Perehinets/shared/database"

type MessageService interface {
	// Search performs a semantic search using the provided query and returns the results.
	Search(vector []float32, limit int, topN int) ([]database.Messages, error)
}
