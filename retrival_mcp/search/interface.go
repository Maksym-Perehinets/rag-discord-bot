package search

type SemanticSearch interface {
	// Search performs a semantic search using the provided query and returns the results.
	Search(query string, limit int, bottomLine float64) []TopMatch
}
