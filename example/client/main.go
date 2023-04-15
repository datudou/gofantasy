package main

import (
	"context"
	"fmt"
	"github.com/gofantasy"
	"time"
)

func main() {
	ctx := context.Background()
	yc, err := gofantasy.
		NewClient().
		WithOptions(gofantasy.WithCache()).
		Yahoo().
		WithAccessToken("")
	if err != nil {
		panic(err)
	}

	//league, err := yc.GetLeague(ctx, "223.l.431")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%+v", league)
	for {
		go func() {
			game, err := yc.GetGame(ctx, "nfl")
			if err != nil {
				panic(err)
			}
			fmt.Printf("%+v", game)
		}()
		time.Sleep(1 * time.Second)
	}
}
