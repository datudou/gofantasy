# GoFantasy
## About
GoFantasy is a go api wrapper for kind of fantasy platforms like Yahoo fantasy, ESPN and so on. 

## Target
Currently, GoFantasy is still on its early stages of development. The following platforms are will be supported:

* Yahoo Fantasy
* ESPN Fantasy

## How to use 
### Yahoo Fantasy  
#### Get your roster

```go
package main

import (
	"context"
	"fmt"
	"github.com/gofantasy"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()
	//Yahoo api call need a access token, you can get it from https://developer.yahoo.com/oauth2/guide/flows_authcode/
	//or you can refer the code in the ./example/auth/main.go to get the accessToken
	accessToken := "Your access token"
	yc, err := gofantasy.
		NewClient().
		Yahoo().LoadAccessToken(accessToken)

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
			fmt.Printf("player: %s, position: %s, eligiblePositions: %s\n", p.Name.Full, p.DisplayPosition, p.EligiblePositions)
		}
	}
}
```

## Write operations

Besides reads, the Yahoo client also supports state-mutating operations:
`SetRoster` / `SetRosterForWeek` (lineups), `AddPlayer` / `DropPlayer` /
`AddDropPlayer` (waivers and free agents), `ProposeTrade` / `RespondToTrade` /
`CancelTransaction` (trades). These mutate league state — applications built on
them should gate every write behind explicit user approval to stay within
Yahoo's API terms, which prohibit unattended automation acting on a user's
behalf.

### ESPN Fantasy

ESPN has no official public API. This client uses the undocumented v3 read API
(`lm-api-reads.fantasy.espn.com`). Private leagues require browser session
cookies `SWID` and `espn_s2` (DevTools → Application → Cookies on espn.com).

```go
package main

import (
	"context"
	"fmt"

	"github.com/gofantasy"
)

func main() {
	ctx := context.Background()
	ec := gofantasy.NewClient().ESPN().WithCookies(
		"{YOUR-SWID}",
		"YOUR_ESPN_S2_COOKIE",
	)

	leagueID := 899513
	season := 2024

	roster, err := ec.GetRoster(ctx, "nfl", season, leagueID, 1)
	if err != nil {
		panic(err)
	}
	for _, e := range roster.Entries {
		fmt.Println(e.PlayerPoolEntry.Player.FullName, e.LineupSlotID)
	}

	teams, err := ec.GetStandings(ctx, "nba", season, leagueID)
	if err != nil {
		panic(err)
	}
	for _, t := range teams {
		fmt.Printf("%s %s (%d-%d)\n", t.Location, t.Nickname, t.Record.Overall.Wins, t.Record.Overall.Losses)
	}
}
```

Sport aliases: `nfl`→`ffl`, `nba`→`fba`, `mlb`→`flb`, `nhl`→`nhl`. Seasons
before 2018 use the `leagueHistory` endpoint automatically.

ESPN does not expose a "list all my leagues" API. Use `DiscoverManagedTeams`
with league IDs from your fantasy.espn.com URL, or call `GetLeague` /
`GetTeams` when you already know the league ID.

Write operations (lineup changes, add/drop) are planned; reads are implemented
in Phase 1.


## Contribution

Contributions are welcome! If you find a bug or want to suggest a new feature, feel free to open an issue or create a pull request.

## Attribution

Fantasy data is provided by Yahoo Fantasy. Applications built on the Yahoo
Fantasy Sports API must display this attribution and comply with Yahoo's API
Access and Use Agreement.

## License

GoFantasy is released under the MIT License. See LICENSE file for details.
