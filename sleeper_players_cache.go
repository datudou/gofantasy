package gofantasy

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type sleeperPlayersCache struct {
	dir string
}

func newSleeperPlayersCache(dir string) *sleeperPlayersCache {
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			dir = os.TempDir()
		} else {
			dir = filepath.Join(home, ".cache", "gofantasy")
		}
	}
	return &sleeperPlayersCache{dir: dir}
}

func (c *sleeperPlayersCache) path(sport string) string {
	return filepath.Join(c.dir, fmt.Sprintf("sleeper_players_%s.json", sport))
}

func (c *sleeperPlayersCache) load(ctx context.Context, sport string, fetch func(context.Context, string) (map[string]json.RawMessage, error)) (map[string]json.RawMessage, error) {
	path := c.path(sport)
	if b, err := os.ReadFile(path); err == nil {
		var players map[string]json.RawMessage
		if json.Unmarshal(b, &players) == nil && len(players) > 0 {
			if info, err := os.Stat(path); err == nil && time.Since(info.ModTime()) < 24*time.Hour {
				return players, nil
			}
		}
	}

	players, err := fetch(ctx, sport)
	if err != nil {
		return nil, err
	}
	_ = os.MkdirAll(c.dir, 0o755)
	if b, err := json.Marshal(players); err == nil {
		_ = os.WriteFile(path, b, 0o644)
	}
	return players, nil
}
