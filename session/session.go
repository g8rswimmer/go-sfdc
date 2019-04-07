// Package session provides handles creation of a Salesforce session
package session

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/g8rswimmer/goforce"
	"github.com/g8rswimmer/goforce/credentials"
)

// Session is the authentication response.  This is used to generate the
// authroization header for the Salesforce API calls.
type Session struct {
	response *sessionPasswordResponse
	config   goforce.Configuration
}

// Formatter is the session interface that formats the
// session information for use in the APIs.
//
// ServiceURL will return the Salesforce instance for the
// service URL.
//
// AuthorizationHeader will add the authorization to the
// HTTP request's header.
//
// Client returns the HTTP client to be used in APIs calls.
type Formatter interface {
	ServiceURL() string
	AuthorizationHeader(*http.Request)
	Client() *http.Client
}
type sessionPasswordResponse struct {
	AccessToken string `json:"access_token"`
	InstanceURL string `json:"instance_url"`
	ID          string `json:"id"`
	TokenType   string `json:"token_type"`
	IssuedAt    string `json:"issued_at"`
	Signature   string `json:"signature"`
}

const oauthEndpoint = "/services/oauth2/token"

// OpenSession is used to authenticate with Salesforce and open a session.  The user will need to
// supply the proper credentails and a HTTP client.
func OpenSession(config goforce.Configuration) (*Session, error) {

	request, err := passwordSessionRequest(config.Credentials)

	if err != nil {
		return nil, err
	}

	response, err := passwordSessionResponse(request, config.Client)
	if err != nil {
		return nil, err
	}

	session := &Session{
		response: response,
		config:   config,
	}

	return session, nil
}

func passwordSessionRequest(creds *credentials.Credentials) (*http.Request, error) {

	oauthURL := creds.URL() + oauthEndpoint

	body, err := creds.Retrieve()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, oauthURL, body)

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

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("session response error: %d %s", response.StatusCode, response.Status)
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

// ServiceURL will return the Salesforce instance for the
// service URL.
func (session *Session) ServiceURL() string {
	return fmt.Sprintf("%s/services/data/v%d.0", session.response.InstanceURL, session.config.Version)
}

// AuthorizationHeader will add the authorization to the
// HTTP request's header.
func (session *Session) AuthorizationHeader(request *http.Request) {
	auth := fmt.Sprintf("%s %s", session.response.TokenType, session.response.AccessToken)
	request.Header.Add("Authorization", auth)
}

// Client returns the HTTP client to be used in APIs calls.
func (session *Session) Client() *http.Client {
	return session.config.Client
}
