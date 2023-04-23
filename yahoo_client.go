package yahoo

import (
	"context"
	"fmt"
	"github.com/gofantasy"
	"net/http"
)

type IYahooClient interface {
	GetLeague(ctx context.Context, leagueID string) (*League, error)
	GetGameKeyBySeason(ctx context.Context, season string) (*Game, error)
	GetUserGames(ctx context.Context, gameKey string) (*Game, error)
	GetUserTeams(ctx context.Context, leagueID string) (*User, error)
	LoadAccessToken() (IYahooClient, error)
	OAuth2(clientID, clientSecret, redirectURL string) IYahooOAuth2
}

type YahooClient struct {
	BaseUrl     string
	BaseClient  *gofantasy.Client
	yahooOAuth2 *yahooOAuth2
}

var _ IYahooClient = &YahooClient{}

func (y *YahooClient) OAuth2(clientID, clientSecret, redirectURL string) IYahooOAuth2 {
	return y.yahooOAuth2.OAuth2(clientID, clientSecret, redirectURL)
}

func (y *YahooClient) GetGameKeyBySeason(ctx context.Context, season string) (*Game, error) {
	//TODO implement me
	panic("implement me")
}

func (*YahooClient) GetUserTeams(ctx context.Context, leagueID string) (*User, error) {
	panic("unimplemented")
}

func (y *YahooClient) LoadAccessToken() (IYahooClient, error) {
	err := readToken("", &y.yahooOAuth2.token)
	if err != nil {
		return nil, err
	}
	accessToken := y.yahooOAuth2.token.AccessToken
	y.BaseClient.Requestor.AuthorizationDecorator = func(req *http.Request) *http.Request {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
		return req
	}
	return y, nil
}

func (y *YahooClient) GetUserGames(ctx context.Context, gameKey string) (*Game, error) {
	cacheKey := fmt.Sprintf("userGame%s", gameKey)
	endpoint := fmt.Sprintf("%s/users;use_login=1/games;games_key=%s", y.BaseUrl, gameKey)
	fc, err := y.getCachedOrFetch(ctx, cacheKey, endpoint, "league")
	if err != nil {
		return nil, err
	}
	return &fc.Users[0].Games[0], nil
}

func (y *YahooClient) GetLeague(ctx context.Context, leagueID string) (*League, error) {
	cacheKey := fmt.Sprintf("league%s", leagueID)
	endpoint := fmt.Sprintf("%s/league/%s", y.BaseUrl, leagueID)
	fc, err := y.getCachedOrFetch(ctx, cacheKey, endpoint, "league")
	if err != nil {
		return nil, err
	}
	return &fc.League, nil
}

func (y *YahooClient) GetLeagueSettings(ctx context.Context, leagueID string) (League, error) {
	//TODO implement me
	panic("implement me")
}

func (y *YahooClient) getCachedOrFetch(ctx context.Context, cacheKey, endpoint string, objType string) (*FantasyContent, error) {
	var fc FantasyContent
	if y.BaseClient.Cache != nil && cacheKey != "" {
		v, exist := y.BaseClient.Cache.Get(cacheKey)
		if exist {
			fmt.Println("----> cache hit!")
			return v.(*FantasyContent), nil
		}
	}
	_, err := y.BaseClient.Requestor.Get(ctx, endpoint, &fc)
	if err != nil {
		return nil, err
	}
	if y.BaseClient.Cache != nil {
		y.BaseClient.Cache.Add(cacheKey, &fc)
	}
	return &fc, nil
}
