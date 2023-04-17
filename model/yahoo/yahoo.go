package yahoo

import (
	"encoding/xml"
	"time"
)

// FantasyContent is the root level response containing the data from a request
// to the fantasy sports API.
type FantasyContent struct {
	XMLName xml.Name `xml:"fantasy_content"`
	League  League   `xml:"league"`
	Users   []User   `xml:"users>user"`
	Games   []Game   `xml:"games>game"`
}

type User struct {
	Guid  string `xml:"guid"`
	Games []Game `xml:"games>game"`
}

type Game struct {
	GameKey string `xml:"game_key"`
	GameID  string `xml:"game_id"`
	Name    string `xml:"name"`
	Code    string `xml:"code"`
	Type    string `xml:"type"`
	URL     string `xml:"url"`
	Season  string `xml:"season"`
}

type League struct {
	LeagueKey             string   `xml:"league_key"`
	LeagueID              string   `xml:"league_id"`
	Name                  string   `xml:"name"`
	URL                   string   `xml:"url"`
	DraftStatus           string   `xml:"draft_status"`
	NumTeams              int      `xml:"num_teams"`
	EditKey               int      `xml:"edit_key"`
	WeeklyDeadline        string   `xml:"weekly_deadline"`
	LeagueUpdateTimestamp int64    `xml:"league_update_timestamp"`
	ScoringType           string   `xml:"scoring_type"`
	CurrentWeek           int      `xml:"current_week"`
	StartWeek             int      `xml:"start_week"`
	EndWeek               int      `xml:"end_week"`
	GameCode              string   `xml:"game_code"`
	IsFinished            int      `xml:"is_finished"`
	Season                int      `xml:"season"`
	Settings              Settings `xml:"settings"`
}

type RosterPosition struct {
	Position string `xml:"position"`
	Count    int    `xml:"count"`
}

type Settings struct {
	DraftType               string           `xml:"draft_type"`
	ScoringType             string           `xml:"scoring_type"`
	UsesPlayoff             bool             `xml:"uses_playoff"`
	PlayoffStartWeek        int              `xml:"playoff_start_week"`
	UsesPlayoffReseeding    bool             `xml:"uses_playoff_reseeding"`
	UsesLockEliminatedTeams bool             `xml:"uses_lock_eliminated_teams"`
	UsesFAAB                bool             `xml:"uses_faab"`
	TradeEndDate            time.Time        `xml:"trade_end_date"`
	TradeRatifyType         string           `xml:"trade_ratify_type"`
	TradeRejectTime         int              `xml:"trade_reject_time"`
	RosterPositions         []RosterPosition `xml:"roster_positions>roster_position"`
	StatCategories          []Stat           `xml:"stat_categories>stats>stat"`
	StatModifiers           []Stat           `xml:"stat_modifiers>stats>stat"`
	Divisions               []Division       `xml:"divisions>division"`
}

type Stat struct {
	StatID       int    `xml:"stat_id"`
	Enabled      int    `xml:"enabled"`
	Name         string `xml:"name"`
	DisplayName  string `xml:"display_name"`
	SortOrder    int    `xml:"sort_order"`
	PositionType string `xml:"position_type"`
	Value        int    `xml:"value"`
}

type Division struct {
	DivisionID int    `xml:"division_id"`
	Name       string `xml:"name"`
}
