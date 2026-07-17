package gofantasy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofantasy/model/sleeper"
)

// ISleeperClient reads Sleeper Fantasy leagues via the public v1 API.
type ISleeperClient interface {
	GetUser(ctx context.Context, username string) (*sleeper.User, error)
	GetUserByID(ctx context.Context, userID string) (*sleeper.User, error)
	GetUserLeagues(ctx context.Context, userID, sport string, season int) ([]sleeper.League, error)
	GetLeague(ctx context.Context, leagueID string) (*sleeper.League, error)
	GetRosters(ctx context.Context, leagueID string) ([]sleeper.Roster, error)
	GetLeagueUsers(ctx context.Context, leagueID string) ([]sleeper.LeagueUser, error)
	GetMatchups(ctx context.Context, leagueID string, week int) ([]sleeper.Matchup, error)
	GetTransactions(ctx context.Context, leagueID string, week int) ([]sleeper.Transaction, error)
	GetState(ctx context.Context, sport string) (*sleeper.State, error)
	GetPlayer(ctx context.Context, sport, playerID string) (*sleeper.Player, error)
	AllPlayers(ctx context.Context, sport string) (map[string]*sleeper.Player, error)
	GetTrendingPlayers(ctx context.Context, sport, trendType string, lookbackHours, limit int) ([]sleeper.Player, error)
	GetProjections(ctx context.Context, sport string, season, week int) (map[string]sleeper.PlayerProjection, error)
	ResolveRoster(ctx context.Context, sport, leagueID string, rosterID int) (*sleeper.ResolvedRoster, error)
	DiscoverManagedTeams(ctx context.Context, userID, sport string, season int) ([]sleeper.ManagedTeam, error)
}

type sleeperClient struct {
	baseURL      string
	baseClient   *client
	playersCache *sleeperPlayersCache
	rateMu       sync.Mutex
	lastReq      time.Time
	playersRaw   map[string]map[string]json.RawMessage
	playersMu    sync.RWMutex
}

var _ ISleeperClient = &sleeperClient{}

const sleeperMinRequestInterval = 700 * time.Millisecond

func (s *sleeperClient) throttle() {
	s.rateMu.Lock()
	defer s.rateMu.Unlock()
	if wait := time.Since(s.lastReq); wait < sleeperMinRequestInterval {
		time.Sleep(sleeperMinRequestInterval - wait)
	}
	s.lastReq = time.Now()
}

func (s *sleeperClient) getJSON(ctx context.Context, path string, into any) error {
	s.throttle()
	endpoint := s.baseURL + path
	cacheKey := md5Hash(endpoint)
	if s.baseClient.cache != nil {
		if v, ok := s.baseClient.cache.Get(ctx, cacheKey); ok {
			if cached, ok := v.(json.RawMessage); ok {
				return json.Unmarshal(cached, into)
			}
		}
	}
	_, err := s.baseClient.requestor.Get(ctx, endpoint, into, jsonDecorator, &jsonDecoder{})
	if err != nil {
		return err
	}
	if s.baseClient.cache != nil {
		if b, mErr := json.Marshal(into); mErr == nil {
			s.baseClient.cache.Set(ctx, cacheKey, json.RawMessage(b))
		}
	}
	return nil
}

func resolveSleeperSport(sport string) (string, error) {
	sport = strings.ToLower(strings.TrimSpace(sport))
	switch sport {
	case "", "nfl", "ffl":
		return "nfl", nil
	case "nba", "fba":
		return "nba", nil
	case "mlb", "flb":
		return "mlb", nil
	case "nhl":
		return "nhl", nil
	case "lcs":
		return "lcs", nil
	default:
		return sport, nil
	}
}

func (s *sleeperClient) GetUser(ctx context.Context, username string) (*sleeper.User, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, fmt.Errorf("username is required")
	}
	var u sleeper.User
	if err := s.getJSON(ctx, "/user/"+username, &u); err != nil {
		return nil, err
	}
	if u.UserID == "" {
		return nil, fmt.Errorf("no Sleeper user found for username %q", username)
	}
	return &u, nil
}

