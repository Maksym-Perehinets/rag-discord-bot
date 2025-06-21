package bot

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "what",
			Description: "Command that will ask RAG AI agent a question and return the answer",
			Options: []*discordgo.ApplicationCommandOption{
				{
					// We want text, so we use ApplicationCommandOptionString
					Type: discordgo.ApplicationCommandOptionString,

					// We'll call this text field "question" internally
					Name: "question",

					// This text is shown to the user in the Discord UI
					Description: "The text you want to provide to the bot",

					// The user MUST provide text for the command to work
					Required: true,
				},
			},
		},
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"what": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "starting working on this... let the magic happen!",
				},
			})

			opt := i.ApplicationCommandData().Options
			optMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(opt))
			for _, opt := range opt {
				optMap[opt.Name] = opt
			}
			userQuestion := optMap["question"].StringValue()

			if err != nil {
				log.Printf("Failed to respond to /what command: %v", err)
			}
			time.Sleep(5 * time.Second) // Simulate processing time
			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &userQuestion,
			})
			if err != nil {
				log.Printf("Failed to edit response for /what command: %v", err)
			}
		},
	}
)

func (b *bot) SetUpSlashCommands(guildID string) {
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
