package search

import "github.com/Maksym-Perehinets/shared/database"

type SemanticSearch interface {
	// Search performs a semantic search using the provided query and returns the results.
	Search(vector []float64, limit int, topN int) []database.Messages
}
