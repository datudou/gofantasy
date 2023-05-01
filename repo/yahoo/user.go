package yahoo

import (
	"context"
	"fmt"
	"strings"

	gff "github.com/gofantasy"
	"github.com/gofantasy/model/yahoo"
	y "github.com/gofantasy/model/yahoo"
)

var GameKeys = map[string]string{
	"nfl": "nfl",
	"nba": "nba",
	"mlb": "mlb",
	"nhl": "nhl",
}

type user struct {
	y.User
	season      string
	gameCode    string
	client      gff.IClient
	queryParams []string
}

func NewUser(client gff.IClient) IUser {
	return &user{
		client: client,
	}
}

func (u *user) Me() *user {
	u.queryParams = append(u.queryParams, "users;use_login=1")
	return u
}

func (u *user) Games(gameKeys ...string) *user {
	if !isValidGameKeys(gameKeys...) {
		return nil
	}
	gcs := strings.Join(gameKeys, ",")
	u.queryParams = append(u.queryParams, fmt.Sprintf("games;games_key=%s", gcs))
	return u
}

func (u *user) Get(ctx context.Context) (*yahoo.FantasyContent, error) {
	endpoint := fmt.Sprintf("%s/%s", gff.YahooBaseURL, strings.Join(u.queryParams, "/"))
	fmt.Println(endpoint)
	return u.client.Get(ctx, endpoint, "")
}

func (u *user) GetTeams(ctx context.Context, opts ...UserOption) ([]*y.Team, error) {
	panic("implement me")
}

func (u *user) WithSeason(season string) *user {
	u.season = season
	return u
}

func (u *user) WithGameCode(gameCode string) *user {
	u.gameCode = gameCode
	return u
}

// isValidGameKeys checks if the provided game keys are valid by looking them up in the GameKeys map.
// If any of the keys are not present in the map, the function returns false. Otherwise, it returns true.
func isValidGameKeys(gameKeys ...string) bool {
	for _, v := range gameKeys {
		if _, ok := GameKeys[v]; !ok {
			return false
		}
	}
	return true
}
