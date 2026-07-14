package main

import (
	"context"
	"fmt"

	"github.com/gofantasy"
)

func main() {
	c := gofantasy.NewClient().Sleeper()
	projs, err := c.GetProjections(context.Background(), "nfl", 2025, 7)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	fmt.Printf("players with projections: %d\n", len(projs))
	n := 0
	for id, p := range projs {
		if pts, ok := p.Points("half_ppr"); ok && pts > 15 {
			fmt.Printf("  player %s: half_ppr=%.1f ppr=%v\n", id, pts, p["pts_ppr"])
			n++
			if n >= 3 {
				break
			}
		}
	}
}
