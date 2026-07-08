package gofantasy

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/gofantasy/model/yahoo"
)

// PlayerSlot assigns a player to a roster position for a SetRoster call.
// Position is a Yahoo position code such as "QB", "BN" (bench), "IL" (injured
// list) or "IR". Benching a player uses "BN".
type PlayerSlot struct {
	PlayerKey string
	Position  string
}

// Valid trade response actions for RespondToTrade.
const (
	TradeActionAccept   = "accept"
	TradeActionReject   = "reject"
	TradeActionAllow    = "allow"    // commissioner: allow a proposed trade
	TradeActionDisallow = "disallow" // commissioner: veto a proposed trade
)

// SetRoster sets the daily lineup for a team on the given date (YYYY-MM-DD).
// Use this for daily-coverage sports (MLB, NBA, NHL). For weekly sports (NFL)
// use SetRosterForWeek.
func (y *yahooClient) SetRoster(ctx context.Context, teamKey string, date string, assignments []PlayerSlot) error {
	if len(assignments) == 0 {
		return fmt.Errorf("no roster assignments provided")
	}
	payload := rosterPayload{
		Roster: rosterBody{
			CoverageType: "date",
			Date:         date,
			Players:      toRosterPlayers(assignments),
		},
	}
	endpoint := fmt.Sprintf("%s/team/%s/roster", y.baseUrl, teamKey)
	_, err := y.writeXML(ctx, http.MethodPut, endpoint, payload)
	return err
}

// SetRosterForWeek sets the lineup for a team for a given week. Use this for
// weekly-coverage sports (NFL).
func (y *yahooClient) SetRosterForWeek(ctx context.Context, teamKey string, week int, assignments []PlayerSlot) error {
	if len(assignments) == 0 {
		return fmt.Errorf("no roster assignments provided")
	}
	payload := rosterPayload{
		Roster: rosterBody{
			CoverageType: "week",
			Week:         week,
			Players:      toRosterPlayers(assignments),
		},
	}
	endpoint := fmt.Sprintf("%s/team/%s/roster", y.baseUrl, teamKey)
	_, err := y.writeXML(ctx, http.MethodPut, endpoint, payload)
	return err
}

// AddPlayer claims a free agent / waiver player onto the given team. faabBid may
// be nil for non-FAAB leagues; for FAAB leagues it is the bid amount.
func (y *yahooClient) AddPlayer(ctx context.Context, leagueKey, teamKey, playerKey string, faabBid *int) (*yahoo.Transaction, error) {
	payload := transactionPayload{
		Transaction: transactionBody{
			Type:    "add",
			FAABBid: faabBid,
			Player: &txnPlayer{
				PlayerKey: playerKey,
				TransactionData: txnPlayerTx{
					Type:               "add",
					DestinationTeamKey: teamKey,
				},
			},
		},
	}
	return y.submitTransaction(ctx, leagueKey, payload)
}

// DropPlayer drops a player from the given team back to the free agent pool.
func (y *yahooClient) DropPlayer(ctx context.Context, leagueKey, teamKey, playerKey string) (*yahoo.Transaction, error) {
	payload := transactionPayload{
		Transaction: transactionBody{
			Type: "drop",
			Player: &txnPlayer{
				PlayerKey: playerKey,
				TransactionData: txnPlayerTx{
					Type:          "drop",
					SourceTeamKey: teamKey,
				},
			},
		},
	}
	return y.submitTransaction(ctx, leagueKey, payload)
}

// AddDropPlayer atomically adds one player and drops another in a single
// transaction. This is the correct call when a roster is full.
func (y *yahooClient) AddDropPlayer(ctx context.Context, leagueKey, teamKey, addKey, dropKey string, faabBid *int) (*yahoo.Transaction, error) {
	payload := transactionPayload{
		Transaction: transactionBody{
			Type:    "add/drop",
			FAABBid: faabBid,
			Players: &txnPlayers{Player: []txnPlayer{
				{
					PlayerKey: addKey,
					TransactionData: txnPlayerTx{
						Type:               "add",
						DestinationTeamKey: teamKey,
					},
				},
				{
					PlayerKey: dropKey,
					TransactionData: txnPlayerTx{
						Type:          "drop",
						SourceTeamKey: teamKey,
					},
				},
			}},
		},
	}
	return y.submitTransaction(ctx, leagueKey, payload)
}

