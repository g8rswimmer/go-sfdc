package credentials

import (
	"io"
	"net/url"
	"strings"
)

type refreshProvider struct {
	creds RefreshCredentials
}

func (provider *refreshProvider) Retrieve() (io.Reader, error) {
	form := url.Values{}
	form.Add("grant_type", string(refreshGrantType))
	form.Add("client_id", provider.creds.ClientID)
	form.Add("client_secret", provider.creds.ClientSecret)
	form.Add("refresh_token", provider.creds.AccessToken)

	return strings.NewReader(form.Encode()), nil
}

func (provider *refreshProvider) URL() string {
	return provider.creds.URL
}

func (provider *refreshProvider) ClientID() string {
	return provider.creds.ClientID
}

func (provider *refreshProvider) ClientSecret() string {
	return provider.creds.ClientSecret
}