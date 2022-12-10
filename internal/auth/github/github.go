package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/t0nyandre/git-notified/pkg/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type Github struct {
	state       string
	oauthConfig *oauth2.Config
}

type githubRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
}

type githubAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func NewGithub() *Github {
	state, err := utils.GenerateRandomState()
	if err != nil {
		log.Fatalf("Could not generate random state: %v", err)
	}

	return &Github{
		state: state,
		oauthConfig: &oauth2.Config{
			ClientID:     os.Getenv("GITHUB_CLIENTID"),
			ClientSecret: os.Getenv("GITHUB_CLIENTSECRET"),
			RedirectURL:  "http://localhost:4000/auth/github/callback",
			Scopes:       []string{"repo:status", "read:user", "user:email"},
			Endpoint:     github.Endpoint,
		},
	}
}

func (oauth *Github) GithubLogin(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=%s&state=%s",
		oauth.oauthConfig.ClientID,
		oauth.oauthConfig.RedirectURL,
		oauth.oauthConfig.Scopes,
		oauth.state)
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

func getAccessToken(input []byte) (*githubAccessTokenResponse, error) {
	req, err := http.NewRequest("POST", github.Endpoint.TokenURL, bytes.NewBuffer(input))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var githubResponse githubAccessTokenResponse

	if err := json.Unmarshal(body, &githubResponse); err != nil {
		return nil, err
	}

	return &githubResponse, nil
}

func (oauth *Github) GithubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	request := githubRequest{
		ClientID:     oauth.oauthConfig.ClientID,
		ClientSecret: oauth.oauthConfig.ClientSecret,
		Code:         code,
	}
	requestJSON, _ := json.Marshal(request)
	oauthData, err := getAccessToken(requestJSON)
	if err != nil {
		log.Fatalf("Could not get access token: %v", err)
	}

	// TODO: Save access token to database and save it to session (cookie)
	// Make a route for getting information from api.github.com
	req, err := http.NewRequest("GET", "https://api.github.com/search/commits?q=author:t0nyandre&per_page=1&sort=committer-date&order=desc", nil)
	if err != nil {
		log.Panicf("API request failed: %v", err)
	}

	authHeaderValue := fmt.Sprintf("token %s", oauthData.AccessToken)
	req.Header.Set("Authorization", authHeaderValue)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panicf("Request failed: %v", err)
	}

	body, _ := ioutil.ReadAll(res.Body)

	fmt.Fprintf(w, string(body))
}
