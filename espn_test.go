package gofantasy

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gofantasy/model/espn"
)

func TestResolveEspnSport(t *testing.T) {
	code, err := resolveEspnSport("nfl")
	if err != nil || code != "ffl" {
		t.Fatalf("nfl -> ffl: got %q err=%v", code, err)
	}
	if _, err := resolveEspnSport("cricket"); err == nil {
		t.Fatal("expected error for unsupported sport")
	}
}

func TestEspnLeagueEndpoint(t *testing.T) {
	endpoint, history, err := espnLeagueEndpoint(EspnBaseURL, "ffl", 2024, 899513, newEspnLeagueConfig(WithEspnViews("mTeam", "mRoster")))
	if err != nil {
		t.Fatal(err)
	}
	if history {
		t.Fatal("2024 should not use leagueHistory")
	}
	if !strings.Contains(endpoint, "/seasons/2024/segments/0/leagues/899513") {
		t.Fatalf("unexpected endpoint: %s", endpoint)
	}
	if !strings.Contains(endpoint, "view=mTeam") || !strings.Contains(endpoint, "view=mRoster") {
		t.Fatalf("missing views: %s", endpoint)
	}

	old, history, err := espnLeagueEndpoint(EspnBaseURL, "ffl", 2015, 899513, newEspnLeagueConfig())
	if err != nil {
		t.Fatal(err)
	}
	if !history {
		t.Fatal("2015 should use leagueHistory")
	}
	if !strings.Contains(old, "/leagueHistory/899513") || !strings.Contains(old, "seasonId=2015") {
		t.Fatalf("unexpected history endpoint: %s", old)
	}
}

func TestTeamOwnedBySWID(t *testing.T) {
	team := espn.Team{Owners: []string{"{11111111-1111-1111-1111-111111111111}"}}
	if !teamOwnedBySWID(team, "{11111111-1111-1111-1111-111111111111}") {
		t.Fatal("expected owner match with braces")
	}
	if teamOwnedBySWID(team, "22222222-2222-2222-2222-222222222222") {
		t.Fatal("unexpected owner match")
	}
}

func TestEspnClientGetLeague(t *testing.T) {
	body, err := os.ReadFile("testdata/espn_league_ffl.json")
	if err != nil {
		t.Fatal(err)
	}

	var gotCookie string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotCookie = r.Header.Get("Cookie")
		if !strings.Contains(r.URL.Path, "/ffl/seasons/2024/segments/0/leagues/899513") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	c := NewClient().ESPN().WithCookies("{SWID}", "espn-s2-value").(*espnClient)
	c.baseURL = srv.URL + "/apis/v3/games"

	lg, err := c.GetLeague(context.Background(), "nfl", 2024, 899513, WithEspnViews("mTeam", "mRoster"))
	if err != nil {
		t.Fatal(err)
	}
	if lg.Settings.Name != "Test League" {
		t.Fatalf("league name: %q", lg.Settings.Name)
	}
	if len(lg.Teams) != 2 {
		t.Fatalf("teams: %d", len(lg.Teams))
	}
	if gotCookie != "SWID={SWID}; espn_s2=espn-s2-value" {
		t.Fatalf("cookie header: %q", gotCookie)
	}
}

func TestEspnClientGetRosterAndFreeAgents(t *testing.T) {
	body, err := os.ReadFile("testdata/espn_league_ffl.json")
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	c := NewClient().ESPN().(*espnClient)
	c.baseURL = srv.URL + "/apis/v3/games"

	roster, err := c.GetRoster(context.Background(), "ffl", 2024, 899513, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(roster.Entries) != 1 || roster.Entries[0].PlayerPoolEntry.Player.FullName != "Patrick Mahomes" {
		t.Fatalf("unexpected roster: %+v", roster)
	}

	players, err := c.GetFreeAgents(context.Background(), "nfl", 2024, 899513)
	if err != nil {
		t.Fatal(err)
	}
	if len(players) != 1 || players[0].FullName != "Free Agent Player" {
		t.Fatalf("unexpected free agents: %+v", players)
	}
}

func TestEspnDiscoverManagedTeams(t *testing.T) {
	body, err := os.ReadFile("testdata/espn_league_ffl.json")
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	c := NewClient().ESPN().WithCookies("{11111111-1111-1111-1111-111111111111}", "s2").(*espnClient)
	c.baseURL = srv.URL + "/apis/v3/games"

	teams, err := c.DiscoverManagedTeams(context.Background(), "nfl", 2024, 899513)
	if err != nil {
		t.Fatal(err)
	}
	if len(teams) != 1 || teams[0].Team.ID != 1 {
		t.Fatalf("unexpected managed teams: %+v", teams)
	}
}

func TestEspnPlayerFilterJSON(t *testing.T) {
	cfg := newEspnPlayerConfig(WithEspnPlayerLimit(25), WithEspnPlayerStatus("FREEAGENT"))
	s, err := cfg.xFantasyFilterJSON()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(s, `"limit":25`) || !strings.Contains(s, `"FREEAGENT"`) {
		t.Fatalf("filter json: %s", s)
	}
}

func TestNormalizeSWID(t *testing.T) {
	if normalizeSWID(" {abc} ") != "ABC" {
		t.Fatal("normalizeSWID failed")
	}
}
