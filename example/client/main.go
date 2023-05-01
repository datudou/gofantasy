package main

import (
	"context"
	"fmt"

	gff "github.com/gofantasy"
	repo "github.com/gofantasy/repo/yahoo"
)

func main() {
	ctx := context.Background()
	lruCache := gff.NewLocalCache(100)
	client, err := gff.
		NewYahooClient().
		WithOptions(gff.WithCache(lruCache)).
		LoadAccessToken("")
	if err != nil {
		panic(err)
	}

	fc, err := repo.NewUser(client).Me().Games("mlb").Get(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(fc.Users[0].Games)

}
