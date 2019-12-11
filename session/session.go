// Package session provides handles creation of a Salesforce session
package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/TuSKan/go-sfdc"
	"github.com/TuSKan/go-sfdc/credentials"
)

// Session is the authentication response.  This is used to generate the
// authroization header for the Salesforce API calls.
type Session struct {
	response *Response
	config   sfdc.Configuration
}

// Clienter interface provides the HTTP client used by the
// the resources.
type Clienter interface {
	Client() *http.Client
}

// InstanceFormatter is the session interface that
// formaters the session instance information used
// by the resources.
//
// InstanceURL will return the Salesforce instance.
//
// AuthorizationHeader will add the authorization to the
// HTTP request's header.
type InstanceFormatter interface {
	InstanceURL() string
	AuthorizationHeader(*http.Request)
	Clienter
}

// ServiceFormatter is the session interface that
// formats the session for service resources.
//
// ServiceURL provides the service URL for resources to
// user.
type ServiceFormatter interface {
	InstanceFormatter
	ServiceURL() string
}

// Response ...
type Response struct {
	AccessToken string `json:"access_token"`
	InstanceURL string `json:"instance_url"`
	ID          string `json:"id"`
	TokenType   string `json:"token_type"`
	IssuedAt    string `json:"issued_at"`
	Signature   string `json:"signature"`
}

const oauthEndpoint = "/services/oauth2/token"

// Open is used to authenticate with Salesforce and open a session.  The user will need to
// supply the proper credentials and a HTTP client.
func Open(config sfdc.Configuration) (*Session, error) {
	if config.Credentials == nil {
		return nil, errors.New("session: configuration crendentials can not be nil")
	}
	if config.Client == nil {
		return nil, errors.New("session: configuration client can not be nil")
	}
	if config.Version <= 0 {
		return nil, errors.New("session: configuration version can not be less than zero")
	}
	request, err := sessionRequest(config.Credentials)

	if err != nil {
		return nil, err
	}

	response, err := sessionResponse(request, config.Client)
	if err != nil {
		return nil, err
	}

	session := &Session{
		response: response,
		config:   config,
	}

	return session, nil
}

// Refresh - refresh session
func Refresh(session *Session) (*Session, error) {
	rcreds := credentials.RefreshCredentials{
		URL:          session.config.Credentials.URL(),
		ClientID:     session.config.Credentials.ClientID(),
		ClientSecret: session.config.Credentials.ClientSecret(),
		AccessToken:  session.response.AccessToken,
	}

	newcred, err := credentials.NewRefreshCredentials(rcreds)
	if err != nil {
		return nil, err
	}

	request, err := sessionRequest(newcred)
	if err != nil {
		return nil, err
	}

	response, err := sessionResponse(request, session.config.Client)
	if err != nil {
		return nil, err
	}

	session = &Session{
		response: response,
		config:   session.config,
	}

	return session, nil
}

// IsValid ...
func IsValid(session *Session) (bool, error) {

	request, err := http.NewRequest(http.MethodGet, session.ServiceURL(), nil)

	if err != nil {
		return false, err
	}

	session.AuthorizationHeader(request)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Accept", "application/json")

	response, err := session.Client().Do(request)
	if err != nil {
		return false, err
	}

	if response.StatusCode != http.StatusOK {
		return false, fmt.Errorf("session response error: %d %s", response.StatusCode, response.Status)
	}
	decoder := json.NewDecoder(response.Body)
	defer response.Body.Close()

	var sessionResponse Response
	err = decoder.Decode(&sessionResponse)
	if err != nil {
		return false, err
	}

	return true, nil
}

func sessionRequest(creds *credentials.Credentials) (*http.Request, error) {

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

func sessionResponse(request *http.Request, client *http.Client) (*Response, error) {
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("session response error: %d %s", response.StatusCode, response.Status)
	}
	decoder := json.NewDecoder(response.Body)
	defer response.Body.Close()

	var sessionResponse Response
	err = decoder.Decode(&sessionResponse)
	if err != nil {
		return nil, err
	}

	return &sessionResponse, nil
}

// InstanceURL will retuern the Salesforce instance
// from the session authentication.
func (session *Session) InstanceURL() string {
	return session.response.InstanceURL
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
