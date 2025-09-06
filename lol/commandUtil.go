package lol

import (
	"github.com/bwmarrin/discordgo"
)

// RegisterLoLCommands registers LoL-related slash commands
func RegisterLoLCommands(s *discordgo.Session, guildID string) error {
	// Lol API status command

	// Summoner lookup command
	_, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, &discordgo.ApplicationCommand{
		Name:        "add-summoner",
		Description: "Lookup summoner info by name",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "Summoner name",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "tag-line",
				Description: "Tagline without #",
				Required:    true,
			},
		},
	})
	if err != nil {
		return err
	}

	_, err = s.ApplicationCommandCreate(s.State.User.ID, guildID, &discordgo.ApplicationCommand{
		Name:        "add-my-matches",
		Description: "Lookup summoner info by name",
	})

	if err != nil {
		return err
	}

	_, err = s.ApplicationCommandCreate(s.State.User.ID, guildID, &discordgo.ApplicationCommand{
		Name:        "add-group",
		Description: "Add player group",
	})

	if err != nil {
		return err
	}

	_, err = s.ApplicationCommandCreate(s.State.User.ID, guildID, &discordgo.ApplicationCommand{
		Name:        "group-stats",
		Description: "Get player group",
	})

	if err != nil {
		return err
	}

	return nil
}

// HandleLoLCommands processes interaction events for LoL commands
func HandleLoLCommands(s *discordgo.Session, i *discordgo.InteractionCreate, riotApiKey string) {
	switch i.ApplicationCommandData().Name {
	case "add-summoner":
		name := i.ApplicationCommandData().Options[0].StringValue()
		tagLine := i.ApplicationCommandData().Options[1].StringValue()
		handleAddSummoner(s, i, name, tagLine, riotApiKey)
	}
}

func HandleLoLDropdowns(s *discordgo.Session, i *discordgo.InteractionCreate, riotApiKey string) {
	switch i.ApplicationCommandData().Name {
	case "add-my-matches":
		handleAddMyMatchesDropdown(s, i)
	case "add-group":
		handleAddGroupDopdown(s, i)
	case "group-stats":
		handleGetGroupStats(s, i)
	}
}
