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
		WithOptions(gofantasy.WithCache(128)).Yahoo().LoadAccessToken("")

	if err != nil {
		panic(err)
	}

	game, err := yc.GetGameKeyBySeason(ctx, "nfl", "2020")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", game)
}
