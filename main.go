package main

import (
	"fmt"
	"os"
	"sheetlingTracker/db"
	"sheetlingTracker/lol"
	"sheetlingTracker/sheetling"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

type UserRecord struct {
	User   string `json:"user"`
	Reason string `json:"reason"`
}

type User struct {
	User string `jsong:"user"`
}

type Tracker struct {
	LastMessageID string `json:"last_message_id"`
}

func main() {
	_ = godotenv.Load()

	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		panic("DISCORD_BOT_TOKEN not set")
	}
	sheetlingChannelId := os.Getenv("SHEETLING_CHANNEL_ID")
	if sheetlingChannelId == "" {
		panic("WATCH_CHANNEL_ID not set")
	}

	guildID := os.Getenv("DISCORD_GUILD_ID")
	if guildID == "" {
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

	db.Init()
	defer db.Conn.Close()

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		handleInteraction(s, i, sheetlingChannelId, riotApiKey)
	})

	err = dg.Open()
	if err != nil {
		panic(err)
	}
	defer dg.Close()
	//clearCommands(dg, guildID)
	// Register commands
	sheetling.SheetlingRegisterCommands(dg, guildID)

	fmt.Println("Bot is running. Press CTRL-C to exit.")
	select {}
}

func handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate, sheetlingChannelId string, riotApiKey string) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return // Ignore non-command interactions
	}

	fmt.Printf("[DEBUG] Handling interaction: %s\n", i.ApplicationCommandData().Name)

	switch i.ApplicationCommandData().Name {
	case "update-sheetling", "find-sheetling":
		sheetling.HandleSheetlingCommands(s, i, riotApiKey, sheetlingChannelId)
	case "summoner", "add-my-matches":
		lol.HandleLoLCommands(s, i, riotApiKey)
	default:
		fmt.Printf("[DEBUG] Unknown command: %s\n", i.ApplicationCommandData().Name)
	}
}
