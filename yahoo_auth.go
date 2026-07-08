package gofantasy

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/yahoo"
	"io"
	"os"
	"time"
)

type yahooOAuth2 struct {
	config       *oauth2.Config
	state        string
	codeVerifier string
	token        oauth2.Token
}

type IYahooOAuth2 interface {
	OAuth2(clientID, clientSecret, redirectURL string) IYahooOAuth2
	GetAuthCodeUrl() (string, error)
	GetAccessToken(code string) error
	SaveToken(path string) error
	LoadAccessToken(path string) (*oauth2.Token, error)
	// Token returns the currently held token (may be the zero value if none
	// has been loaded/exchanged yet).
	Token() *oauth2.Token
	// SetToken installs an already-obtained token (e.g. loaded from a
	// database rather than a file) for subsequent use/refresh.
	SetToken(t *oauth2.Token)
	// RefreshIfNeeded refreshes the held token when it is expired, using
	// YAHOO_CLIENT_ID / YAHOO_REDIRECT_URL from the environment, and returns
	// the (possibly unchanged) current token.
	RefreshIfNeeded() (*oauth2.Token, error)
}

var _ IYahooOAuth2 = &yahooOAuth2{}

func (y *yahooOAuth2) RefreshToken() error {
	redirectURL := os.Getenv("YAHOO_REDIRECT_URL")
	clientID := os.Getenv("YAHOO_CLIENT_ID")
	clientSecret := os.Getenv("YAHOO_CLIENT_SECRET")
	oc := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"openid"},
		Endpoint:     yahoo.Endpoint,
	}
	token := &oauth2.Token{
		RefreshToken: y.token.RefreshToken,
	}
	// Create a context
	ctx := context.Background()

	// Create a TokenSource using the configuration and the provided token
	tokenSource := oc.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		fmt.Printf("Error refreshing token: %v\n", err)
		return err
	}
	y.token = *newToken
	return nil
}

func (y *yahooOAuth2) LoadAccessToken(path string) (*oauth2.Token, error) {
	if err := readToken(path, &y.token); err != nil {
		return nil, err
	}
	if y.token.Expiry.UTC().Before(time.Now().UTC()) {
		if err := y.RefreshToken(); err != nil {
			return nil, err
		}
	}
	return &y.token, nil
}

// Token returns the currently held token.
func (y *yahooOAuth2) Token() *oauth2.Token {
	return &y.token
}

// SetToken installs an already-obtained token for subsequent use/refresh.
func (y *yahooOAuth2) SetToken(t *oauth2.Token) {
	if t != nil {
		y.token = *t
	}
}

// RefreshIfNeeded refreshes the held token when it is expired.
func (y *yahooOAuth2) RefreshIfNeeded() (*oauth2.Token, error) {
	if y.token.Expiry.IsZero() || y.token.Expiry.UTC().After(time.Now().UTC()) {
		return &y.token, nil
	}
	if err := y.RefreshToken(); err != nil {
		return nil, err
	}
	return &y.token, nil
}

func (y *yahooOAuth2) SaveToken(path string) error {
	return saveToken(&y.token, path)
}

func (y *yahooOAuth2) GetAccessToken(code string) error {
	ctx := context.Background()
	o := oauth2.SetAuthURLParam("code_verifier", y.codeVerifier)
	token, err := y.getAccessToken(ctx, code, o)
	if err != nil {
		return err
	}
	y.token = *token
	return nil
}

func (y *yahooOAuth2) getAccessToken(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	token, err := y.config.Exchange(ctx, code, opts...)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (y *yahooOAuth2) GetAuthCodeUrl() (string, error) {
	sha2 := sha256.New()
	io.WriteString(sha2, y.codeVerifier)
	codeChallenge := base64.RawURLEncoding.EncodeToString(sha2.Sum(nil))
	o1 := oauth2.SetAuthURLParam("code_challenge_method", "S256")
	o2 := oauth2.SetAuthURLParam("code_challenge", codeChallenge)
	authCodeUrl := y.config.AuthCodeURL(y.state, o1, o2)
	return authCodeUrl, nil
}

func (y *yahooOAuth2) OAuth2(clientID, clientSecret, redirectURL string) IYahooOAuth2 {
	codeVerifier, err := randomBytesInHex(32) // 64 character string here
	if err != nil {
		return nil
	}
	return &yahooOAuth2{
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
}

func saveToken(token *oauth2.Token, path string) error {
	if path == "" {
		path = os.Getenv("HOME") + YahooTokenPath
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
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("opening token file %s: %w", path, err)
	}
	defer f.Close()
	byteValue, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("reading token file %s: %w", path, err)
	}
	if err := json.Unmarshal(byteValue, t); err != nil {
		return fmt.Errorf("parsing token file %s: %w", path, err)
	}
	return nil
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
