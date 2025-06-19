package bot

import (
	"github.com/bwmarrin/discordgo"
	"testing"
)

// mockSession is a minimal mock for discordgo.Session for testing intents
// and handler registration.
type mockSession struct {
	AddHandlerCalled bool
	Identify         discordgo.Identify
}

func (m *mockSession) AddHandler(handler interface{}) {
	m.AddHandlerCalled = true
}

// botWithMockSession embeds bot and allows swapping the session for testing.
type botWithMockSession struct {
	*bot
	sessionMock *mockSession
}

func (b *botWithMockSession) Session() *discordgo.Session {
	// Return a pointer to a discordgo.Session with Identify field mapped to mock
	return &discordgo.Session{Identify: b.sessionMock.Identify}
}

func TestRegisterHandlerRegistersHandler(t *testing.T) {
	mockSess := &mockSession{}
	b := &botWithMockSession{sessionMock: mockSess}
	handler := func(s *discordgo.Session, m *discordgo.MessageCreate) {}

	// Should not panic or error
	b.RegisterHandler(handler)
	if !mockSess.AddHandlerCalled {
		t.Errorf("Expected AddHandler to be called")
	}
}

func TestRegisterHandlerSetsIntents(t *testing.T) {
	mockSess := &mockSession{}
	b := &botWithMockSession{sessionMock: mockSess}
	intent := discordgo.IntentGuildMessages
	handler := func(s *discordgo.Session, m *discordgo.MessageCreate) {}
	b.RegisterHandler(handler, intent)

	if b.sessionMock.Identify.Intents&intent == 0 {
		t.Errorf("Expected intent to be set on session")
	}
}
