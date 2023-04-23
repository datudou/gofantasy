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
		WithAccessToken("") // if pass "" , it will read token object from the file saved at  ~/.config/gofantasy/yahoo_token.json

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
