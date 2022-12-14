package github

import (
	"log"

	"github.com/t0nyandre/git-notified/pkg/utils"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

type Provider struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	CallbackURL  string `json:"callback_url"`
	State        string
	AuthURL      string
	TokenURL     string
	ProfileURL   string
	EmailURL     string
	Config       *oauth2.Config
}

func New(ctx context.Context, clientID, clientSecret, callbackURL string, scopes ...string) *Provider {
	logger, ok := ctx.Value("logger").(*zap.SugaredLogger)
	if !ok {
		log.Fatalf("Could not get logger from context")
	}

	state, err := utils.GenerateRandomState()
	if err != nil {
		logger.Fatalw("Could not generate random state", "error", err)
	}

	provider := &Provider{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		CallbackURL:  callbackURL,
		State:        state,
		AuthURL:      "https://github.com/login/oauth/authorize",
		TokenURL:     "https://github.com/login/oauth/access_token",
		ProfileURL:   "https://api.github.com/user",
		EmailURL:     "https://api.github.com/user/emails",
	}
	provider.newConfig(scopes)

	return provider
}

func (p *Provider) newConfig(scopes []string) *oauth2.Config {
	c := &oauth2.Config{
		ClientID:     p.ClientID,
		ClientSecret: p.ClientSecret,
		RedirectURL:  p.CallbackURL,
		Scopes:       []string{},
		Endpoint: oauth2.Endpoint{
			TokenURL: p.TokenURL,
		},
	}

	for _, scope := range scopes {
		c.Scopes = append(c.Scopes, scope)
	}

	c.Endpoint.AuthURL = c.AuthCodeURL(p.State)

	return c
}

func (p *Provider) getAccessToken(code string) (string, error) {
	token, err := p.Config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return "", err
	}

	return token.AccessToken, nil
}
