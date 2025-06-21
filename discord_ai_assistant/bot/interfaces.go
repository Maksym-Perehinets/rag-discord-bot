package bot

import "github.com/bwmarrin/discordgo"

type Bot interface {
	// Session returns the discordgo session
	Session() *discordgo.Session

	// Run starts the bot and connects to Discord
	Run() func()

	// Close closes the discordgo session
	Close()

	// GetChannels returns all channels in a guild with more than 0 messages
	GetChannels(guildID string) []*discordgo.Channel

	// GetGuilds returns all guilds the bot is in
	GetGuilds() []*discordgo.UserGuild

	// RegisterHandler registers a handler for a specific event
	RegisterHandler(handler interface{}, intents ...discordgo.Intent)

	SetUpSlashCommands()
}
