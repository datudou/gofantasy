package gofantasy

import (
	"encoding/json"
	"fmt"
)

type espnLeagueConfig struct {
	views           []string
	scoringPeriodID *int
	matchupPeriodID *int
	xFantasyFilter  any
}

// EspnLeagueOption configures a league fetch (views, scoring period, filters).
type EspnLeagueOption func(*espnLeagueConfig)

// WithEspnViews requests one or more ESPN league views on the same call.
func WithEspnViews(views ...string) EspnLeagueOption {
	return func(c *espnLeagueConfig) {
		c.views = append(c.views, views...)
	}
}

// WithScoringPeriodID scopes roster/live views to a scoring period.
func WithScoringPeriodID(period int) EspnLeagueOption {
	return func(c *espnLeagueConfig) {
		c.scoringPeriodID = &period
	}
}

// WithMatchupPeriodID scopes matchup views to a matchup period.
func WithMatchupPeriodID(period int) EspnLeagueOption {
	return func(c *espnLeagueConfig) {
		c.matchupPeriodID = &period
	}
}

// WithXFantasyFilter sets the X-Fantasy-Filter header (used for player pool queries).
func WithXFantasyFilter(filter any) EspnLeagueOption {
	return func(c *espnLeagueConfig) {
		c.xFantasyFilter = filter
	}
}

type espnPlayerConfig struct {
	limit  int
	offset int
	status []string
}

// EspnPlayerOption configures free-agent / player-pool queries.
type EspnPlayerOption func(*espnPlayerConfig)

// WithEspnPlayerLimit caps how many players are returned (default 50).
func WithEspnPlayerLimit(limit int) EspnPlayerOption {
	return func(c *espnPlayerConfig) {
		c.limit = limit
	}
}

// WithEspnPlayerOffset skips the first N players in the pool query.
func WithEspnPlayerOffset(offset int) EspnPlayerOption {
	return func(c *espnPlayerConfig) {
		c.offset = offset
	}
}

// WithEspnPlayerStatus filters pool entries (e.g. "FREEAGENT", "WAIVERS").
func WithEspnPlayerStatus(status ...string) EspnPlayerOption {
	return func(c *espnPlayerConfig) {
		c.status = append(c.status, status...)
	}
}

func newEspnLeagueConfig(opts ...EspnLeagueOption) espnLeagueConfig {
	cfg := espnLeagueConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

func newEspnPlayerConfig(opts ...EspnPlayerOption) espnPlayerConfig {
	cfg := espnPlayerConfig{
		limit:  50,
		offset: 0,
		status: []string{"FREEAGENT", "WAIVERS"},
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

func (c espnPlayerConfig) xFantasyFilterJSON() (string, error) {
	filter := map[string]any{
		"players": map[string]any{
			"filterStatus": map[string]any{
				"value": c.status,
			},
			"limit":  c.limit,
			"offset": c.offset,
		},
	}
	b, err := json.Marshal(filter)
	if err != nil {
		return "", fmt.Errorf("marshal X-Fantasy-Filter: %w", err)
	}
	return string(b), nil
}

func resolveEspnSport(sport string) (string, error) {
	if code, ok := EspnGameCodes[sport]; ok {
		return code, nil
	}
	return "", fmt.Errorf("unsupported sport %q (use nfl/nba/mlb/nhl or ffl/fba/flb/nhl)", sport)
}
