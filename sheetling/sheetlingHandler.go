package sheetling

import (
	utils "sheetlingTracker/util"

	"github.com/bwmarrin/discordgo"
)

func handleUpdateSheetlings(s *discordgo.Session, i *discordgo.InteractionCreate, channelID string) {
	utils.Respond(s, i)
	msg := updateSheetlings(s, channelID)
	utils.EditResponse(s, i, msg)
}

func handleFindSheetling(s *discordgo.Session, i *discordgo.InteractionCreate, query string) {
	utils.Respond(s, i)
	msg := findSheetling(query)
	utils.EditResponse(s, i, msg)
}
