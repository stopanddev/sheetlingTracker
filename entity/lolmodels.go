package entity

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
	AllInPings                     int           `json:"allInPings"`
	AssistMePings                  int           `json:"assistMePings"`
	Assists                        int           `json:"assists"`
	BaronKills                     int           `json:"baronKills"`
	BountyLevel                    int           `json:"bountyLevel"`
	ChampExperience                int           `json:"champExperience"`
	ChampLevel                     int           `json:"champLevel"`
	ChampionId                     int           `json:"championId"`
	ChampionName                   string        `json:"championName"`
	CommandPings                   int           `json:"commandPings"`
	ChampionTransform              int           `json:"championTransform"`
	ConsumablesPurchased           int           `json:"consumablesPurchased"`
	Challenges                     ChallengesDto `json:"challenges"`
	DamageDealtToBuildings         int           `json:"damageDealtToBuildings"`
	DamageDealtToObjectives        int           `json:"damageDealtToObjectives"`
	DamageDealtToTurrets           int           `json:"damageDealtToTurrets"`
	DamageSelfMitigated            int           `json:"damageSelfMitigated"`
	Deaths                         int           `json:"deaths"`
	DetectorWardsPlaced            int           `json:"detectorWardsPlaced"`
	DoubleKills                    int           `json:"doubleKills"`
	DragonKills                    int           `json:"dragonKills"`
	EligibleForProgression         bool          `json:"eligibleForProgression"`
	EnemyMissingPings              int           `json:"enemyMissingPings"`
	EnemyVisionPings               int           `json:"enemyVisionPings"`
	FirstBloodAssist               bool          `json:"firstBloodAssist"`
	FirstBloodKill                 bool          `json:"firstBloodKill"`
	FirstTowerAssist               bool          `json:"firstTowerAssist"`
	FirstTowerKill                 bool          `json:"firstTowerKill"`
	GameEndedInEarlySurrender      bool          `json:"gameEndedInEarlySurrender"`
	GameEndedInSurrender           bool          `json:"gameEndedInSurrender"`
	HoldPings                      int           `json:"holdPings"`
	GetBackPings                   int           `json:"getBackPings"`
	GoldEarned                     int           `json:"goldEarned"`
	GoldSpent                      int           `json:"goldSpent"`
	IndividualPosition             string        `json:"individualPosition"`
	InhibitorKills                 int           `json:"inhibitorKills"`
	InhibitorTakedowns             int           `json:"inhibitorTakedowns"`
	InhibitorsLost                 int           `json:"inhibitorsLost"`
	Item0                          int           `json:"item0"`
	Item1                          int           `json:"item1"`
	Item2                          int           `json:"item2"`
	Item3                          int           `json:"item3"`
	Item4                          int           `json:"item4"`
	Item5                          int           `json:"item5"`
	Item6                          int           `json:"item6"`
	ItemsPurchased                 int           `json:"itemsPurchased"`
	KillingSprees                  int           `json:"killingSprees"`
	Kills                          int           `json:"kills"`
	Lane                           string        `json:"lane"`
	LargestCriticalStrike          int           `json:"largestCriticalStrike"`
	LargestKillingSpree            int           `json:"largestKillingSpree"`
	LargestMultiKill               int           `json:"largestMultiKill"`
	LongestTimeSpentLiving         int           `json:"longestTimeSpentLiving"`
	MagicDamageDealt               int           `json:"magicDamageDealt"`
	MagicDamageDealtToChampions    int           `json:"magicDamageDealtToChampions"`
	MagicDamageTaken               int           `json:"magicDamageTaken"`
	NeutralMinionsKilled           int           `json:"neutralMinionsKilled"`
	NeedVisionPings                int           `json:"needVisionPings"`
	NexusKills                     int           `json:"nexusKills"`
	NexusTakedowns                 int           `json:"nexusTakedowns"`
	NexusLost                      int           `json:"nexusLost"`
	ObjectivesStolen               int           `json:"objectivesStolen"`
	ObjectivesStolenAssists        int           `json:"objectivesStolenAssists"`
	OnMyWayPings                   int           `json:"onMyWayPings"`
	ParticipantId                  int           `json:"participantId"`
	PlayerScore0                   int           `json:"playerScore0"`
	PlayerScore1                   int           `json:"playerScore1"`
	PlayerScore2                   int           `json:"playerScore2"`
	PlayerScore3                   int           `json:"playerScore3"`
	PlayerScore4                   int           `json:"playerScore4"`
	PlayerScore5                   int           `json:"playerScore5"`
	PlayerScore6                   int           `json:"playerScore6"`
	PlayerScore7                   int           `json:"playerScore7"`
	PlayerScore8                   int           `json:"playerScore8"`
	PlayerScore9                   int           `json:"playerScore9"`
	PlayerScore10                  int           `json:"playerScore10"`
	PlayerScore11                  int           `json:"playerScore11"`
	PentaKills                     int           `json:"pentaKills"`
	Perks                          PerksDto      `json:"perks"`
	PhysicalDamageDealt            int           `json:"physicalDamageDealt"`
	PhysicalDamageDealtToChampions int           `json:"physicalDamageDealtToChampions"`
	PhysicalDamageTaken            int           `json:"physicalDamageTaken"`
	Placement                      int           `json:"placement"`
	PlayerAugment1                 int           `json:"playerAugment1"`
	PlayerAugment2                 int           `json:"playerAugment2"`
	PlayerAugment3                 int           `json:"playerAugment3"`
	PlayerAugment4                 int           `json:"playerAugment4"`
	PlayerSubteamId                int           `json:"playerSubteamId"`
	PushPings                      int           `json:"pushPings"`
	ProfileIcon                    int           `json:"profileIcon"`
	Puuid                          string        `json:"puuid"`
	QuadraKills                    int           `json:"quadraKills"`
	RiotIdGameName                 string        `json:"riotIdGameName"`
	RiotIdTagline                  string        `json:"riotIdTagline"`
	Role                           string        `json:"role"`
	SightWardsBoughtInGame         int           `json:"sightWardsBoughtInGame"`
	Spell1Casts                    int           `json:"spell1Casts"`
	Spell2Casts                    int           `json:"spell2Casts"`
	Spell3Casts                    int           `json:"spell3Casts"`
	Spell4Casts                    int           `json:"spell4Casts"`
	SubteamPlacement               int           `json:"subteamPlacement"`
	Summoner1Casts                 int           `json:"summoner1Casts"`
	Summoner1Id                    int           `json:"summoner1Id"`
	Summoner2Casts                 int           `json:"summoner2Casts"`
	Summoner2Id                    int           `json:"summoner2Id"`
	SummonerId                     string        `json:"summonerId"`
	SummonerLevel                  int           `json:"summonerLevel"`
	SummonerName                   string        `json:"summonerName"`
	TeamEarlySurrendered           bool          `json:"teamEarlySurrendered"`
	TeamId                         int           `json:"teamId"`
	TeamPosition                   string        `json:"teamPosition"`
	TimeCCingOthers                int           `json:"timeCCingOthers"`
	TimePlayed                     int           `json:"timePlayed"`
	TotalAllyJungleMinionsKilled   int           `json:"totalAllyJungleMinionsKilled"`
	TotalDamageDealt               int           `json:"totalDamageDealt"`
	TotalDamageDealtToChampions    int           `json:"totalDamageDealtToChampions"`
	TotalDamageShieldedOnTeammates int           `json:"totalDamageShieldedOnTeammates"`
	TotalDamageTaken               int           `json:"totalDamageTaken"`
	TotalEnemyJungleMinionsKilled  int           `json:"totalEnemyJungleMinionsKilled"`
	TotalHeal                      int           `json:"totalHeal"`
	TotalHealsOnTeammates          int           `json:"totalHealsOnTeammates"`
	TotalMinionsKilled             int           `json:"totalMinionsKilled"`
	TotalTimeCCDealt               int           `json:"totalTimeCCDealt"`
	TotalTimeSpentDead             int           `json:"totalTimeSpentDead"`
	TotalUnitsHealed               int           `json:"totalUnitsHealed"`
	TripleKills                    int           `json:"tripleKills"`
	TrueDamageDealt                int           `json:"trueDamageDealt"`
	TrueDamageDealtToChampions     int           `json:"trueDamageDealtToChampions"`
	TrueDamageTaken                int           `json:"trueDamageTaken"`
	TurretKills                    int           `json:"turretKills"`
	TurretTakedowns                int           `json:"turretTakedowns"`
	TurretsLost                    int           `json:"turretsLost"`
	UnrealKills                    int           `json:"unrealKills"`
	VisionScore                    int           `json:"visionScore"`
	VisionClearedPings             int           `json:"visionClearedPings"`
	VisionWardsBoughtInGame        int           `json:"visionWardsBoughtInGame"`
	WardsKilled                    int           `json:"wardsKilled"`
	WardsPlaced                    int           `json:"wardsPlaced"`
	Win                            bool          `json:"win"`
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
