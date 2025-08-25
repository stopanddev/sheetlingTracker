package utils

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func Respond(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Loading",
		},
	})
}

func EditResponse(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
	if err != nil {
		fmt.Println("Failed to edit response:", err)
	}
}
