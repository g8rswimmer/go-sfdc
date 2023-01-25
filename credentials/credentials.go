package credentials

import (
	"crypto/rsa"
	"errors"
	"io"
)

// PasswordCredentials is a structure for the OAuth credentials
// that are needed to authenticate with a Salesforce org.
//
// URL is the login URL used, examples would be https://test.salesforce.com or https://login.salesforce.com
//
// Username is the Salesforce user name for logging into the org.
//
// Password is the Salesforce password for the user.
//
// ClientID is the client ID from the connected application.
//
// ClientSecret is the client secret from the connected application.
type PasswordCredentials struct {
	URL          string
	Username     string
	Password     string
	ClientID     string
	ClientSecret string
}

type JwtCredentials struct {
	URL            string
	ClientId       string // the client id as defined in the connected app in SalesForce
	ClientUsername string
	ClientKey      *rsa.PrivateKey // the client RSA key uploaded for authentication in the ConnectedApp
}

// Credentials is the structure that contains all of the
// information for creating a session.
type Credentials struct {
	provider Provider
}

// Provider is the interface that is able to provide the
// session creator with all of the valid information.
//
// Retrieve will return the reader for the HTTP request body.
//
// URL is the URL base for the session endpoint.
type Provider interface {
	Retrieve() (io.Reader, error)
	URL() string
}

type grantType string

const (
	passwordGrantType grantType = "password"
	jwtGrantType      grantType = "urn:ietf:params:oauth:grant-type:jwt-bearer"
)

// Retrieve will return the reader for the HTTP request body.
func (creds *Credentials) Retrieve() (io.Reader, error) {
	return creds.provider.Retrieve()
}

// URL is the URL base for the session endpoint.
func (creds *Credentials) URL() string {
	return creds.provider.URL()
}

// NewCredentials will create a credential with the custom provider.
func NewCredentials(provider Provider) (*Credentials, error) {
	if provider == nil {
		return nil, errors.New("credentials: the provider can not be nil")
	}
	return &Credentials{
		provider: provider,
	}, nil
}

// NewPasswordCredentials will create a crendential with the password credentials.
func NewPasswordCredentials(creds PasswordCredentials) (*Credentials, error) {
	if err := validatePasswordCredentials(creds); err != nil {
		return nil, err
	}
	return &Credentials{
		provider: &passwordProvider{
			creds: creds,
		},
	}, nil
}

// NewJWTCredentials weill create a credntial with all required info about generating a JWT claims parameter
func NewJWTCredentials(creds JwtCredentials) (*Credentials, error) {
	if err := validateJWTCredentials(creds); err != nil {
		return nil, err
	}
	return &Credentials{
		provider: &jwtProvider{
			creds: creds,
		},
	}, nil
}

func validatePasswordCredentials(cred PasswordCredentials) error {
	switch {
	case len(cred.URL) == 0:
		return errors.New("credentials: password credential's URL can not be empty")
	case len(cred.Username) == 0:
		return errors.New("credentials: password credential's username can not be empty")
	case len(cred.Password) == 0:
		return errors.New("credentials: password credential's password can not be empty")
	case len(cred.ClientID) == 0:
		return errors.New("credentials: password credential's client ID can not be empty")
	case len(cred.ClientSecret) == 0:
		return errors.New("credentials: password credential's client secret can not be empty")
	}
	return nil
}

func validateJWTCredentials(cred JwtCredentials) error {
	switch {
	case len(cred.URL) == 0:
		return errors.New("URL cannot be empty")
	case cred.ClientKey == nil:
		return errors.New("client key cannot be empty")
	case len(cred.ClientUsername) == 0:
		return errors.New("client username cannot be empty")
	case len(cred.ClientId) == 0:
		return errors.New("client id cannot be empty")
	}
	return nil

}
