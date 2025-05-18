package bot

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type bot struct {
	session *discordgo.Session
}

var botInstance *bot

func (b *bot) Session() *discordgo.Session {
	return b.session
}

func (b *bot) Close() {
	err := b.session.Close()
	if err != nil {
		log.Printf("Error closing Discord session: %v", err)
	}
}

func (b *bot) Run() func() {
	err := b.session.Open()
	if err != nil {
		panic(err)
	}

	return func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sc
		b.Close()
	}
}

func StartBot() *bot {
	if botInstance != nil {
		return botInstance
	}

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("No token provided. Please set DISCORD_BOT_TOKEN environment variable or config.")
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	botInstance = &bot{session: session}

	return botInstance
}
