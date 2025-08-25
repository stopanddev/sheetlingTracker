package utils

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func EditResponse(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
	if err != nil {
		fmt.Println("Failed to edit response:", err)
	}
}
