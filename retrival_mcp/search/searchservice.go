package search

import (
	"gorm.io/gorm"
)

type searchService struct {
	db *gorm.DB
}

type TopMatch struct {
	messageID string
}

func (s *searchService) Search(query string, limit int, topN int) []TopMatch, {

	return nil
}
