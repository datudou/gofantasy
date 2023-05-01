package gofantasy

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofantasy/model/yahoo"
)

type IYahooClient interface {
	// GetLeague(ctx context.Context, leagueID string) (*yahoo.League, error)
	// GetGameBySeason(ctx context.Context, gameCode string, season string) (*[]yahoo.Game, error)
	// GetUserAttendGames(ctx context.Context, gameKey ...string) ([]yahoo.Game, error)
	// GetUserManagedTeams(ctx context.Context, gameKey ...string) ([]*yahoo.Team, error)
	// GetUserRoster(ctx context.Context, teamKey string) (*yahoo.Roster, error)
	// OAuth2(clientID, clientSecret, redirectURL string) IYahooOAuth2
	// LoadAccessToken(path string) (IYahooClient, error)
}

type yahooClient struct {
	baseUrl     string
	yahooOAuth2 *yahooOAuth2
	Client
}

var defaultHTTPClient = &http.Client{
	Timeout: time.Second * 10,
	Transport: &http.Transport{
		TLSHandshakeTimeout: 10 * time.Second,
	},
}

func NewYahooClient(opts ...ClientOption) *yahooClient {

	r := &requestor{
		httpClient: defaultHTTPClient,
	}
	c := &yahooClient{
		yahooOAuth2: &yahooOAuth2{},
		baseUrl:     YahooBaseURL,
	}
	c.Requestor = r
	return c
}

func (c *yahooClient) WithOptions(opts ...ClientOption) *yahooClient {
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (y *yahooClient) SetCache(c ICache) {
	y.Cache = c
}

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

// GetGameBySeason
//
//	@Description:
//	@receiver y
//	@param ctx
//	@param gameCode
//	@param season
//	@return *[]yahoo.Game
//	@return error
// func (y *yahooClient) GetGameBySeason(ctx context.Context, gameCode string, season string) (*[]yahoo.Game, error) {

// 	endpoint := fmt.Sprintf("%s/games;game_codes=%s;seasons=%s", y.baseUrl, gameCode, season)
// 	fc, err := y.get(ctx, endpoint, "")
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(fc.Games) == 0 {
// 		return nil, fmt.Errorf("no games found for gameCode %s and season %s", gameCode, season)

// 	}
// 	return &fc.Games, nil
// }

// GetUserRoster
//
//	@Description:
//	@receiver y
//	@param ctx
//	@param teamKey
//	@return *yahoo.Roster
//	@return error
// func (y *yahooClient) GetUserRoster(ctx context.Context, teamKey string) (*yahoo.Roster, error) {

// 	endpoint := fmt.Sprintf("%s/team/%s/roster/players", y.baseUrl, teamKey)
// 	fc, err := y.get(ctx, endpoint, "")
// 	if err != nil {
// 		return nil, err
// 	}
// 	if &fc.Team == nil {
// 		return nil, fmt.Errorf("no roster found for teamKey %s", teamKey)
// 	}
// 	return &fc.Team.Roster, nil
// }

// GetUserAttendGames
//
//	@Description: get user attend games
//	@receiver y
//	@param ctx
//	@param gameKeys
//	@return []yahoo.Game
//	@return error
func (y *yahooClient) GetUserAttendGames(ctx context.Context, gameKeys ...string) ([]yahoo.Game, error) {
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
// func (y *yahooClient) GetUserManagedTeams(ctx context.Context, gameKeys ...string) ([]*yahoo.Team, error) {
// 	if !isValidGameKeys(gameKeys...) {
// 		return nil, fmt.Errorf("invalid gameCodes %v", gameKeys)
// 	}
// 	gcs := strings.Join(gameKeys, ",")

// 	endpoint := fmt.Sprintf("%s/users;use_login=1/games;game_keys=%s/teams", y.baseUrl, gcs)
// 	fc, err := y.get(ctx, endpoint, "")
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(fc.Users) <= 0 {
// 		return nil, fmt.Errorf("no teams found for gameCodes %v", gameKeys)
// 	}
// 	if len(fc.Users[0].Games) <= 0 {
// 		return nil, fmt.Errorf("no games found for gameCodes %v", gameKeys)
// 	}
// 	games := fc.Users[0].Games
// 	var teams []*yahoo.Team
// 	for _, v := range games {
// 		teams = append(teams, v.Teams...)
// 	}

// 	return teams, nil
// }

// LoadAccessToken

// @Description:
// @receiver y
// @param path
// @return IYahooClient
// @return error
func (y *yahooClient) LoadAccessToken(path string) (*yahooClient, error) {
	if path == "" {
		path = os.Getenv("HOME") + YahooTokenPath
	}
	token, err := y.yahooOAuth2.LoadAccessToken(path)
	if err != nil {
		return nil, err
	}
	y.Requestor.AuthorizationDecorator = func(req *http.Request) *http.Request {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
		return req
	}
	return y, nil
}

// GetLeague
//
//	@Description:GetLeague retrieves information about a specific Yahoo fantasy league with the given leagueID.
//	@receiver y
//	@param ctx
//	@param leagueID
//	@return *yahoo.League
//	@return error

// func (y *yahooClient) GetLeague(ctx context.Context, leagueID string) (*yahoo.League, error) {
// 	endpoint := fmt.Sprintf("%s/league/%s", y.baseUrl, leagueID)
// 	fc, err := y.get(ctx, endpoint, "league")
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &fc.League, nil
// }

func md5Hash(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func (y *yahooClient) Get(ctx context.Context, endpoint string, objType string) (*yahoo.FantasyContent, error) {
	var fc yahoo.FantasyContent
	if y.Cache != nil {
		v, exist := y.Cache.Get(ctx, md5Hash(endpoint))
		if exist {
			fmt.Println("cache hit")
			return v.(*yahoo.FantasyContent), nil
		} else {
			fmt.Printf("cache not exist for %s", endpoint)
		}
	}
	_, err := y.Requestor.Get(ctx, endpoint, &fc, xmlDecorator, &xmlDecoder{})
	if err != nil {
		return nil, err
	}

	if y.Cache != nil {
		y.Cache.Set(ctx, md5Hash(endpoint), &fc)
	}
	return &fc, nil
}
