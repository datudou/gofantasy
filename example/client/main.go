package main

import (
	"context"
	"fmt"
	"github.com/gofantasy"
	"time"
)

func main() {
	ctx := context.Background()
	lruCache := gofantasy.NewLocalCache(100)
	yc, err := gofantasy.
		NewClient().WithOptions(gofantasy.WithCache(lruCache)).
		Yahoo().LoadAccessToken("")

	if err != nil {
		panic(err)
	}

	teams, err := yc.GetUserManagedTeams(ctx, "mlb")
	if err != nil {
		panic(err)
	}

	for {
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
		time.Sleep(time.Second * 1)
	}

}
