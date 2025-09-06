package lol

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sheetlingTracker/db"
	"sheetlingTracker/entity"
	"sheetlingTracker/entityUtils.go"
	utils "sheetlingTracker/util"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v4"
)

func getsummonerPuuid(s *discordgo.Session, i *discordgo.InteractionCreate, name string, tagLine string, riotApiKey string) (entity.Summoner, error) {

	name = strings.ReplaceAll(name, " ", "")
	url := fmt.Sprintf("https://americas.api.riotgames.com/riot/account/v1/accounts/by-riot-id/%s/%s", name, tagLine)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		utils.EditResponse(s, i, "Failed to create request.")
		return entity.Summoner{}, err
	}
	req.Header.Set("X-Riot-Token", riotApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.EditResponse(s, i, "Failed to contact Riot API.")
		return entity.Summoner{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		utils.EditResponse(s, i, fmt.Sprintf("Summoner `%s` not found.", name))
		return entity.Summoner{}, err
	}
	if resp.StatusCode != 200 {
		utils.EditResponse(s, i, fmt.Sprintf("Riot API error: %d", resp.StatusCode))
		return entity.Summoner{}, err
	}
	var summoner entity.Summoner
	err = json.NewDecoder(resp.Body).Decode(&summoner)
	if err != nil {
		utils.EditResponse(s, i, "Failed to decode summoner data.")
		return entity.Summoner{}, err
	}
	return summoner, err
}

func summonerLookup(s *discordgo.Session, i *discordgo.InteractionCreate, name string) (entity.Summoner, error) {
	var summoner entity.Summoner

	err := db.Conn.QueryRow(context.Background(), `
		SELECT p.player_name, p.tag_line, p.puuid
		FROM server_players p
		WHERE p.player_name = $1
		ORDER BY p.player_name
		LIMIT 1
	`, name).Scan(&summoner.GameName, &summoner.TagLine, &summoner.Puuid)
	if err != nil {
		utils.EditResponse(s, i, err.Error())
		return entity.Summoner{}, err
	}

	return summoner, nil
}

func insertSummoner(s *discordgo.Session, i *discordgo.InteractionCreate, name string, tagLine string, riotApiKey string) {
	summoner, err := getsummonerPuuid(s, i, name, tagLine, riotApiKey)
	if err != nil {
		msg := fmt.Sprintf("Summoner %s #%s does not exist", name, tagLine)
		utils.EditResponse(s, i, msg)
		return
	}

	ctx := context.Background()
	query := `
		INSERT INTO server_players (
			player_name, tag_line, puuid
		)
		VALUES ($1, $2, $3)
		ON CONFLICT (player_name, tag_line) DO NOTHING;
	`

	_, err = db.Conn.Exec(ctx, query, summoner.GameName, summoner.TagLine, summoner.Puuid)
	if err != nil {
		utils.EditResponse(s, i, err.Error())
		return
	}

	utils.EditResponse(s, i, fmt.Sprintf("Successfully added %s", name))

}
func getMatchIds(s *discordgo.Session, i *discordgo.InteractionCreate, puuid string, riotApiKey string) ([]string, error) {
	url := fmt.Sprintf("https://americas.api.riotgames.com/lol/match/v5/matches/by-puuid/%s/ids?type=normal&start=0&count=30", puuid)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		utils.Respond(s, i)
		return []string{}, err
	}
	req.Header.Set("X-Riot-Token", riotApiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.Respond(s, i)
		return []string{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		utils.Respond(s, i)
		return []string{}, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []string{}, err
	}
	var results []string
	err = json.Unmarshal(body, &results)
	if err != nil {
		return []string{}, err
	}
	return results, nil
}