func (s *sleeperClient) GetUserByID(ctx context.Context, userID string) (*sleeper.User, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	var u sleeper.User
	if err := s.getJSON(ctx, "/user/"+userID, &u); err != nil {
		return nil, err
	}
	if u.UserID == "" {
		return nil, fmt.Errorf("no Sleeper user found for id %q", userID)
	}
	return &u, nil
}

func (s *sleeperClient) GetUserLeagues(ctx context.Context, userID, sport string, season int) ([]sleeper.League, error) {
	sport, err := resolveSleeperSport(sport)
	if err != nil {
		return nil, err
	}
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	if season <= 0 {
		season = time.Now().Year()
	}
	var leagues []sleeper.League
	path := fmt.Sprintf("/user/%s/leagues/%s/%d", userID, sport, season)
	if err := s.getJSON(ctx, path, &leagues); err != nil {
		return nil, err
	}
	return leagues, nil
}

func (s *sleeperClient) GetLeague(ctx context.Context, leagueID string) (*sleeper.League, error) {
	leagueID = strings.TrimSpace(leagueID)
	if leagueID == "" {
		return nil, fmt.Errorf("league_id is required")
	}
	var lg sleeper.League
	if err := s.getJSON(ctx, "/league/"+leagueID, &lg); err != nil {
		return nil, err
	}
	if lg.LeagueID == "" {
		lg.LeagueID = leagueID
	}
	return &lg, nil
}

func (s *sleeperClient) GetRosters(ctx context.Context, leagueID string) ([]sleeper.Roster, error) {
	leagueID = strings.TrimSpace(leagueID)
	if leagueID == "" {
		return nil, fmt.Errorf("league_id is required")
	}
	var rosters []sleeper.Roster
	if err := s.getJSON(ctx, "/league/"+leagueID+"/rosters", &rosters); err != nil {
		return nil, err
	}
	return rosters, nil
}

