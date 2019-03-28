package credentials

import (
	"io"
	"net/url"
	"strings"
)

type passwordProvider struct {
	creds PasswordCredentails
}

func (provider *passwordProvider) Retrieve() (io.Reader, error) {
	form := url.Values{}
	form.Add("grant_type", string(passwordGrantType))
	form.Add("username", provider.creds.Username)
	form.Add("password", provider.creds.Password)
	form.Add("client_id", provider.creds.ClientID)
	form.Add("client_secret", provider.creds.ClientSecret)

	return strings.NewReader(form.Encode()), nil
}

func (provider *passwordProvider) URL() string {
	return provider.creds.URL
}
