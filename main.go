package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sheetlingTracker/lol"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/sahilm/fuzzy"
)

type UserRecord struct {
	User   string `json:"user"`
	Reason string `json:"reason"`
}

type User struct {
	User string `jsong:"user"`
}

var dataDir = "data"
var recordsFile = dataDir + "/records.json"
var lastMsgFile = dataDir + "/last_message.json"
var trackedUserFile = dataDir + "/tracked_users.json"

type Tracker struct {
	LastMessageID string `json:"last_message_id"`
}

func main() {
	_ = godotenv.Load()

	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		panic("DISCORD_BOT_TOKEN not set")
	}
	watchChannelID := os.Getenv("WATCH_CHANNEL_ID")
	if watchChannelID == "" {
		panic("WATCH_CHANNEL_ID not set")
	}

	guildID := os.Getenv("DISCORD_GUILD_ID")
	if watchChannelID == "" {
		panic("GUID_ID NOT SET")
	}

	riotApiKey := os.Getenv("RIOT_API_KEY")
	if riotApiKey == "" {
		panic("RIOT_API_KEY NOT SET")
	}
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		handleInteraction(s, i, watchChannelID, riotApiKey)
	})

	err = dg.Open()
	if err != nil {
		panic(err)
	}
	defer dg.Close()
	//clearCommands(dg, guildID)
	// Register commands
	registerCommands(dg, guildID)

	fmt.Println("Bot is running. Press CTRL-C to exit.")
	select {}
}

