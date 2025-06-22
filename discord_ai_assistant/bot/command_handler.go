package bot

import (
	"context"
	"github.com/Maksym-Perehinets/discord_ai_assistant/ai_client"
	"github.com/bwmarrin/discordgo"
	"log"
)

func AddCommandHandlerForQuery(
	ctx context.Context,
	commandName string,
	logicHandlerFunction func(
		ctx context.Context,
		request ai_client.ChatRequest,
	) (
		*ai_client.ChatResponse,
		error),
) map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {

	handlerFunction := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "starting working on this... let the magic happen!",
			},
		})
		if err != nil {
			log.Printf("Failed to respond to /what command: %v", err)
		}

		response, err := s.InteractionResponse(i.Interaction)
		if err != nil {
			log.Printf("Failed to get interaction response: %v", err)
		}

		log.Printf(response.ID)

		opt := i.ApplicationCommandData().Options
		optMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(opt))
		for _, opt := range opt {
			optMap[opt.Name] = opt
		}
		userQuestion := optMap["question"].StringValue()

		resp, err := logicHandlerFunction(ctx, ai_client.ChatRequest{
			UserID:    i.Interaction.Member.User.ID,
			MessageID: response.ID,
			Query: []ai_client.ChatMessage{
				{
					Role:    "user",
					Content: userQuestion,
				},
			},
		})

		if err != nil {
			log.Printf("Failed to process query: %v", err)
			errorResponse := "An error occurred while processing your request. Please try again later."

			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &errorResponse,
			})
			if err != nil {
				log.Fatalf("Failed to edit response for /what command: %v", err)
			}
			return
		}

		responseContent := resp.Answer.Content

		th, err := s.MessageThreadStart(i.Interaction.ChannelID, response.ID, "response thread", 0)
		if err != nil {
			log.Printf("Failed to create thread for response: %v", err)
		}

		for _, chunk := range splitStringIfNeeded(responseContent) {
			_, err = s.ChannelMessageSend(th.ID, chunk)
			if err != nil {
				log.Printf("Failed to send message in thread: %v", err)
			}
		}

		//_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		//	Content: &responseContent,
		//})
		//if err != nil {
		//	log.Printf("Failed to edit response for /what command: %v", err)
		//}
	}
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		commandName: handlerFunction,
	}
}
