package database_service

import (
	"fmt"
	"github.com/Maksym-Perehinets/shared/database"
	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const DefaultMinSimilarity = 0.8

type messageService struct {
	db *gorm.DB
}

func NewMessageService(db *gorm.DB) MessageService {
	return &messageService{
		db: db,
	}
}

// Search performs a semantic search using the provided vector and returns the results.
// no input validation is required, as the vector is already validated in the vectorizer service
func (s *messageService) Search(vector pgvector.Vector, limit int, minSimilarityInput float64) ([]database.Messages, error) {

	maxCosineDistance := DefaultMinSimilarity
	if minSimilarityInput > 0.0 && minSimilarityInput <= 1.0 {
		maxCosineDistance = minSimilarityInput
	}

	if maxCosineDistance < 0.0 {
		maxCosineDistance = 0.0
	} else if maxCosineDistance > 1.0 {
		maxCosineDistance = 1.0
	}

	var results []database.Messages

	query := s.db.Model(&database.Messages{})

	query = query.Where("vectorized_message <=> ?::vector <= ?", vector, maxCosineDistance)

	query = query.Order(clause.OrderBy{
		Expression: clause.Expr{
			SQL:  "vectorized_message <=> ?::vector ASC",
			Vars: []interface{}{vector},
		},
	})

	if limit > 0 {
		query = query.Limit(limit)
	}

	// Execute the query
	if err := query.Find(&results).Error; err != nil {
		return nil, fmt.Errorf("error searching messages: %w", err)
	}

	return results, nil
}