func getSingleMatch(s *discordgo.Session, i *discordgo.InteractionCreate, matchId string, riotApiKey string) (entity.MatchDto, error) {
	url := fmt.Sprintf("https://americas.api.riotgames.com/lol/match/v5/matches/%s", matchId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		utils.EditResponse(s, i, err.Error())
		return entity.MatchDto{}, err
	}
	req.Header.Set("X-Riot-Token", riotApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.EditResponse(s, i, err.Error())
		return entity.MatchDto{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		utils.EditResponse(s, i, "404 Error getSingleMatch")
		return entity.MatchDto{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.EditResponse(s, i, err.Error())
		return entity.MatchDto{}, err
	}
	var result entity.MatchDto
	err = json.Unmarshal(body, &result)
	if err != nil {
		utils.EditResponse(s, i, err.Error())
		return entity.MatchDto{}, err
	}

	return result, nil
}

func getMatches(s *discordgo.Session, i *discordgo.InteractionCreate, puuId string, riotApiKey string) {
	matchIds, err := getMatchIds(s, i, puuId, riotApiKey)
	var matchList []entity.MatchDto
	if err != nil {
		utils.EditResponse(s, i, err.Error())
		return
	}
	for _, matchId := range matchIds {
		match, err := getSingleMatch(s, i, matchId, riotApiKey)
		if err != nil {
			utils.EditResponse(s, i, err.Error())
			return
		}
		matchList = append(matchList, match)
	}
	var validMatches = filterMatches(matchList)
	for _, match := range validMatches {
		matchRecord, matchDetails := entityUtils.ParseMatchDto(match)
		insertMatches(s, i, matchRecord, matchDetails)
	}

	utils.EditResponse(s, i, "Finished adding records")

}

func filterMatches(matchList []entity.MatchDto) []entity.MatchDto {
	var validMatches []entity.MatchDto

	for _, match := range matchList {
		if match.Info.GameMode == "CLASSIC" {
			validMatches = append(validMatches, match)
		}
	}

	return validMatches
}

func insertMatches(s *discordgo.Session, i *discordgo.InteractionCreate, matchRecord entity.MatchRecordDTO, matchDetails []entity.MatchPlayerDetailsDTO) {
	batch := &pgx.Batch{}
	ctx := context.Background()
	query := `
		INSERT INTO matches (
			match_id, player1, player2, player3, player4, player5,
			player6, player7, player8, player9, player10,
			end_of_game_result, game_creation, game_duration
		)
		VALUES ($1, $2, $3, $4, $5, $6,
		        $7, $8, $9, $10, $11,
		        $12, $13, $14)
		ON CONFLICT (match_id) DO NOTHING;
	`

	_, err := db.Conn.Exec(ctx, query,
		matchRecord.MatchID,
		matchRecord.Player1, matchRecord.Player2, matchRecord.Player3, matchRecord.Player4, matchRecord.Player5,
		matchRecord.Player6, matchRecord.Player7, matchRecord.Player8, matchRecord.Player9, matchRecord.Player10,
		matchRecord.EndOfGameResult,
		matchRecord.GameCreation,
		matchRecord.GameDuration,
	)
	if err != nil {
		utils.EditResponse(s, i, err.Error())
		return
	}

	for _, d := range matchDetails {
		batch.Queue(`
			INSERT INTO match_player_details (
				match_id, all_in_pings, assist_me_pings, assists, baron_kills, bounty_level,
				champ_experience, champ_level, champion_id, champion_name, command_pings,
				champion_transform, consumables_purchased, damage_dealt_to_buildings, damage_dealt_to_objectives,
				damage_dealt_to_turrets, damage_self_mitigated, deaths, detector_wards_placed,
				double_kills, dragon_kills, enemy_missing_pings, enemy_vision_pings, first_blood_assist,
				first_blood_kill, first_tower_assist, first_tower_kill, game_ended_in_early_surrender,
				game_ended_in_surrender, hold_pings, get_back_pings, gold_earned, gold_spent,
				individual_position, inhibitor_kills, inhibitor_takedowns, inhibitors_lost,
				item0, item1, item2, item3, item4, item5, item6, items_purchased, killing_sprees,
				kills, lane, largest_critical_strike, largest_killing_spree, largest_multi_kill,
				longest_time_spent_living, magic_damage_dealt, magic_damage_dealt_to_champions,
				magic_damage_taken, neutral_minions_killed, need_vision_pings, nexus_kills,
				nexus_takedowns, nexus_lost, objectives_stolen, objectives_stolen_assists,
				on_my_way_pings, participant_id, player_score0, player_score1, player_score2,
				player_score3, player_score4, player_score5, player_score6, player_score7,
				player_score8, player_score9, player_score10, player_score11, penta_kills,
				physical_damage_dealt, physical_damage_dealt_to_champions, physical_damage_taken,
				placement, player_augment1, player_augment2, player_augment3, player_augment4,
				player_subteam_id, push_pings, profile_icon, puuid, quadra_kills, riot_id_game_name,
				riot_id_tagline, role, sight_wards_bought_in_game, spell1_casts, spell2_casts,
				spell3_casts, spell4_casts, subteam_placement, summoner1_casts, summoner1_id,
				summoner2_casts, summoner2_id, summoner_id, summoner_level, summoner_name,
				team_early_surrendered, team_id, team_position, time_ccing_others, time_played,
				total_ally_jungle_minions_killed, total_damage_dealt, total_damage_dealt_to_champions,
				total_damage_shielded_on_teammates, total_damage_taken, total_enemy_jungle_minions_killed,
				total_heal, total_heals_on_teammates, total_minions_killed, total_time_cc_dealt,
				total_time_spent_dead, total_units_healed, triple_kills, true_damage_dealt,
				true_damage_dealt_to_champions, true_damage_taken, turret_kills, turret_takedowns,
				turrets_lost, unreal_kills, vision_score, vision_cleared_pings, vision_wards_bought_in_game,
				wards_killed, wards_placed, win
			)
			VALUES (
				$1, $2, $3, $4, $5, $6,
				$7, $8, $9, $10, $11,
				$12, $13, $14, $15, $16,
				$17, $18, $19, $20, $21,
				$22, $23, $24, $25, $26,
				$27, $28, $29, $30, $31,
				$32, $33, $34, $35, $36,
				$37, $38, $39, $40, $41,
				$42, $43, $44, $45, $46,
				$47, $48, $49, $50, $51,
				$52, $53, $54, $55, $56,
				$57, $58, $59, $60, $61,
				$62, $63, $64, $65, $66,
				$67, $68, $69, $70, $71,
				$72, $73, $74, $75, $76,
				$77, $78, $79, $80, $81,
				$82, $83, $84, $85, $86,
				$87, $88, $89, $90, $91,
				$92, $93, $94, $95, $96,
				$97, $98, $99, $100, $101,
				$102, $103, $104, $105, $106,
				$107, $108, $109, $110, $111,
				$112, $113, $114, $115, $116,
				$117, $118, $119, $120, $121,
				$122, $123, $124, $125, $126,
				$127, $128, $129, $130, $131,
				$132, $133, $134, $135, $136,
				$137
			)
			ON CONFLICT (match_id, participant_id) DO NOTHING;
			`,
			// Your 141 values go here: d.MatchID, d.AllInPings, ..., d.Win
			// You can generate this with a script (or let me do it for you)
			d.MatchID, d.AllInPings, d.AssistMePings, d.Assists, d.BaronKills, d.BountyLevel,
			d.ChampExperience, d.ChampLevel, d.ChampionId, d.ChampionName, d.CommandPings,
			d.ChampionTransform, d.ConsumablesPurchased, d.DamageDealtToBuildings, d.DamageDealtToObjectives,
			d.DamageDealtToTurrets, d.DamageSelfMitigated, d.Deaths, d.DetectorWardsPlaced,
			d.DoubleKills, d.DragonKills, d.EnemyMissingPings, d.EnemyVisionPings, d.FirstBloodAssist,
			d.FirstBloodKill, d.FirstTowerAssist, d.FirstTowerKill, d.GameEndedInEarlySurrender,
			d.GameEndedInSurrender, d.HoldPings, d.GetBackPings, d.GoldEarned, d.GoldSpent,
			d.IndividualPosition, d.InhibitorKills, d.InhibitorTakedowns, d.InhibitorsLost,
			d.Item0, d.Item1, d.Item2, d.Item3, d.Item4, d.Item5, d.Item6, d.ItemsPurchased, d.KillingSprees,
			d.Kills, d.Lane, d.LargestCriticalStrike, d.LargestKillingSpree, d.LargestMultiKill,
			d.LongestTimeSpentLiving, d.MagicDamageDealt, d.MagicDamageDealtToChampions,
			d.MagicDamageTaken, d.NeutralMinionsKilled, d.NeedVisionPings, d.NexusKills,
			d.NexusTakedowns, d.NexusLost, d.ObjectivesStolen, d.ObjectivesStolenAssists,
			d.OnMyWayPings, d.ParticipantId, d.PlayerScore0, d.PlayerScore1, d.PlayerScore2,
			d.PlayerScore3, d.PlayerScore4, d.PlayerScore5, d.PlayerScore6, d.PlayerScore7,
			d.PlayerScore8, d.PlayerScore9, d.PlayerScore10, d.PlayerScore11, d.PentaKills,
			d.PhysicalDamageDealt, d.PhysicalDamageDealtToChampions, d.PhysicalDamageTaken,
			d.Placement, d.PlayerAugment1, d.PlayerAugment2, d.PlayerAugment3, d.PlayerAugment4,
			d.PlayerSubteamId, d.PushPings, d.ProfileIcon, d.Puuid, d.QuadraKills, d.RiotIdGameName,
			d.RiotIdTagline, d.Role, d.SightWardsBoughtInGame, d.Spell1Casts, d.Spell2Casts,
			d.Spell3Casts, d.Spell4Casts, d.SubteamPlacement, d.Summoner1Casts, d.Summoner1Id,
			d.Summoner2Casts, d.Summoner2Id, d.SummonerId, d.SummonerLevel, d.SummonerName,
			d.TeamEarlySurrendered, d.TeamId, d.TeamPosition, d.TimeCCingOthers, d.TimePlayed,
			d.TotalAllyJungleMinionsKilled, d.TotalDamageDealt, d.TotalDamageDealtToChampions,
			d.TotalDamageShieldedOnTeammates, d.TotalDamageTaken, d.TotalEnemyJungleMinionsKilled,
			d.TotalHeal, d.TotalHealsOnTeammates, d.TotalMinionsKilled, d.TotalTimeCCDealt,
			d.TotalTimeSpentDead, d.TotalUnitsHealed, d.TripleKills, d.TrueDamageDealt,
			d.TrueDamageDealtToChampions, d.TrueDamageTaken, d.TurretKills, d.TurretTakedowns,
			d.TurretsLost, d.UnrealKills, d.VisionScore, d.VisionClearedPings, d.VisionWardsBoughtInGame,
			d.WardsKilled, d.WardsPlaced, d.Win,
		)
	}

	br := db.Conn.SendBatch(context.Background(), batch)
	defer br.Close()

	for range matchDetails {
		if _, err := br.Exec(); err != nil {
			utils.EditResponse(s, i, err.Error())
			return
		}
	}

	utils.EditResponse(s, i, "Matches added")
}

func addGroups(s *discordgo.Session, i *discordgo.InteractionCreate, names []string, riotApiKey string) {
	placeholders := make([]string, len(names))
	args := make([]interface{}, len(names))
	ctx := context.Background()
	for i, name := range names {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = name
	}

	query := fmt.Sprintf("SELECT puuid, player_name FROM server_players WHERE player_name IN (%s)", strings.Join(placeholders, ", "))

	rows, err := db.Conn.Query(context.Background(), query, args...)
	if err != nil {
		utils.EditResponse(s, i, err.Error())
		return
	}

	defer rows.Close()

	var results []string
	var groupName string
	for rows.Next() {
		var value string
		var name string
		if err := rows.Scan(&value, &name); err != nil {
			utils.EditResponse(s, i, err.Error())
			continue
		}
		results = append(results, value)
		groupName += name + ", "
	}

	if rows.Err() != nil {
		utils.EditResponse(s, i, rows.Err().Error())
	}

	var groupID int
	err = db.Conn.QueryRow(ctx, "SELECT nextval('group_id_seq')").Scan(&groupID)
	if err != nil {
		utils.EditResponse(s, i, err.Error())
	}

	batch := &pgx.Batch{}

	for _, puuid := range results {
		batch.Queue(`
			INSERT INTO groups (group_id, puuid, group_name)
			VALUES ($1, $2, $3)
		`, groupID, puuid, groupName)
	}

	br := db.Conn.SendBatch(ctx, batch)
	defer br.Close()
	utils.EditResponse(s, i, "Group Added")
}

func getGroups(s *discordgo.Session, i *discordgo.InteractionCreate, name string, riotApiKey string) {
	rows, err := db.Conn.Query(context.Background(), "SELECT puuid from groups WHERE group_name = $1", name)

	if err != nil {
		utils.EditResponse(s, i, err.Error())
		return
	}
	defer rows.Close()
	var puuids []string
	for rows.Next() {
		var puuid string
		if err = rows.Scan(&puuid); err != nil {
			utils.EditResponse(s, i, err.Error())
			return
		}

		puuids = append(puuids, puuid)
	}
	count := len(puuids)

	placeholders := make([]string, count)
	args := make([]interface{}, count)
	for i, id := range puuids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM (
			SELECT match_id
			FROM match_player_details
			WHERE puuid IN (%s) AND win = true
			GROUP BY match_id
			HAVING COUNT(DISTINCT puuid) = %d
		) as wins;
	`, strings.Join(placeholders, ", "), count)

	var winCount int
	err = db.Conn.QueryRow(context.Background(), query, args...).Scan(&winCount)
	if err != nil {
		utils.EditResponse(s, i, err.Error())
		return
	}

	query = fmt.Sprintf(`
			SELECT COUNT(*)
			FROM (
				SELECT match_id
				FROM match_player_details
				WHERE puuid IN (%s) AND win = false
				GROUP BY match_id
				HAVING COUNT(DISTINCT puuid) = %d
			) as losses;
		`, strings.Join(placeholders, ", "), count)

	var lossCount int
	err = db.Conn.QueryRow(context.Background(), query, args...).Scan(&lossCount)
	if err != nil {
		utils.EditResponse(s, i, err.Error())
		return
	}

	msg := fmt.Sprintf("%s has: %d wins and %d losses", name, winCount, lossCount)
	utils.EditResponse(s, i, msg)
}
