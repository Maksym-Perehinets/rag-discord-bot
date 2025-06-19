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

// botWithMockSession allows swapping the session for testing.
type botWithMockSession struct {
	sessionMock *mockSession
}

func (b *botWithMockSession) Session() *discordgo.Session {
	return &discordgo.Session{Identify: b.sessionMock.Identify}
}

func (b *botWithMockSession) RegisterHandler(handler interface{}, intents ...discordgo.Intent) {
	b.sessionMock.AddHandler(handler)
	for _, intent := range intents {
		b.sessionMock.Identify.Intents |= intent
	}
}

func TestMockRegisterHandlerRegistersHandler(t *testing.T) {
	mockSess := &mockSession{}
	b := &botWithMockSession{sessionMock: mockSess}
	handler := func(s *discordgo.Session, m *discordgo.MessageCreate) {}

	b.sessionMock.AddHandler(handler)
	if !mockSess.AddHandlerCalled {
		t.Errorf("Expected AddHandler to be called")
	}
}

func TestMockRegisterHandlerSetsIntents(t *testing.T) {
	mockSess := &mockSession{}
	b := &botWithMockSession{sessionMock: mockSess}
	intent := discordgo.IntentGuildMessages
	handler := func(s *discordgo.Session, m *discordgo.MessageCreate) {}
	b.sessionMock.AddHandler(handler)
	b.sessionMock.Identify.Intents |= intent

	if b.sessionMock.Identify.Intents&intent == 0 {
		t.Errorf("Expected intent to be set on session")
	}
}
