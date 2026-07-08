package gofantasy

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gofantasy/model/yahoo"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"strings"
)

type IYahooClient interface {
	GetLeague(ctx context.Context, leagueKey string) (*yahoo.League, error)
	GetLeagueStandings(ctx context.Context, leagueKey string) ([]yahoo.Team, error)
	GetLeagueScoreboard(ctx context.Context, leagueKey string, weeks ...int) (*yahoo.Scoreboard, error)
	GetLeagueTeams(ctx context.Context, leagueKey string) ([]yahoo.Team, error)
	GetLeagueTransactions(ctx context.Context, leagueKey string, opts ...TransactionOption) ([]yahoo.Transaction, error)
	GetLeagueDraftResults(ctx context.Context, leagueKey string) ([]yahoo.DraftResult, error)
	GetLeaguePlayers(ctx context.Context, leagueKey string, opts ...PlayerOption) ([]yahoo.Player, error)

	GetGameBySeason(ctx context.Context, gameCode string, season string) ([]yahoo.Game, error)
	GetUserAttendGames(ctx context.Context, gameKey ...string) ([]yahoo.Game, error)
	GetUserManagedTeams(ctx context.Context, gameKey ...string) ([]*yahoo.Team, error)
	// GetAllUserGames returns every fantasy game (one entry per sport per
	// season) the logged-in user has ever participated in, across all
	// years — not just the "current" game for each sport. Use this to
	// discover past seasons' game keys (Game.GameKey), then pass them to
	// GetUserManagedTeamsForGames to fetch that season's teams.
	GetAllUserGames(ctx context.Context) ([]yahoo.Game, error)
	// GetUserManagedTeamsForGames returns the given games with their Teams
	// populated. Unlike GetUserManagedTeams, gameKeys here may be literal
	// current-game aliases ("nfl") or specific numeric season game keys
	// (e.g. "423" for NFL 2023) as returned by GetAllUserGames — there is
	// no isValidGameKeys restriction.
	GetUserManagedTeamsForGames(ctx context.Context, gameKeys ...string) ([]yahoo.Game, error)

	GetUserRoster(ctx context.Context, teamKey string) (*yahoo.Roster, error)
	GetTeam(ctx context.Context, teamKey string) (*yahoo.Team, error)
	GetTeamStandings(ctx context.Context, teamKey string) (*yahoo.TeamStandings, error)
	GetTeamMatchups(ctx context.Context, teamKey string, weeks ...int) ([]yahoo.Matchup, error)
	GetTeamStats(ctx context.Context, teamKey string, opts ...StatOption) (*yahoo.TeamStats, error)

	GetPlayer(ctx context.Context, playerKey string) (*yahoo.Player, error)
	GetPlayerStats(ctx context.Context, playerKey string, opts ...StatOption) (*yahoo.PlayerStats, error)

	// Write operations. These mutate league state and should be gated behind
	// explicit user approval by callers.
	SetRoster(ctx context.Context, teamKey string, date string, assignments []PlayerSlot) error
	SetRosterForWeek(ctx context.Context, teamKey string, week int, assignments []PlayerSlot) error
	AddPlayer(ctx context.Context, leagueKey, teamKey, playerKey string, faabBid *int) (*yahoo.Transaction, error)
	DropPlayer(ctx context.Context, leagueKey, teamKey, playerKey string) (*yahoo.Transaction, error)
	AddDropPlayer(ctx context.Context, leagueKey, teamKey, addKey, dropKey string, faabBid *int) (*yahoo.Transaction, error)
	ProposeTrade(ctx context.Context, leagueKey, traderTeamKey, tradeeTeamKey string, send, receive []string, note string) (*yahoo.Transaction, error)
	RespondToTrade(ctx context.Context, transactionKey, action string, voteAgainst *bool) error
	CancelTransaction(ctx context.Context, transactionKey string) error

	OAuth2(clientID, clientSecret, redirectURL string) IYahooOAuth2
	LoadAccessToken(path string) (IYahooClient, error)
	// WithAccessToken configures the client to use an already-obtained OAuth2
	// token (e.g. persisted per-user in a database) instead of a token file.
	// If the token is expired it is refreshed using YAHOO_CLIENT_ID /
	// YAHOO_REDIRECT_URL from the environment.
	WithAccessToken(token *oauth2.Token) (IYahooClient, error)
	// Token returns the token currently in use (reflecting any refresh
	// performed by WithAccessToken/LoadAccessToken), so callers can persist
	// it again after a refresh.
	Token() *oauth2.Token
}

