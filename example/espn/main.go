package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gofantasy"
)

func main() {
	swid := os.Getenv("ESPN_SWID")
	s2 := os.Getenv("ESPN_S2")
	leagueID := 899513
	season := 2024
	teamID := 1

	ctx := context.Background()
	ec := gofantasy.NewClient().ESPN().WithCookies(swid, s2)

	roster, err := ec.GetRoster(ctx, "nfl", season, leagueID, teamID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Roster entries: %d\n", len(roster.Entries))
	for _, e := range roster.Entries {
		fmt.Printf("- %s (slot %d)\n", e.PlayerPoolEntry.Player.FullName, e.LineupSlotID)
	}
}
