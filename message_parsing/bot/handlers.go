package bot

import "github.com/bwmarrin/discordgo"

func (b *bot) RegisterHandler(handler interface{}, intents ...discordgo.Intent) {
	b.Session().AddHandler(handler)
	for _, intent := range intents {
		b.Session().Identify.Intents |= intent
	}
}
