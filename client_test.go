package gofantasy

import (
	"encoding/xml"
	"strings"
	"testing"
)

func TestIsValidGameKeys(t *testing.T) {
	if !isValidGameKeys("mlb", "nba") {
		t.Fatalf("expected mlb,nba to be valid")
	}
	if isValidGameKeys("mlb", "cricket") {
		t.Fatalf("expected cricket to be invalid")
	}
}

func TestMD5HashDeterministic(t *testing.T) {
	a := md5Hash("https://example.com/x")
	b := md5Hash("https://example.com/x")
	if a != b {
		t.Fatalf("md5Hash not deterministic: %s != %s", a, b)
	}
	if a == md5Hash("https://example.com/y") {
		t.Fatalf("md5Hash collision for different inputs")
	}
}

func TestPlayerOptionEncode(t *testing.T) {
	f := newPlayerFilter(
		WithPlayerStatus("FA"),
		WithPlayerPosition("QB"),
	)
	got := f.encode()
	if !strings.Contains(got, ";status=FA") || !strings.Contains(got, ";position=QB") {
		t.Fatalf("unexpected player filter encoding: %q", got)
	}
}

func TestAddDropPayloadXML(t *testing.T) {
	bid := 7
	payload := transactionPayload{
		Transaction: transactionBody{
			Type:    "add/drop",
			FAABBid: &bid,
			Players: &txnPlayers{Player: []txnPlayer{
				{PlayerKey: "422.p.1", TransactionData: txnPlayerTx{Type: "add", DestinationTeamKey: "422.l.1.t.2"}},
				{PlayerKey: "422.p.9", TransactionData: txnPlayerTx{Type: "drop", SourceTeamKey: "422.l.1.t.2"}},
			}},
		},
	}
	out, err := xml.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	for _, want := range []string{
		"<fantasy_content>", "<type>add/drop</type>", "<faab_bid>7</faab_bid>",
		"<destination_team_key>422.l.1.t.2</destination_team_key>",
		"<source_team_key>422.l.1.t.2</source_team_key>",
	} {
		if !strings.Contains(s, want) {
			t.Fatalf("payload missing %q; got:\n%s", want, s)
		}
	}
}

func TestSingleAddPayloadOmitsPlayersWrapper(t *testing.T) {
	payload := transactionPayload{
		Transaction: transactionBody{
			Type:   "add",
			Player: &txnPlayer{PlayerKey: "422.p.1", TransactionData: txnPlayerTx{Type: "add", DestinationTeamKey: "422.l.1.t.2"}},
		},
	}
	out, err := xml.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(out), "<players>") {
		t.Fatalf("single add should not emit a <players> wrapper; got:\n%s", string(out))
	}
}
