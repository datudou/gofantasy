package gofantasy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gofantasy/model/espn"
)

// IEspnClient reads (and eventually writes) ESPN Fantasy leagues via the v3 API.
type IEspnClient interface {
	// WithCookies configures SWID and espn_s2 session cookies required for
	// private leagues. Either cookie may be empty for public-league reads.
	WithCookies(swid, espnS2 string) IEspnClient
	Cookies() (swid, espnS2 string)

	// GetLeague fetches a league with optional views and query parameters.
	GetLeague(ctx context.Context, sport string, season, leagueID int, opts ...EspnLeagueOption) (*espn.League, error)

	GetTeams(ctx context.Context, sport string, season, leagueID int) ([]espn.Team, error)
	GetTeam(ctx context.Context, sport string, season, leagueID, teamID int) (*espn.Team, error)
	GetRoster(ctx context.Context, sport string, season, leagueID, teamID int, scoringPeriodID ...int) (*espn.Roster, error)
	GetStandings(ctx context.Context, sport string, season, leagueID int) ([]espn.Team, error)
	GetMatchups(ctx context.Context, sport string, season, leagueID int, matchupPeriod ...int) ([]espn.ScheduleEntry, error)
	GetTransactions(ctx context.Context, sport string, season, leagueID int) ([]espn.Transaction, error)
	GetFreeAgents(ctx context.Context, sport string, season, leagueID int, opts ...EspnPlayerOption) ([]espn.Player, error)

	// DiscoverManagedTeams returns teams owned by the authenticated SWID across
	// the given league IDs. ESPN has no documented "list all my leagues" API;
	// callers typically collect league IDs from the fantasy.espn.com URL.
	DiscoverManagedTeams(ctx context.Context, sport string, season int, leagueIDs ...int) ([]espn.ManagedTeam, error)
}

type espnClient struct {
	baseURL    string
	baseClient *client
	swid       string
	espnS2     string
}

var _ IEspnClient = &espnClient{}

func (e *espnClient) WithCookies(swid, espnS2 string) IEspnClient {
	clone := *e
	clone.swid = strings.TrimSpace(swid)
	clone.espnS2 = strings.TrimSpace(espnS2)
	return &clone
}

func (e *espnClient) Cookies() (swid, espnS2 string) {
	return e.swid, e.espnS2
}

func (e *espnClient) GetLeague(ctx context.Context, sport string, season, leagueID int, opts ...EspnLeagueOption) (*espn.League, error) {
	cfg := newEspnLeagueConfig(opts...)
	return e.fetchLeague(ctx, sport, season, leagueID, cfg)
}

func (e *espnClient) GetTeams(ctx context.Context, sport string, season, leagueID int) ([]espn.Team, error) {
	lg, err := e.GetLeague(ctx, sport, season, leagueID, WithEspnViews("mTeam"))
	if err != nil {
		return nil, err
	}
	if len(lg.Teams) == 0 {
		return nil, fmt.Errorf("no teams found for league %d season %d", leagueID, season)
	}
	return lg.Teams, nil
}

func (e *espnClient) GetTeam(ctx context.Context, sport string, season, leagueID, teamID int) (*espn.Team, error) {
	teams, err := e.GetTeams(ctx, sport, season, leagueID)
	if err != nil {
		return nil, err
	}
	for i := range teams {
		if teams[i].ID == teamID {
			return &teams[i], nil
		}
	}
	return nil, fmt.Errorf("team %d not found in league %d", teamID, leagueID)
}

func (e *espnClient) GetRoster(ctx context.Context, sport string, season, leagueID, teamID int, scoringPeriodID ...int) (*espn.Roster, error) {
	opts := []EspnLeagueOption{WithEspnViews("mRoster")}
	if len(scoringPeriodID) > 0 {
		opts = append(opts, WithScoringPeriodID(scoringPeriodID[0]))
	}
	lg, err := e.GetLeague(ctx, sport, season, leagueID, opts...)
	if err != nil {
		return nil, err
	}
	for i := range lg.Teams {
		if lg.Teams[i].ID == teamID {
			if lg.Teams[i].Roster == nil {
				return nil, fmt.Errorf("no roster returned for team %d (private league may need cookies)", teamID)
			}
			return lg.Teams[i].Roster, nil
		}
	}
	return nil, fmt.Errorf("team %d not found in league %d", teamID, leagueID)
}

