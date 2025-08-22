package lol

import (
	"github.com/bwmarrin/discordgo"
)

// RegisterLoLCommands registers LoL-related slash commands
func RegisterLoLCommands(s *discordgo.Session, guildID string) error {
	// Lol API status command
	_, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, &discordgo.ApplicationCommand{
		Name:        "lol-status",
		Description: "Get Riot Games API status for a region",
		Options:     []*discordgo.ApplicationCommandOption{},
	})
	if err != nil {
		return err
	}

	// Summoner lookup command
	_, err = s.ApplicationCommandCreate(s.State.User.ID, guildID, &discordgo.ApplicationCommand{
		Name:        "summoner",
		Description: "Lookup summoner info by name",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "Summoner name",
				Required:    true,
			},
		},
	})
	if err != nil {
		return err
	}

	// Game Streak
	_, err = s.ApplicationCommandCreate(s.State.User.ID, guildID, &discordgo.ApplicationCommand{
		Name:        "duo-history",
		Description: "Lookup how you and partner's games have gone",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name1",
				Description: "Summoner name",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name2",
				Description: "Summoner name",
				Required:    true,
			},
		},
	})
	if err != nil {
		return err
	}

	// Find censored user in your game
	_, err = s.ApplicationCommandCreate(s.State.User.ID, guildID, &discordgo.ApplicationCommand{
		Name:        "find-censored",
		Description: "Lookup summoner info by name",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "your-name",
				Description: "Summoner name",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "their-name",
				Description: "Summoner name",
				Required:    true,
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

// HandleLoLCommands processes interaction events for LoL commands
func HandleLoLCommands(s *discordgo.Session, i *discordgo.InteractionCreate, riotApiKey string) {
	switch i.ApplicationCommandData().Name {
	case "lol-status":
		handleLoLStatus(s, i)
	case "summoner":
		name := i.ApplicationCommandData().Options[0].StringValue()
		handleSummonerLookup(s, i, name, riotApiKey)
	case "duo-history":
		name1 := i.ApplicationCommandData().Options[0].StringValue()
		name2 := i.ApplicationCommandData().Options[1].StringValue()
		handleMatchHistory(s, i, name1, name2, riotApiKey)
	case "find-censored":
		yourname := i.ApplicationCommandData().Options[0].StringValue()
		theirname := i.ApplicationCommandData().Options[1].StringValue()
		handleFindCensored(s, i, theirname, riotApiKey, yourname)
	}
}
