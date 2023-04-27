package yahoo

import (
	"encoding/xml"
	"time"
)

// FantasyContent is the root level response containing the data from a request
// to the fantasy sports API.
type FantasyContent struct {
	XMLName xml.Name `xml:"fantasy_content,omitempty"`
	League  League   `xml:"league,omitempty"`
	Users   []User   `xml:"users>user,omitempty"`
	Games   []Game   `xml:"games>game,omitempty"`
	Team    Team     `xml:"team,omitempty"`
}

type Team struct {
	TeamKey               string     `xml:"team_key,omitempty"`
	TeamID                int        `xml:"team_id,omitempty"`
	Name                  string     `xml:"name,omitempty"`
	IsOwnedByCurrentLogin int        `xml:"is_owned_by_current_login,omitempty"`
	URL                   string     `xml:"url,omitempty"`
	TeamLogos             []TeamLogo `xml:"team_logos>team_logo,omitempty"`
	WaiverPriority        int        `xml:"waiver_priority,omitempty"`
	FAABBalance           int        `xml:"faab_balance,omitempty"`
	NumberOfMoves         int        `xml:"number_of_moves,omitempty"`
	NumberOfTrades        int        `xml:"number_of_trades,omitempty"`
	RosterAdds            RosterAdds `xml:"roster_adds,omitempty"`
	ClinchedPlayoffs      int        `xml:"clinched_playoffs,omitempty"`
	LeagueScoringType     string     `xml:"league_scoring_type,omitempty"`
	HasDraftGrade         int        `xml:"has_draft_grade,omitempty"`
	AuctionBudgetTotal    int        `xml:"auction_budget_total,omitempty"`
	AuctionBudgetSpent    int        `xml:"auction_budget_spent,omitempty"`
	Managers              []Manager  `xml:"managers>manager,omitempty"`
	Roster                Roster     `xml:"roster,omitempty"`
}

type TeamLogo struct {
	Size string `xml:"size,omitempty"`
	URL  string `xml:"url,omitempty"`
}

type RosterAdds struct {
	CoverageType  string `xml:"coverage_type,omitempty"`
	CoverageValue int    `xml:"coverage_value,omitempty"`
	Value         int    `xml:"value,omitempty"`
}

type Manager struct {
	ManagerID      int    `xml:"manager_id,omitempty"`
	Nickname       string `xml:"nickname,omitempty"`
	Guid           string `xml:"guid,omitempty"`
	IsCommissioner int    `xml:"is_commissioner,omitempty"`
	Email          string `xml:"email,omitempty"`
	ImageURL       string `xml:"image_url,omitempty"`
	FELOScore      int    `xml:"felo_score,omitempty"`
	FELOTier       string `xml:"felo_tier,omitempty"`
}

type User struct {
	Guid  string `xml:"guid,omitempty"`
	Games []Game `xml:"games>game,omitempty"`
}

type Game struct {
	GameKey string  `xml:"game_key,omitempty"`
	GameID  string  `xml:"game_id,omitempty"`
	Name    string  `xml:"name,omitempty"`
	Code    string  `xml:"code,omitempty"`
	Type    string  `xml:"type,omitempty"`
	URL     string  `xml:"url,omitempty"`
	Season  string  `xml:"season,omitempty"`
	Teams   []*Team `xml:"teams>team,omitempty"`
}

type League struct {
	LeagueKey             string   `xml:"league_key,omitempty"`
	LeagueID              string   `xml:"league_id,omitempty"`
	Name                  string   `xml:"name,omitempty"`
	URL                   string   `xml:"url,omitempty"`
	DraftStatus           string   `xml:"draft_status,omitempty"`
	NumTeams              int      `xml:"num_teams,omitempty"`
	EditKey               int      `xml:"edit_key,omitempty"`
	WeeklyDeadline        string   `xml:"weekly_deadline,omitempty"`
	LeagueUpdateTimestamp int64    `xml:"league_update_timestamp,omitempty"`
	ScoringType           string   `xml:"scoring_type,omitempty"`
	CurrentWeek           int      `xml:"current_week,omitempty"`
	StartWeek             int      `xml:"start_week,omitempty"`
	EndWeek               int      `xml:"end_week,omitempty"`
	GameCode              string   `xml:"game_code,omitempty"`
	IsFinished            int      `xml:"is_finished,omitempty"`
	Season                int      `xml:"season,omitempty"`
	Settings              Settings `xml:"settings,omitempty"`
}

