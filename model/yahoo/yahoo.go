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
	Team    Team     `xml:"team"`
}

type Team struct {
	TeamKey               string     `xml:"team_key"`
	TeamID                int        `xml:"team_id"`
	Name                  string     `xml:"name"`
	IsOwnedByCurrentLogin int        `xml:"is_owned_by_current_login"`
	URL                   string     `xml:"url"`
	TeamLogos             []TeamLogo `xml:"team_logos>team_logo"`
	WaiverPriority        int        `xml:"waiver_priority"`
	FAABBalance           int        `xml:"faab_balance"`
	NumberOfMoves         int        `xml:"number_of_moves"`
	NumberOfTrades        int        `xml:"number_of_trades"`
	RosterAdds            RosterAdds `xml:"roster_adds"`
	ClinchedPlayoffs      int        `xml:"clinched_playoffs"`
	LeagueScoringType     string     `xml:"league_scoring_type"`
	HasDraftGrade         int        `xml:"has_draft_grade"`
	AuctionBudgetTotal    int        `xml:"auction_budget_total"`
	AuctionBudgetSpent    int        `xml:"auction_budget_spent"`
	Managers              []Manager  `xml:"managers>manager"`
	Roster                Roster     `xml:"roster"`
}

type TeamLogo struct {
	Size string `xml:"size"`
	URL  string `xml:"url"`
}

type RosterAdds struct {
	CoverageType  string `xml:"coverage_type"`
	CoverageValue int    `xml:"coverage_value"`
	Value         int    `xml:"value"`
}

type Manager struct {
	ManagerID      int    `xml:"manager_id"`
	Nickname       string `xml:"nickname"`
	Guid           string `xml:"guid"`
	IsCommissioner int    `xml:"is_commissioner"`
	Email          string `xml:"email"`
	ImageURL       string `xml:"image_url"`
	FELOScore      int    `xml:"felo_score"`
	FELOTier       string `xml:"felo_tier"`
}

type User struct {
	Guid  string `xml:"guid"`
	Games []Game `xml:"games>game"`
}

type Game struct {
	GameKey string  `xml:"game_key"`
	GameID  string  `xml:"game_id"`
	Name    string  `xml:"name"`
	Code    string  `xml:"code"`
	Type    string  `xml:"type"`
	URL     string  `xml:"url"`
	Season  string  `xml:"season"`
	Teams   []*Team `xml:"teams>team"`
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

type Roster struct {
	CoverageType string   `xml:"coverage_type"`
	Date         string   `xml:"date"`
	Players      []Player `xml:"players>player"`
}

type Name struct {
	Full       string `xml:"full"`
	First      string `xml:"first"`
	Last       string `xml:"last"`
	AsciiFirst string `xml:"ascii_first"`
	AsciiLast  string `xml:"ascii_last"`
}

type EligiblePositions struct {
	Position []string `xml:"position"`
}

type SelectedPosition struct {
	CoverageType string `xml:"coverage_type"`
	Date         string `xml:"date"`
	Position     string `xml:"position"`
}

type Player struct {
	PlayerKey             string            `xml:"player_key"`
	PlayerID              int               `xml:"player_id"`
	Name                  Name              `xml:"name"`
	EditorialPlayerKey    string            `xml:"editorial_player_key"`
	EditorialTeamKey      string            `xml:"editorial_team_key"`
	EditorialTeamFullName string            `xml:"editorial_team_full_name"`
	EditorialTeamAbbr     string            `xml:"editorial_team_abbr"`
	UniformNumber         int               `xml:"uniform_number"`
	DisplayPosition       string            `xml:"display_position"`
	ImageURL              string            `xml:"image_url"`
	IsUndroppable         int               `xml:"is_undroppable"`
	PositionType          string            `xml:"position_type"`
	EligiblePositions     EligiblePositions `xml:"eligible_positions"`
	HasPlayerNotes        int               `xml:"has_player_notes"`
	SelectedPosition      SelectedPosition  `xml:"selected_position"`
}
