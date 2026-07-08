package espn

// PlayerPoolEntry wraps player metadata as returned inside roster entries and
// free-agent pool responses.
type PlayerPoolEntry struct {
	ID               int     `json:"id"`
	Player           Player  `json:"player"`
	Status           string  `json:"status"`
	AppliedStatTotal float64 `json:"appliedStatTotal"`
}

// Player is the core athlete record from ESPN fantasy payloads.
type Player struct {
	ID                  int    `json:"id"`
	FullName            string `json:"fullName"`
	FirstName           string `json:"firstName"`
	LastName            string `json:"lastName"`
	DefaultPositionID   int    `json:"defaultPositionId"`
	EligibleSlots       []int  `json:"eligibleSlots"`
	ProTeamID           int    `json:"proTeamId"`
	Active              bool   `json:"active"`
	Injured             bool   `json:"injured"`
	InjuryStatus        string `json:"injuryStatus"`
	Ownership           *Ownership `json:"ownership,omitempty"`
	Stats               []StatSplit `json:"stats,omitempty"`
}

type Ownership struct {
	PercentOwned   float64 `json:"percentOwned"`
	PercentStarted float64 `json:"percentStarted"`
}

type StatSplit struct {
	SeasonID        int                `json:"seasonId"`
	StatSourceID    int                `json:"statSourceId"`
	StatSplitTypeID int                `json:"statSplitTypeId"`
	Stats           map[string]float64 `json:"stats"`
}
