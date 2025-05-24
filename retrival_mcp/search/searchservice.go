package search

import (
	"github.com/Maksym-Perehinets/retrival_mcp/database_service"
	"github.com/Maksym-Perehinets/retrival_mcp/message"
	"github.com/Maksym-Perehinets/shared/database"
	"github.com/Maksym-Perehinets/shared/vectorizer"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
	"log"
)

type searchService struct {
	db      *gorm.DB
	session *discordgo.Session
}

type TopMatch struct {
	MessageID string  `json:"message_id"`
	Content   string  `json:"content"`
	Score     float64 `json:"score,omitempty"` // Optional score field
}

func NewSearchService(db *gorm.DB, session *discordgo.Session) SemanticSearch {
	return &searchService{
		db:      db,
		session: session,
	}
}

// Search combines all required actions to perform a semantic search.
func (s *searchService) Search(query string, limit int, bottomLine float64) []TopMatch {
	vectorizerService := vectorizer.NewVectorizer()

	vector, err := vectorizerService.VectorizeMessage(query)
	if err != nil {
		log.Printf("Error vectorizing message: %v", err)
		return nil
	}

	searchServ := database_service.NewMessageService(s.db)
	messagesEmbeddings, err := searchServ.Search(database.ToPgVector(vector), limit, bottomLine)
	if err != nil {
		log.Printf("Error searching messages: %v", err)
		return nil
	}

	discordBotService := message.NewMessageService(s.session)
	messages, err := discordBotService.GetMessages(messagesEmbeddings)
	if err != nil {
		log.Printf("Error getting messages: %v", err)
		return nil
	}

	var topMatches []TopMatch
	for _, m := range messages {
		topMatches = append(topMatches, TopMatch{
			MessageID: m.ID,
			Content:   m.Content,
		})
	}

	return topMatches
}
