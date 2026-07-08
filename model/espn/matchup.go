package espn

// Matchup is one head-to-head pairing for a matchup period.
type Matchup struct {
	ID              int          `json:"id"`
	MatchupPeriodID int          `json:"matchupPeriodId"`
	Home            MatchupTeam  `json:"home"`
	Away            MatchupTeam  `json:"away"`
	Winner          string       `json:"winner"`
}

type MatchupTeam struct {
	TeamID      int     `json:"teamId"`
	TotalPoints float64 `json:"totalPoints"`
}

// ScheduleEntry is the schedule array shape returned with mSchedule/mMatchup views.
type ScheduleEntry struct {
	ID              int         `json:"id"`
	MatchupPeriodID int         `json:"matchupPeriodId"`
	Home            MatchupTeam `json:"home"`
	Away            MatchupTeam `json:"away"`
	Winner          string      `json:"winner"`
}
