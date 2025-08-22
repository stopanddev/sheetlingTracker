package sheetling

import (
	"encoding/json"
	"fmt"
	"os"
	"sheetlingTracker/entity"
	utils "sheetlingTracker/util"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sahilm/fuzzy"
)

var DataDir = "data"
var RecordsFile = DataDir + "/records.json"
var LastMsgFile = DataDir + "/last_message.json"
var TrackedUserFile = DataDir + "/tracked_users.json"

func updateSheetlings(s *discordgo.Session, channelID string) string {
	records, err := loadRecords()
	if err != nil {
		return "Failed to load records."
	}

	lastID := loadLastMessageID()
	messages, err := s.ChannelMessages(channelID, 100, "", "", lastID)
	if err != nil {
		return "Failed to fetch messages."
	}

	newestID := lastID
	updatedUsers := make(map[string]bool)
	for _, msg := range messages {
		if msg.Author.Bot {
			continue
		}
		user, reason := extractUserReason(msg.Content)
		if user == "" || reason == "" {
			continue
		}

		lowerUser := strings.ToLower(user)
		oldRec, exists := records[lowerUser]
		if !exists || oldRec.Reason != reason {
			records[lowerUser] = entity.UserRecord{User: user, Reason: reason}
			updatedUsers[lowerUser] = true
		}

		if msg.ID > newestID {
			newestID = msg.ID
		}
	}

	count := len(updatedUsers)
	if count > 0 {
		err = saveRecords(records)
		if err != nil {
			return "Failed to save records."
		}
		saveLastMessageID(newestID)
		return fmt.Sprintf("Updated records with %d new entries.", count)
	} else {
		return "No new records found."
	}
}

func findCensoredName(s *discordgo.Session, i *discordgo.InteractionCreate, query string) string {
	records, err := loadRecords()
	if err != nil {
		return "Filed to load tracked players."
	}

	var usernames []string
	for k := range records {
		usernames = append(usernames, strings.ToLower(k))
	}

	matches := fuzzy.Find(strings.ToLower(query), usernames)
	fmt.Println(matches)
	if len(matches) == 0 {
		return "Ther are no close matches."
	}

	resp := "**Closest matches:**\n"
	for i, match := range matches {
		fmt.Println(i)
		if i >= 3 {
			break
		}
		rec := records[match.Str]
		resp += fmt.Sprintf("• User: **%s**  Offense: **%s**\n", rec.User, rec.Reason)
	}

	return resp
}

func addTrackedUser(s *discordgo.Session, i *discordgo.InteractionCreate, query string) string {
	records, err := loadTrackedPlayers()
	if err != nil {
		return "Failed to get or create tracked players file."
	}
	exists := false
	for _, record := range records {
		if strings.EqualFold(strings.ToLower(query), strings.ToLower(record.User)) {
			exists = true
		}
	}

	if exists {
		return "User already followed"
	}
	lowerQuery := strings.ToLower(query)
	records[lowerQuery] = entity.User{User: lowerQuery}
	err = saveTrackedUser(records)
	if err != nil {
		return "Failed To Add Player"
	}

	return "User Added"
}

// Utility functions

func extractUserReason(content string) (string, string) {
	lines := strings.Split(content, "\n")
	var user, reason string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToLower(line), "user:") {
			user = strings.TrimSpace(line[5:])
		} else if strings.HasPrefix(strings.ToLower(line), "reason:") {
			reason = strings.TrimSpace(line[7:])
		}
	}
	return user, reason
}

func loadRecords() (map[string]entity.UserRecord, error) {
	records := make(map[string]entity.UserRecord)
	data, err := os.ReadFile(RecordsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return records, nil // no file yet, return empty map
		}
		return nil, err
	}
	err = json.Unmarshal(data, &records)
	return records, err
}

func loadTrackedPlayers() (map[string]entity.User, error) {
	records := make(map[string]entity.User)
	data, err := os.ReadFile(TrackedUserFile)
	if err != nil {
		if os.IsNotExist(err) {
			return records, nil // no file yet, return empty map
		}
		return nil, err
	}
	err = json.Unmarshal(data, &records)
	return records, err
}

func saveTrackedUser(records map[string]entity.User) error {
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	err = os.MkdirAll(DataDir, 0755)
	if err != nil {
		return err
	}
	return os.WriteFile(TrackedUserFile, data, 0644)
}

func saveRecords(records map[string]entity.UserRecord) error {
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	err = os.MkdirAll(DataDir, 0755)
	if err != nil {
		return err
	}
	return os.WriteFile(RecordsFile, data, 0644)
}

func loadLastMessageID() string {
	data, err := os.ReadFile(LastMsgFile)
	if err != nil {
		return ""
	}
	var t entity.Tracker
	err = json.Unmarshal(data, &t)
	if err != nil {
		return ""
	}
	return t.LastMessageID
}

func saveLastMessageID(id string) {
	t := entity.Tracker{LastMessageID: id}
	data, err := json.Marshal(t)
	if err != nil {
		return
	}
	_ = os.MkdirAll(DataDir, 0755)
	_ = os.WriteFile(LastMsgFile, data, 0644)
}

func clearCommands(s *discordgo.Session, guildID string) {
	cmds, err := s.ApplicationCommands(s.State.User.ID, guildID)
	if err != nil {
		fmt.Println("Failed to fetch commands:", err)
		return
	}
	for _, cmd := range cmds {
		err := s.ApplicationCommandDelete(s.State.User.ID, guildID, cmd.ID)
		if err != nil {
			fmt.Printf("Failed to delete command %s: %v\n", cmd.Name, err)
		} else {
			fmt.Println("Deleted command:", cmd.Name)
		}
	}
}

func deleteSheetUser(s *discordgo.Session, i *discordgo.InteractionCreate, username string) error {
	records, err := loadRecords()
	if err != nil {
		return err
	}

	lowerUser := strings.ToLower(username)

	// Check if user exists and delete if found
	if _, exists := records[lowerUser]; exists {
		delete(records, lowerUser)
	} else {
		msg := fmt.Sprintf("user %s not found in records", username)
		findCensoredName(s, i, username)
		utils.Respond(s, i, msg)
	}

	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	err = os.MkdirAll(DataDir, 0755)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("user %s removed from shit list", username)
	utils.Respond(s, i, msg)
	return os.WriteFile(RecordsFile, data, 0644)
}

func deleteTrackedUserRecord(s *discordgo.Session, i *discordgo.InteractionCreate, username string) error {
	records, err := loadTrackedPlayers()
	if err != nil {
		return err
	}

	lowerUser := strings.ToLower(username)

	// Check if user exists and delete if found
	if _, exists := records[lowerUser]; exists {
		delete(records, lowerUser)
	} else {
		msg := fmt.Sprintf("user %s not found in records", username)
		findCensoredName(s, i, username)
		utils.Respond(s, i, msg)
	}

	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	err = os.MkdirAll(DataDir, 0755)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("user %s removed from tracked list", username)
	utils.Respond(s, i, msg)
	return os.WriteFile(TrackedUserFile, data, 0644)
}
