package espn

// Transaction is a league transaction (add, drop, trade, waiver).
type Transaction struct {
	ID              int                `json:"id"`
	Type            string             `json:"type"`
	Status          string             `json:"status"`
	MemberID        string             `json:"memberId"`
	TeamID          int                `json:"teamId"`
	ScoringPeriodID int                `json:"scoringPeriodId"`
	Items           []TransactionItem  `json:"items"`
}

type TransactionItem struct {
	Type              string `json:"type"`
	TeamID            int    `json:"teamId"`
	PlayerID          int    `json:"playerId"`
	FromTeamID        int    `json:"fromTeamId"`
	ToTeamID          int    `json:"toTeamId"`
	FromLineupSlotID  int    `json:"fromLineupSlotId"`
	ToLineupSlotID    int    `json:"toLineupSlotId"`
}
