package gofantasy

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestResolveSleeperSport(t *testing.T) {
	code, err := resolveSleeperSport("nfl")
	if err != nil || code != "nfl" {
		t.Fatalf("nfl: %q err=%v", code, err)
	}
	code, _ = resolveSleeperSport("nba")
	if code != "nba" {
		t.Fatalf("nba: %q", code)
	}
}

func TestSleeperGetUser(t *testing.T) {
	body, _ := os.ReadFile("testdata/sleeper_user.json")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/user/testuser") {
			t.Fatalf("path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	c := NewClient().Sleeper().(*sleeperClient)
	c.baseURL = srv.URL

	u, err := c.GetUser(context.Background(), "testuser")
	if err != nil {
		t.Fatal(err)
	}
	if u.UserID != "user123" || u.Username != "testuser" {
		t.Fatalf("user: %+v", u)
	}
}

func TestSleeperGetLeagueAndRosters(t *testing.T) {
	leagueBody, _ := os.ReadFile("testdata/sleeper_league.json")
	rosterBody, _ := os.ReadFile("testdata/sleeper_rosters.json")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.Contains(r.URL.Path, "/rosters"):
			_, _ = w.Write(rosterBody)
		default:
			_, _ = w.Write(leagueBody)
		}
	}))
	defer srv.Close()

	c := NewClient().Sleeper().(*sleeperClient)
	c.baseURL = srv.URL

	lg, err := c.GetLeague(context.Background(), "123456789")
	if err != nil || lg.Name != "Test League" {
		t.Fatalf("league: %+v err=%v", lg, err)
	}
	rosters, err := c.GetRosters(context.Background(), "123456789")
	if err != nil || len(rosters) != 1 || rosters[0].RosterID != 1 {
		t.Fatalf("rosters: %+v err=%v", rosters, err)
	}
}

func TestSleeperDiscoverManagedTeams(t *testing.T) {
	leagues := `[{"league_id":"123456789","name":"Test","season":"2024","sport":"nfl","status":"in_season"}]`
	rosterBody, _ := os.ReadFile("testdata/sleeper_rosters.json")
	users := `[{"user_id":"user123","display_name":"Test","metadata":{"team_name":"My Team"}}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.Contains(r.URL.Path, "/leagues/"):
			_, _ = w.Write([]byte(leagues))
		case strings.Contains(r.URL.Path, "/rosters"):
			_, _ = w.Write(rosterBody)
		case strings.Contains(r.URL.Path, "/users"):
			_, _ = w.Write([]byte(users))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer srv.Close()

	c := NewClient().Sleeper().(*sleeperClient)
	c.baseURL = srv.URL

	teams, err := c.DiscoverManagedTeams(context.Background(), "user123", "nfl", 2024)
	if err != nil {
		t.Fatal(err)
	}
	if len(teams) != 1 || teams[0].Roster.RosterID != 1 {
		t.Fatalf("teams: %+v", teams)
	}
}

func TestSleeperGetProjections(t *testing.T) {
	body, _ := os.ReadFile("testdata/sleeper_projections.json")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/projections/nfl/regular/2025/1") {
			t.Fatalf("path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	c := NewClient().Sleeper().(*sleeperClient)
	c.baseURL = srv.URL

	projs, err := c.GetProjections(context.Background(), "nfl", 2025, 1)
	if err != nil {
		t.Fatal(err)
	}

	p, ok := projs["9509"]
	if !ok {
		t.Fatalf("player 9509 missing: %v", projs)
	}
	if pts, ok := p.Points("half_ppr"); !ok || pts != 14.6 {
		t.Fatalf("half_ppr points = %v (ok=%v), want 14.6", pts, ok)
	}
	if pts, ok := p.Points("ppr"); !ok || pts != 17.3 {
		t.Fatalf("ppr points = %v (ok=%v), want 17.3", pts, ok)
	}

	// A player with only ADP noise still decodes; one with a non-numeric
	// field keeps its numeric stats instead of failing the whole response.
	if _, ok := projs["6462"]; !ok {
		t.Fatal("ADP-only entry should survive")
	}
	junk, ok := projs["junk"]
	if !ok {
		t.Fatal("entry with mixed value types should survive")
	}
	if pts, ok := junk.Points("half_ppr"); !ok || pts != 3.2 {
		t.Fatalf("junk half_ppr = %v (ok=%v), want 3.2", pts, ok)
	}
	if _, ok := junk["note"]; ok {
		t.Fatal("non-numeric stat should be dropped")
	}

	// Invalid args fail fast instead of hitting the API.
	if _, err := c.GetProjections(context.Background(), "nfl", 0, 1); err == nil {
		t.Fatal("season 0 should error")
	}
}