// ProposeTrade proposes a trade from traderTeamKey to tradeeTeamKey. send is the
// list of player keys traderTeamKey gives up; receive is the list of player keys
// traderTeamKey wants from tradeeTeamKey. note is an optional message.
func (y *yahooClient) ProposeTrade(ctx context.Context, leagueKey, traderTeamKey, tradeeTeamKey string, send, receive []string, note string) (*yahoo.Transaction, error) {
	if len(send) == 0 && len(receive) == 0 {
		return nil, fmt.Errorf("a trade must include at least one player")
	}
	players := make([]txnPlayer, 0, len(send)+len(receive))
	for _, pk := range send {
		players = append(players, txnPlayer{
			PlayerKey: pk,
			TransactionData: txnPlayerTx{
				Type:               "pending_trade",
				SourceTeamKey:      traderTeamKey,
				DestinationTeamKey: tradeeTeamKey,
			},
		})
	}
	for _, pk := range receive {
		players = append(players, txnPlayer{
			PlayerKey: pk,
			TransactionData: txnPlayerTx{
				Type:               "pending_trade",
				SourceTeamKey:      tradeeTeamKey,
				DestinationTeamKey: traderTeamKey,
			},
		})
	}
	payload := transactionPayload{
		Transaction: transactionBody{
			Type:          "pending_trade",
			TraderTeamKey: traderTeamKey,
			TradeeTeamKey: tradeeTeamKey,
			TradeNote:     note,
			Players:       &txnPlayers{Player: players},
		},
	}
	return y.submitTransaction(ctx, leagueKey, payload)
}

// RespondToTrade accepts, rejects, allows or disallows a pending trade.
// action must be one of the TradeAction* constants. voteAgainst is only used by
// commissioners in leagues that vote on trades.
func (y *yahooClient) RespondToTrade(ctx context.Context, transactionKey, action string, voteAgainst *bool) error {
	switch action {
	case TradeActionAccept, TradeActionReject, TradeActionAllow, TradeActionDisallow:
	default:
		return fmt.Errorf("invalid trade action %q", action)
	}
	body := transactionBody{
		TransactionKey: transactionKey,
		Type:           "pending_trade",
		Action:         action,
	}
	payload := transactionPayload{Transaction: body}
	endpoint := fmt.Sprintf("%s/transaction/%s", y.baseUrl, transactionKey)
	_, err := y.writeXML(ctx, http.MethodPut, endpoint, payload)
	return err
}

// CancelTransaction cancels a pending transaction (e.g. a waiver claim or a
// trade you proposed) identified by its transaction key.
func (y *yahooClient) CancelTransaction(ctx context.Context, transactionKey string) error {
	endpoint := fmt.Sprintf("%s/transaction/%s", y.baseUrl, transactionKey)
	_, err := y.baseClient.requestor.Delete(ctx, endpoint, nil, xmlDecorator, &xmlDecoder{})
	return err
}

func (y *yahooClient) submitTransaction(ctx context.Context, leagueKey string, payload transactionPayload) (*yahoo.Transaction, error) {
	endpoint := fmt.Sprintf("%s/league/%s/transactions", y.baseUrl, leagueKey)
	fc, err := y.writeXML(ctx, http.MethodPost, endpoint, payload)
	if err != nil {
		return nil, err
	}
	return &fc.Transaction, nil
}

func (y *yahooClient) writeXML(ctx context.Context, method, endpoint string, payload any) (*yahoo.FantasyContent, error) {
	body, err := xml.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}
	var fc yahoo.FantasyContent
	switch method {
	case http.MethodPost:
		_, err = y.baseClient.requestor.Post(ctx, endpoint, bytes.NewReader(body), &fc, xmlDecorator, &xmlDecoder{})
	case http.MethodPut:
		_, err = y.baseClient.requestor.Put(ctx, endpoint, bytes.NewReader(body), &fc, xmlDecorator, &xmlDecoder{})
	default:
		return nil, fmt.Errorf("unsupported write method %s", method)
	}
	if err != nil {
		return nil, err
	}
	return &fc, nil
}

func toRosterPlayers(assignments []PlayerSlot) []rosterPlayer {
	players := make([]rosterPlayer, 0, len(assignments))
	for _, a := range assignments {
		players = append(players, rosterPlayer{
			PlayerKey: a.PlayerKey,
			Position:  a.Position,
		})
	}
	return players
}
