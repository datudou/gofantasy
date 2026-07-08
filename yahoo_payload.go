package gofantasy

import "encoding/xml"

// The Yahoo Fantasy API expects write operations to be submitted as XML bodies
// rooted at <fantasy_content>. The structs below model those request bodies.
// They are intentionally kept separate from the response models in
// model/yahoo because the request shapes differ (e.g. transaction_data on a
// player, coverage_type on a roster) from what Yahoo returns.

// rosterPayload is the body for PUT /team/{team_key}/roster.
type rosterPayload struct {
	XMLName xml.Name   `xml:"fantasy_content"`
	Roster  rosterBody `xml:"roster"`
}

type rosterBody struct {
	CoverageType string         `xml:"coverage_type"`
	Date         string         `xml:"date,omitempty"`
	Week         int            `xml:"week,omitempty"`
	Players      []rosterPlayer `xml:"players>player"`
}

type rosterPlayer struct {
	PlayerKey string `xml:"player_key"`
	Position  string `xml:"position"`
}

// transactionPayload is the body for POST /league/{league_key}/transactions
// (add, drop, add/drop, pending_trade) as well as
// PUT /transaction/{transaction_key} (accept/reject/allow/disallow a trade).
type transactionPayload struct {
	XMLName     xml.Name        `xml:"fantasy_content"`
	Transaction transactionBody `xml:"transaction"`
}

type transactionBody struct {
	Type           string `xml:"type"`
	FAABBid        *int   `xml:"faab_bid,omitempty"`
	TraderTeamKey  string `xml:"trader_team_key,omitempty"`
	TradeeTeamKey  string `xml:"tradee_team_key,omitempty"`
	TradeNote      string `xml:"trade_note,omitempty"`
	TransactionKey string `xml:"transaction_key,omitempty"`
	Action         string `xml:"action,omitempty"`
	VoterTeamKey   string `xml:"voter_team_key,omitempty"`
	// Player is used for single-player add or drop transactions.
	// Note: encoding/xml does not honour omitempty on nested "parent>child"
	// paths (it still emits an empty wrapper), so Players is modelled as a
	// pointer wrapper to guarantee the <players> element is dropped entirely
	// when unused.
	Player *txnPlayer `xml:"player,omitempty"`
	// Players is used for multi-player transactions (add/drop, pending_trade).
	Players *txnPlayers `xml:"players,omitempty"`
}

type txnPlayers struct {
	Player []txnPlayer `xml:"player"`
}

type txnPlayer struct {
	PlayerKey       string      `xml:"player_key"`
	TransactionData txnPlayerTx `xml:"transaction_data"`
}

type txnPlayerTx struct {
	Type               string `xml:"type"`
	SourceTeamKey      string `xml:"source_team_key,omitempty"`
	DestinationTeamKey string `xml:"destination_team_key,omitempty"`
}