func (e *espnClient) GetStandings(ctx context.Context, sport string, season, leagueID int) ([]espn.Team, error) {
	lg, err := e.GetLeague(ctx, sport, season, leagueID, WithEspnViews("mTeam", "mStandings"))
	if err != nil {
		return nil, err
	}
	if len(lg.Teams) == 0 {
		return nil, fmt.Errorf("no standings found for league %d", leagueID)
	}
	return lg.Teams, nil
}

func (e *espnClient) GetMatchups(ctx context.Context, sport string, season, leagueID int, matchupPeriod ...int) ([]espn.ScheduleEntry, error) {
	opts := []EspnLeagueOption{WithEspnViews("mMatchup", "mScoreboard")}
	if len(matchupPeriod) > 0 {
		opts = append(opts, WithMatchupPeriodID(matchupPeriod[0]))
	}
	lg, err := e.GetLeague(ctx, sport, season, leagueID, opts...)
	if err != nil {
		return nil, err
	}
	if len(lg.Schedule) == 0 {
		return nil, fmt.Errorf("no matchups found for league %d", leagueID)
	}
	return lg.Schedule, nil
}

func (e *espnClient) GetTransactions(ctx context.Context, sport string, season, leagueID int) ([]espn.Transaction, error) {
	lg, err := e.GetLeague(ctx, sport, season, leagueID, WithEspnViews("mTransactions2"))
	if err != nil {
		return nil, err
	}
	return lg.Transactions, nil
}

func (e *espnClient) GetFreeAgents(ctx context.Context, sport string, season, leagueID int, opts ...EspnPlayerOption) ([]espn.Player, error) {
	cfg := newEspnPlayerConfig(opts...)
	filterJSON, err := cfg.xFantasyFilterJSON()
	if err != nil {
		return nil, err
	}
	var filter any
	if err := json.Unmarshal([]byte(filterJSON), &filter); err != nil {
		return nil, err
	}

	lg, err := e.GetLeague(ctx, sport, season, leagueID,
		WithEspnViews("kona_player_info"),
		WithXFantasyFilter(filter),
	)
	if err != nil {
		return nil, err
	}

	players := make([]espn.Player, 0, len(lg.Players))
	for _, entry := range lg.Players {
		players = append(players, entry.Player)
	}
	if len(players) == 0 {
		for _, team := range lg.Teams {
			if team.Roster == nil {
				continue
			}
			for _, entry := range team.Roster.Entries {
				players = append(players, entry.PlayerPoolEntry.Player)
			}
		}
	}
	return players, nil
}

func (e *espnClient) DiscoverManagedTeams(ctx context.Context, sport string, season int, leagueIDs ...int) ([]espn.ManagedTeam, error) {
	if e.swid == "" {
		return nil, fmt.Errorf("SWID cookie is required to discover managed teams")
	}
	if len(leagueIDs) == 0 {
		return nil, fmt.Errorf("at least one league ID is required")
	}
	code, err := resolveEspnSport(sport)
	if err != nil {
		return nil, err
	}

	var out []espn.ManagedTeam
	for _, leagueID := range leagueIDs {
		teams, err := e.GetTeams(ctx, sport, season, leagueID)
		if err != nil {
			return nil, fmt.Errorf("league %d: %w", leagueID, err)
		}
		for _, team := range teams {
			if teamOwnedBySWID(team, e.swid) {
				out = append(out, espn.ManagedTeam{
					Sport:    code,
					Season:   season,
					LeagueID: leagueID,
					Team:     team,
				})
			}
		}
	}
	return out, nil
}

