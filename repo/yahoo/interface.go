package yahoo

import (
	"context"

	"github.com/gofantasy/model/yahoo"
)

type UserOption func(*user) IUser

type IUser interface {
	GetTeams(ctx context.Context, opts ...UserOption) ([]*yahoo.Team, error)
	Me() *user
	Games(gameKeys ...string) *user
	Get(ctx context.Context) (*yahoo.FantasyContent, error)

	// WithSeason(season string) UserOption
	// WithGame(gameCode string) UserOption
}
