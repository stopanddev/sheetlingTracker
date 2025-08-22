package lol

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sheetlingTracker/entity"
	utils "sheetlingTracker/util"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sahilm/fuzzy"
)

func lolStatus(s *discordgo.Session, i *discordgo.InteractionCreate) string {
	start := time.Now()
	fmt.Println("[DEBUG] /lolstatus called with region: NA1")
	msg := " "
	// Step 1: Defer the response
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		msg = fmt.Sprintf("Failed to defer response after %v: %v\n", time.Since(start), err)
		return msg
	}
	fmt.Printf("[DEBUG] Deferred response in %v\n", time.Since(start))

	apiKey := os.Getenv("RIOT_API_KEY")
	if apiKey == "" {
		return "Riot API key is not set."
	}

	url := "https://na1.api.riotgames.com/lol/status/v4/platform-data"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		msg = "Failed to create request."
		return msg
	}

	req.Header.Set("X-Riot-Token", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		msg = "Failed to reach Riot API: " + err.Error()
		return msg
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusTooManyRequests {
			msg = "Rate limit exceeded. Please try again later."
			return msg
		}
		msg = fmt.Sprintf("Riot API returned status %d .", resp.StatusCode)
		return msg
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
		return "Failed to parse Riot response."
	}

	msg = fmt.Sprintf(" **Status for %s (%s):**\n", result.Name, result.ID)
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

	return msg
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

func summonerLookup(s *discordgo.Session, i *discordgo.InteractionCreate, name string, riotApiKey string) (entity.Summoner, error) {

	name = strings.ReplaceAll(name, " ", "")
	url := fmt.Sprintf("https://americas.api.riotgames.com/riot/account/v1/accounts/by-riot-id/%s/NA1", name)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		utils.Respond(s, i, "Failed to create request.")
		return entity.Summoner{}, err
	}
	req.Header.Set("X-Riot-Token", riotApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.Respond(s, i, "Failed to contact Riot API.")
		return entity.Summoner{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		utils.Respond(s, i, fmt.Sprintf("Summoner `%s` not found.", name))
		return entity.Summoner{}, err
	}
	if resp.StatusCode != 200 {
		utils.Respond(s, i, fmt.Sprintf("Riot API error: %d", resp.StatusCode))
		return entity.Summoner{}, err
	}
	var summoner entity.Summoner
	err = json.NewDecoder(resp.Body).Decode(&summoner)
	if err != nil {
		utils.Respond(s, i, "Failed to decode summoner data.")
		return entity.Summoner{}, err
	}
	return summoner, err
}

func matchHistory(s *discordgo.Session, i *discordgo.InteractionCreate, name1 string, name2 string, riotApiKey string) string {

	var match_info []entity.MatchDto

	player1, err := summonerLookup(s, i, name1, riotApiKey)
	if err != nil {
		return "Failed to find player 1"
	}
	player2, err := summonerLookup(s, i, name2, riotApiKey)
	if err != nil {
		return "Failed to find player 2"
	}

	player1_matches, err := handleFindMatches(s, i, player1.Puuid, riotApiKey)
	if err != nil {
		return "Player 1 has no matches"
	}

	player2_matches, err := handleFindMatches(s, i, player2.Puuid, riotApiKey)
	if err != nil {
		return "Player 2 has no matches"
	}

	common_matches := findCommonMatches(player1_matches, player2_matches)
	if err != nil {
		return "Comparing matches failed"
	}

	if len(common_matches) == 0 {
		return "No matches played together"
	} else {
		for _, id := range common_matches {
			temp_match_info, err := handleFindMatchInfo(s, i, id, riotApiKey)
			if err != nil {
				fmt.Println("Error getting match info")
			}
			match_info = append(match_info, temp_match_info)
		}
	}
	total, victories := countVictories(match_info, player1.Puuid)
	msg := fmt.Sprintf("%s and %s have won **%d** out of **%d** games played together", player1.GameName, player2.GameName, victories, total)
	return msg
}

