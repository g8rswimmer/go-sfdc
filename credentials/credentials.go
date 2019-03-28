package credentials

import "io"

// PasswordCredentails is a structure for the OAuth credentials
// that are needed to authenticate with a Salesforce org.
//
// URL is the login URL used, either https://test.salesforce.com or https://login.salesforce.com
//
// Username is the Salesforce user name for logging into the org.
//
// Password is the Salesforce password for the user.
//
// ClientID is the client ID from the connected application.
//
// ClientSecret is the client secret from the connected application.
type PasswordCredentails struct {
	URL          string
	Username     string
	Password     string
	ClientID     string
	ClientSecret string
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
func NewCredentials(provider Provider) *Credentials {
	return &Credentials{
		provider: provider,
	}
}

// NewPasswordCredentials will create a crendential with the password credentials.
func NewPasswordCredentials(creds PasswordCredentails) *Credentials {
	return &Credentials{
		provider: &passwordProvider{
			creds: creds,
		},
	}
}
