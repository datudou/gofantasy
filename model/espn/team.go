package espn

// Team is a fantasy team with optional roster and record populated depending
// on the requested league views.
type Team struct {
	ID       int      `json:"id"`
	Abbrev   string   `json:"abbrev"`
	Location string   `json:"location"`
	Nickname string   `json:"nickname"`
	Owners   []string `json:"owners"`
	Record   *Record  `json:"record,omitempty"`
	Roster   *Roster  `json:"roster,omitempty"`
	ValuesByStat map[string]float64 `json:"valuesByStat,omitempty"`
}

// TeamInfo is the lightweight team shape returned by some league views.
// Deprecated: prefer Team; kept for backward compatibility with early stubs.
type TeamInfo = Team

// Record holds season-to-date team results.
type Record struct {
	Overall OverallRecord `json:"overall"`
}

type OverallRecord struct {
	Wins         int     `json:"wins"`
	Losses       int     `json:"losses"`
	Ties         int     `json:"ties"`
	PointsFor    float64 `json:"pointsFor"`
	PointsAgainst float64 `json:"pointsAgainst"`
}

// Roster is a team's player list for a scoring period.
type Roster struct {
	Entries []RosterEntry `json:"entries"`
}

type RosterEntry struct {
	PlayerID         int              `json:"playerId"`
	LineupSlotID     int              `json:"lineupSlotId"`
	Status           string           `json:"status"`
	PlayerPoolEntry  PlayerPoolEntry  `json:"playerPoolEntry"`
}

// ManagedTeam ties a league to one team the authenticated member owns.
type ManagedTeam struct {
	Sport    string `json:"sport"`
	Season   int    `json:"season"`
	LeagueID int    `json:"leagueId"`
	Team     Team   `json:"team"`
}
