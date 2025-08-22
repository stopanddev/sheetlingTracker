package sheetling

import (
	"sheetlingTracker/lol"

	"github.com/bwmarrin/discordgo"
)

func SheetlingRegisterCommands(s *discordgo.Session, guildID string) {
	_, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, &discordgo.ApplicationCommand{
		Name:        "update",
		Description: "Scan channel and update user records",
	})
	if err != nil {
		panic(err)
	}

	_, err = s.ApplicationCommandCreate(s.State.User.ID, guildID, &discordgo.ApplicationCommand{
		Name:        "find",
		Description: "Find a user in the records",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "query",
				Description: "Username to search for",
				Required:    true,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	_, err = s.ApplicationCommandCreate(s.State.User.ID, guildID, &discordgo.ApplicationCommand{
		Name:        "track-user",
		Description: "Track your games to spot shitling ",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "query",
				Description: "Username to add",
				Required:    true,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	_, err = s.ApplicationCommandCreate(s.State.User.ID, guildID, &discordgo.ApplicationCommand{
		Name:        "delete-user",
		Description: "Delete a user in the records",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "query",
				Description: "Username to delete",
				Required:    true,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	_, err = s.ApplicationCommandCreate(s.State.User.ID, guildID, &discordgo.ApplicationCommand{
		Name:        "delete-tracked-user",
		Description: "Delete tracked user",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "query",
				Description: "Username to delete",
				Required:    true,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	err = lol.RegisterLoLCommands(s, guildID)
	if err != nil {
		panic(err)
	}
}

func HandleSheetlingCommands(s *discordgo.Session, i *discordgo.InteractionCreate, riotApiKey string, watchChannelID string) {
	switch i.ApplicationCommandData().Name {
	case "update":
		handleUpdateSheetlings(s, i, watchChannelID)
	case "find":
		query := i.ApplicationCommandData().Options[0].StringValue()
		handleFind(s, i, query)
	case "track-user":
		query := i.ApplicationCommandData().Options[0].StringValue()
		handleAddTrackedUser(s, i, query)
	case "delete-user":
		query := i.ApplicationCommandData().Options[0].StringValue()
		handleDeleteUserRecord(s, i, query)
	case "delete-tracked-user":
		query := i.ApplicationCommandData().Options[0].StringValue()
		handleDeleteTrackedUserRecord(s, i, query)
	}
}
