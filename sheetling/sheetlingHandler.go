package sheetling

import (
	utils "sheetlingTracker/util"

	"github.com/bwmarrin/discordgo"
)

func handleUpdateSheetlings(s *discordgo.Session, i *discordgo.InteractionCreate, channelID string) {
	msg := updateSheetlings(s, channelID)
	utils.EditResponse(s, i, msg)
}

func handleFindSheetling(s *discordgo.Session, i *discordgo.InteractionCreate, query string) {
	msg := findSheetling(query)
	utils.EditResponse(s, i, msg)
}

func handleAddTrackedUser(s *discordgo.Session, i *discordgo.InteractionCreate, query string) {
	msg := addTrackedUser(s, i, query)
	utils.EditResponse(s, i, msg)
}
