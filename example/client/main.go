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
		WithOptions(gofantasy.WithCache(128)).
		Yahoo().
		WithAccessToken("")

	if err != nil {
		panic(err)
	}

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