func handleFindMatches(s *discordgo.Session, i *discordgo.InteractionCreate, puuid string, riotApiKey string) ([]string, error) {
	url := fmt.Sprintf("https://americas.api.riotgames.com/lol/match/v5/matches/by-puuid/%s/ids?start=0&count=30", puuid)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		utils.Respond(s, i, "Failed to create request.")
		return nil, err
	}
	req.Header.Set("X-Riot-Token", riotApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.Respond(s, i, "Failed to contact Riot API.")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		fmt.Println("Match not found.")
		fmt.Println(resp.Body)
		return nil, err
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Riot API error: %d", resp.StatusCode)
		return nil, err
	}

	var matches []string
	err = json.NewDecoder(resp.Body).Decode(&matches)
	if err != nil {
		utils.Respond(s, i, "Failed to decode summoner data.")
		fmt.Println("ERROR DECODE")
		fmt.Println(resp.Body)
		return nil, err
	}

	return matches, err
}

func handleFindMatchInfo(s *discordgo.Session, i *discordgo.InteractionCreate, matchID string, riotApiKey string) (entity.MatchDto, error) {

	url := fmt.Sprintf("https://americas.api.riotgames.com/lol/match/v5/matches/%s", matchID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		utils.Respond(s, i, "Failed to create request.")
		return entity.MatchDto{}, err
	}
	req.Header.Set("X-Riot-Token", riotApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.Respond(s, i, "Failed to contact Riot API.")
		return entity.MatchDto{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		utils.Respond(s, i, "Match not found.")
		return entity.MatchDto{}, err
	}
	if resp.StatusCode != 200 {
		utils.Respond(s, i, fmt.Sprintf("Riot API error: %d", resp.StatusCode))
		return entity.MatchDto{}, err
	}
	var match_info entity.MatchDto
	err = json.NewDecoder(resp.Body).Decode(&match_info)
	if err != nil {
		utils.Respond(s, i, "Failed to decode summoner data.")
		fmt.Println(resp.Body)
		return entity.MatchDto{}, err
	}

	return match_info, err
}

func findCensored(s *discordgo.Session, i *discordgo.InteractionCreate, search string, riotApiKey string, name string) string {
	user, err := summonerLookup(s, i, name, riotApiKey)
	if err != nil {
		return "Can't find user in handleFindCensored"
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	some_matches, err := handleFindMatches(s, i, user.Puuid, riotApiKey)
	if err != nil {
		return "handleFindMatches in handleFindCensored errored out"
	}
	var usernames []string
	usernameSet := make(map[string]struct{})

	for _, match := range some_matches {
		curr_match, err := handleFindMatchInfo(s, i, match, riotApiKey)
		if err != nil {
			return "handleFindMatchInfo in handleFindCensored errored"
		}

		for _, p := range curr_match.Info.Participants {
			name := strings.ToLower(p.RiotIdGameName)
			if _, exists := usernameSet[name]; !exists {
				usernameSet[name] = struct{}{}
				usernames = append(usernames, name)
			}
		}
	}
	matches := fuzzy.Find(strings.ToLower(search), usernames)
	if len(matches) == 0 {
		return "No close matches."
	}

	resp := "**Closest matches:**\n"
	for i, match := range matches {
		if i >= 9 {
			break
		}
		resp += fmt.Sprintf("• User: **%s**\n", match.Str)
	}

	return resp
}

func findCommonMatches(a, b []string) []string {
	matchSet := make(map[string]struct{})
	for _, id := range a {
		matchSet[id] = struct{}{}
	}

	var common []string
	for _, id := range b {
		if _, exists := matchSet[id]; exists {
			common = append(common, id)
		}
	}
	return common
}

func countVictories(matches []entity.MatchDto, puuid string) (total int, victories int) {
	for _, match := range matches {
		total = total + 1

		for _, participant := range match.Info.Participants {
			if participant.Puuid == puuid && participant.Win {
				victories = victories + 1
				break
			}
		}
	}
	return total, victories
}