var GameKeys = map[string]string{
	"nfl": "nfl",
	"nba": "nba",
	"mlb": "mlb",
	"nhl": "nhl",
}

type yahooClient struct {
	baseUrl     string
	baseClient  *client
	yahooOAuth2 *yahooOAuth2
}

var _ IYahooClient = &yahooClient{}

// OAuth2
//
//	@Description: returns an instance of yahooOAuth2
//	@receiver y
//	@param clientID
//	@param clientSecret
//	@param redirectURL
//	@return IYahooOAuth2
func (y *yahooClient) OAuth2(clientID, clientSecret, redirectURL string) IYahooOAuth2 {
	return y.yahooOAuth2.OAuth2(clientID, clientSecret, redirectURL)
}

// GetGameBySeason returns the Yahoo games that match the given gameCode/season.
// gameCode is e.g. "mlb", "nfl", "nba", "nhl"; season is the four-digit year.
func (y *yahooClient) GetGameBySeason(ctx context.Context, gameCode string, season string) ([]yahoo.Game, error) {
	endpoint := fmt.Sprintf("%s/games;game_codes=%s;seasons=%s", y.baseUrl, gameCode, season)
	fc, err := y.get(ctx, endpoint, "")
	if err != nil {
		return nil, err
	}
	if len(fc.Games) == 0 {
		return nil, fmt.Errorf("no games found for gameCode %s and season %s", gameCode, season)
	}
	return fc.Games, nil
}

// GetUserRoster returns the current roster (list of players) for the given team key.
func (y *yahooClient) GetUserRoster(ctx context.Context, teamKey string) (*yahoo.Roster, error) {
	endpoint := fmt.Sprintf("%s/team/%s/roster/players", y.baseUrl, teamKey)
	fc, err := y.get(ctx, endpoint, "")
	if err != nil {
		return nil, err
	}
	if fc.Team.TeamKey == "" {
		return nil, fmt.Errorf("no roster found for teamKey %s", teamKey)
	}
	return &fc.Team.Roster, nil
}

// GetUserAttendGames
//
//	@Description: get user attend games
//	@receiver y
//	@param ctx
//	@param gameKeys
//	@return []yahoo.Game
//	@return error
func (y *yahooClient) GetUserAttendGames(ctx context.Context, gameKeys ...string) ([]yahoo.Game, error) {
	if !isValidGameKeys(gameKeys...) {
		return nil, fmt.Errorf("invalid gameCodes %v", gameKeys)
	}
	gcs := strings.Join(gameKeys, ",")

	endpoint := fmt.Sprintf("%s/users;use_login=1/games;games_key=%s", y.baseUrl, gcs)
	fc, err := y.get(ctx, endpoint, "")
	if err != nil {
		return nil, err
	}
	if len(fc.Users) <= 0 {
		return nil, fmt.Errorf("no games found for gameCodes %v", gameKeys)
	}

	return fc.Users[0].Games, nil
}

