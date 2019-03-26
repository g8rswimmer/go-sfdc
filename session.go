package goforce

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

type grantType string

type SessionPasswordCredentials struct {
	URL          string
	Username     string
	Password     string
	ClientID     string
	ClientSecret string
}

type Session struct {
	response *sessionPasswordResponse
}

type sessionPasswordResponse struct {
	AccessToken string `json:"access_token"`
	InstanceURL string `json:"instance_url"`
	ID          string `json:"id"`
	TokenType   string `json:"token_type"`
	IssuedAd    string `json:"issued_at"`
	Signature   string `json:"signature"`
}

const (
	passwordGrantType grantType = "password"
)
const oauthService = "/services/oauth2/token"

func NewPasswordSession(credentials SessionPasswordCredentials, client *http.Client) (*Session, error) {

	request, err := passwordSessionRequest(credentials)

	if err != nil {
		return nil, err
	}

	response, err := passwordSessionResponse(request, client)
	if err != nil {
		return nil, err
	}

	session := &Session{
		response: response,
	}

	return session, nil
}

func passwordSessionRequest(credentials SessionPasswordCredentials) (*http.Request, error) {
	form := url.Values{}
	form.Add("grant_type", string(passwordGrantType))
	form.Add("username", credentials.Username)
	form.Add("password", credentials.Password)
	form.Add("client_id", credentials.ClientID)
	form.Add("client_secret", credentials.ClientSecret)

	oauthURL := credentials.URL + oauthService

	request, err := http.NewRequest(http.MethodPost, oauthURL, strings.NewReader(form.Encode()))

	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Accept", "application/json")
	return request, nil
}

func passwordSessionResponse(request *http.Request, client *http.Client) (*sessionPasswordResponse, error) {
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(response.Body)
	defer response.Body.Close()

	var sessionResponse sessionPasswordResponse
	err = decoder.Decode(&sessionResponse)
	if err != nil {
		return nil, err
	}

	return &sessionResponse, nil
}
