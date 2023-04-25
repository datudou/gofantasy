# GoFantasy
## About
GoFantasy is a go api wrapper for kind of fantasy platforms like Yahoo fantasy, ESPN and so on. 

## Target
Currently, GoFantasy is still on its early stages of development. The following platforms are will be supported:

* Yahoo Fantasy
* ESPN Fantasy

## How to use 
### Yahoo Fantasy  
#### Get You roster 
```go
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
			fmt.Printf("player: %s, position: %s, eligiblePositions: %s\n", p.Name.Full, p.DisplayPosition, p.EligiblePositions)
		}
	}
}
```

## Contribution

Contributions are welcome! If you find a bug or want to suggest a new feature, feel free to open an issue or create a pull request.

## License

GoFantasy is released under the MIT License. See LICENSE file for details.
