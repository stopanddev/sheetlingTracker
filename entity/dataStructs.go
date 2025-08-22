package entity

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
