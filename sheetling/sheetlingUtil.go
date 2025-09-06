package sheetling

import (
	"sheetlingTracker/lol"

	"github.com/bwmarrin/discordgo"
)

func SheetlingRegisterCommands(s *discordgo.Session, guildID string) {
	_, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, &discordgo.ApplicationCommand{
		Name:        "find-sheetling",
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
		Name:        "update-sheetling",
		Description: "Track your games to spot shitling ",
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
	case "update-sheetling":
		handleUpdateSheetlings(s, i, watchChannelID)
	case "find-sheetling":
		query := i.ApplicationCommandData().Options[0].StringValue()
		handleFindSheetling(s, i, query)
	}
}
