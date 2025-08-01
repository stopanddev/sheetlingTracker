package lol

type Summoner struct {
	Puuid    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

type MatchDto struct {
	Metadata MetadataDto `json:"metadata"`
	Info     InfoDto     `json:"info"`
}

type MetadataDto struct {
	DataVersion  string   `json:"dataVersion"`
	MatchID      string   `json:"matchId"`
	Participants []string `json:"participants"`
}

type InfoDto struct {
	EndOfGameResult    string           `json:"endOfGameResult"`
	GameCreation       int64            `json:"gameCreation"`
	GameDuration       int64            `json:"gameDuration"`
	GameEndTimestamp   int64            `json:"gameEndTimestamp"`
	GameID             int64            `json:"gameId"`
	GameMode           string           `json:"gameMode"`
	GameName           string           `json:"gameName"`
	GameStartTimestamp int64            `json:"gameStartTimestamp"`
	GameType           string           `json:"gameType"`
	GameVersion        string           `json:"gameVersion"`
	MapID              int              `json:"mapId"`
	Participants       []ParticipantDto `json:"participants"`
	PlatformID         string           `json:"platformId"`
	QueueID            int              `json:"queueId"`
	Teams              []TeamDto        `json:"teams"`
	TournamentCode     string           `json:"tournamentCode"`
}

type ParticipantDto struct {
	Assists            int           `json:"assists"`
	ChampExperience    int           `json:"champExperience"`
	ChampLevel         int           `json:"champLevel"`
	ChampionID         int           `json:"championId"`
	ChampionName       string        `json:"championName"`
	Deaths             int           `json:"deaths"`
	Kills              int           `json:"kills"`
	TeamID             int           `json:"teamId"`
	TeamPosition       string        `json:"teamPosition"`
	TotalDamageDealt   int           `json:"totalDamageDealt"`
	TotalDamageTaken   int           `json:"totalDamageTaken"`
	TotalMinionsKilled int           `json:"totalMinionsKilled"`
	SummonerName       string        `json:"summonerName"`
	Puuid              string        `json:"puuid"`
	Win                bool          `json:"win"`
	Perks              PerksDto      `json:"perks"`
	Challenges         ChallengesDto `json:"challenges"`
}

type ChallengesDto struct {
	KillParticipation float64 `json:"killParticipation"`
	Kda               float64 `json:"kda"`
	DamagePerMinute   float64 `json:"damagePerMinute"`
	// Add more fields as needed
}

type MissionsDto struct {
	PlayerScore0 int `json:"playerScore0"`
	// Add additional fields as needed
}

type PerksDto struct {
	StatPerks PerkStatsDto   `json:"statPerks"`
	Styles    []PerkStyleDto `json:"styles"`
}

type PerkStatsDto struct {
	Defense int `json:"defense"`
	Flex    int `json:"flex"`
	Offense int `json:"offense"`
}

type PerkStyleDto struct {
	Description string                  `json:"description"`
	Selections  []PerkStyleSelectionDto `json:"selections"`
	Style       int                     `json:"style"`
}

type PerkStyleSelectionDto struct {
	Perk int `json:"perk"`
	Var1 int `json:"var1"`
	Var2 int `json:"var2"`
	Var3 int `json:"var3"`
}

type TeamDto struct {
	Bans       []BanDto      `json:"bans"`
	Objectives ObjectivesDto `json:"objectives"`
	TeamID     int           `json:"teamId"`
	Win        bool          `json:"win"`
}

type BanDto struct {
	ChampionID int `json:"championId"`
	PickTurn   int `json:"pickTurn"`
}

type ObjectivesDto struct {
	Baron      ObjectiveDto `json:"baron"`
	Champion   ObjectiveDto `json:"champion"`
	Dragon     ObjectiveDto `json:"dragon"`
	Horde      ObjectiveDto `json:"horde"`
	Inhibitor  ObjectiveDto `json:"inhibitor"`
	RiftHerald ObjectiveDto `json:"riftHerald"`
	Tower      ObjectiveDto `json:"tower"`
}

type ObjectiveDto struct {
	First bool `json:"first"`
	Kills int  `json:"kills"`
}
