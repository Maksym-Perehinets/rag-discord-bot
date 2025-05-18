package vectorizer

import (
	"context"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/azure"
	"log"
	"os"
)

type openAIClient struct {
	client openai.Client
}

const (
	openAIAPIVersion = "2023-05-15"
)

// NewVectorizer initializes a new OpenAI client with the provided API key.
func NewVectorizer() Vectorizer {
	client := openai.NewClient(
		azure.WithEndpoint(
			os.Getenv("AZURE_OPENAI_ENDPOINT"),
			openAIAPIVersion,
		),
		azure.WithAPIKey(
			os.Getenv("AZURE_OPENAI_API_KEY"),
		),
	)
	return &openAIClient{
		client: client,
	}
}

// VectorizeMessage takes a message and returns its vector representation.
func (o *openAIClient) VectorizeMessage(message string) ([]float64, error) {
	log.Printf("Vectorizing message: %s", message)
	r, err := o.client.Embeddings.New(context.TODO(), openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String(message),
		},
		Model:      os.Getenv("AZURE_OPENAI_EMBEDDING_MODEL"),
		Dimensions: openai.Int(1024),
	})
	if err != nil {
		return nil, err
	}
	return r.Data[0].Embedding, nil
}