// GetUserManagedTeams
//
//	@Description:
//	@receiver y
//	@param ctx
//	@param gameKeys
//	@return []*yahoo.Team
//	@return error
func (y *yahooClient) GetUserManagedTeams(ctx context.Context, gameKeys ...string) ([]*yahoo.Team, error) {
	if !isValidGameKeys(gameKeys...) {
		return nil, fmt.Errorf("invalid gameCodes %v", gameKeys)
	}
	gcs := strings.Join(gameKeys, ",")

	endpoint := fmt.Sprintf("%s/users;use_login=1/games;game_keys=%s/teams", y.baseUrl, gcs)
	fc, err := y.get(ctx, endpoint, "")
	if err != nil {
		return nil, err
	}

	if len(fc.Users) <= 0 {
		return nil, fmt.Errorf("no teams found for gameCodes %v", gameKeys)
	}
	if len(fc.Users[0].Games) <= 0 {
		return nil, fmt.Errorf("no games found for gameCodes %v", gameKeys)
	}
	games := fc.Users[0].Games
	var teams []*yahoo.Team
	for _, v := range games {
		teams = append(teams, v.Teams...)
	}

	return teams, nil
}

// GetAllUserGames returns every fantasy game the logged-in user has ever
// participated in (all sports, all seasons), unfiltered. See
// IYahooClient.GetAllUserGames.
func (y *yahooClient) GetAllUserGames(ctx context.Context) ([]yahoo.Game, error) {
	endpoint := fmt.Sprintf("%s/users;use_login=1/games", y.baseUrl)
	fc, err := y.get(ctx, endpoint, "")
	if err != nil {
		return nil, err
	}
	if len(fc.Users) == 0 {
		return nil, fmt.Errorf("no games found for the current user")
	}
	return fc.Users[0].Games, nil
}

// GetUserManagedTeamsForGames returns the given games with their Teams
// populated. See IYahooClient.GetUserManagedTeamsForGames.
func (y *yahooClient) GetUserManagedTeamsForGames(ctx context.Context, gameKeys ...string) ([]yahoo.Game, error) {
	if len(gameKeys) == 0 {
		return nil, fmt.Errorf("at least one game key is required")
	}
	gcs := strings.Join(gameKeys, ",")

	endpoint := fmt.Sprintf("%s/users;use_login=1/games;game_keys=%s/teams", y.baseUrl, gcs)
	fc, err := y.get(ctx, endpoint, "")
	if err != nil {
		return nil, err
	}
	if len(fc.Users) == 0 {
		return nil, fmt.Errorf("no games found for gameKeys %v", gameKeys)
	}
	return fc.Users[0].Games, nil
}

// LoadAccessToken
//
//	@Description:
//	@receiver y
//	@param path
//	@return IYahooClient
//	@return error
func (y *yahooClient) LoadAccessToken(path string) (IYahooClient, error) {
	if path == "" {
		path = os.Getenv("HOME") + YahooTokenPath
	}
	token, err := y.yahooOAuth2.LoadAccessToken(path)
	if err != nil {
		return nil, err
	}
	y.baseClient.requestor.AuthorizationDecorator = func(req *http.Request) *http.Request {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
		return req
	}
	return y, nil
}

// WithAccessToken configures the client to use an already-obtained OAuth2
// token instead of loading one from a file. See IYahooClient.WithAccessToken.
func (y *yahooClient) WithAccessToken(token *oauth2.Token) (IYahooClient, error) {
	if token == nil {
		return nil, fmt.Errorf("token is nil")
	}
	y.yahooOAuth2.SetToken(token)
	refreshed, err := y.yahooOAuth2.RefreshIfNeeded()
	if err != nil {
		return nil, err
	}
	y.baseClient.requestor.AuthorizationDecorator = func(req *http.Request) *http.Request {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", refreshed.AccessToken))
		return req
	}
	return y, nil
}

// Token returns the token currently in use.
func (y *yahooClient) Token() *oauth2.Token {
	return y.yahooOAuth2.Token()
}

// GetLeague returns the high-level metadata and settings for a league.
func (y *yahooClient) GetLeague(ctx context.Context, leagueKey string) (*yahoo.League, error) {
	endpoint := fmt.Sprintf("%s/league/%s", y.baseUrl, leagueKey)
	fc, err := y.get(ctx, endpoint, "league")
	if err != nil {
		return nil, err
	}
	if fc.League.LeagueKey == "" {
		return nil, fmt.Errorf("no league found for leagueKey %s", leagueKey)
	}
	return &fc.League, nil
}

