package main

import (
	"context"
	"fmt"
	"github.com/gofantasy"
)

func main() {
	ctx := context.Background()
	yc, err := gofantasy.
		NewClient().
		Yahoo().LoadAccessToken("")

	if err != nil {
		panic(err)
	}

	teams, err := yc.GetUserManagedTeams(ctx, "mlb")
	if err != nil {
		panic(err)
	}

	for _, t := range teams {
		roster, err := yc.GetUserRoster(ctx, t.TeamKey)
		if err != nil {
			fmt.Printf("error getting roster for team %s: %s\n", t.TeamKey, err)
		}
		for _, p := range roster.Players {
			fmt.Printf("player: %s, position: %s, eligiblePositions: %s\n",
				p.Name.Full, p.DisplayPosition, p.EligiblePositions)
		}
	}
}
