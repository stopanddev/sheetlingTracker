package sheetling

import (
	"context"
	"fmt"
	"sheetlingTracker/db"
	"sheetlingTracker/entity"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v4"
)

func updateSheetlings(s *discordgo.Session, channelId string) string {

	lastId := loadLastSheetlingMessageId(channelId)
	messages, err := s.ChannelMessages(channelId, 100, "", lastId, "")
	if err != nil {
		return "Failed to fetch messages."
	}

	newestId := lastId
	records := make(map[string]entity.UserRecord)
	for _, msg := range messages {
		if msg.Author.Bot {
			continue
		}
		user, reason := extractUserReason(msg.Content)
		if user == "" || reason == "" {
			continue
		}

		lowerUser := strings.ToLower(user)
		records[lowerUser] = entity.UserRecord{User: user, Reason: reason}
		if msg.ID > newestId {
			newestId = msg.ID
		}
	}

	count := len(records)
	if count > 0 {
		batch := &pgx.Batch{}

		for _, record := range records {
			batch.Queue("INSERT INTO sheetlings (name, reason) VALUES ($1, $2)", record.User, record.Reason)
		}
		batchRequest := db.Conn.SendBatch(context.Background(), batch)
		defer batchRequest.Close()
		fmt.Println(records)
		for range records {
			if _, err := batchRequest.Exec(); err != nil {
				return "Error inserting shitling update"
			}
		}
		result := saveLastMessageId(newestId, channelId)
		if result != "" {
			return result
		}
		return fmt.Sprintf("Updated records with %d new entries.", count)
	} else {
		return "No new records found."
	}
}

func findSheetling(query string) string {
	rows, err := db.Conn.Query(context.Background(), `
		(
			SELECT name, reason
			FROM sheetlings
			WHERE name = $1
			LIMIT 1
		)
		UNION ALL
		(
			SELECT name, reason
			FROM sheetlings
			WHERE name ILIKE '%' || $1 || '%'
			  AND name <> $1
			LIMIT 2
		)
	`, query)

	if err != nil {
		return fmt.Sprintf("Error searching sheetlings: %v", err)
	}
	defer rows.Close()
	var results []string
	for rows.Next() {
		var name, reason string
		if err := rows.Scan(&name, &reason); err != nil {
			return fmt.Sprintf("Error reading result: %v", err)
		}
		results = append(results, fmt.Sprintf("User: %s\nReason: %s", name, reason))
	}

	if len(results) == 0 {
		return fmt.Sprintf("No sheetling found matching \"%s\"", query)
	}

	return strings.Join(results, "\n\n")
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

func loadLastSheetlingMessageId(sheetlingChannelId string) string {
	var msg entity.MessageCheckpointTable

	err := db.Conn.QueryRow(context.Background(), `
		SELECT m.channel_id, m.message_id, mt.name
		FROM message_checkpoint m
		JOIN message_type mt ON m.message_type = mt.id
		WHERE m.channel_id = $1
		ORDER BY m.channel_id
		LIMIT 1
	`, sheetlingChannelId).Scan(&msg.ChannelId, &msg.MessageId, &msg.MessageType)
	if err != nil {
		return ""
	}

	return msg.MessageId
}

func saveLastMessageId(messageId string, channelId string) string {
	cmdTag, err := db.Conn.Exec(context.Background(), `
		UPDATE message_checkpoint
		SET message_id = $1
		WHERE channel_id = $2
	`, messageId, channelId)

	if err != nil {
		return fmt.Sprintf("Issue saving last message Id %d", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Sprintf("No rows updated for save last message Id %s", channelId)
	}

	return ""
}