// GetLeagueStandings returns the teams in a league ordered by their standings.
// Each returned Team carries a populated TeamStandings.
func (y *yahooClient) GetLeagueStandings(ctx context.Context, leagueKey string) ([]yahoo.Team, error) {
	endpoint := fmt.Sprintf("%s/league/%s/standings", y.baseUrl, leagueKey)
	fc, err := y.get(ctx, endpoint, "")
	if err != nil {
		return nil, err
	}
	if len(fc.League.Standings) == 0 {
		return nil, fmt.Errorf("no standings found for leagueKey %s", leagueKey)
	}
	return fc.League.Standings, nil
}

// GetLeagueScoreboard returns the matchups (and box scores) for the given week(s).
// When no weeks are supplied Yahoo returns the current week.
func (y *yahooClient) GetLeagueScoreboard(ctx context.Context, leagueKey string, weeks ...int) (*yahoo.Scoreboard, error) {
	endpoint := fmt.Sprintf("%s/league/%s/scoreboard", y.baseUrl, leagueKey)
	if w := formatWeeks(weeks); w != "" {
		endpoint += ";week=" + w
	}
	fc, err := y.get(ctx, endpoint, "")
	if err != nil {
		return nil, err
	}
	return &fc.League.Scoreboard, nil
}

// GetLeagueTeams returns every team in a league.
func (y *yahooClient) GetLeagueTeams(ctx context.Context, leagueKey string) ([]yahoo.Team, error) {
	endpoint := fmt.Sprintf("%s/league/%s/teams", y.baseUrl, leagueKey)
	fc, err := y.get(ctx, endpoint, "")
	if err != nil {
		return nil, err
	}
	if len(fc.League.Teams) == 0 {
		return nil, fmt.Errorf("no teams found for leagueKey %s", leagueKey)
	}
	return fc.League.Teams, nil
}

// GetLeagueTransactions returns transactions (adds/drops/trades) for a league.
// Use TransactionOption to filter by type, team or count.
func (y *yahooClient) GetLeagueTransactions(ctx context.Context, leagueKey string, opts ...TransactionOption) ([]yahoo.Transaction, error) {
	filter := newTransactionFilter(opts...)
	endpoint := fmt.Sprintf("%s/league/%s/transactions%s", y.baseUrl, leagueKey, filter.encode())
	fc, err := y.get(ctx, endpoint, "")
	if err != nil {
		return nil, err
	}
	return fc.League.Transactions, nil
}

// GetLeagueDraftResults returns the draft picks for a league.
func (y *yahooClient) GetLeagueDraftResults(ctx context.Context, leagueKey string) ([]yahoo.DraftResult, error) {
	endpoint := fmt.Sprintf("%s/league/%s/draftresults", y.baseUrl, leagueKey)
	fc, err := y.get(ctx, endpoint, "")
	if err != nil {
		return nil, err
	}
	return fc.League.DraftResults, nil
}

// GetLeaguePlayers returns players in the league filtered by the supplied PlayerOptions.
// By default Yahoo returns 25 players starting at offset 0; use WithPlayerPagination to paginate.
func (y *yahooClient) GetLeaguePlayers(ctx context.Context, leagueKey string, opts ...PlayerOption) ([]yahoo.Player, error) {
	filter := newPlayerFilter(opts...)
	endpoint := fmt.Sprintf("%s/league/%s/players%s", y.baseUrl, leagueKey, filter.encode())
	fc, err := y.get(ctx, endpoint, "")
	if err != nil {
		return nil, err
	}
	return fc.League.Players, nil
}

