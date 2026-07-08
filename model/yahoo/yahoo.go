package yahoo

import (
	"encoding/xml"
)

// FantasyContent is the root level response containing the data from a request
// to the fantasy sports API.
type FantasyContent struct {
	XMLName     xml.Name    `xml:"fantasy_content,omitempty"`
	League      League      `xml:"league,omitempty"`
	Users       []User      `xml:"users>user,omitempty"`
	Games       []Game      `xml:"games>game,omitempty"`
	Team        Team        `xml:"team,omitempty"`
	Player      Player      `xml:"player,omitempty"`
	Players     []Player    `xml:"players>player,omitempty"`
	Transaction Transaction `xml:"transaction,omitempty"`
}

type Team struct {
	TeamKey               string        `xml:"team_key,omitempty"`
	TeamID                int           `xml:"team_id,omitempty"`
	Name                  string        `xml:"name,omitempty"`
	IsOwnedByCurrentLogin int           `xml:"is_owned_by_current_login,omitempty"`
	URL                   string        `xml:"url,omitempty"`
	TeamLogos             []TeamLogo    `xml:"team_logos>team_logo,omitempty"`
	WaiverPriority        int           `xml:"waiver_priority,omitempty"`
	FAABBalance           int           `xml:"faab_balance,omitempty"`
	NumberOfMoves         int           `xml:"number_of_moves,omitempty"`
	NumberOfTrades        int           `xml:"number_of_trades,omitempty"`
	RosterAdds            RosterAdds    `xml:"roster_adds,omitempty"`
	ClinchedPlayoffs      int           `xml:"clinched_playoffs,omitempty"`
	LeagueScoringType     string        `xml:"league_scoring_type,omitempty"`
	HasDraftGrade         int           `xml:"has_draft_grade,omitempty"`
	AuctionBudgetTotal    int           `xml:"auction_budget_total,omitempty"`
	AuctionBudgetSpent    int           `xml:"auction_budget_spent,omitempty"`
	Managers              []Manager     `xml:"managers>manager,omitempty"`
	Roster                Roster        `xml:"roster,omitempty"`
	TeamStandings         TeamStandings `xml:"team_standings,omitempty"`
	TeamStats             TeamStats     `xml:"team_stats,omitempty"`
	TeamPoints            TeamPoints    `xml:"team_points,omitempty"`
	TeamProjectedPoints   TeamPoints    `xml:"team_projected_points,omitempty"`
	WinProbability        float64       `xml:"win_probability,omitempty"`
	Matchups              []Matchup     `xml:"matchups>matchup,omitempty"`
}

// TeamStandings holds standings information for a single team.
type TeamStandings struct {
	Rank          int           `xml:"rank,omitempty"`
	PlayoffSeed   int           `xml:"playoff_seed,omitempty"`
	OutcomeTotals OutcomeTotals `xml:"outcome_totals,omitempty"`
	GamesBack     string        `xml:"games_back,omitempty"`
	PointsFor     float64       `xml:"points_for,omitempty"`
	PointsAgainst float64       `xml:"points_against,omitempty"`
	Streak        Streak        `xml:"streak,omitempty"`
}

// OutcomeTotals holds win/loss/tie counts.
type OutcomeTotals struct {
	Wins       int     `xml:"wins,omitempty"`
	Losses     int     `xml:"losses,omitempty"`
	Ties       int     `xml:"ties,omitempty"`
	Percentage float64 `xml:"percentage,omitempty"`
}

// Streak holds the current win/loss streak.
type Streak struct {
	Type  string `xml:"type,omitempty"`
	Value int    `xml:"value,omitempty"`
}

// TeamStats holds the stats for a single team for a given coverage.
type TeamStats struct {
	CoverageType string `xml:"coverage_type,omitempty"`
	Week         int    `xml:"week,omitempty"`
	Season       int    `xml:"season,omitempty"`
	Date         string `xml:"date,omitempty"`
	Stats        []Stat `xml:"stats>stat,omitempty"`
}

// TeamPoints holds total fantasy points scored by a team for a coverage.
type TeamPoints struct {
	CoverageType string  `xml:"coverage_type,omitempty"`
	Week         int     `xml:"week,omitempty"`
	Season       int     `xml:"season,omitempty"`
	Total        float64 `xml:"total,omitempty"`
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
	LeagueKey             string        `xml:"league_key,omitempty"`
	LeagueID              string        `xml:"league_id,omitempty"`
	Name                  string        `xml:"name,omitempty"`
	URL                   string        `xml:"url,omitempty"`
	DraftStatus           string        `xml:"draft_status,omitempty"`
	NumTeams              int           `xml:"num_teams,omitempty"`
	EditKey               string        `xml:"edit_key,omitempty"`
	WeeklyDeadline        string        `xml:"weekly_deadline,omitempty"`
	LeagueUpdateTimestamp int64         `xml:"league_update_timestamp,omitempty"`
	ScoringType           string        `xml:"scoring_type,omitempty"`
	CurrentWeek           int           `xml:"current_week,omitempty"`
	StartWeek             int           `xml:"start_week,omitempty"`
	EndWeek               int           `xml:"end_week,omitempty"`
	GameCode              string        `xml:"game_code,omitempty"`
	IsFinished            int           `xml:"is_finished,omitempty"`
	Season                int           `xml:"season,omitempty"`
	Settings              Settings      `xml:"settings,omitempty"`
	Standings             []Team        `xml:"standings>teams>team,omitempty"`
	Scoreboard            Scoreboard    `xml:"scoreboard,omitempty"`
	Teams                 []Team        `xml:"teams>team,omitempty"`
	Players               []Player      `xml:"players>player,omitempty"`
	Transactions          []Transaction `xml:"transactions>transaction,omitempty"`
	DraftResults          []DraftResult `xml:"draft_results>draft_result,omitempty"`
}