type RosterPosition struct {
	Position string `xml:"position,omitempty"`
	Count    int    `xml:"count,omitempty"`
}

type Settings struct {
	DraftType               string           `xml:"draft_type,omitempty"`
	ScoringType             string           `xml:"scoring_type,omitempty"`
	UsesPlayoff             bool             `xml:"uses_playoff,omitempty"`
	PlayoffStartWeek        int              `xml:"playoff_start_week,omitempty"`
	UsesPlayoffReseeding    bool             `xml:"uses_playoff_reseeding,omitempty"`
	UsesLockEliminatedTeams bool             `xml:"uses_lock_eliminated_teams,omitempty"`
	UsesFAAB                bool             `xml:"uses_faab,omitempty"`
	TradeEndDate            time.Time        `xml:"trade_end_date,omitempty"`
	TradeRatifyType         string           `xml:"trade_ratify_type,omitempty"`
	TradeRejectTime         int              `xml:"trade_reject_time,omitempty"`
	RosterPositions         []RosterPosition `xml:"roster_positions>roster_position,omitempty"`
	StatCategories          []Stat           `xml:"stat_categories>stats>stat,omitempty"`
	StatModifiers           []Stat           `xml:"stat_modifiers>stats>stat,omitempty"`
	Divisions               []Division       `xml:"divisions>division,omitempty"`
}

type Stat struct {
	StatID       int    `xml:"stat_id,omitempty"`
	Enabled      int    `xml:"enabled,omitempty"`
	Name         string `xml:"name,omitempty"`
	DisplayName  string `xml:"display_name,omitempty"`
	SortOrder    int    `xml:"sort_order,omitempty"`
	PositionType string `xml:"position_type,omitempty"`
	Value        int    `xml:"value,omitempty"`
}

type Division struct {
	DivisionID int    `xml:"division_id,omitempty"`
	Name       string `xml:"name,omitempty"`
}

type Roster struct {
	CoverageType string   `xml:"coverage_type,omitempty"`
	Date         string   `xml:"date,omitempty"`
	Players      []Player `xml:"players>player,omitempty"`
}

type Name struct {
	Full       string `xml:"full,omitempty"`
	First      string `xml:"first,omitempty"`
	Last       string `xml:"last,omitempty"`
	AsciiFirst string `xml:"ascii_first,omitempty"`
	AsciiLast  string `xml:"ascii_last,omitempty"`
}

type EligiblePositions struct {
	Position []string `xml:"position,omitempty"`
}

type SelectedPosition struct {
	CoverageType string `xml:"coverage_type,omitempty"`
	Date         string `xml:"date,omitempty"`
	Position     string `xml:"position,omitempty"`
}

type Player struct {
	PlayerKey             string            `xml:"player_key,omitempty"`
	PlayerID              int               `xml:"player_id,omitempty"`
	Name                  Name              `xml:"name,omitempty"`
	EditorialPlayerKey    string            `xml:"editorial_player_key,omitempty"`
	EditorialTeamKey      string            `xml:"editorial_team_key,omitempty"`
	EditorialTeamFullName string            `xml:"editorial_team_full_name,omitempty"`
	EditorialTeamAbbr     string            `xml:"editorial_team_abbr,omitempty"`
	UniformNumber         int               `xml:"uniform_number,omitempty"`
	DisplayPosition       string            `xml:"display_position,omitempty"`
	ImageURL              string            `xml:"image_url,omitempty"`
	IsUndroppable         int               `xml:"is_undroppable,omitempty"`
	PositionType          string            `xml:"position_type,omitempty"`
	EligiblePositions     EligiblePositions `xml:"eligible_positions,omitempty"`
	HasPlayerNotes        int               `xml:"has_player_notes,omitempty"`
	SelectedPosition      SelectedPosition  `xml:"selected_position,omitempty"`
}
