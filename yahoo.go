package gofantasy

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/yahoo"
	"io"
	"net/http"
	"os"
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

type yahooOAuth2 struct {
	config       *oauth2.Config
	state        string
	codeVerifier string
}

type IYahooClient interface {
	WithAccessToken(accessToken string) (IYahooClient, error)
	OAuth2(clientID, clientSecret, redirectURL string) IYahooClient
	GetGame(ctx context.Context, gameKey string) (*Game, error)
	GetLeague(ctx context.Context, leagueID string) (*League, error)
	GetLeagueSettings(ctx context.Context, leagueID string) (*League, error)
	GetAuthCodeUrl() (string, error)
	GetAccessToken(code string) IYahooClient
	SaveToken(path string) error
}

type yahooClient struct {
	baseUrl     string
	baseClient  *client
	yahooOAuth2 *yahooOAuth2
	token       oauth2.Token
}

var _ IYahooClient = &yahooClient{}

func (y *yahooClient) SaveToken(path string) error {
	return saveToken(&y.token, path)
}

func (y *yahooClient) GetAccessToken(code string) IYahooClient {
	ctx := context.Background()
	o := oauth2.SetAuthURLParam("code_verifier", y.yahooOAuth2.codeVerifier)
	token, err := y.yahooOAuth2.config.Exchange(ctx, code, o)
	if err != nil {
		fmt.Printf("Error authorizing token: %s\n", err)
		return nil
	}
	y.token = *token
	return y
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

func (y *yahooClient) WithAccessToken(accessToken string) (IYahooClient, error) {
	if accessToken == "" {
		err := readToken("", &y.token)
		if err != nil {
			return nil, err
		}
		accessToken = y.token.AccessToken
	}
	y.baseClient.requestor.authorizationDecorator = func(req *http.Request) *http.Request {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
		return req
	}
	return y, nil
}

func (y *yahooClient) OAuth2(clientID, clientSecret, redirectURL string) IYahooClient {
	codeVerifier, err := randomBytesInHex(32) // 64 character string here
	if err != nil {
		return nil
	}
	yo2 := &yahooOAuth2{
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
	y.yahooOAuth2 = yo2
	return y
}

func (y *yahooClient) GetGame(ctx context.Context, gameKey string) (*Game, error) {
	cacheKey := fmt.Sprintf("game%s", gameKey)
	endpoint := fmt.Sprintf("%s/game/%s", y.baseUrl, gameKey)
	fc, err := y.getCachedOrFetch(ctx, cacheKey, endpoint, "game")
	if err != nil {
		return nil, err
	}
	return &fc.Game, nil
}

func (y *yahooClient) GetLeague(ctx context.Context, leagueID string) (*League, error) {
	cacheKey := fmt.Sprintf("league%s", leagueID)
	endpoint := fmt.Sprintf("%s/league/%s", y.baseUrl, leagueID)
	fc, err := y.getCachedOrFetch(ctx, cacheKey, endpoint, "league")
	if err != nil {
		return nil, err
	}
	return &fc.League, nil
}

func (y *yahooClient) GetLeagueSettings(ctx context.Context, leagueID string) (*League, error) {
	//TODO implement me
	panic("implement me")
}

func (y *yahooClient) getCachedOrFetch(ctx context.Context, cacheKey, endpoint string, objType string) (*FantasyContent, error) {
	var fc FantasyContent
	if y.baseClient.cache != nil {
		v, exist := y.baseClient.cache.Get(cacheKey)
		if exist {
			fmt.Println("----> cache hit!")
			return v.(*FantasyContent), nil
		}
	}
	_, err := y.baseClient.requestor.Get(ctx, endpoint, &fc)
	if err != nil {
		return nil, err
	}
	if y.baseClient.cache != nil {
		y.baseClient.cache.Add(cacheKey, &fc)
	}
	return &fc, nil
}

// use in yahoo oauth2
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

func saveToken(token *oauth2.Token, path string) error {
	if path == "" {
		path = os.Getenv("HOME") + "/.config/gofantasy/yahoo_token.json"
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		if !os.IsExist(err) {
			err := os.Mkdir(os.Getenv("HOME")+"/.config/gofantasy", 0755)
			if err != nil {
				return err
			}
			f, err = os.Create(path)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	defer f.Close()
	t, _ := json.Marshal(token)
	_, err = f.Write(t)
	if err != nil {
		return err
	}
	return nil
}

func readToken(path string, t *oauth2.Token) error {
	if path == "" {
		path = os.Getenv("HOME") + "/.config/gofantasy/yahoo_token.json"
	}
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	byteValue, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(byteValue), t)
	if err != nil {
		panic(err)
	}
	return nil
}
