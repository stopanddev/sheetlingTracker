package lol

import (
	utils "sheetlingTracker/util"

	"github.com/bwmarrin/discordgo"
)

func handleAddSummoner(s *discordgo.Session, i *discordgo.InteractionCreate, name string, tagLine string, riotApiKey string) {
	utils.Respond(s, i)
	insertSummoner(s, i, name, tagLine, riotApiKey)
}

func hanleAddMatches(s *discordgo.Session, i *discordgo.InteractionCreate, name string, riotApiKey string) {
	utils.Respond(s, i)
	puuid, err := summonerLookup(s, i, name)
	if err != nil {
		utils.EditResponse(s, i, "Failed to find player")
		return
	}
	getMatches(s, i, puuid.Puuid, riotApiKey)
}
