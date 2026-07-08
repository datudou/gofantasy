package yahoo

import (
	"encoding/xml"
	"testing"
)

func TestUnmarshalLeagueSettingsTradeEndDate(t *testing.T) {
	const raw = `<?xml version="1.0" encoding="UTF-8"?>
<fantasy_content xmlns="http://fantasysports.yahooapis.com/fantasy/v2/base.rng">
  <league>
    <league_key>428.l.36687</league_key>
    <settings>
      <trade_end_date>2023-03-10</trade_end_date>
      <uses_playoff>1</uses_playoff>
    </settings>
  </league>
</fantasy_content>`
	var fc FantasyContent
	if err := xml.Unmarshal([]byte(raw), &fc); err != nil {
		t.Fatalf("unmarshal settings: %v", err)
	}
	if fc.League.Settings.TradeEndDate != "2023-03-10" || !fc.League.Settings.UsesPlayoff {
		t.Fatalf("unexpected settings: %+v", fc.League.Settings)
	}
}

func TestUnmarshalLeagueStandings(t *testing.T) {
	const raw = `<?xml version="1.0" encoding="UTF-8"?>
<fantasy_content xmlns="http://fantasysports.yahooapis.com/fantasy/v2/base.rng">
  <league>
    <league_key>428.l.36687</league_key>
    <standings>
      <teams>
        <team>
          <team_key>428.l.36687.t.15</team_key>
          <name>Mine</name>
          <team_standings>
            <rank>1</rank>
            <outcome_totals><wins>45</wins><losses>30</losses><ties>0</ties></outcome_totals>
            <points_for>5123.4</points_for>
            <points_against>4980.1</points_against>
          </team_standings>
        </team>
      </teams>
    </standings>
  </league>
</fantasy_content>`
	var fc FantasyContent
	if err := xml.Unmarshal([]byte(raw), &fc); err != nil {
		t.Fatalf("unmarshal standings: %v", err)
	}
	if len(fc.League.Standings) != 1 || fc.League.Standings[0].TeamStandings.Rank != 1 {
		t.Fatalf("unexpected standings: %+v", fc.League.Standings)
	}
}

func TestUnmarshalLeagueEditKeyDate(t *testing.T) {
	const raw = `<?xml version="1.0" encoding="UTF-8"?>
<fantasy_content xmlns="http://fantasysports.yahooapis.com/fantasy/v2/base.rng">
  <league>
    <league_key>428.l.36687</league_key>
    <name>Test</name>
    <edit_key>2014-11-13</edit_key>
    <is_finished>1</is_finished>
    <draft_status>postdraft</draft_status>
  </league>
</fantasy_content>`
	var fc FantasyContent
	if err := xml.Unmarshal([]byte(raw), &fc); err != nil {
		t.Fatalf("unmarshal league: %v", err)
	}
	if fc.League.LeagueKey != "428.l.36687" || fc.League.IsFinished != 1 {
		t.Fatalf("unexpected league: %+v", fc.League)
	}
}

func TestUnmarshalMatchupEmptyWinProbability(t *testing.T) {
	const raw = `<?xml version="1.0" encoding="UTF-8"?>
<fantasy_content xmlns="http://fantasysports.yahooapis.com/fantasy/v2/base.rng">
  <team>
    <team_key>428.l.36687.t.15</team_key>
    <matchups count="1">
      <matchup>
        <week>1</week>
        <status>postevent</status>
        <teams count="2">
          <team>
            <team_key>428.l.36687.t.15</team_key>
            <win_probability></win_probability>
            <team_points><total>100.5</total></team_points>
          </team>
          <team>
            <team_key>428.l.36687.t.16</team_key>
            <team_points><total>90.0</total></team_points>
          </team>
        </teams>
      </matchup>
    </matchups>
  </team>
</fantasy_content>`
	var fc FantasyContent
	if err := xml.Unmarshal([]byte(raw), &fc); err != nil {
		t.Fatalf("unmarshal matchups: %v", err)
	}
	if len(fc.Team.Matchups) != 1 {
		t.Fatalf("expected 1 matchup, got %d", len(fc.Team.Matchups))
	}
}
