package espn

type League struct {
	GameId          int             `json:"gameId"`
	ID              int             `json:"id"`
	Members         []Member        `json:"members"`
	ScoringPeriodId int             `json:"scoringPeriodId"`
	SeasonId        int             `json:"seasonId"`
	SegmentId       int             `json:"segmentId"`
	Settings        Settings        `json:"settings"`
	Status          LeagueStatus    `json:"status"`
	Teams           []Team          `json:"teams"`
	Schedule        []ScheduleEntry `json:"schedule"`
	DraftDetail     DraftDetail     `json:"draftDetail"`
	Transactions    []Transaction   `json:"transactions"`
	Players         []PlayerPoolEntry `json:"players,omitempty"`
}

// Member is a league member (owner/manager).
type Member struct {
	ID              string `json:"id"`
	DisplayName     string `json:"displayName"`
	IsLeagueManager bool   `json:"isLeagueManager"`
}

// Memeber is a deprecated misspelling kept as an alias.
type Memeber = Member

type Settings struct {
	Name           string         `json:"name"`
	Size           int            `json:"size"`
	RosterSettings RosterSettings `json:"rosterSettings"`
	ScoringSettings ScoringSettings `json:"scoringSettings"`
}

type RosterSettings struct {
	LineupSlotCounts map[string]int `json:"lineupSlotCounts"`
}

type ScoringSettings struct {
	ScoringItems []ScoringItem `json:"scoringItems"`
}

type ScoringItem struct {
	StatID int     `json:"statId"`
	Points float64 `json:"points"`
}

type LeagueStatus struct {
	CurrentMatchupPeriod int  `json:"currentMatchupPeriod"`
	IsActive             bool `json:"isActive"`
	LatestScoringPeriod  int  `json:"latestScoringPeriod"`
	IsComplete           bool `json:"isComplete"`
}

type Link struct {
	Language   string   `json:"language"`
	Rel        []string `json:"rel"`
	Href       string   `json:"href"`
	Text       string   `json:"text"`
	ShortText  string   `json:"shortText"`
	IsExternal bool     `json:"isExternal"`
	IsPremium  bool     `json:"isPremium"`
}

type DraftDetail struct {
	CompleteDate int64  `json:"completeDate"`
	Drafted      bool   `json:"drafted"`
	InProgress   bool   `json:"inProgress"`
	Picks        []Pick `json:"picks"`
}

type Pick struct {
	AutoDraftTypeId   int  `json:"autoDraftTypeId"`
	BidAmount         int  `json:"bidAmount"`
	Id                int  `json:"id"`
	Keeper            bool `json:"keeper"`
	LineupSlotId      int  `json:"lineupSlotId"`
	NominatingTeamId  int  `json:"nominatingTeamId"`
	OverallPickNumber int  `json:"overallPickNumber"`
	PlayerId          int  `json:"playerId"`
	ReservedForKeeper bool `json:"reservedForKeeper"`
	RoundId           int  `json:"roundId"`
	RoundPickNumber   int  `json:"roundPickNumber"`
	TeamId            int  `json:"teamId"`
	TradeLocked       bool `json:"tradeLocked"`
}

type GameDetail struct {
	DraftDetail     DraftDetail `json:"draftDetail"`
	GameId          int         `json:"gameId"`
	Id              int         `json:"id"`
	ScoringPeriodId int         `json:"scoringPeriodId"`
	SeasonId        int         `json:"seasonId"`
	SegmentId       int         `json:"segmentId"`
	Settings        Settings    `json:"settings"`
	Status          Status      `json:"status"`
}

type DraftSettings struct {
	AuctionBudget      int    `json:"auctionBudget"`
	AvailableDate      int64  `json:"availableDate"`
	Date               int64  `json:"date"`
	IsTradingEnabled   bool   `json:"isTradingEnabled"`
	KeeperCount        int    `json:"keeperCount"`
	KeeperCountFuture  int    `json:"keeperCountFuture"`
	KeeperDeadlineDate int64  `json:"keeperDeadlineDate"`
	KeeperOrderType    string `json:"keeperOrderType"`
	LeagueSubType      string `json:"leagueSubType"`
	OrderType          string `json:"orderType"`
	PickOrder          []int  `json:"pickOrder"`
	TimePerSelection   int    `json:"timePerSelection"`
	Type               string `json:"type"`
}

type Status struct {
	ActivatedDate            int64          `json:"activatedDate"`
	CreatedAsLeagueType      int            `json:"createdAsLeagueType"`
	CurrentLeagueType        int            `json:"currentLeagueType"`
	CurrentMatchupPeriod     int            `json:"currentMatchupPeriod"`
	FinalScoringPeriod       int            `json:"finalScoringPeriod"`
	FirstScoringPeriod       int            `json:"firstScoringPeriod"`
	IsActive                 bool           `json:"isActive"`
	IsExpired                bool           `json:"isExpired"`
	IsFull                   bool           `json:"isFull"`
	IsPlayoffMatchupEdited   bool           `json:"isPlayoffMatchupEdited"`
	IsToBeDeleted            bool           `json:"isToBeDeleted"`
	IsViewable               bool           `json:"isViewable"`
	IsWaiverOrderEdited      bool           `json:"isWaiverOrderEdited"`
	LatestScoringPeriod      int            `json:"latestScoringPeriod"`
	PreviousSeasons          []int          `json:"previousSeasons"`
	StandingsUpdateDate      int64          `json:"standingsUpdateDate"`
	TeamsJoined              int            `json:"teamsJoined"`
	TransactionScoringPeriod int            `json:"transactionScoringPeriod"`
	WaiverLastExecutionDate  int64          `json:"waiverLastExecutionDate"`
	WaiverNextExecutionDate  int64          `json:"waiverNextExecutionDate"`
	WaiverProcessStatus      map[string]int `json:"waiverProcessStatus"`
}

type Season struct {
	Abbrev string `json:"abbrev"`
	Active bool   `json:"active"`
	//CurrentScoringPeriod ScoringPeriod `json:"currentScoringPeriod"`
	Display      bool   `json:"display"`
	DisplayOrder int    `json:"displayOrder"`
	EndDate      int64  `json:"endDate"`
	GameId       int    `json:"gameId"`
	Id           int    `json:"id"`
	Name         string `json:"name"`
	StartDate    int64  `json:"startDate"`
}
