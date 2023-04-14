package gofantasy

import "context"

type IEspnClient interface {
	GetGame(ctx context.Context, gameKey string) (*Game, error)
	GetLeague(ctx context.Context, leagueID string) (*League, error)
	GetLeagueSettings(ctx context.Context, leagueID string) (*League, error)
}

type espnClient struct {
	baseUrl    string
	baseClient *client
}

var _ IEspnClient = &espnClient{}

func (e espnClient) GetLeagueSettings(ctx context.Context, leagueID string) (*League, error) {
	//TODO implement me
	panic("implement me")
}

func (e espnClient) GetLeague(ctx context.Context, leagueID string) (*League, error) {
	//TODO implement me
	panic("implement me")
}

func (e espnClient) GetGames(ctx context.Context, gameKey string) (*interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (e espnClient) GetGame(ctx context.Context, gameKey string) (*Game, error) {
	//TODO implement me
	panic("implement me")
}