func registerCommands(s *discordgo.Session, guildID string) {
	_, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, &discordgo.ApplicationCommand{
		Name:        "update",
		Description: "Scan channel and update user records",
	})
	if err != nil {
		panic(err)
	}

	_, err = s.ApplicationCommandCreate(s.State.User.ID, guildID, &discordgo.ApplicationCommand{
		Name:        "find",
		Description: "Find a user in the records",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "query",
				Description: "Username to search for",
				Required:    true,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	_, err = s.ApplicationCommandCreate(s.State.User.ID, guildID, &discordgo.ApplicationCommand{
		Name:        "track-user",
		Description: "Track your games to spot shitling ",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "query",
				Description: "Username to add",
				Required:    true,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	_, err = s.ApplicationCommandCreate(s.State.User.ID, guildID, &discordgo.ApplicationCommand{
		Name:        "delete-user",
		Description: "Delete a user in the records",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "query",
				Description: "Username to delete",
				Required:    true,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	_, err = s.ApplicationCommandCreate(s.State.User.ID, guildID, &discordgo.ApplicationCommand{
		Name:        "delete-tracked-user",
		Description: "Delete tracked user",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "query",
				Description: "Username to delete",
				Required:    true,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	err = lol.RegisterLoLCommands(s, guildID)
	if err != nil {
		panic(err)
	}
}

func handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate, watchChannelID string, riotApiKey string) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return // Ignore non-command interactions
	}

	fmt.Printf("[DEBUG] Handling interaction: %s\n", i.ApplicationCommandData().Name)

	switch i.ApplicationCommandData().Name {
	case "update":
		handleUpdate(s, i, watchChannelID)
	case "find":
		query := i.ApplicationCommandData().Options[0].StringValue()
		handleFind(s, i, query)
	case "track-user":
		query := i.ApplicationCommandData().Options[0].StringValue()
		handleAddTrackedUser(s, i, query)
	case "delete-user":
		query := i.ApplicationCommandData().Options[0].StringValue()
		handleDeleteUserRecord(s, i, query)
	case "delete-tracked-user":
		query := i.ApplicationCommandData().Options[0].StringValue()
		handleDeleteTrackedUserRecord(s, i, query)
	case "lol-status", "summoner", "duo-history", "find-censored":
		lol.HandleLoLCommands(s, i, riotApiKey)
	default:
		fmt.Printf("[DEBUG] Unknown command: %s\n", i.ApplicationCommandData().Name)
	}
}

func handleUpdate(s *discordgo.Session, i *discordgo.InteractionCreate, channelID string) {
	records, err := loadRecords()
	if err != nil {
		respond(s, i, "Failed to load records.")
		return
	}

	lastID := loadLastMessageID()
	messages, err := s.ChannelMessages(channelID, 100, "", "", lastID)
	if err != nil {
		respond(s, i, "Failed to fetch messages.")
		return
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
			records[lowerUser] = UserRecord{User: user, Reason: reason}
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
			respond(s, i, "Failed to save records.")
			return
		}
		saveLastMessageID(newestID)
		respond(s, i, fmt.Sprintf("Updated records with %d new entries.", count))
	} else {
		respond(s, i, "No new records found.")
	}
}

func handleFind(s *discordgo.Session, i *discordgo.InteractionCreate, query string) {
	records, err := loadRecords()
	if err != nil {
		respond(s, i, "Filed to load tracked players.")
		return
	}

	var usernames []string
	for k := range records {
		usernames = append(usernames, strings.ToLower(k))
	}

	matches := fuzzy.Find(strings.ToLower(query), usernames)
	fmt.Println(matches)
	if len(matches) == 0 {
		respond(s, i, "Ther are no close matches.")
		return
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

	respond(s, i, resp)
}

func handleAddTrackedUser(s *discordgo.Session, i *discordgo.InteractionCreate, query string) {
	records, err := loadTrackedPlayers()
	if err != nil {
		respond(s, i, "Fialed to get or create tracked players file.")
		return
	}
	exists := false
	for _, record := range records {
		if strings.EqualFold(strings.ToLower(query), strings.ToLower(record.User)) {
			exists = true
		}
	}

	if exists {
		return
	}
	lowerQuery := strings.ToLower(query)
	records[lowerQuery] = User{User: lowerQuery}
	err = saveTrackedUser(records)
	if err != nil {
		respond(s, i, "Failed To Add Player")
	}

	respond(s, i, "User Added")
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

func respond(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	fmt.Println("In Respond")
	fmt.Println(msg)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
}

func loadRecords() (map[string]UserRecord, error) {
	records := make(map[string]UserRecord)
	data, err := os.ReadFile(recordsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return records, nil // no file yet, return empty map
		}
		return nil, err
	}
	err = json.Unmarshal(data, &records)
	return records, err
}

func loadTrackedPlayers() (map[string]User, error) {
	records := make(map[string]User)
	data, err := os.ReadFile(trackedUserFile)
	if err != nil {
		if os.IsNotExist(err) {
			return records, nil // no file yet, return empty map
		}
		return nil, err
	}
	err = json.Unmarshal(data, &records)
	return records, err
}

func saveTrackedUser(records map[string]User) error {
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	err = os.MkdirAll(dataDir, 0755)
	if err != nil {
		return err
	}
	return os.WriteFile(trackedUserFile, data, 0644)
}

func saveRecords(records map[string]UserRecord) error {
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	err = os.MkdirAll(dataDir, 0755)
	if err != nil {
		return err
	}
	return os.WriteFile(recordsFile, data, 0644)
}

func loadLastMessageID() string {
	data, err := os.ReadFile(lastMsgFile)
	if err != nil {
		return ""
	}
	var t Tracker
	err = json.Unmarshal(data, &t)
	if err != nil {
		return ""
	}
	return t.LastMessageID
}

func saveLastMessageID(id string) {
	t := Tracker{LastMessageID: id}
	data, err := json.Marshal(t)
	if err != nil {
		return
	}
	_ = os.MkdirAll(dataDir, 0755)
	_ = os.WriteFile(lastMsgFile, data, 0644)
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

func handleDeleteUserRecord(s *discordgo.Session, i *discordgo.InteractionCreate, username string) error {
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
		handleFind(s, i, username)
		respond(s, i, msg)
	}

	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	err = os.MkdirAll(dataDir, 0755)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("user %s removed from shit list", username)
	respond(s, i, msg)
	return os.WriteFile(recordsFile, data, 0644)
}

func handleDeleteTrackedUserRecord(s *discordgo.Session, i *discordgo.InteractionCreate, username string) error {
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
		handleFind(s, i, username)
		respond(s, i, msg)
	}

	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	err = os.MkdirAll(dataDir, 0755)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("user %s removed from tracked list", username)
	respond(s, i, msg)
	return os.WriteFile(trackedUserFile, data, 0644)
}
