package bot

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

func (b *bot) SetUpSlashCommands(
	guildID string,
	commands []*discordgo.ApplicationCommand,
	commandHandlers map[string]func(
		s *discordgo.Session,
		i *discordgo.InteractionCreate),
) {
	s := b.Session()
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
		log.Printf("Registered command: /%s", cmd.Name)
	}

}
