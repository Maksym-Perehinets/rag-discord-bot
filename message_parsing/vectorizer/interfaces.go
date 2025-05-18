package vectorizer

type Vectorizer interface {
	// VectorizeMessage takes a message and returns its vector representation.
	VectorizeMessage(message string) ([]float64, error)
}
