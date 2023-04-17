package gofantasy

import (
	"context"
	"crypto"
	"fmt"
	"github.com/gofantasy/model/yahoo"
	"io"
	"net/http"
	"os"
)

type IYahooClient interface {
	GetLeague(ctx context.Context, leagueID string) (*yahoo.League, error)
	GetGameKeyBySeason(ctx context.Context, gameCode string, season string) (*[]yahoo.Game, error)
	GetUserGames(ctx context.Context, gameKey string) (*yahoo.Game, error)
	GetUserTeams(ctx context.Context, leagueID string) (*yahoo.User, error)
	OAuth2(clientID, clientSecret, redirectURL string) IYahooOAuth2
	LoadAccessToken(path string) (IYahooClient, error)
}

type yahooClient struct {
	baseUrl     string
	baseClient  *client
	yahooOAuth2 *yahooOAuth2
}

var _ IYahooClient = &yahooClient{}

func (y *yahooClient) OAuth2(clientID, clientSecret, redirectURL string) IYahooOAuth2 {
	return y.yahooOAuth2.OAuth2(clientID, clientSecret, redirectURL)
}

func (y *yahooClient) GetGameKeyBySeason(ctx context.Context, gameCode string, season string) (*[]yahoo.Game, error) {

	endpoint := fmt.Sprintf("%s/games;game_codes=%s;seasons=%s", y.baseUrl, gameCode, season)
	fc, err := y.get(ctx, endpoint, "")
	if err != nil {
		return nil, err
	}
	if len(fc.Games) == 0 {
		return nil, fmt.Errorf("no games found for gameCode %s and season %s", gameCode, season)

	}
	return &fc.Games, nil
}

func (*yahooClient) GetUserTeams(ctx context.Context, leagueID string) (*yahoo.User, error) {
	//endpoint := fmt.Sprintf("%s/games;game_codes=%s;seasons=%s", y.baseUrl, gameCode, season)
	//fc, err := y.get(ctx, endpoint, "")
	//if err != nil {
	//	return nil, err
	//}
	//if len(fc.Games) == 0 {
	//	return nil, fmt.Errorf("no games found for gameCode %s and season %s", gameCode, season)
	//
	//}
	//return &fc.Games, nil
	panic("unimplemented")
}

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

func (y *yahooClient) GetUserGames(ctx context.Context, gameKey string) (*yahoo.Game, error) {
	endpoint := fmt.Sprintf("%s/users;use_login=1/games;games_key=%s", y.baseUrl, gameKey)
	fc, err := y.get(ctx, endpoint, "")
	if err != nil {
		return nil, err
	}
	return &fc.Users[0].Games[0], nil
}

func (y *yahooClient) GetLeague(ctx context.Context, leagueID string) (*yahoo.League, error) {
	endpoint := fmt.Sprintf("%s/league/%s", y.baseUrl, leagueID)
	fc, err := y.get(ctx, endpoint, "league")
	if err != nil {
		return nil, err
	}
	return &fc.League, nil
}

func (y *yahooClient) GetLeagueSettings(ctx context.Context, leagueID string) (yahoo.League, error) {
	panic("implement me")
}

func (y *yahooClient) get(ctx context.Context, endpoint string, objType string) (*yahoo.FantasyContent, error) {
	var fc yahoo.FantasyContent
	if y.baseClient.cache != nil {
		v, exist := y.baseClient.cache.Get(endpoint)
		if exist {
			return v.(*yahoo.FantasyContent), nil
		}
	}
	_, err := y.baseClient.requestor.Get(ctx, endpoint, &fc, xmlDecorator, &xmlDecoder{})
	if err != nil {
		return nil, err
	}

	if y.baseClient.cache != nil {
		y.baseClient.cache.Add(endpoint, &fc)
	}
	return &fc, nil
}

func md5(str string) string {
	w := crypto.MD5.New()
	io.WriteString(w, str)
	md5str := fmt.Sprintf("%x", w.Sum(nil))
	return md5str
}
