package gofantasy

const (
	// YahooBaseURL is the base URL for all calls to Yahoo's fantasy sports API
	YahooBaseURL   = "https://fantasysports.yahooapis.com/fantasy/v2"
	YahooTokenPath = "/.config/gofantasy/yahoo_token.json"

	// EspnBaseURL is the read API base for ESPN Fantasy v3 (2018+ seasons).
	EspnBaseURL = "https://lm-api-reads.fantasy.espn.com/apis/v3/games"
)

// EspnGameCodes maps common sport aliases to ESPN game path segments.
var EspnGameCodes = map[string]string{
	"nfl": "ffl",
	"nba": "fba",
	"mlb": "flb",
	"nhl": "nhl",
	"ffl": "ffl",
	"fba": "fba",
	"flb": "flb",
}
