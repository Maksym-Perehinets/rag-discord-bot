package bot

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

func (b *bot) GetChannels(guildID string) []*discordgo.Channel {
	log.Printf("Getting channels for guild %s", guildID)
	var channels []*discordgo.Channel

	elements, err := b.session.GuildChannels(guildID)
	if err != nil {
		log.Fatalf("failed to get guild channels: %v", err)
	}

	for _, element := range elements {
		if element.Type == discordgo.ChannelTypeGuildVoice || element.Type == discordgo.ChannelTypeGuildStageVoice || element.Type == discordgo.ChannelTypeGuildCategory {
			continue
		}
		channels = append(channels, element)
	}
	return channels
}

func (b *bot) GetGuilds() []*discordgo.UserGuild {
	log.Printf("Getting guilds")
	guilds, err := b.Session().UserGuilds(100, "", "", true)
	if err != nil {
		log.Printf("Error getting guilds: %v", err)
	}
	return guilds
}
