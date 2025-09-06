package lol

import (
	"context"
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
			utils.EditResponse(s, i, err.Error())
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

func handleGetGroupStats(s *discordgo.Session, i *discordgo.InteractionCreate) {
	rows, err := db.Conn.Query(context.Background(), `SELECT distinct group_name FROM groups`)
	if err != nil {
		utils.EditResponse(s, i, "Error querying player groups:")
		return
	}
	defer rows.Close()
	var options []discordgo.SelectMenuOption
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			utils.EditResponse(s, i, err.Error())
		}

		options = append(options, discordgo.SelectMenuOption{
			Label:       name,
			Value:       name,
			Description: "Get Group Stats " + name,
		})
	}
	if len(options) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No groups found.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	min := 1
	max := 1
	if len(options) < 5 {
		max = len(options)
	}
	selectMenu := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			&discordgo.SelectMenu{
				CustomID:    "get_group",
				Placeholder: "Choose a group...",
				Options:     options,
				MinValues:   &min,
				MaxValues:   max,
			},
		},
	}
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    "Get group:",
			Components: []discordgo.MessageComponent{selectMenu},
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		utils.EditResponse(s, i, err.Error())
	}
}

func handleAddGroupDopdown(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
			Description: "Add to group " + name,
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

	min := 2
	max := 5
	if len(options) < 5 {
		max = len(options)
	}
	selectMenu := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			&discordgo.SelectMenu{
				CustomID:    "add_group",
				Placeholder: "Choose a summoner...",
				Options:     options,
				MinValues:   &min,
				MaxValues:   max,
			},
		},
	}
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    "Add players to group:",
			Components: []discordgo.MessageComponent{selectMenu},
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		utils.EditResponse(s, i, err.Error())
	}
}

func HanleAddMatches(s *discordgo.Session, i *discordgo.InteractionCreate, name string, riotApiKey string) {
	utils.Respond(s, i)
	puuid, err := summonerLookup(s, i, name)
	if err != nil {
		utils.EditResponse(s, i, "Failed to find player")
		return
	}
	getMatches(s, i, puuid.Puuid, riotApiKey)
}

func HanleAddGroup(s *discordgo.Session, i *discordgo.InteractionCreate, names []string, riotApiKey string) {
	utils.Respond(s, i)
	addGroups(s, i, names, riotApiKey)

}

func HandleGetGroup(s *discordgo.Session, i *discordgo.InteractionCreate, names string, riotApiKey string) {
	utils.Respond(s, i)
	var name = names + " "
	getGroups(s, i, name, riotApiKey)
}
