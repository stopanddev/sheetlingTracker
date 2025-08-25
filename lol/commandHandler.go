package lol

import (
	"fmt"
	utils "sheetlingTracker/util"

	"github.com/bwmarrin/discordgo"
)

func handleSummonerLookup(s *discordgo.Session, i *discordgo.InteractionCreate, name string, riotApiKey string) {
	player, err := summonerLookup(s, i, name, riotApiKey)
	if err != nil {
		utils.EditResponse(s, i, "Failed to find plyaer")
	} else {
		msg := fmt.Sprintf("%s was found", player)
		utils.EditResponse(s, i, msg)
	}
}

func handleLoLStatus(s *discordgo.Session, i *discordgo.InteractionCreate) {
	msg := lolStatus(s, i)
	utils.EditResponse(s, i, msg)
}

func handleMatchHistory(s *discordgo.Session, i *discordgo.InteractionCreate, name1 string, name2 string, riotApiKey string) {
	msg := matchHistory(s, i, name1, name2, riotApiKey)
	utils.EditResponse(s, i, msg)
}

func handleFindCensored(s *discordgo.Session, i *discordgo.InteractionCreate, theirName string, yourname string, riotApiKey string) {
	msg := findCensored(s, i, theirName, riotApiKey, yourname)
	utils.EditResponse(s, i, msg)
}