func (e *espnClient) fetchLeague(ctx context.Context, sport string, season, leagueID int, cfg espnLeagueConfig) (*espn.League, error) {
	code, err := resolveEspnSport(sport)
	if err != nil {
		return nil, err
	}
	endpoint, history, err := espnLeagueEndpoint(e.baseURL, code, season, leagueID, cfg)
	if err != nil {
		return nil, err
	}

	var lg espn.League
	if history {
		var leagues []espn.League
		if err := e.getJSON(ctx, endpoint, cfg, &leagues); err != nil {
			return nil, err
		}
		if len(leagues) == 0 {
			return nil, fmt.Errorf("no league data for league %d season %d", leagueID, season)
		}
		lg = leagues[0]
	} else {
		if err := e.getJSON(ctx, endpoint, cfg, &lg); err != nil {
			return nil, err
		}
	}
	if lg.ID == 0 {
		lg.ID = leagueID
	}
	return &lg, nil
}

func espnLeagueEndpoint(baseURL, sport string, season, leagueID int, cfg espnLeagueConfig) (string, bool, error) {
	var raw string
	history := season < 2018
	if history {
		raw = fmt.Sprintf("%s/%s/leagueHistory/%d", baseURL, sport, leagueID)
	} else {
		raw = fmt.Sprintf("%s/%s/seasons/%d/segments/0/leagues/%d", baseURL, sport, season, leagueID)
	}

	u, err := url.Parse(raw)
	if err != nil {
		return "", false, err
	}
	q := u.Query()
	if history {
		q.Set("seasonId", fmt.Sprintf("%d", season))
	}
	for _, view := range cfg.views {
		q.Add("view", view)
	}
	if cfg.scoringPeriodID != nil {
		q.Set("scoringPeriodId", fmt.Sprintf("%d", *cfg.scoringPeriodID))
	}
	if cfg.matchupPeriodID != nil {
		q.Set("matchupPeriodId", fmt.Sprintf("%d", *cfg.matchupPeriodID))
	}
	u.RawQuery = q.Encode()
	return u.String(), history, nil
}

func (e *espnClient) getJSON(ctx context.Context, endpoint string, cfg espnLeagueConfig, into any) error {
	cacheKey := md5Hash(endpoint)
	if e.baseClient.cache != nil {
		if v, ok := e.baseClient.cache.Get(ctx, cacheKey); ok {
			if cached, ok := v.(json.RawMessage); ok {
				return json.Unmarshal(cached, into)
			}
		}
	}

	reqDecorator := espnRequestDecorator(e.swid, e.espnS2, cfg.xFantasyFilter)
	_, err := e.baseClient.requestor.Get(ctx, endpoint, into, reqDecorator, &jsonDecoder{})
	if err != nil {
		return err
	}

	if e.baseClient.cache != nil {
		// Re-fetch through marshal for cache storage (Get already decoded into `into`).
		if b, mErr := json.Marshal(into); mErr == nil {
			e.baseClient.cache.Set(ctx, cacheKey, json.RawMessage(b))
		}
	}
	return nil
}

func espnRequestDecorator(swid, espnS2 string, xFantasyFilter any) requestDecorator {
	return func(req *http.Request) *http.Request {
		req = jsonDecorator(req)
		if swid != "" || espnS2 != "" {
			req.Header.Set("Cookie", fmt.Sprintf("SWID=%s; espn_s2=%s", swid, espnS2))
		}
		if xFantasyFilter != nil {
			if b, err := json.Marshal(xFantasyFilter); err == nil {
				req.Header.Set("X-Fantasy-Filter", string(b))
			}
		}
		return req
	}
}

func teamOwnedBySWID(team espn.Team, swid string) bool {
	norm := normalizeSWID(swid)
	if norm == "" {
		return false
	}
	for _, owner := range team.Owners {
		if normalizeSWID(owner) == norm {
			return true
		}
	}
	return false
}

func normalizeSWID(s string) string {
	s = strings.TrimSpace(strings.ToUpper(s))
	s = strings.TrimPrefix(s, "{")
	s = strings.TrimSuffix(s, "}")
	return s
}
