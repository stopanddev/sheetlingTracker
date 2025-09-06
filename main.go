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
	if i.Type == discordgo.InteractionMessageComponent {
		switch i.MessageComponentData().CustomID {
		case "add_summoner_matches":
			selectedName := i.MessageComponentData().Values[0]
			lol.HanleAddMatches(s, i, selectedName, riotApiKey)
		case "add_group":
			groupMembers := i.MessageComponentData().Values
			lol.HanleAddGroup(s, i, groupMembers, riotApiKey)
		}
	} else {
		switch i.ApplicationCommandData().Name {
		case "update-sheetling", "find-sheetling":
			sheetling.HandleSheetlingCommands(s, i, riotApiKey, sheetlingChannelId)
		case "add-summoner":
			lol.HandleLoLCommands(s, i, riotApiKey)
		case "add-my-matches", "add-group":
			lol.HandleLoLDropdowns(s, i, riotApiKey)
		default:
			fmt.Printf("[DEBUG] Unknown command: %s\n", i.ApplicationCommandData().Name)
		}
	}
}
