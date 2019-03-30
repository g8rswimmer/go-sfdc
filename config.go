package goforce

import (
	"net/http"

	"github.com/g8rswimmer/goforce/credentials"
)

// Configuration is the structure for goforce sessions.
//
// Credentials is the credentials that will be used to form a session.
//
// Client is the HTTP client that will be used.
//
// Version is the Salesforce version for the APIs
type Configuration struct {
	Credentials *credentials.Credentials
	Client      *http.Client
	Version     int
}