func (s *sleeperClient) GetLeagueUsers(ctx context.Context, leagueID string) ([]sleeper.LeagueUser, error) {
	leagueID = strings.TrimSpace(leagueID)
	if leagueID == "" {
		return nil, fmt.Errorf("league_id is required")
	}
	var users []sleeper.LeagueUser
	if err := s.getJSON(ctx, "/league/"+leagueID+"/users", &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (s *sleeperClient) GetMatchups(ctx context.Context, leagueID string, week int) ([]sleeper.Matchup, error) {
	leagueID = strings.TrimSpace(leagueID)
	if leagueID == "" {
		return nil, fmt.Errorf("league_id is required")
	}
	if week <= 0 {
		return nil, fmt.Errorf("week is required")
	}
	var matchups []sleeper.Matchup
	if err := s.getJSON(ctx, fmt.Sprintf("/league/%s/matchups/%d", leagueID, week), &matchups); err != nil {
		return nil, err
	}
	return matchups, nil
}

func (s *sleeperClient) GetTransactions(ctx context.Context, leagueID string, week int) ([]sleeper.Transaction, error) {
	leagueID = strings.TrimSpace(leagueID)
	if leagueID == "" {
		return nil, fmt.Errorf("league_id is required")
	}
	if week <= 0 {
		return nil, fmt.Errorf("week is required")
	}
	var txns []sleeper.Transaction
	if err := s.getJSON(ctx, fmt.Sprintf("/league/%s/transactions/%d", leagueID, week), &txns); err != nil {
		return nil, err
	}
	return txns, nil
}

// GetProjections returns Sleeper's projected stat lines for one regular-season
// week, keyed by Sleeper player ID. This is the endpoint Sleeper's own app
// uses (v1/projections) — not formally documented, so the decoding tolerates
// unknown value shapes by keeping only numeric stats. Entries that carry no
// numeric stats at all are dropped.
func (s *sleeperClient) GetProjections(ctx context.Context, sport string, season, week int) (map[string]sleeper.PlayerProjection, error) {
	sport, err := resolveSleeperSport(sport)
	if err != nil {
		return nil, err
	}
	if season <= 0 || week <= 0 {
		return nil, fmt.Errorf("season and week must be positive (got season=%d week=%d)", season, week)
	}
	var raw map[string]map[string]any
	path := fmt.Sprintf("/projections/%s/regular/%d/%d", sport, season, week)
	if err := s.getJSON(ctx, path, &raw); err != nil {
		return nil, err
	}
	out := make(map[string]sleeper.PlayerProjection, len(raw))
	for id, stats := range raw {
		proj := make(sleeper.PlayerProjection, len(stats))
		for k, v := range stats {
			if f, ok := v.(float64); ok {
				proj[k] = f
			}
		}
		if len(proj) > 0 {
			out[id] = proj
		}
	}
	return out, nil
}

func (s *sleeperClient) GetState(ctx context.Context, sport string) (*sleeper.State, error) {
	sport, err := resolveSleeperSport(sport)
	if err != nil {
		return nil, err
	}
	var st sleeper.State
	if err := s.getJSON(ctx, "/state/"+sport, &st); err != nil {
		return nil, err
	}
	return &st, nil
}

func (s *sleeperClient) loadPlayersRaw(ctx context.Context, sport string) (map[string]json.RawMessage, error) {
	sport, err := resolveSleeperSport(sport)
	if err != nil {
		return nil, err
	}
	s.playersMu.RLock()
	if s.playersRaw != nil {
		if m, ok := s.playersRaw[sport]; ok && len(m) > 0 {
			s.playersMu.RUnlock()
			return m, nil
		}
	}
	s.playersMu.RUnlock()

	fetch := func(ctx context.Context, sport string) (map[string]json.RawMessage, error) {
		var players map[string]json.RawMessage
		if err := s.getJSON(ctx, "/players/"+sport, &players); err != nil {
			return nil, err
		}
		return players, nil
	}
	players, err := s.playersCache.load(ctx, sport, fetch)
	if err != nil {
		return nil, err
	}
	s.playersMu.Lock()
	if s.playersRaw == nil {
		s.playersRaw = make(map[string]map[string]json.RawMessage)
	}
	s.playersRaw[sport] = players
	s.playersMu.Unlock()
	return players, nil
}

// AllPlayers decodes the full player pool for a sport, keyed by Sleeper
// player ID. The pool is large (thousands of entries, 24h disk cache
// underneath) — callers should derive and cache whatever index they need
// rather than calling this per lookup.
func (s *sleeperClient) AllPlayers(ctx context.Context, sport string) (map[string]*sleeper.Player, error) {
	raw, err := s.loadPlayersRaw(ctx, sport)
	if err != nil {
		return nil, err
	}
	out := make(map[string]*sleeper.Player, len(raw))
	for id, b := range raw {
		var p sleeper.Player
		if err := json.Unmarshal(b, &p); err != nil {
			continue // tolerate malformed pool entries
		}
		if p.PlayerID == "" {
			p.PlayerID = id
		}
		out[id] = &p
	}
	return out, nil
}

func (s *sleeperClient) GetPlayer(ctx context.Context, sport, playerID string) (*sleeper.Player, error) {
	players, err := s.loadPlayersRaw(ctx, sport)
	if err != nil {
		return nil, err
	}
	raw, ok := players[playerID]
	if !ok {
		return nil, fmt.Errorf("player %s not found", playerID)
	}
	var p sleeper.Player
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, err
	}
	if p.PlayerID == "" {
		p.PlayerID = playerID
	}
	if p.FullName == "" {
		p.FullName = strings.TrimSpace(p.FirstName + " " + p.LastName)
	}
	return &p, nil
}

func (s *sleeperClient) GetTrendingPlayers(ctx context.Context, sport, trendType string, lookbackHours, limit int) ([]sleeper.Player, error) {
	sport, err := resolveSleeperSport(sport)
	if err != nil {
		return nil, err
	}
	if trendType == "" {
		trendType = "add"
	}
	if limit <= 0 {
		limit = 25
	}
	if lookbackHours <= 0 {
		lookbackHours = 24
	}
	var trending []struct {
		PlayerID string `json:"player_id"`
		Count    int    `json:"count"`
	}
	path := fmt.Sprintf("/players/%s/trending/%s?lookback_hours=%d&limit=%d", sport, trendType, lookbackHours, limit)
	if err := s.getJSON(ctx, path, &trending); err != nil {
		return nil, err
	}
	out := make([]sleeper.Player, 0, len(trending))
	for _, t := range trending {
		p, err := s.GetPlayer(ctx, sport, t.PlayerID)
		if err != nil {
			continue
		}
		out = append(out, *p)
	}
	return out, nil
}

func (s *sleeperClient) ResolveRoster(ctx context.Context, sport, leagueID string, rosterID int) (*sleeper.ResolvedRoster, error) {
	rosters, err := s.GetRosters(ctx, leagueID)
	if err != nil {
		return nil, err
	}
	var roster *sleeper.Roster
	for i := range rosters {
		if rosters[i].RosterID == rosterID {
			roster = &rosters[i]
			break
		}
	}
	if roster == nil {
		return nil, fmt.Errorf("roster %d not found in league %s", rosterID, leagueID)
	}
	starterSet := make(map[string]bool, len(roster.Starters))
	for _, id := range roster.Starters {
		starterSet[id] = true
	}
	reserveSet := make(map[string]bool, len(roster.Reserve))
	for _, id := range roster.Reserve {
		reserveSet[id] = true
	}

	resolved := &sleeper.ResolvedRoster{
		RosterID: roster.RosterID,
		OwnerID:  roster.OwnerID,
		LeagueID: roster.LeagueID,
		Settings: roster.Settings,
	}
	for _, id := range roster.Starters {
		p, _ := s.GetPlayer(ctx, sport, id)
		resolved.Starters = append(resolved.Starters, sleeper.ResolvedPlayer{PlayerID: id, Slot: "starter", Player: p})
	}
	for _, id := range roster.Players {
		if starterSet[id] || reserveSet[id] {
			continue
		}
		p, _ := s.GetPlayer(ctx, sport, id)
		resolved.Bench = append(resolved.Bench, sleeper.ResolvedPlayer{PlayerID: id, Slot: "bench", Player: p})
	}
	for _, id := range roster.Reserve {
		p, _ := s.GetPlayer(ctx, sport, id)
		resolved.Reserve = append(resolved.Reserve, sleeper.ResolvedPlayer{PlayerID: id, Slot: "reserve", Player: p})
	}
	return resolved, nil
}

func (s *sleeperClient) DiscoverManagedTeams(ctx context.Context, userID, sport string, season int) ([]sleeper.ManagedTeam, error) {
	sport, err := resolveSleeperSport(sport)
	if err != nil {
		return nil, err
	}
	leagues, err := s.GetUserLeagues(ctx, userID, sport, season)
	if err != nil {
		return nil, err
	}
	var out []sleeper.ManagedTeam
	for _, lg := range leagues {
		rosters, err := s.GetRosters(ctx, lg.LeagueID)
		if err != nil {
			return nil, fmt.Errorf("league %s: %w", lg.LeagueID, err)
		}
		users, err := s.GetLeagueUsers(ctx, lg.LeagueID)
		if err != nil {
			return nil, fmt.Errorf("league %s users: %w", lg.LeagueID, err)
		}
		userByID := make(map[string]sleeper.LeagueUser, len(users))
		for _, u := range users {
			userByID[u.UserID] = u
		}
		for _, r := range rosters {
			if r.OwnerID != userID {
				continue
			}
			seasonNum := season
			if n, err := strconv.Atoi(lg.Season); err == nil {
				seasonNum = n
			}
			out = append(out, sleeper.ManagedTeam{
				Sport:      sport,
				Season:     seasonNum,
				League:     lg,
				Roster:     r,
				LeagueUser: userByID[userID],
			})
		}
	}
	return out, nil
}

// Ensure jsonDecorator is used
var _ = func(req *http.Request) *http.Request { return jsonDecorator(req) }
