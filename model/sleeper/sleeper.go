package sleeper

// User is a Sleeper account.
type User struct {
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Avatar      string `json:"avatar"`
}

// League is a Sleeper fantasy league.
type League struct {
	LeagueID string         `json:"league_id"`
	Name     string         `json:"name"`
	Season   string         `json:"season"`
	Sport    string         `json:"sport"`
	Status   string         `json:"status"`
	Settings LeagueSettings `json:"settings"`
}

type LeagueSettings struct {
	NumTeams int `json:"num_teams"`
}

// LeagueUser is a manager in a league.
type LeagueUser struct {
	UserID      string            `json:"user_id"`
	DisplayName string            `json:"display_name"`
	Metadata    map[string]string `json:"metadata"`
}

// TeamName returns the fantasy team name when set in metadata.
func (u LeagueUser) TeamName() string {
	if u.Metadata != nil {
		if n := u.Metadata["team_name"]; n != "" {
			return n
		}
	}
	return u.DisplayName
}

// Roster is a team's players and record in a league.
type Roster struct {
	RosterID int            `json:"roster_id"`
	OwnerID  string         `json:"owner_id"`
	LeagueID string         `json:"league_id"`
	Players  []string       `json:"players"`
	Starters []string       `json:"starters"`
	Reserve  []string       `json:"reserve"`
	Taxi     []string       `json:"taxi"`
	Settings RosterSettings `json:"settings"`
}

type RosterSettings struct {
	Wins           int     `json:"wins"`
	Losses         int     `json:"losses"`
	Ties           int     `json:"ties"`
	Fpts           float64 `json:"fpts"`
	FptsAgainst    float64 `json:"fpts_against"`
	WaiverPosition int     `json:"waiver_position"`
}

// ResolvedRoster joins roster slots with player metadata.
type ResolvedRoster struct {
	RosterID int              `json:"roster_id"`
	OwnerID  string           `json:"owner_id"`
	LeagueID string           `json:"league_id"`
	Starters []ResolvedPlayer `json:"starters"`
	Bench    []ResolvedPlayer `json:"bench"`
	Reserve  []ResolvedPlayer `json:"reserve"`
	Settings RosterSettings   `json:"settings"`
}

type ResolvedPlayer struct {
	PlayerID string  `json:"player_id"`
	Slot     string  `json:"slot"`
	Player   *Player `json:"player,omitempty"`
}

// Player is a single athlete from the Sleeper player pool.
type Player struct {
	PlayerID         string   `json:"player_id"`
	FirstName        string   `json:"first_name"`
	LastName         string   `json:"last_name"`
	FullName         string   `json:"full_name"`
	Position         string   `json:"position"`
	FantasyPositions []string `json:"fantasy_positions"`
	Team             string   `json:"team"`
	Status           string   `json:"status"`
	InjuryStatus     string   `json:"injury_status"`
	Active           bool     `json:"active"`
	Number           int      `json:"number"`
}

// Matchup is one team's score for a week.
type Matchup struct {
	RosterID  int      `json:"roster_id"`
	MatchupID int      `json:"matchup_id"`
	Points    float64  `json:"points"`
	Starters  []string `json:"starters"`
	Players   []string `json:"players"`
}

// Transaction is a waiver, trade, or free-agent move.
type Transaction struct {
	TransactionID string           `json:"transaction_id"`
	Type          string           `json:"type"`
	Status        string           `json:"status"`
	Leg           int              `json:"leg"`
	RosterIDs     []int            `json:"roster_ids"`
	Adds          map[string]int   `json:"adds"`
	Drops         map[string]int   `json:"drops"`
	WaiverBudget  []map[string]any `json:"waiver_budget"`
}

// State is the current NFL/NBA season week from Sleeper.
type State struct {
	Season     string `json:"season"`
	SeasonType string `json:"season_type"`
	Week       int    `json:"week"`
	Leg        int    `json:"leg"`
}

// ManagedTeam is a user's roster in a league (for team discovery).
type ManagedTeam struct {
	Sport      string     `json:"sport"`
	Season     int        `json:"season"`
	League     League     `json:"league"`
	Roster     Roster     `json:"roster"`
	LeagueUser LeagueUser `json:"league_user"`
}

// PlayerProjection is one player's projected stat line for a week, as served
// by Sleeper's projections endpoint. Keys are Sleeper stat names — the
// projected fantasy totals live in "pts_ppr", "pts_half_ppr" and "pts_std",
// alongside component stats (rush_yd, rec, fgm, …). The endpoint is the one
// Sleeper's own app uses; it is not formally documented, so unknown or
// missing keys must be tolerated.
type PlayerProjection map[string]float64

// Points returns the projected fantasy points for a scoring format
// ("ppr", "half_ppr" or "std") and whether the projection carries it.
func (p PlayerProjection) Points(scoring string) (float64, bool) {
	v, ok := p["pts_"+scoring]
	return v, ok
}
