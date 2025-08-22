package sheetling

import (
	utils "sheetlingTracker/util"

	"github.com/bwmarrin/discordgo"
)

func handleUpdateSheetlings(s *discordgo.Session, i *discordgo.InteractionCreate, channelID string) {
	msg := updateSheetlings(s, channelID)
	utils.Respond(s, i, msg)
}

func handleFindCensoredNames(s *discordgo.Session, i *discordgo.InteractionCreate, query string) {
	msg := findCensoredName(s, i, query)
	utils.Respond(s, i, msg)
}

func handleAddTrackedUser(s *discordgo.Session, i *discordgo.InteractionCreate, query string) {
	msg := addTrackedUser(s, i, query)
	utils.Respond(s, i, msg)
}

func handleDeleteSheetUser(s *discordgo.Session, i *discordgo.InteractionCreate, query string) {
	msg := deleteSheetUser(s, i, query)
	utils.Respond(s, i, msg.Error())
}

func handleDeleteTrackedUser(s *discordgo.Session, i *discordgo.InteractionCreate, query string) {
	msg := deleteTrackedUserRecord(s, i, query)
	utils.Respond(s, i, msg.Error())
}
