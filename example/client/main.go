package main

import (
	"context"
	"fmt"
	"github.com/gofantasy"
)

func main() {
	ctx := context.Background()
	yc, err := gofantasy.NewClient().Yahoo().WithAccessToken("")
	if err != nil {
		panic(err)
	}
	game, err := yc.GetGame(ctx, "nfl")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", game)
	league, err := yc.GetLeague(ctx, "223.l.431")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", league)
}