// Scoreboard contains the matchups for a particular week of a league.
type Scoreboard struct {
	Week     int       `xml:"week,omitempty"`
	Matchups []Matchup `xml:"matchups>matchup,omitempty"`
}

// Matchup is a single head-to-head matchup between two teams in a week.
type Matchup struct {
	Week                int     `xml:"week,omitempty"`
	WeekStart           string  `xml:"week_start,omitempty"`
	WeekEnd             string  `xml:"week_end,omitempty"`
	Status              string  `xml:"status,omitempty"`
	IsPlayoffs          int     `xml:"is_playoffs,omitempty"`
	IsConsolation       int     `xml:"is_consolation,omitempty"`
	IsTied              int     `xml:"is_tied,omitempty"`
	WinnerTeamKey       string  `xml:"winner_team_key,omitempty"`
	IsMatchupRecapAvail int     `xml:"is_matchup_recap_available,omitempty"`
	MatchupRecapURL     string  `xml:"matchup_recap_url,omitempty"`
	MatchupRecapTitle   string  `xml:"matchup_recap_title,omitempty"`
	MatchupGrades       []Grade `xml:"matchup_grades>matchup_grade,omitempty"`
	Teams               []Team  `xml:"teams>team,omitempty"`
}

// Grade represents a per-team grade for a matchup.
type Grade struct {
	TeamKey string `xml:"team_key,omitempty"`
	Grade   string `xml:"grade,omitempty"`
}

// Transaction represents a single transaction (add/drop/trade/commish action).
type Transaction struct {
	TransactionKey string              `xml:"transaction_key,omitempty"`
	TransactionID  int                 `xml:"transaction_id,omitempty"`
	Type           string              `xml:"type,omitempty"`
	Status         string              `xml:"status,omitempty"`
	Timestamp      int64               `xml:"timestamp,omitempty"`
	FAABBid        int                 `xml:"faab_bid,omitempty"`
	TraderTeamKey  string              `xml:"trader_team_key,omitempty"`
	TradeeTeamKey  string              `xml:"tradee_team_key,omitempty"`
	Players        []TransactionPlayer `xml:"players>player,omitempty"`
}

// TransactionPlayer is a player associated with a transaction. It embeds Player
// and adds a TransactionData section describing the move (add/drop).
type TransactionPlayer struct {
	Player
	TransactionData []TransactionData `xml:"transaction_data,omitempty"`
}

// TransactionData describes a single add/drop action for a transaction player.
type TransactionData struct {
	Type                string `xml:"type,omitempty"`
	SourceType          string `xml:"source_type,omitempty"`
	SourceTeamKey       string `xml:"source_team_key,omitempty"`
	SourceTeamName      string `xml:"source_team_name,omitempty"`
	DestinationType     string `xml:"destination_type,omitempty"`
	DestinationTeamKey  string `xml:"destination_team_key,omitempty"`
	DestinationTeamName string `xml:"destination_team_name,omitempty"`
}

// DraftResult is one pick in the draft.
type DraftResult struct {
	Pick      int    `xml:"pick,omitempty"`
	Round     int    `xml:"round,omitempty"`
	Cost      int    `xml:"cost,omitempty"`
	TeamKey   string `xml:"team_key,omitempty"`
	PlayerKey string `xml:"player_key,omitempty"`
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
	TradeEndDate            string           `xml:"trade_end_date,omitempty"`
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
	Status                string            `xml:"status,omitempty"`
	StatusFull            string            `xml:"status_full,omitempty"`
	InjuryNote            string            `xml:"injury_note,omitempty"`
	OwnershipType         string            `xml:"ownership>ownership_type,omitempty"`
	OwnerTeamKey          string            `xml:"ownership>owner_team_key,omitempty"`
	OwnerTeamName         string            `xml:"ownership>owner_team_name,omitempty"`
	PercentOwned          PercentOwned      `xml:"percent_owned,omitempty"`
	PlayerStats           PlayerStats       `xml:"player_stats,omitempty"`
	PlayerPoints          PlayerPoints      `xml:"player_points,omitempty"`
}

// PercentOwned reports what percent of leagues the player is owned in.
type PercentOwned struct {
	CoverageType string  `xml:"coverage_type,omitempty"`
	Week         int     `xml:"week,omitempty"`
	Value        float64 `xml:"value,omitempty"`
	Delta        float64 `xml:"delta,omitempty"`
}

// PlayerStats holds a single player's stats for a given coverage.
type PlayerStats struct {
	CoverageType string `xml:"coverage_type,omitempty"`
	Week         int    `xml:"week,omitempty"`
	Season       int    `xml:"season,omitempty"`
	Date         string `xml:"date,omitempty"`
	Stats        []Stat `xml:"stats>stat,omitempty"`
}

// PlayerPoints holds a single player's total fantasy points for a coverage.
type PlayerPoints struct {
	CoverageType string  `xml:"coverage_type,omitempty"`
	Week         int     `xml:"week,omitempty"`
	Season       int     `xml:"season,omitempty"`
	Total        float64 `xml:"total,omitempty"`
}
