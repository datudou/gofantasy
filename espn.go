package gofantasy

// import (
// 	"context"
// 	"github.com/gofantasy/model/espn"
// )

// type IEspnClient interface {
// 	GetLeague(ctx context.Context, leagueID string) (*espn.League, error)
// }

// type espnClient struct {
// 	baseUrl    string
// 	baseClient *client
// }

// var _ IEspnClient = &espnClient{}

// func (e *espnClient) GetLeague(ctx context.Context, leagueID string) (*espn.League, error) {

// }
// func (e *espnClient) get(ctx context.Context, endpoint string) (*eahoo, error) {
// 	if e.baseClient.cache != nil {
// 		v, exist := e.baseClient.cache.Get(ctx, endpoint)
// 		if exist {
// 		}
// 	}
// 	_, err := e.baseClient.requestor.Get(ctx, endpoint, &fc, jsonDecorator, &jsonDecoder{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	if e.baseClient.cache != nil {
// 		e.baseClient.cache.Add(ctx, endpoint, &fc)
// 	}
// 	return &fc, nil
// }
