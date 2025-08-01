package lol

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

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
	}
}

func handleLoLStatus(s *discordgo.Session, i *discordgo.InteractionCreate) {
	fmt.Println("I AM HERE")
	start := time.Now()
	fmt.Println("[DEBUG] /lolstatus called with region: NA1")

	// Step 1: Defer the response
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		fmt.Printf("Failed to defer response after %v: %v\n", time.Since(start), err)
		return
	}
	fmt.Printf("[DEBUG] Deferred response in %v\n", time.Since(start))

	apiKey := os.Getenv("RIOT_API_KEY")
	if apiKey == "" {
		editResponse(s, i, "Riot API key is not set.")
		return
	}

	url := fmt.Sprintf("https://na1.api.riotgames.com/lol/status/v4/platform-data")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		respond(s, i, "Failed to create request.")
		return
	}

	req.Header.Set("X-Riot-Token", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		editResponse(s, i, "Failed to reach Riot API: "+err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusTooManyRequests {
			editResponse(s, i, "Rate limit exceeded. Please try again later.")
			return
		}
		editResponse(s, i, fmt.Sprintf("Riot API returned status %d .", resp.StatusCode))
		return
	}

	var result struct {
		Name      string `json:"name"`
		ID        string `json:"id"`
		Incidents []struct {
			Titles []struct {
				Locale  string `json:"locale"`
				Content string `json:"content"`
			} `json:"titles"`
		} `json:"incidents"`
		Maintenances []struct {
			Titles []struct {
				Locale  string `json:"locale"`
				Content string `json:"content"`
			} `json:"titles"`
		} `json:"maintenances"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		editResponse(s, i, "Failed to parse Riot response.")
		return
	}

	msg := fmt.Sprintf(" **Status for %s (%s):**\n", result.Name, result.ID)
	if len(result.Maintenances) == 0 && len(result.Incidents) == 0 {
		msg += "No incidents or maintenance reported."
	} else {
		for _, m := range result.Maintenances {
			msg += "\n **Maintenance:** " + findEnglishTitle(m.Titles)
		}
		for _, inc := range result.Incidents {
			msg += "\n **Incident:** " + findEnglishTitle(inc.Titles)
		}
	}

	editResponse(s, i, msg)
}

func findEnglishTitle(titles []struct {
	Locale  string `json:"locale"`
	Content string `json:"content"`
}) string {
	for _, title := range titles {
		if title.Locale == "en_US" {
			return title.Content
		}
	}
	return "(no English title available)"
}

func handleSummonerLookup(s *discordgo.Session, i *discordgo.InteractionCreate, name string, riotApiKey string) (Summoner, error) {

	name = strings.ReplaceAll(name, " ", "")
	url := fmt.Sprintf("https://americas.api.riotgames.com/riot/account/v1/accounts/by-riot-id/%s/NA1", name)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		respond(s, i, "Failed to create request.")
		return Summoner{}, err
	}
	req.Header.Set("X-Riot-Token", riotApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		respond(s, i, "Failed to contact Riot API.")
		return Summoner{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		respond(s, i, fmt.Sprintf("Summoner `%s` not found.", name))
		return Summoner{}, err
	}
	if resp.StatusCode != 200 {
		respond(s, i, fmt.Sprintf("Riot API error: %d", resp.StatusCode))
		return Summoner{}, err
	}
	var summoner Summoner
	err = json.NewDecoder(resp.Body).Decode(&summoner)
	if err != nil {
		respond(s, i, "Failed to decode summoner data.")
		fmt.Println(resp.Body)
		return Summoner{}, err
	}

	return summoner, err
}

func handleMatchHistory(s *discordgo.Session, i *discordgo.InteractionCreate, name1 string, name2 string, riotApiKey string) {

	var match_info []MatchDto
	player1, err := handleSummonerLookup(s, i, name1, riotApiKey)
	if err != nil {
		respond(s, i, "Failed to find player 1")
	}
	fmt.Println("Found Plyaer 1")
	player2, err := handleSummonerLookup(s, i, name2, riotApiKey)
	if err != nil {
		respond(s, i, "Failed to find player 1")
	}

	fmt.Println("Found Plyaer 2")
	player1_matches, err := handleFindMatches(s, i, player1.Puuid, riotApiKey)
	fmt.Println(player1_matches)
	if err != nil {
		respond(s, i, "Player 1 has no matches")
	}

	player2_matches, err := handleFindMatches(s, i, player2.Puuid, riotApiKey)
	if err != nil {
		respond(s, i, "Player 2 has no matches")
	}

	fmt.Println(player2_matches)
	common_matches := FindCommonMatches(player1_matches, player2_matches)
	if err != nil {
		respond(s, i, "Comparing matches failed")
	}

	if len(common_matches) == 0 {
		msg := "No matches played together"
		respond(s, i, msg)
	} else {
		fmt.Printf("WE IN THE ELSE")
		for _, id := range common_matches {
			fmt.Printf("Match ID %s", id)
			temp_match_info, err := handleFindMatcheInfo(s, i, id, riotApiKey)
			if err != nil {
				respond(s, i, "Error getting match info")
			}
			match_info = append(match_info, temp_match_info)
		}
	}

	for i, match := range match_info {
		fmt.Printf("Match #%d\n", i+1)
		fmt.Printf("Match ID: %s\n", match.Metadata.MatchID)
		fmt.Printf("Game Mode: %s\n", match.Info.GameMode)
		fmt.Printf("Game Duration: %d seconds\n", match.Info.GameDuration)
		fmt.Println("Participants:")

		for _, p := range match.Info.Participants {
			fmt.Printf(" - %s (%s): %d Kills, %d Deaths, %d Assists\n",
				p.SummonerName, p.ChampionName, p.Kills, p.Deaths, p.Assists)
		}

		fmt.Println(strings.Repeat("-", 40))
	}
}

func handleFindMatches(s *discordgo.Session, i *discordgo.InteractionCreate, puuid string, riotApiKey string) ([]string, error) {
	fmt.Printf("THE PUUID %s", puuid)
	url := fmt.Sprintf("https://americas.api.riotgames.com/lol/match/v5/matches/by-puuid/%s", puuid)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		respond(s, i, "Failed to create request.")
		return nil, err
	}
	req.Header.Set("X-Riot-Token", riotApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		respond(s, i, "Failed to contact Riot API.")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		respond(s, i, "Match not found.")
		return nil, err
	}
	if resp.StatusCode != 200 {
		respond(s, i, fmt.Sprintf("Riot API error: %d", resp.StatusCode))
		return nil, err
	}
	var matches []string
	err = json.NewDecoder(resp.Body).Decode(&matches)
	if err != nil {
		respond(s, i, "Failed to decode summoner data.")
		fmt.Println(resp.Body)
		return nil, err
	}

	return matches, err
}

func handleFindMatcheInfo(s *discordgo.Session, i *discordgo.InteractionCreate, matchID string, riotApiKey string) (MatchDto, error) {

	url := fmt.Sprintf("https://americas.api.riotgames.com/lol/match/v5/matches/%s", matchID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		respond(s, i, "Failed to create request.")
		return MatchDto{}, err
	}
	req.Header.Set("X-Riot-Token", riotApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		respond(s, i, "Failed to contact Riot API.")
		return MatchDto{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		respond(s, i, fmt.Sprintf("Match not found."))
		return MatchDto{}, err
	}
	if resp.StatusCode != 200 {
		respond(s, i, fmt.Sprintf("Riot API error: %d", resp.StatusCode))
		return MatchDto{}, err
	}
	var match_info MatchDto
	err = json.NewDecoder(resp.Body).Decode(&match_info)
	if err != nil {
		respond(s, i, "Failed to decode summoner data.")
		fmt.Println(resp.Body)
		return MatchDto{}, err
	}

	return match_info, err
}

func respond(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
}

func editResponse(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
	if err != nil {
		fmt.Println("Failed to edit response:", err)
	}
}

func FindCommonMatches(a, b []string) []string {
	fmt.Printf("-----------")
	matchSet := make(map[string]struct{})
	for _, id := range a {
		fmt.Printf("%s", id)
		matchSet[id] = struct{}{}
	}

	var common []string
	for _, id := range b {
		fmt.Printf("%s", id)
		if _, exists := matchSet[id]; exists {
			common = append(common, id)
		}
	}
	return common
}