// GetTeam returns the metadata for a single team.
func (y *yahooClient) GetTeam(ctx context.Context, teamKey string) (*yahoo.Team, error) {
	endpoint := fmt.Sprintf("%s/team/%s", y.baseUrl, teamKey)
	fc, err := y.get(ctx, endpoint, "")
	if err != nil {
		return nil, err
	}
	if fc.Team.TeamKey == "" {
		return nil, fmt.Errorf("no team found for teamKey %s", teamKey)
	}
	return &fc.Team, nil
}

// GetTeamStandings returns the standings detail (rank, W/L/T, points) for a single team.
func (y *yahooClient) GetTeamStandings(ctx context.Context, teamKey string) (*yahoo.TeamStandings, error) {
	endpoint := fmt.Sprintf("%s/team/%s/standings", y.baseUrl, teamKey)
	fc, err := y.get(ctx, endpoint, "")
	if err != nil {
		return nil, err
	}
	return &fc.Team.TeamStandings, nil
}

// GetTeamMatchups returns the matchups a team has played; pass week numbers to narrow it down.
func (y *yahooClient) GetTeamMatchups(ctx context.Context, teamKey string, weeks ...int) ([]yahoo.Matchup, error) {
	endpoint := fmt.Sprintf("%s/team/%s/matchups", y.baseUrl, teamKey)
	if w := formatWeeks(weeks); w != "" {
		endpoint += ";weeks=" + w
	}
	fc, err := y.get(ctx, endpoint, "")
	if err != nil {
		return nil, err
	}
	return fc.Team.Matchups, nil
}

// GetTeamStats returns the stats for a team. Defaults to season-to-date when no StatOption is supplied.
func (y *yahooClient) GetTeamStats(ctx context.Context, teamKey string, opts ...StatOption) (*yahoo.TeamStats, error) {
	filter := newStatFilter(opts...)
	endpoint := fmt.Sprintf("%s/team/%s/stats%s", y.baseUrl, teamKey, filter.encode())
	fc, err := y.get(ctx, endpoint, "")
	if err != nil {
		return nil, err
	}
	return &fc.Team.TeamStats, nil
}

// GetPlayer returns the metadata for a single player.
func (y *yahooClient) GetPlayer(ctx context.Context, playerKey string) (*yahoo.Player, error) {
	endpoint := fmt.Sprintf("%s/player/%s", y.baseUrl, playerKey)
	fc, err := y.get(ctx, endpoint, "")
	if err != nil {
		return nil, err
	}
	if fc.Player.PlayerKey == "" {
		return nil, fmt.Errorf("no player found for playerKey %s", playerKey)
	}
	return &fc.Player, nil
}

// GetPlayerStats returns the stats for a player. Defaults to season-to-date when no StatOption is supplied.
func (y *yahooClient) GetPlayerStats(ctx context.Context, playerKey string, opts ...StatOption) (*yahoo.PlayerStats, error) {
	filter := newStatFilter(opts...)
	endpoint := fmt.Sprintf("%s/player/%s/stats%s", y.baseUrl, playerKey, filter.encode())
	fc, err := y.get(ctx, endpoint, "")
	if err != nil {
		return nil, err
	}
	return &fc.Player.PlayerStats, nil
}

func (y *yahooClient) get(ctx context.Context, endpoint string, objType string) (*yahoo.FantasyContent, error) {
	var fc yahoo.FantasyContent
	cacheKey := md5Hash(endpoint)
	if y.baseClient.cache != nil {
		if v, ok := y.baseClient.cache.Get(ctx, cacheKey); ok {
			if cached, ok := v.(*yahoo.FantasyContent); ok {
				return cached, nil
			}
		}
	}
	_, err := y.baseClient.requestor.Get(ctx, endpoint, &fc, xmlDecorator, &xmlDecoder{})
	if err != nil {
		return nil, err
	}

	if y.baseClient.cache != nil {
		y.baseClient.cache.Set(ctx, cacheKey, &fc)
	}
	return &fc, nil
}

func isValidGameKeys(gameKeys ...string) bool {
	for _, v := range gameKeys {
		if _, ok := GameKeys[v]; !ok {
			return false
		}
	}
	return true
}

func md5Hash(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
