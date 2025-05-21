package database_service

import (
	"fmt"
	"github.com/Maksym-Perehinets/shared/database"
	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type messageService struct {
	db *gorm.DB
}

func NewMessageService(db *gorm.DB) MessageService {
	return &messageService{
		db: db,
	}
}

// Search performs a semantic search using the provided query and returns the results.
func (s *messageService) Search(vector []float32, limit int, topN int) ([]database.Messages, error) {
	var results []database.Messages
	// TODO REWRITE HOW topN WORKS FUCK again
	// Ensure "vectorized_message" is your correct vector column name.
	// Use pq.Array to wrap the []float64 slice.
	// Use an explicit ::vector cast in the SQL.
	query := s.db.Order(clause.OrderBy{
		Expression: clause.Expr{
			SQL:  "vectorized_message <=> ?::vector",
			Vars: []interface{}{pgvector.NewVector(vector)}, // Pass pq.Array(vector)
		},
	})

	if topN > 0 {
		query = query.Limit(topN)
	}

	r := query.Find(&results)
	if r.Error != nil {
		return nil, fmt.Errorf("error searching messages: %w", r.Error)
	}

	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}
