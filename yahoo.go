package gofantasy

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/yahoo"
	"io"
	"net/http"
	"time"
)

// FantasyContent is the root level response containing the data from a request
// to the fantasy sports API.
type FantasyContent struct {
	XMLName xml.Name `xml:"fantasy_content"`
	League  League   `xml:"league"`
	Game    Game     `xml:"game"`
}

type Game struct {
	GameKey string `xml:"game_key"`
	GameID  string `xml:"game_id"`
	Name    string `xml:"name"`
	Code    string `xml:"code"`
	Type    string `xml:"type"`
	URL     string `xml:"url"`
	Season  string `xml:"season"`
}

type League struct {
	LeagueKey             string   `xml:"league_key"`
	LeagueID              string   `xml:"league_id"`
	Name                  string   `xml:"name"`
	URL                   string   `xml:"url"`
	DraftStatus           string   `xml:"draft_status"`
	NumTeams              int      `xml:"num_teams"`
	EditKey               int      `xml:"edit_key"`
	WeeklyDeadline        string   `xml:"weekly_deadline"`
	LeagueUpdateTimestamp int64    `xml:"league_update_timestamp"`
	ScoringType           string   `xml:"scoring_type"`
	CurrentWeek           int      `xml:"current_week"`
	StartWeek             int      `xml:"start_week"`
	EndWeek               int      `xml:"end_week"`
	GameCode              string   `xml:"game_code"`
	IsFinished            int      `xml:"is_finished"`
	Season                int      `xml:"season"`
	Settings              Settings `xml:"settings"`
}

type RosterPosition struct {
	Position string `xml:"position"`
	Count    int    `xml:"count"`
}

type Settings struct {
	DraftType               string           `xml:"draft_type"`
	ScoringType             string           `xml:"scoring_type"`
	UsesPlayoff             bool             `xml:"uses_playoff"`
	PlayoffStartWeek        int              `xml:"playoff_start_week"`
	UsesPlayoffReseeding    bool             `xml:"uses_playoff_reseeding"`
	UsesLockEliminatedTeams bool             `xml:"uses_lock_eliminated_teams"`
	UsesFAAB                bool             `xml:"uses_faab"`
	TradeEndDate            time.Time        `xml:"trade_end_date"`
	TradeRatifyType         string           `xml:"trade_ratify_type"`
	TradeRejectTime         int              `xml:"trade_reject_time"`
	RosterPositions         []RosterPosition `xml:"roster_positions>roster_position"`
	StatCategories          []Stat           `xml:"stat_categories>stats>stat"`
	StatModifiers           []Stat           `xml:"stat_modifiers>stats>stat"`
	Divisions               []Division       `xml:"divisions>division"`
}

type Stat struct {
	StatID       int    `xml:"stat_id"`
	Enabled      int    `xml:"enabled"`
	Name         string `xml:"name"`
	DisplayName  string `xml:"display_name"`
	SortOrder    int    `xml:"sort_order"`
	PositionType string `xml:"position_type"`
	Value        int    `xml:"value"`
}

type Division struct {
	DivisionID int    `xml:"division_id"`
	Name       string `xml:"name"`
}

const (
	// YahooBaseURL is the base URL for all calls to Yahoo's fantasy sports API
	YahooBaseURL = "https://fantasysports.yahooapis.com/fantasy/v2"
)

type YahooOAuth2 struct {
	config       *oauth2.Config
	state        string
	codeVerifier string
}

func generateState() string {
	b := make([]byte, 128)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	state := base64.URLEncoding.EncodeToString(b)
	return state
}

func randomBytesInHex(count int) (string, error) {
	buf := make([]byte, count)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return "", fmt.Errorf(" Could not generate %d random bytes: %v", count, err)
	}

	return hex.EncodeToString(buf), nil
}

type IYahooClient interface {
	WithAccessToken(accessToken string) IYahooClient
	WithOAuth2(clientID, clientSecret, redirectURL string) IYahooClient
	GetGame(ctx context.Context, gameKey string) (*Game, error)
	GetLeague(ctx context.Context, leagueID string) (*League, error)
	GetLeagueSettings(ctx context.Context, leagueID string) (*League, error)
	GetAuthCodeUrl() (string, error)
	GetAccessToken(code string) (*oauth2.Token, error)
}

type yahooClient struct {
	baseUrl     string
	requestor   *requestor
	yahooOAuth2 YahooOAuth2
}

var _ IYahooClient = &yahooClient{}

func (y *yahooClient) GetAccessToken(code string) (*oauth2.Token, error) {
	ctx := context.Background()
	o := oauth2.SetAuthURLParam("code_verifier", y.yahooOAuth2.codeVerifier)
	token, err := y.yahooOAuth2.config.Exchange(ctx, code, o)
	if err != nil {
		return nil, fmt.Errorf("Error authorizing token: %s\n", err)
	}
	return token, nil
}

func (y *yahooClient) GetAuthCodeUrl() (string, error) {
	sha2 := sha256.New()
	io.WriteString(sha2, y.yahooOAuth2.codeVerifier)
	codeChallenge := base64.RawURLEncoding.EncodeToString(sha2.Sum(nil))
	o1 := oauth2.SetAuthURLParam("code_challenge_method", "S256")
	o2 := oauth2.SetAuthURLParam("code_challenge", codeChallenge)
	authCodeUrl := y.yahooOAuth2.config.AuthCodeURL(y.yahooOAuth2.state, o1, o2)
	return authCodeUrl, nil
}

func (y *yahooClient) WithAccessToken(accessToken string) IYahooClient {
	y.requestor.authorizationDecorator = func(req *http.Request) *http.Request {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
		return req
	}
	return y
}

func (y *yahooClient) WithOAuth2(clientID, clientSecret, redirectURL string) IYahooClient {
	codeVerifier, err := randomBytesInHex(32) // 64 character string here
	if err != nil {
		return nil
	}
	yo2 := &YahooOAuth2{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"openid"},
			Endpoint:     yahoo.Endpoint,
		},
		state:        generateState(),
		codeVerifier: codeVerifier,
	}
	y.yahooOAuth2 = *yo2
	return y
}

func (y *yahooClient) GetGame(ctx context.Context, gameKey string) (*Game, error) {
	var fc FantasyContent
	_, err := y.requestor.Get(ctx, fmt.Sprintf("%s/game/%s", y.baseUrl, gameKey), &fc)
	if err != nil {
		return nil, err
	}
	return &fc.Game, nil
}

func (y *yahooClient) GetLeague(ctx context.Context, leagueID string) (*League, error) {
	panic("implement me")
}

func (y *yahooClient) GetLeagueSettings(ctx context.Context, leagueID string) (*League, error) {
	//TODO implement me
	panic("implement me")
}