package lol

import (
	"context"
	"fmt"
	"sheetlingTracker/db"
	utils "sheetlingTracker/util"

	"github.com/bwmarrin/discordgo"
)

func handleAddSummoner(s *discordgo.Session, i *discordgo.InteractionCreate, name string, tagLine string, riotApiKey string) {
	utils.Respond(s, i)
	insertSummoner(s, i, name, tagLine, riotApiKey)
}

func handleAddMyMatchesDropdown(s *discordgo.Session, i *discordgo.InteractionCreate) {
	rows, err := db.Conn.Query(context.Background(), `SELECT player_name FROM server_players`)
	if err != nil {
		utils.EditResponse(s, i, "Error querying players:")
		return
	}
	defer rows.Close()

	var options []discordgo.SelectMenuOption
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			continue
		}
		options = append(options, discordgo.SelectMenuOption{
			Label:       name,
			Value:       name,
			Description: "Look up matches for " + name,
		})
	}

	if len(options) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No summoners found. Use `/add-summoner` to add one first.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	min := 1

	selectMenu := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			&discordgo.SelectMenu{
				CustomID:    "add_summoner_matches",
				Placeholder: "Choose a summoner...",
				Options:     options,
				MinValues:   &min,
				MaxValues:   1,
			},
		},
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    "Select a summoner to fetch matches:",
			Components: []discordgo.MessageComponent{selectMenu},
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	})
}

func HanleAddMatches(s *discordgo.Session, i *discordgo.InteractionCreate, name string, riotApiKey string) {
	utils.Respond(s, i)
	puuid, err := summonerLookup(s, i, name)
	if err != nil {
		utils.EditResponse(s, i, "Failed to find player")
		return
	}
	fmt.Println("Calling add match")
	getMatches(s, i, puuid.Puuid, riotApiKey)
}
