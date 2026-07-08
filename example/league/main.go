package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gofantasy"
)

// This example demonstrates the league/team/player-level helpers added on top of
// the existing GetUserManagedTeams / GetUserRoster flow.
//
// Run example/auth first to mint a token, then set LEAGUE_KEY to one of your
// leagues (e.g. "mlb.l.123456") before running this program.
func main() {
	leagueKey := os.Getenv("LEAGUE_KEY")
	if leagueKey == "" {
		fmt.Println("set LEAGUE_KEY to the league you want to inspect (e.g. mlb.l.123456)")
		os.Exit(1)
	}

	ctx := context.Background()
	yc, err := gofantasy.NewClient().
		WithOptions(gofantasy.WithCache(gofantasy.NewLocalCache(256))).
		Yahoo().
		LoadAccessToken("")
	if err != nil {
		panic(err)
	}

	league, err := yc.GetLeague(ctx, leagueKey)
	if err != nil {
		panic(err)
	}
	fmt.Printf("League %s (%s) — current week %d\n\n", league.Name, league.LeagueKey, league.CurrentWeek)

	standings, err := yc.GetLeagueStandings(ctx, leagueKey)
	if err != nil {
		panic(err)
	}
	fmt.Println("Standings:")
	for _, t := range standings {
		s := t.TeamStandings
		fmt.Printf("  %2d. %-25s  %d-%d-%d  PF %.1f / PA %.1f\n",
			s.Rank, t.Name, s.OutcomeTotals.Wins, s.OutcomeTotals.Losses, s.OutcomeTotals.Ties,
			s.PointsFor, s.PointsAgainst)
	}

	scoreboard, err := yc.GetLeagueScoreboard(ctx, leagueKey)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nScoreboard for week %d:\n", scoreboard.Week)
	for _, m := range scoreboard.Matchups {
		if len(m.Teams) != 2 {
			continue
		}
		a, b := m.Teams[0], m.Teams[1]
		fmt.Printf("  %s (%.1f)  vs  %s (%.1f)\n", a.Name, a.TeamPoints.Total, b.Name, b.TeamPoints.Total)
	}

	txns, err := yc.GetLeagueTransactions(ctx, leagueKey,
		gofantasy.WithTransactionType("add,drop"),
		gofantasy.WithTransactionCount(5),
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nLast %d transactions:\n", len(txns))
	for _, tx := range txns {
		fmt.Printf("  [%s] %s\n", tx.Type, tx.TransactionKey)
	}

	free, err := yc.GetLeaguePlayers(ctx, leagueKey,
		gofantasy.WithPlayerStatus("FA"),
		gofantasy.WithPlayerSort("AR"),
		gofantasy.WithPlayerPagination(0, 5),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("\nTop available free agents:")
	for _, p := range free {
		fmt.Printf("  %s (%s) — %s\n", p.Name.Full, p.DisplayPosition, p.EditorialTeamAbbr)
	}
}
