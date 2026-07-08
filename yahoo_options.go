package gofantasy

import (
	"fmt"
	"strconv"
	"strings"
)

// PlayerOption configures filters when listing players in a league.
// See https://developer.yahoo.com/fantasysports/guide/players-collection.html
type PlayerOption func(*playerFilter)

type playerFilter struct {
	params map[string]string
}

func (f *playerFilter) encode() string {
	if len(f.params) == 0 {
		return ""
	}
	var sb strings.Builder
	for k, v := range f.params {
		sb.WriteString(";")
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(v)
	}
	return sb.String()
}

func newPlayerFilter(opts ...PlayerOption) *playerFilter {
	f := &playerFilter{params: map[string]string{}}
	for _, o := range opts {
		o(f)
	}
	return f
}

// WithPlayerPosition filters players by position (e.g. "QB", "RB", "PG").
func WithPlayerPosition(pos string) PlayerOption {
	return func(f *playerFilter) { f.params["position"] = pos }
}

// WithPlayerStatus filters players by status: "A" (available), "FA" (free agents),
// "W" (waivers), "T" (taken), "K" (keepers).
func WithPlayerStatus(status string) PlayerOption {
	return func(f *playerFilter) { f.params["status"] = status }
}

// WithPlayerSearch filters players whose name matches the given query.
func WithPlayerSearch(q string) PlayerOption {
	return func(f *playerFilter) { f.params["search"] = q }
}

// WithPlayerSort sets the sort order. Valid values include "AR" (actual rank),
// "PTS", "OR" (overall rank), "NAME", "OWN" (% owned).
func WithPlayerSort(sort string) PlayerOption {
	return func(f *playerFilter) { f.params["sort"] = sort }
}

// WithPlayerPagination sets the start and count for paging through players.
// Yahoo returns at most 25 players per page; use start to advance.
func WithPlayerPagination(start, count int) PlayerOption {
	return func(f *playerFilter) {
		f.params["start"] = strconv.Itoa(start)
		f.params["count"] = strconv.Itoa(count)
	}
}

// TransactionOption configures filters when listing league transactions.
type TransactionOption func(*transactionFilter)

type transactionFilter struct {
	params map[string]string
}

func (f *transactionFilter) encode() string {
	if len(f.params) == 0 {
		return ""
	}
	var sb strings.Builder
	for k, v := range f.params {
		sb.WriteString(";")
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(v)
	}
	return sb.String()
}

func newTransactionFilter(opts ...TransactionOption) *transactionFilter {
	f := &transactionFilter{params: map[string]string{}}
	for _, o := range opts {
		o(f)
	}
	return f
}

// WithTransactionType filters by type: "add", "drop", "commish", "trade".
// Multiple types may be comma-separated.
func WithTransactionType(t string) TransactionOption {
	return func(f *transactionFilter) { f.params["type"] = t }
}

// WithTransactionTeamKey limits transactions to those involving the given team.
func WithTransactionTeamKey(teamKey string) TransactionOption {
	return func(f *transactionFilter) { f.params["team_key"] = teamKey }
}

// WithTransactionCount limits the number of returned transactions.
func WithTransactionCount(count int) TransactionOption {
	return func(f *transactionFilter) { f.params["count"] = strconv.Itoa(count) }
}

// StatOption configures the coverage when fetching team or player stats.
type StatOption func(*statFilter)

type statFilter struct {
	params map[string]string
}

func (f *statFilter) encode() string {
	if len(f.params) == 0 {
		return ""
	}
	var sb strings.Builder
	for k, v := range f.params {
		sb.WriteString(";")
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(v)
	}
	return sb.String()
}

func newStatFilter(opts ...StatOption) *statFilter {
	f := &statFilter{params: map[string]string{}}
	for _, o := range opts {
		o(f)
	}
	return f
}

// WithStatSeason fetches stats aggregated for the given season.
func WithStatSeason(season int) StatOption {
	return func(f *statFilter) {
		f.params["type"] = "season"
		f.params["season"] = strconv.Itoa(season)
	}
}

// WithStatWeek fetches stats for the given week (head-to-head leagues).
func WithStatWeek(week int) StatOption {
	return func(f *statFilter) {
		f.params["type"] = "week"
		f.params["week"] = strconv.Itoa(week)
	}
}

// WithStatDate fetches stats for the given date (YYYY-MM-DD, daily leagues).
func WithStatDate(date string) StatOption {
	return func(f *statFilter) {
		f.params["type"] = "date"
		f.params["date"] = date
	}
}

// WithStatLastWeek fetches stats for the last week.
func WithStatLastWeek() StatOption {
	return func(f *statFilter) { f.params["type"] = "lastweek" }
}

// formatWeeks turns a list of weeks into the comma-separated form Yahoo expects.
// Returns an empty string when no weeks are provided.
func formatWeeks(weeks []int) string {
	if len(weeks) == 0 {
		return ""
	}
	parts := make([]string, 0, len(weeks))
	for _, w := range weeks {
		parts = append(parts, fmt.Sprintf("%d", w))
	}
	return strings.Join(parts, ",")
}
